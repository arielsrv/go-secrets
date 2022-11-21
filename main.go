package main

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:

import (
	"log"

	"github.com/go-chassis/go-archaius"

	"github.com/arielsrv/go-secrets/secrets"
)

func main() {
	secretsStore := secrets.NewSecretsStore()
	value := secretsStore.Get("PETS_PETS-API_GITLAB_TOKEN")
	if value.Err != nil {
		log.Fatal(value.Err)
	}
	log.Print(value.String())

	err := archaius.Init(
		archaius.WithENVSource(),
	)

	if err != nil {
		log.Fatal(err)
	}

	archaiusStore := secrets.NewArchaiusSecretStore()
	localValue := archaiusStore.Get("PATH")
	if localValue.Err != nil {
		log.Fatal(localValue.Err)
	}
	log.Print(localValue.Value)
}
