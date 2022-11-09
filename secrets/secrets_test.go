package secrets_test

import (
	"context"
	"testing"

	"github.com/arielsrv/go-secrets/secrets"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretClient struct {
	mock.Mock
}

func (m *MockSecretClient) GetSecretValue(context.Context, *secretsmanager.GetSecretValueInput, ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called()
	return args.Get(0).(*secretsmanager.GetSecretValueOutput), args.Error(1)
}

func TestSecretService_Get(t *testing.T) {
	secretService := secrets.NewSecretService()
	secretClient := new(MockSecretClient)
	secretClient.On("GetSecretValue").Return(GetSecretValue())
	secretService.Client = secretClient

	actual, err := secretService.Get("my_key")
	assert.NoError(t, err)
	assert.Equal(t, "my_secret", actual)
}

func GetSecretValue() (*secretsmanager.GetSecretValueOutput, error) {
	output := new(secretsmanager.GetSecretValueOutput)
	output.SecretString = aws.String("my_secret")
	return output, nil
}
