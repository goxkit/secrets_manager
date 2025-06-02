// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package secretsmanager provides a unified abstraction layer for accessing various
// secret management services such as AWS Secrets Manager, HashiCorp Vault, and others.
// It offers a consistent interface for applications to retrieve sensitive information
// without being tightly coupled to specific secret management implementations.
//
// This package follows the adapter pattern, allowing applications to switch between
// different secret management providers with minimal code changes. This is particularly
// useful for applications that need to work across different cloud environments or
// deployment scenarios.
package secretsmanager

import "context"

type (
	// SecretClient defines the interface for interacting with secret providers.
	// This interface allows for different implementations such as AWS Secrets Manager,
	// HashiCorp Vault, or other secret management solutions to be used interchangeably.
	//
	// Implementations of this interface should handle connection management,
	// authentication, and any provider-specific behaviors required to access secrets.
	SecretClient interface {
		// LoadSecrets loads all secrets from the provider into memory.
		// This method should be called during application initialization to
		// prefetch secrets and improve performance by reducing the need for
		// repeated calls to the external service.
		//
		// The context can be used to control timeouts or cancellation of the
		// secret loading operation.
		//
		// Returns an error if any secrets cannot be loaded successfully.
		LoadSecrets(ctx context.Context) error

		// GetSecret retrieves a specific secret value by its key.
		// The implementation determines how the key is mapped to the actual secret
		// in the underlying secret management service.
		//
		// The context can be used to control timeouts or cancellation of the operation.
		//
		// Returns:
		//   - The secret value as a string if found
		//   - An error if the key doesn't exist or if there's a problem accessing the secret
		GetSecret(ctx context.Context, key string) (string, error)
	}
)
