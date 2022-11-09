package main

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://aws.github.io/aws-sdk-go-v2/docs/getting-started/

import (
	"log"

	"github.com/arielsrv/go-secrets/secrets"
)

func main() {
	secretService := secrets.NewSecretService()
	value, err := secretService.Get("secretin")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(value)
}
