package api

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/redshift-datasource/pkg/redshift/models"
)

type API struct {
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
	return &API{secretsmanager.New(session), settings}, nil
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
