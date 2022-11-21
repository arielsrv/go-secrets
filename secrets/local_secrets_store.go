package secrets

import (
	"fmt"

	"github.com/go-chassis/go-archaius"
)

type LocalSecretStore struct {
}

type ArchaiusSecretStore struct {
	LocalSecretStore
}

func NewArchaiusSecretStore() *ArchaiusSecretStore {
	return &ArchaiusSecretStore{}
}

func (l LocalSecretStore) Get(key string) *SecretDto {
	secretDto := new(SecretDto)
	secretDto.Key = key
	secretDto.Value = archaius.GetString(key, "")

	if secretDto.Value == "" {
		secretDto.Err = fmt.Errorf("missing secret: %s", key)
	}

	return secretDto
}
