// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package aws provides an AWS Secrets Manager implementation of the SecretClient interface.
// It enables applications to retrieve secrets stored in AWS Secrets Manager using
// a consistent API defined by the secretsmanager package.
package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/goxkit/configs"
	"github.com/goxkit/logging"
	"go.uber.org/zap"

	sm "github.com/goxkit/secretsmanager"
)

// awsSecretClient is an implementation of the SecretClient interface that uses
// AWS Secrets Manager to store and retrieve secrets. It maintains an in-memory
// cache of secrets to minimize API calls and improve performance.
type awsSecretClient struct {
	logger      logging.Logger
	client      *secretsmanager.Client
	appSecretId string            // The AWS Secrets Manager secret identifier
	secrets     map[string]string // In-memory cache of secret key-value pairs
}

// NewAwsSecretClient creates a new instance of AWS Secrets Manager client.
//
// It initializes the AWS configuration using the default credential providers chain,
// and prepares the secret identifier based on the application environment and secret key.
// The secret ID format follows the pattern: "{environment}/{secretKey}".
//
// Parameters:
//   - cfgs: Application configuration containing environment, secret key, and logger
//
// Returns:
//   - A SecretClient interface implementation for AWS Secrets Manager
//   - An error if AWS configuration cannot be loaded
func NewAwsSecretClient(cfgs *configs.Configs) (sm.SecretClient, error) {
	logger := cfgs.Logger

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error get aws configs from env", zap.Error(err))
		return nil, err
	}

	// Format the secret ID using environment and app secret key
	appSecretId := fmt.Sprintf("%s/%s", cfgs.AppConfigs.Environment.ToString(), cfgs.AppConfigs.SecretKey)

	return &awsSecretClient{
		logger:      logger,
		client:      secretsmanager.NewFromConfig(awsCfg),
		appSecretId: appSecretId,
		secrets:     make(map[string]string),
	}, nil
}

// LoadSecrets retrieves all secrets from AWS Secrets Manager for the configured secret ID.
//
// This method makes an API call to AWS Secrets Manager to fetch the secret value as a JSON blob,
// then unmarshals it into an in-memory map of string keys to string values. This approach
// enables fast access to secrets without requiring repeated calls to AWS for each secret lookup.
//
// The method should be called during application initialization to ensure secrets are available
// when needed. If the secret values change in AWS Secrets Manager, the application would need
// to be restarted or this method called again to refresh the cached values.
//
// Parameters:
//   - ctx: Context for controlling the request lifecycle
//
// Returns:
//   - An error if the secret cannot be fetched or parsed
func (c *awsSecretClient) LoadSecrets(ctx context.Context) error {
	// Call AWS Secrets Manager API to get the secret value
	res, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &c.appSecretId,
	})

	if err != nil {
		c.logger.Error("error to get secret", zap.Error(err))
		return err
	}

	// Reset the cache before loading new values
	c.secrets = map[string]string{}

	// Parse the secret JSON data into our cache map
	err = json.Unmarshal(res.SecretBinary, &c.secrets)
	if err != nil {
		c.logger.Error("error get secret from aws", zap.Error(err))
		return err
	}

	return nil
}

// GetSecret retrieves a specific secret value by its key from the in-memory cache.
//
// This method performs a lookup in the in-memory cache that was populated by LoadSecrets.
// It's designed to be fast and efficient, avoiding repeated calls to AWS Secrets Manager
// for each secret retrieval. The method will return an error if the requested key does
// not exist in the cache.
//
// Note that the context parameter is not used in this implementation since the lookup
// is performed on the local cache, but it's included to satisfy the SecretClient interface.
//
// Parameters:
//   - ctx: Context (not used in this implementation)
//   - key: The secret key to look up
//
// Returns:
//   - The secret value as a string if found
//   - An error if the key doesn't exist in the cache
func (c *awsSecretClient) GetSecret(_ context.Context, key string) (string, error) {
	value, ok := c.secrets[key]
	if !ok {
		return "", errors.New("secret was not found")
	}

	return value, nil
}
