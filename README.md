> This package provides a Secrets client adapter with some features

# ⚡️ Usage

## Configuration

```go
package main

import (
	"gitlab.com/arielsrv/go-secrets"
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
```