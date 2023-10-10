package main

import (
	"log"

	_ "github.com/arielsrv/go-config/autoload"
	"github.com/arielsrv/go-config/env"
	"gitlab.com/iskaypetcom/digital/sre/tools/dev/go-mq-producer/secrets"
)

func main() {
	secretStore := secrets.NewAWSSecretBuilder().
		Build()

	secretDto := secretStore.Get(env.Get("key"))

	if secretDto.Err != nil {
		log.Println(secretDto.Err)
	}

	log.Println(secretDto.Value)
}
