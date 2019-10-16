package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
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
	service, err := New(ctx, config)
	if err != nil {
		JSON, _ := json.Marshal(config)
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create service from config %s", JSON))
	}
	toolbox.Dump(config)
	singletonEnvKey = envKey
	singleton = service
	return singleton, nil
}
