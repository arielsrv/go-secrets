package secrets

type ISecretStore interface {
	Get(key string) *SecretDto
}

type SecretDto struct {
	Key   string
	Value string
	Err   error
}

func (s *SecretDto) String() string {
	return s.Value
}
