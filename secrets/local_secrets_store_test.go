package secrets_test

import (
	"testing"

	"github.com/arielsrv/go-secrets/secrets"
	"github.com/go-chassis/go-archaius"
	"github.com/stretchr/testify/assert"
)

func TestLocalSecretStore_Get(t *testing.T) {
	t.Setenv("my_key", "value")
	err := archaius.Init(archaius.WithENVSource())
	assert.NoError(t, err)

	secretsStore := secrets.NewArchaiusSecretStore()
	actual := secretsStore.Get("my_key")

	assert.NotNil(t, actual)
	assert.NoError(t, actual.Err)
	assert.Equal(t, "value", actual.Value)
}

func TestLocalSecretStore_Get_Err(t *testing.T) {
	err := archaius.Init(archaius.WithENVSource())
	assert.NoError(t, err)

	secretsStore := secrets.NewArchaiusSecretStore()
	actual := secretsStore.Get("missing_key")

	assert.NotNil(t, actual)
	assert.Error(t, actual.Err)
}
