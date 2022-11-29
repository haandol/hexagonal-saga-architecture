package internal

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/haandol/hexagonal/pkg/util"
)

type DatabaseSecrets struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Name     string `json:"dbname" validate:"required"`
	Host     string `json:"host" validate:"required"`
}

// GetSecretString fetch secret values as string from AWS Secrets Manager using aws-sdk-go-v2
func GetSecretsWithID(awsCfg *aws.Config, secretID string) (*DatabaseSecrets, error) {
	// Create the Secrets Manager client
	client := secretsmanager.NewFromConfig(*awsCfg)

	// Build the request with its input parameters
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}
	// Get the secret value
	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var secrets DatabaseSecrets
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return nil, err
	}
	if err := util.ValidateStruct(secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}
