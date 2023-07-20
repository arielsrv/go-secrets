package main

import (
	"github.com/arielsrv/go-secrets/secrets"
	"log"
)

func main() {
	secretStore := secrets.NewAWSSecretBuilder().
		Build()

	secretDto := secretStore.Get("cache.password")

	if secretDto.Err != nil {
		log.Println(secretDto.Err)
	}

	log.Println(secretDto.Value)
}
