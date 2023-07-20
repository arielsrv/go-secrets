package secrets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/store"
	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/patrickmn/go-cache"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSAdapter interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

type AWSSecretStore struct {
	AWSAdapter
	AppCache *go_cache.GoCacheStore
}

type Options func(f *AWSSecretStore)

func NewAWSSecretStore(opts ...Options) *AWSSecretStore {
	p := new(AWSSecretStore)

	p.AppCache = go_cache.NewGoCache(cache.New(1*time.Hour, 10*time.Minute))

	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		opt(p)
	}

	return p
}

func WithSecretClient(adapter AWSAdapter) Options {
	return func(s *AWSSecretStore) {
		s.AWSAdapter = adapter
	}
}

func (c AWSSecretStore) Get(key string) *SecretDto {
	return c.GetWithContext(context.Background(), key)
}

func (c AWSSecretStore) GetWithContext(ctx context.Context, key string) *SecretDto {
	secretDto, found, err := c.getFromCache(ctx, key)
	if err != nil {
		secretDto.Err = err
		return secretDto
	}

	if !found {
		secretDto = c.getFromAWS(ctx, key)
		err = c.put(ctx, key, secretDto)
		if err != nil {
			secretDto.Err = err
			return secretDto
		}
	}

	return secretDto
}

func (c AWSSecretStore) getFromCache(ctx context.Context, key string) (*SecretDto, bool, error) {
	value, err := c.AppCache.Get(ctx, key)
	if err != nil {
		if !errors.Is(err, &store.NotFound{}) {
			return nil, false, err
		}
	}

	if value == nil {
		return nil, false, nil
	}

	secretDto, ok := value.(*SecretDto)
	if !ok {
		return nil, false, fmt.Errorf("invalid cache casting for key: %s", key)
	}

	return secretDto, true, nil
}

func (c AWSSecretStore) getFromAWS(ctx context.Context, key string) *SecretDto {
	secretDto := new(SecretDto)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	secretDto.Key = *input.SecretId

	result, err := c.GetSecretValue(ctx, input)
	if err != nil {
		secretDto.Err = err
		return secretDto
	}

	// Decrypts secret using the associated KMS key.
	secretDto.Value = *result.SecretString

	return secretDto
}

func (c AWSSecretStore) put(ctx context.Context, key string, secretDto *SecretDto) error {
	err := c.AppCache.Set(ctx, key, secretDto)
	if err != nil {
		return err
	}
	return nil
}
