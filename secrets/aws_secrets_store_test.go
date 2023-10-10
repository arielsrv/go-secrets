package secrets_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-mq-producer/secrets"

	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/patrickmn/go-cache"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSecretStore struct {
	mock.Mock
	secrets map[string]*string
}

func NewMockSecretStore(values map[string]*string) secrets.AWSSecretStore {
	return secrets.AWSSecretStore{
		AWSAdapter: &MockSecretStore{
			secrets: values,
		},
		AppCache: go_cache.NewGoCache(cache.New(1*time.Hour, 10*time.Minute)),
	}
}

func (m *MockSecretStore) GetSecretValue(_ context.Context, input *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if m.secrets[aws.ToString(input.SecretId)] == nil {
		return nil, errors.New("secret not found")
	}

	return &secretsmanager.GetSecretValueOutput{
		SecretString: m.secrets[aws.ToString(input.SecretId)],
	}, nil
}

func TestSecretService_Get(t *testing.T) {
	secretService := NewMockSecretStore(map[string]*string{
		"my_key": aws.String("my_secret"),
	})

	actual := secretService.Get("my_key")
	assert.NoError(t, actual.Err)
	assert.Equal(t, "my_secret", actual.String())

	actual = secretService.Get("my_key")
	assert.NoError(t, actual.Err)
	assert.Equal(t, "my_secret", actual.String())
}

func TestSecretService_GetErr(t *testing.T) {
	secretService := NewMockSecretStore(map[string]*string{})

	actual := secretService.Get("my_key")
	assert.Error(t, actual.Err)
}
