package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grafana/sqlds/v3"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/aws/aws-sdk-go/service/redshift/redshiftiface"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice/redshiftdataapiserviceiface"
	"github.com/aws/aws-sdk-go/service/redshiftserverless"
	"github.com/aws/aws-sdk-go/service/redshiftserverless/redshiftserverlessiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-aws-sdk/pkg/sql/api"
	awsModels "github.com/grafana/grafana-aws-sdk/pkg/sql/models"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

type API struct {
	DataClient                 redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI
	SecretsClient              secretsmanageriface.SecretsManagerAPI
	ManagementClient           redshiftiface.RedshiftAPI
	ServerlessManagementClient redshiftserverlessiface.RedshiftServerlessAPI
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

	sess, err := sessionCache.GetSession(awsds.SessionConfig{
		Settings:      redshiftSettings.AWSDatasourceSettings,
		HTTPClient:    httpClient,
		UserAgentName: aws.String("Redshift"),
	})
	if err != nil {
		return nil, err
	}

	return &API{
		DataClient:                 redshiftdataapiservice.New(sess),
		SecretsClient:              secretsmanager.New(sess),
		ManagementClient:           redshift.New(sess),
		ServerlessManagementClient: redshiftserverless.New(sess),
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
	redshiftInput := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Sql:               aws.String(input.Query),
		WithEvent:         aws.Bool(c.settings.WithEvent),
		WorkgroupName:     commonInput.WorkgroupName,
	}
	output, err := c.DataClient.ExecuteStatementWithContext(ctx, redshiftInput)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", api.ExecuteError, err)
	}

	return &api.ExecuteQueryOutput{ID: *output.Id}, nil
}

// GetQueryID always returns not found. To actually check if the query has been called requires calling ListStatements, which can lead to timeouts
// when there are many statements to page through
func (c *API) GetQueryID(ctx context.Context, query string, args ...interface{}) (bool, string, error) {
	return false, "", nil
}

func (c *API) Status(ctx aws.Context, output *api.ExecuteQueryOutput) (*api.ExecuteQueryStatus, error) {
	statusResp, err := c.DataClient.DescribeStatementWithContext(ctx, &redshiftdataapiservice.DescribeStatementInput{
		Id: aws.String(output.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", api.StatusError, err)
	}

	if statusResp.Error != nil && *statusResp.Error != "" {
		return nil, fmt.Errorf("%w: %v", api.ExecuteError, *statusResp.Error)
	}

	var finished bool
	state := *statusResp.Status
	switch state {
	case redshiftdataapiservice.StatusStringFailed,
		redshiftdataapiservice.StatusStringAborted:
		finished = true
	case redshiftdataapiservice.StatusStringFinished:
		finished = true
	default:
		finished = false
	}

	return &api.ExecuteQueryStatus{
		ID:       output.ID,
		State:    state,
		Finished: finished,
	}, nil
}

func (c *API) CancelQuery(ctx context.Context, options sqlds.Options, queryID string) error {
	return c.Stop(&api.ExecuteQueryOutput{ID: queryID})
}

func (c *API) Stop(output *api.ExecuteQueryOutput) error {
	_, err := c.DataClient.CancelStatement(&redshiftdataapiservice.CancelStatementInput{
		Id: &output.ID,
	})
	// ignore finished query error
	if err != nil && !strings.Contains(err.Error(), "Could not cancel a query that is already in FINISHED state") {
		return fmt.Errorf("%w: %v", err, api.StopError)
	}
	return nil
}

func (c *API) Regions(aws.Context) ([]string, error) {
	// List from https://docs.aws.amazon.com/general/latest/gr/redshift-service.html
	return []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-2",
		"ap-southeast-3",
		"ap-southeast-4",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"ca-west-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-south-2",
		"eu-north-1",
		"eu-central-2",
		"il-central-1",
		"me-south-1",
		"me-central-1",
		"sa-east-1",
		"us-gov-east-1",
		"us-gov-west-1",
	}, nil
}

func (c *API) Databases(ctx aws.Context, options sqlds.Options) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.ListDatabasesInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.ListDatabasesWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, sc := range out.Databases {
			if sc != nil {
				res = append(res, *sc)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Schemas(ctx aws.Context, options sqlds.Options) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.ListSchemasInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		WorkgroupName:     commonInput.WorkgroupName,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.DataClient.ListSchemasWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, sc := range out.Schemas {
			if sc != nil {
				res = append(res, *sc)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Tables(ctx aws.Context, options sqlds.Options) ([]string, error) {
	schema := options["schema"]
	// We use the "public" schema by default if not specified
	if schema == "" {
		schema = "public"
	}
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.ListTablesInput{
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
		out, err := c.DataClient.ListTablesWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, t := range out.Tables {
			if t.Name != nil {
				res = append(res, *t.Name)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Columns(ctx aws.Context, options sqlds.Options) ([]string, error) {
	schema, table := options["schema"], options["table"]
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.DescribeTableInput{
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
		out, err := c.DataClient.DescribeTableWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		for _, c := range out.ColumnList {
			if c.Name != nil {
				res = append(res, *c.Name)
			}
		}
		if input.NextToken == nil {
			isFinished = true
		}
	}
	return res, nil
}

func (c *API) Secrets(ctx aws.Context) ([]models.ManagedSecret, error) {
	input := &secretsmanager.ListSecretsInput{
		Filters: []*secretsmanager.Filter{
			{
				// Only secrets with the tag RedshiftQueryOwner can be used
				// https://docs.aws.amazon.com/redshift/latest/mgmt/query-editor.html#query-cluster-configure
				Key:    aws.String(secretsmanager.FilterNameStringTypeTagKey),
				Values: []*string{aws.String("RedshiftQueryOwner")},
			},
		},
	}
	isFinished := false
	redshiftSecrets := []models.ManagedSecret{}
	for !isFinished {
		out, err := c.SecretsClient.ListSecretsWithContext(ctx, input)
		if err != nil {
			return nil, err
		}
		input.NextToken = out.NextToken
		if input.NextToken == nil {
			isFinished = true
		}
		for _, s := range out.SecretList {
			if s.ARN == nil || s.Name == nil {
				continue
			}
			redshiftSecrets = append(redshiftSecrets, models.ManagedSecret{
				ARN:  *s.ARN,
				Name: *s.Name,
			})
		}
	}
	return redshiftSecrets, nil
}

func (c *API) Secret(ctx aws.Context, options sqlds.Options) (*models.RedshiftSecret, error) {
	arn := options["secretARN"]
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(arn),
	}
	out, err := c.SecretsClient.GetSecretValueWithContext(ctx, input)
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
	out, err := c.ManagementClient.DescribeClusters(&redshift.DescribeClustersInput{})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("missing clusters content")
	}
	res := []models.RedshiftCluster{}
	for _, r := range out.Clusters {
		if r != nil && r.ClusterIdentifier != nil && r.Endpoint != nil && r.Endpoint.Address != nil && r.Endpoint.Port != nil && r.DBName != nil {
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
	out, err := c.ServerlessManagementClient.ListWorkgroups(&redshiftserverless.ListWorkgroupsInput{})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, fmt.Errorf("missing workgroups content")
	}
	res := []models.RedshiftWorkgroup{}
	for _, r := range out.Workgroups {
		if r != nil && r.WorkgroupName != nil && r.Endpoint != nil && r.Endpoint.Address != nil && r.Endpoint.Port != nil {
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
