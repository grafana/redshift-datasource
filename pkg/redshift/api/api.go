package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-aws-sdk/pkg/awsauth"
	"github.com/grafana/redshift-datasource/pkg/redshift/api/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshiftdata"
	redshiftdatatypes "github.com/aws/aws-sdk-go-v2/service/redshiftdata/types"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	secretsmanagertypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

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
	ServerlessManagementClient types.ServerlessAPIClient
	settings                   *models.RedshiftDataSourceSettings
}

func New(ctx context.Context, settings awsModels.Settings) (api.AWSAPI, error) {
	redshiftSettings := settings.(*models.RedshiftDataSourceSettings)

	httpClientProvider := sdkhttpclient.NewProvider()
	// TODO: Context needs to be added, see https://github.com/grafana/oss-plugin-partnerships/issues/648
	httpClientOptions, err := redshiftSettings.Config.HTTPClientOptions(ctx)

	cfg := backend.GrafanaConfigFromContext(ctx)
	httpClientOptions.Middlewares = append(httpClientOptions.Middlewares, sdkhttpclient.ResponseLimitMiddleware(cfg.ResponseLimit()))

	if err != nil {
		backend.Logger.Error("failed to create HTTP client options", "error", err.Error())
		return nil, err
	}
	httpClient, err := httpClientProvider.New(httpClientOptions)
	if err != nil {
		backend.Logger.Error("failed to create HTTP client", "error", err.Error())
		return nil, err
	}

	region := redshiftSettings.Region
	if region == "" || region == "default" {
		region = redshiftSettings.DefaultRegion
	}

	awsCfg, err := awsauth.NewConfigProvider().GetConfig(ctx, awsauth.Settings{
		LegacyAuthType:     redshiftSettings.AuthType,
		AccessKey:          redshiftSettings.AccessKey,
		SecretKey:          redshiftSettings.SecretKey,
		Region:             region,
		CredentialsProfile: redshiftSettings.Profile,
		AssumeRoleARN:      redshiftSettings.AssumeRoleARN,
		Endpoint:           redshiftSettings.Endpoint,
		ExternalID:         redshiftSettings.ExternalID,
		UserAgent:          "Redshift",
		HTTPClient:         httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &API{
		DataClient:                 redshiftdata.NewFromConfig(awsCfg),
		SecretsClient:              secretsmanager.NewFromConfig(awsCfg),
		ManagementClient:           redshift.NewFromConfig(awsCfg),
		ServerlessManagementClient: redshiftserverless.NewFromConfig(awsCfg),
		settings:                   redshiftSettings,
	}, nil
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
		Database: aws.String(c.settings.Database),
	}
	switch {
	// Serverless + Temporary credential
	case c.settings.UseServerless && !c.settings.UseManagedSecret:
		res.WorkgroupName = aws.String(c.settings.WorkgroupName)
	// Serverless + Managed Secret
	case c.settings.UseServerless && c.settings.UseManagedSecret:
		res.WorkgroupName = aws.String(c.settings.WorkgroupName)
		res.SecretARN = aws.String(c.settings.ManagedSecret.ARN)
	// Provisioned + Temporary credential
	case !c.settings.UseServerless && !c.settings.UseManagedSecret:
		res.ClusterIdentifier = aws.String(c.settings.ClusterIdentifier)
		res.DbUser = aws.String(c.settings.DBUser)
	// Provisioned + Managed Secret
	case !c.settings.UseServerless && c.settings.UseManagedSecret:
		res.ClusterIdentifier = aws.String(c.settings.ClusterIdentifier)
		res.SecretARN = aws.String(c.settings.ManagedSecret.ARN)
	}
	return res
}

func (c *API) Execute(ctx context.Context, input *api.ExecuteQueryInput) (*api.ExecuteQueryOutput, error) {
	commonInput := c.apiInput()
	redshiftInput := &redshiftdata.ExecuteStatementInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Sql:               aws.String(input.Query),
		WithEvent:         aws.Bool(c.settings.WithEvent),
		WorkgroupName:     commonInput.WorkgroupName,
	}
	output, err := c.DataClient.ExecuteStatement(ctx, redshiftInput, func(options *redshiftdata.Options) {
		if c.settings.Region != "" {
			options.Region = c.settings.Region
		} else {
			options.Region = c.settings.DefaultRegion
		}
	})
	if err != nil {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.ErrorExecute, err), false)
	}

	return &api.ExecuteQueryOutput{ID: *output.Id}, nil
}

// GetQueryID always returns not found. To actually check if the query has been called requires calling ListStatements, which can lead to timeouts
// when there are many statements to page through
func (c *API) GetQueryID(_ context.Context, _ string, _ ...interface{}) (bool, string, error) {
	return false, "", nil
}

func (c *API) Status(ctx context.Context, output *api.ExecuteQueryOutput) (*api.ExecuteQueryStatus, error) {
	statusResp, err := c.DataClient.DescribeStatement(ctx, &redshiftdata.DescribeStatementInput{
		Id: aws.String(output.ID),
	})
	if err != nil {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.ErrorStatus, err), false)
	}

	if statusResp.Error != nil && *statusResp.Error != "" {
		return nil, errorsource.DownstreamError(fmt.Errorf("%w: %v", api.ErrorExecute, *statusResp.Error), false)
	}

	var finished bool
	switch statusResp.Status {
	case redshiftdatatypes.StatusStringFailed,
		redshiftdatatypes.StatusStringAborted,
		redshiftdatatypes.StatusStringFinished:
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
	_, err := c.DataClient.CancelStatement(context.TODO(), &redshiftdata.CancelStatementInput{
		Id: &output.ID,
	})
	// ignore finished query error
	if err != nil && !strings.Contains(err.Error(), "Could not cancel a query that is already in FINISHED state") {
		return errorsource.DownstreamError(fmt.Errorf("%w: %v", err, api.ErrorStop), false)
	}
	return nil
}

func (c *API) Regions(context.Context) ([]string, error) {
	// This is not used. If regions are out of date, update them in the @grafana/aws-sdk-react package
	return []string{}, nil
}

func (c *API) Databases(ctx context.Context, _ sqlds.Options) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdata.ListDatabasesInput{
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
	input := &redshiftdata.ListSchemasInput{
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
	input := &redshiftdata.ListTablesInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		SchemaPattern:     aws.String(schema),
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
	input := &redshiftdata.DescribeTableInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Schema:            aws.String(schema),
		Table:             aws.String(table),
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
	input := &secretsmanager.ListSecretsInput{
		Filters: []secretsmanagertypes.Filter{
			{
				// Only secrets with the tag RedshiftQueryOwner can be used
				// https://docs.aws.amazon.com/redshift/latest/mgmt/query-editor.html#query-cluster-configure
				Key:    secretsmanagertypes.FilterNameStringTypeTagKey,
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
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
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

func (c *API) Clusters(ctx context.Context) ([]models.RedshiftCluster, error) {
	out, err := c.ManagementClient.DescribeClusters(ctx, &redshift.DescribeClustersInput{})
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

func (c *API) Workgroups(ctx context.Context) ([]models.RedshiftWorkgroup, error) {
	out, err := c.ServerlessManagementClient.ListWorkgroups(ctx, &redshiftserverless.ListWorkgroupsInput{})
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
