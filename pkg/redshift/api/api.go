package api

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice"
	"github.com/aws/aws-sdk-go/service/redshiftdataapiservice/redshiftdataapiserviceiface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

type API struct {
	Client        redshiftdataapiserviceiface.RedshiftDataAPIServiceAPI
	SecretsClient secretsmanageriface.SecretsManagerAPI
	settings      *models.RedshiftDataSourceSettings
}

func New(sessionCache *awsds.SessionCache, settings *models.RedshiftDataSourceSettings) (*API, error) {
	region := settings.DefaultRegion
	if settings.Region != "" {
		region = settings.Region
	}
	session, err := sessionCache.GetSession(region, settings.AWSDatasourceSettings)
	if err != nil {
		return nil, err
	}
	return &API{redshiftdataapiservice.New(session), secretsmanager.New(session), settings}, nil
}

type apiInput struct {
	ClusterIdentifier *string
	Database          *string
	DbUser            *string
	SecretARN         *string
}

func (c *API) apiInput() apiInput {
	res := apiInput{
		ClusterIdentifier: aws.String(c.settings.ClusterIdentifier),
		Database:          aws.String(c.settings.Database),
	}
	if c.settings.UseManagedSecret {
		res.SecretARN = aws.String(c.settings.ManagedSecret.ARN)
	} else {
		res.DbUser = aws.String(c.settings.DBUser)
	}
	return res
}

func (c *API) Execute(ctx aws.Context, query string) (*redshiftdataapiservice.ExecuteStatementOutput, error) {
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.ExecuteStatementInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Sql:               aws.String(query),
	}

	return c.Client.ExecuteStatementWithContext(ctx, input)
}

func (c *API) ListSchemas(ctx aws.Context) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.ListSchemasInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.Client.ListSchemasWithContext(ctx, input)
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

func (c *API) ListTables(ctx aws.Context, schema string) ([]string, error) {
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
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.Client.ListTablesWithContext(ctx, input)
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

func (c *API) ListColumns(ctx aws.Context, schema, table string) ([]string, error) {
	commonInput := c.apiInput()
	input := &redshiftdataapiservice.DescribeTableInput{
		ClusterIdentifier: commonInput.ClusterIdentifier,
		Database:          commonInput.Database,
		DbUser:            commonInput.DbUser,
		SecretArn:         commonInput.SecretARN,
		Schema:            aws.String(schema),
		Table:             aws.String(table),
	}
	isFinished := false
	res := []string{}
	for !isFinished {
		out, err := c.Client.DescribeTableWithContext(ctx, input)
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

func (c *API) ListSecrets(ctx aws.Context) ([]models.ManagedSecret, error) {
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

func (c *API) GetSecret(ctx aws.Context, arn string) (*models.RedshiftSecret, error) {
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
