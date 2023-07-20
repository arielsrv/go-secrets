package secrets

import (
	"context"
	"time"
)

type Store interface {
	Get(key string) *SecretDto
	GetWithContext(ctx context.Context, key string) *SecretDto
}

type SecretDto struct {
	Key       string
	Value     string
	ExpiresIn time.Duration
	Err       error
}

func (s *SecretDto) String() string {
	return s.Value
}
