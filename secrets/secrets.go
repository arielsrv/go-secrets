package secrets

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type ISecretService interface {
	Get(key string) (string, error)
}

type SecretService struct {
	config aws.Config
	Client ISecretClient
}

type ISecretClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type SecretClient struct {
	client *secretsmanager.Client
}

func (s *SecretClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return s.client.GetSecretValue(ctx, params, optFns...)
}

func NewSecretService() *SecretService {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	client := secretsmanager.NewFromConfig(config)

	return &SecretService{
		config: config,
		Client: &SecretClient{
			client: client,
		},
	}
}

func (s SecretService) Get(key string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := s.Client.GetSecretValue(context.TODO(), input)
	if err != nil {
		// For a list of exceptions thrown, see
		// https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
		return "", err
	}

	// Decrypts secret using the associated KMS key.
	var secretString = *result.SecretString
	return secretString, nil
}
