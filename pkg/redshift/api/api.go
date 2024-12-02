package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/grafana/redshift-datasource/pkg/redshift/api/types"
	"strings"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	redshiftV2 "github.com/aws/aws-sdk-go-v2/service/redshift"
	redshiftdataV2 "github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshifttypesV2 "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	redshiftserverlessV2 "github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	secretsmanagerV2 "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagerV2types "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	awsModels "github.com/grafana/grafana-aws-sdk/pkg/sql/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
	"github.com/grafana/sqlds/v4"
)

type API struct {
	DataClient                 types.RedshiftDataClient
	SecretsClient              types.RedshiftSecretsClient
	ManagementClient           types.RedshiftManagementClient
	ServerlessManagementClient types.ServlessAPIClient
	settings                   *models.RedshiftDataSourceSettings
}

func New(ctx context.Context, sessionCache *awsds.SessionCache, settings awsModels.Settings) (api.AWSAPI, error) {
	redshiftSettings := settings.(*models.RedshiftDataSourceSettings)

	httpClientProvider := sdkhttpclient.NewProvider()
	// TODO: Context needs to be added, see https://github.com/grafana/oss-plugin-partnerships/issues/648
	httpClientOptions, err := redshiftSettings.Config.HTTPClientOptions(ctx)
	if err != nil {
		backend.Logger.Error("failed to create HTTP client options", "error", err.Error())
		return nil, err
	}
	httpClient, err := httpClientProvider.New(httpClientOptions)
	if err != nil {
		backend.Logger.Error("failed to create HTTP client", "error", err.Error())
		return nil, err
	}

	authSettings := awsds.ReadAuthSettings(ctx)
	sess, err := sessionCache.GetSessionWithAuthSettings(awsds.GetSessionConfig{
		Settings:      redshiftSettings.AWSDatasourceSettings,
		HTTPClient:    httpClient,
		UserAgentName: awsV2.String("Redshift"),
	}, *authSettings)
	if err != nil {
		return nil, err
	}

	provider := &SessionCredentialsProvider{sess}
	return &API{
		DataClient:                 redshiftdataV2.New(redshiftdataV2.Options{Credentials: provider}),
		SecretsClient:              secretsmanagerV2.New(secretsmanagerV2.Options{Credentials: provider}),
		ManagementClient:           redshiftV2.New(redshiftV2.Options{Credentials: provider}),
		ServerlessManagementClient: redshiftserverlessV2.New(redshiftserverlessV2.Options{Credentials: provider}),
		settings:                   redshiftSettings,
	}, nil
}

type SessionCredentialsProvider struct {
	session *session.Session
}

func (scp *SessionCredentialsProvider) Retrieve(_ context.Context) (awsV2.Credentials, error) {
	creds := awsV2.Credentials{}
	v1creds, err := scp.session.Config.Credentials.Get()
	if err != nil {
		return creds, err
	}
	creds.AccessKeyID = v1creds.AccessKeyID
	creds.SecretAccessKey = v1creds.SecretAccessKey
	creds.SessionToken = v1creds.SessionToken
	return creds, nil
}

type apiInput struct {
	ClusterIdentifier *string
	WorkgroupName     *string
	Database          *string
	DbUser            *string
	SecretARN         *string
}

func (c *API) apiInput() apiInput {
	res := apiInput{
		Database: awsV2.String(c.settings.Database),
	}
	switch {
	// Serverless + Temporary credential
	case c.settings.UseServerless && !c.settings.UseManagedSecret:
		res.WorkgroupName = awsV2.String(c.settings.WorkgroupName)
	// Serverless + Managed Secret
	case c.settings.UseServerless && c.settings.UseManagedSecret:
		res.WorkgroupName = awsV2.String(c.settings.WorkgroupName)
		res.SecretARN = awsV2.String(c.settings.ManagedSecret.ARN)
	// Provisioned + Temporary credential
	case !c.settings.UseServerless && !c.settings.UseManagedSecret:
		res.ClusterIdentifier = awsV2.String(c.settings.ClusterIdentifier)
		res.DbUser = awsV2.String(c.settings.DBUser)
	// Provisioned + Managed Secret
	case !c.settings.UseServerless && c.settings.UseManagedSecret:
		res.ClusterIdentifier = awsV2.String(c.settings.ClusterIdentifier)
		res.SecretARN = awsV2.String(c.settings.ManagedSecret.ARN)
	}
	return res
}

func (c *API) Execute(ctx context.Context, input *api.ExecuteQueryInput) (*api.ExecuteQueryOutput, error) {
	commonInput := c.apiInput()
	redshiftInput := &redshiftdataV2.ExecuteStatementInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Sql:               awsV2.String(input.Query),
		WithEvent:         awsV2.Bool(c.settings.WithEvent),
		WorkgroupName:     commonInput.WorkgroupName,
	}
	output, err := c.DataClient.ExecuteStatement(ctx, redshiftInput)
	if err != nil {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.ExecuteError, err), false)
	}

	return &api.ExecuteQueryOutput{ID: *output.Id}, nil
}

// GetQueryID always returns not found. To actually check if the query has been called requires calling ListStatements, which can lead to timeouts
// when there are many statements to page through
func (c *API) GetQueryID(_ context.Context, _ string, _ ...interface{}) (bool, string, error) {
	return false, "", nil
}

