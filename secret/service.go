package secret

import (
	"context"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"smirror/config"
	"smirror/secret/kms"
	"smirror/secret/kms/aws"
	"smirror/secret/kms/gcp"
)

type Service interface {
	Init(ctx context.Context, service afs.Service, resources []*config.Resource) error

	StorageOpts(ctx context.Context, esource *config.Resource) ([]storage.Option, error)
}

type service struct {
	sourceScheme string
}

func (s service) Kms(service afs.Service) (kms.Service, error) {
	switch s.sourceScheme {
	case gs.Scheme:
		return gcp.New(service), nil
	case s3.Scheme:
		return aws.New()
	}
	return nil, fmt.Errorf("unsupported scheme: %v", s.sourceScheme)
}

func (s *service) Init(ctx context.Context, service afs.Service, resources []*config.Resource) error {
	kmsService, err := s.Kms(service)
	if err != nil {
		return err
	}
	for i := range resources {
		resource := resources[i]
		if resource == nil {
			continue
		}
		if resource.Credentials != nil {
			data, err := kmsService.Decrypt(ctx, &resource.Credentials.Secret)
			if err != nil {
				return err
			}
			resources[i].Credentials.Auth = decodeBase64IfNeeded(data)
		}
		if resource.CustomKey != nil {
			data, err := kmsService.Decrypt(ctx, &resource.CustomKey.Secret)
			if err != nil {
				return err
			}
			data = decodeBase64IfNeeded(data)

			if resources[i].CustomKey.AES256Key, err = option.NewAES256Key(data); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s service) StorageOpts(ctx context.Context, resource *config.Resource) ([]storage.Option, error) {
	var result = make([]storage.Option, 0)
	if resource == nil {
		return result, nil
	}
	if resource.CustomKey != nil && resource.CustomKey.AES256Key != nil {
		result = append(result, resource.CustomKey.AES256Key)
	}
	if resource.URL == "" {
		return result, nil
	}
	scheme := url.Scheme(resource.URL, file.Scheme)
	if resource.Credentials != nil {
		switch scheme {
		case gs.Scheme:
			auth, err := gs.NewJwtConfig(resource.Credentials.Auth)
			if err != nil {
				return nil, err
			}
			result = append(result, auth)
		case s3.Scheme:
			auth, err := s3.NewAuthConfig(resource.Credentials.Auth)
			if err != nil {
				return nil, err
			}
			result = append(result, auth)
		default:
			//do nothing init should take care of validating supported URL scheme
		}
	}
	return result, nil
}

//New creates a new secret service
func New(sourceScheme string) Service {
	return &service{
		sourceScheme: sourceScheme,
	}
}
