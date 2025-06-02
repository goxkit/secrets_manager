# Secrets Manager

[![Go Reference](https://pkg.go.dev/badge/github.com/goxkit/secretsmanager.svg)](https://pkg.go.dev/github.com/goxkit/secretsmanager)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

The `secretsmanager` package provides a unified abstraction layer for accessing various secret management services like AWS Secrets Manager, HashiCorp Vault, and others. It offers a consistent interface for applications to retrieve sensitive information without being tightly coupled to specific secret management implementations.

This package is part of the [Goxkit](https://github.com/goxkit) toolkit, designed to simplify Go application development with standardized interfaces and implementations for common services.

## Features

- **Provider-agnostic Interface**: Consistent API regardless of the underlying secrets provider
- **In-memory Caching**: Improved performance by reducing calls to external services
- **Context Support**: All operations support context for timeout and cancellation handling
- **Easy Implementation**: Simple interface that can be extended to support various secret providers

## Supported Providers

- **AWS Secrets Manager**: Full implementation available
- **HashiCorp Vault**: Coming soon
- More providers to be added in future releases

## Installation

```bash
go get github.com/goxkit/secretsmanager
```

## Usage

### Basic Usage with AWS Secrets Manager

```go
package main

import (
	"context"
	"log"

	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/goxkit/secretsmanager/aws"
)

func main() {
	// Initialize application configurations
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		MQTT().
		Build()
	if err != nil {
		cfgs.Logger.Fatal(err.Error())
	}

	// Create a new AWS Secrets Manager client
	secretClient, err := aws.NewAwsSecretClient(cfgs)
	if err != nil {
		log.Fatalf("Failed to create AWS Secrets Manager client: %v", err)
	}

	// Load all secrets during application startup
	ctx := context.Background()
	if err := secretClient.LoadSecrets(ctx); err != nil {
		log.Fatalf("Failed to load secrets: %v", err)
	}

	// Retrieve a specific secret
	dbPassword, err := secretClient.GetSecret(ctx, "DB_PASSWORD")
	if err != nil {
		log.Fatalf("Failed to get DB password: %v", err)
	}

	// Use the secret value
	log.Println("Database connection established successfully")
}
```

### Secret Format in AWS Secrets Manager

Secrets in AWS Secrets Manager should be stored as JSON objects with key-value pairs. For example:

```json
{
  "DB_USERNAME": "admin",
  "DB_PASSWORD": "secure-password",
  "API_KEY": "your-api-key"
}
```

The package will parse this JSON and make each key-value pair available through the `GetSecret` method.

## Implementing a New Provider

To implement a new secret provider, create a new package that implements the `SecretClient` interface:

```go
type SecretClient interface {
    LoadSecrets(ctx context.Context) error
    GetSecret(ctx context.Context, key string) (string, error)
}
```

Each implementation should:

1. Handle connection management and authentication with the provider
2. Implement secret caching for performance optimization
3. Properly handle errors and logging
4. Follow the context pattern for operation lifecycle management

## Best Practices

- Call `LoadSecrets` during application initialization
- Use environment-specific secret identifiers
- Handle errors gracefully, especially for missing secrets
- Consider implementing a fallback mechanism for critical secrets

## License

This package is available under the MIT License. See the [LICENSE](LICENSE) file for more information.

## Contributing

Contributions are welcome! If you're interested in adding a new secret provider or enhancing the existing functionality, please feel free to submit a pull request.
