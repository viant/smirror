package cron

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
)

var singleton Service
var singletonEnvKey string

//NewFromEnv returns new service for env key
func NewFromEnv(ctx context.Context, envKey string) (Service, error) {
	if singleton != nil && envKey == singletonEnvKey {
		return singleton, nil
	}
	config, err := NewConfigFromEnv(ctx, envKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config from env key "+envKey)
	}
	service, err := New(ctx, config, afs.New())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create service from config %v", config))
	}
	singletonEnvKey = envKey
	singleton = service
	return singleton, nil
}
