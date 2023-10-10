package main

import (
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-mq-producer/secrets"
	"log"
)

func main() {
	secretStore := secrets.NewAWSSecretBuilder(
		secrets.WithRegion("us-east-1"),
		secrets.WithPartition("aws"),
		secrets.WithEndpointURL("http://localhost:4566")).
		Build()

	secretDto := secretStore.Get("cache.password")

	if secretDto.Err != nil {
		log.Println(secretDto.Err)
	}

	log.Println(secretDto.Value)
}
