package secret

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/cred"
	"github.com/viant/toolbox/kms"
	"github.com/viant/toolbox/kms/gcp"
)

func New(ctx context.Context, config *Config) (*cred.Config, error) {
	return newGcpCredConfig(ctx, config)
}

func newGcpCredConfig(ctx context.Context, config *Config) (*cred.Config, error) {
	decryptRequest := createDecryptRequest(config)
	decoderFactory := toolbox.NewJSONDecoderFactory()
	credConfig := &cred.Config{}
	kmsService := gcp.GetService()
	err := kmsService.Decode(ctx, decryptRequest, decoderFactory, credConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}
	return credConfig, nil
}

func createDecryptRequest(config *Config) *kms.DecryptRequest {
	decryptRequest := &kms.DecryptRequest{}
	decryptRequest.Key = config.Key
	decryptRequest.Parameter = config.Parameter
	decryptRequest.Resource = &kms.Resource{URL: config.URL}
	return decryptRequest
}
