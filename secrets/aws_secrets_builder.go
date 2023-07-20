package secrets

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSSecretBuilder struct {
	partitionID string
	endpointURL string
	region      string
}

func NewAWSSecretBuilder(opts ...BuilderOptions) *AWSSecretBuilder {
	b := new(AWSSecretBuilder)

	for i := 0; i < len(opts); i++ {
		opt := opts[i]
		opt(b)
	}

	return b
}

type BuilderOptions func(f *AWSSecretBuilder)

func WithPartition(partitionID string) BuilderOptions {
	return func(s *AWSSecretBuilder) {
		s.partitionID = partitionID
	}
}

func WithEndpointURL(endpointURL string) BuilderOptions {
	return func(s *AWSSecretBuilder) {
		s.endpointURL = endpointURL
	}
}

func WithRegion(region string) BuilderOptions {
	return func(s *AWSSecretBuilder) {
		s.region = region
	}
}

var (
	once        sync.Once
	secretStore *AWSSecretStore
)

func (b AWSSecretBuilder) Build() Store {
	once.Do(func() {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if b.endpointURL != "" {
				return aws.Endpoint{
					PartitionID:   b.partitionID,
					URL:           b.endpointURL,
					SigningRegion: b.region,
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		awsCfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolverWithOptions(customResolver),
		)

		if err != nil {
			log.Fatal(err)
		}

		secretStore = NewAWSSecretStore(
			WithSecretClient(secretsmanager.NewFromConfig(awsCfg)),
		)
	})

	return secretStore
}
