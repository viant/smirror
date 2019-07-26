package mirror

import (
	"fmt"
	"github.com/pkg/errors"
)

var singleton Service
var singletonEnvKey string

//NewFromEnv returns new service for env key
func NewFromEnv(envKey string) (Service, error) {
	if singleton != nil && envKey == singletonEnvKey {
		return singleton, nil
	}
	config, err := NewConfigFromEnv(envKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create config from env key "+envKey)
	}
	service, err := New(config)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create service from config %v", config))
	}
	singletonEnvKey = envKey
	singleton = service
	return singleton, nil
}