func (c *API) Status(ctx context.Context, output *api.ExecuteQueryOutput) (*api.ExecuteQueryStatus, error) {
	statusResp, err := c.DataClient.DescribeStatement(ctx, &redshiftdataV2.DescribeStatementInput{
		Id: awsV2.String(output.ID),
	})
	if err != nil {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.StatusError, err), false)
	}

	if statusResp.Error != nil && *statusResp.Error != "" {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.ExecuteError, *statusResp.Error), false)
	}

	var finished bool
	switch statusResp.Status {
	case redshifttypesV2.StatusStringFailed,
		redshifttypesV2.StatusStringAborted,
		redshifttypesV2.StatusStringFinished:
		finished = true
	}

	return &api.ExecuteQueryStatus{
		ID:       output.ID,
		State:    string(statusResp.Status),
		Finished: finished,
	}, nil
}

func (c *API) CancelQuery(_ context.Context, _ sqlds.Options, queryID string) error {
	return c.Stop(&api.ExecuteQueryOutput{ID: queryID})
}

func (c *API) Stop(output *api.ExecuteQueryOutput) error {
	_, err := c.DataClient.CancelStatement(context.TODO(), &redshiftdataV2.CancelStatementInput{
		Id: &output.ID,
	})
	// ignore finished query error
	if err != nil && !strings.Contains(err.Error(), "Could not cancel a query that is already in FINISHED state") {
		return errorsource.DownstreamError(fmt.Errorf("%w: %v", err, api.StopError), false)
	}
	return nil
}

func (c *API) Regions(context.Context) ([]string, error) {
	// This is not used. If regions are out of date, update them in the @grafana/aws-sdk-react package
	return []string{}, nil
}

func (c *API) Databases(ctx context.Context, _ sqlds.Options) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataV2.ListDatabasesInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.ListDatabases(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		res = append(res, out.Databases...)
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Schemas(ctx context.Context, _ sqlds.Options) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataV2.ListSchemasInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.ListSchemas(ctx, input)
		if err != nil {
			return nil, err
		}
		res = append(res, out.Schemas...)
		input.NextToken = out.NextToken
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Tables(ctx context.Context, options sqlds.Options) ([]string, error) {
	schema := options["schema"]
	// We use the "public" schema by default if not specified
	if schema == "" {
		schema = "public"
	}
	commonInput := c.apiInput()
	input := &redshiftdataV2.ListTablesInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		SchemaPattern:     awsV2.String(schema),
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.ListTables(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, table := range out.Tables {
			if table.Name != nil {
				res = append(res, *table.Name)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Columns(ctx context.Context, options sqlds.Options) ([]string, error) {
	schema, table := options["schema"], options["table"]
	commonInput := c.apiInput()
	input := &redshiftdataV2.DescribeTableInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Schema:            awsV2.String(schema),
		Table:             awsV2.String(table),
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.DescribeTable(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, column := range out.ColumnList {
			if column.Name != nil {
				res = append(res, *column.Name)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Secrets(ctx context.Context) ([]models.ManagedSecret, error) {
	input := &secretsmanagerV2.ListSecretsInput{
		Filters: []secretsmanagerV2types.Filter{
			{
				// Only secrets with the tag RedshiftQueryOwner can be used
				// https://docs.aws.amazon.com/redshift/latest/mgmt/query-editor.html#query-cluster-configure
				Key:    secretsmanagerV2types.FilterNameStringTypeTagKey,
				Values: []string{"RedshiftQueryOwner"},
			},
		},
	}
	isFinished := false
	redshiftSecrets := []models.ManagedSecret{}
	for !isFinished {
		out, err := c.SecretsClient.ListSecrets(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		if input.NextToken == nil {
			isFinished = true
		}
		for _, secret := range out.SecretList {
			if secret.ARN == nil || secret.Name == nil {
				continue
			}
			redshiftSecrets = append(redshiftSecrets, models.ManagedSecret{
				ARN:  *secret.ARN,
				Name: *secret.Name,
			})
		}
	}
	return redshiftSecrets, nil
}

func (c *API) Secret(ctx context.Context, options sqlds.Options) (*models.RedshiftSecret, error) {
	arn := options["secretARN"]
	input := &secretsmanagerV2.GetSecretValueInput{
		SecretId: awsV2.String(arn),
	}
	out, err := c.SecretsClient.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("missing secret content")
	}
	res := &models.RedshiftSecret{}
	err = json.Unmarshal([]byte(*out.SecretString), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *API) Clusters() ([]models.RedshiftCluster, error) {
	out, err := c.ManagementClient.DescribeClusters(context.TODO(), &redshiftV2.DescribeClustersInput{})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("missing clusters content")
	}
	res := []models.RedshiftCluster{}
	for _, r := range out.Clusters {
		if r.ClusterIdentifier != nil && r.Endpoint != nil && r.Endpoint.Address != nil && r.Endpoint.Port != nil && r.DBName != nil {
			res = append(res, models.RedshiftCluster{
				ClusterIdentifier: *r.ClusterIdentifier,
				Endpoint: models.RedshiftEndpoint{
					Address: *r.Endpoint.Address,
					Port:    *r.Endpoint.Port,
				},
				Database: *r.DBName,
			})
		}
	}
	return res, nil
}

func (c *API) Workgroups() ([]models.RedshiftWorkgroup, error) {
	out, err := c.ServerlessManagementClient.ListWorkgroups(context.TODO(), &redshiftserverlessV2.ListWorkgroupsInput{})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("missing workgroups content")
	}
	res := []models.RedshiftWorkgroup{}
	for _, r := range out.Workgroups {
		if r.WorkgroupName != nil && r.Endpoint != nil && r.Endpoint.Address != nil && r.Endpoint.Port != nil {
			res = append(res, models.RedshiftWorkgroup{
				WorkgroupName: *r.WorkgroupName,
				Endpoint: models.RedshiftEndpoint{
					Address: *r.Endpoint.Address,
					Port:    *r.Endpoint.Port,
				},
			})
		}
	}
	return res, nil
}
