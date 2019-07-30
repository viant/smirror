package secret

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/cred"
	"github.com/viant/toolbox/kms"
	"github.com/viant/toolbox/kms/aws"
	"github.com/viant/toolbox/kms/gcp"
)

//New creates a new secret service
func New(ctx context.Context, config *Config) (*cred.Config, error) {
	switch config.Provider {
	case "aws":
		return newAwsCredConfig(ctx, config)
	case "gcp":
		return newGcpCredConfig(ctx, config)
	default:
		return nil, fmt.Errorf("unsupported secret  provider: %v", config.Provider)
	}
}

func newAwsCredConfig(ctx context.Context, config *Config) (*cred.Config, error) {
	decryptRequest := createDecryptRequest(config)
	decoderFactory := toolbox.NewJSONDecoderFactory()
	credConfig := &cred.Config{}
	kmsService, err := aws.New()
	if err != nil {
		return nil, err
	}
	err = kmsService.Decode(ctx, decryptRequest, decoderFactory, credConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	toolbox.Dump(credConfig)
	return credConfig, nil
}

func newGcpCredConfig(ctx context.Context, config *Config) (*cred.Config, error) {
	decryptRequest := createDecryptRequest(config)
	decoderFactory := toolbox.NewJSONDecoderFactory()
	credConfig := &cred.Config{}
	kmsService := gcp.New()
	toolbox.Dump(decryptRequest)
	err := kmsService.Decode(ctx, decryptRequest, decoderFactory, credConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	return credConfig, nil
}

func createDecryptRequest(config *Config) *kms.DecryptRequest {
	decryptRequest := &kms.DecryptRequest{}
	decryptRequest.Key = config.Key
	decryptRequest.Resource = &kms.Resource{URL: config.URL}
	decryptRequest.Parameter = config.Parameter
	return decryptRequest
}
