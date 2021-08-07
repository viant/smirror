package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/option"
	"github.com/viant/afs/storage"
	"github.com/viant/afs/url"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"smirror/auth"
	"smirror/config"
	"smirror/secret/kms"
	"smirror/secret/kms/aws"
	"smirror/secret/kms/gcp"
)

//Service represents kms service
type Service interface {
	//Decrypt decrypts secrets
	Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error)
	//Load initialises resources
	Init(ctx context.Context, service afs.Service, resources []*config.Resource) error

	//StorageOpts returns storage option for supplied resource
	StorageOpts(ctx context.Context, resource *config.Resource) ([]storage.Option, error)
}

type service struct {
	sourceScheme string
	fs           afs.Service
}

func (s service) Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error) {
	kmsService, err := s.Kms(s.fs)
	if err != nil {
		return nil, err
	}
	data, err := kmsService.Decrypt(ctx, secret)
	if err != nil {
		return nil, err
	}
	data = decodeBase64IfNeeded(data)
	return data, err
}

//Kms returns kms service
func (s service) Kms(service afs.Service) (kms.Service, error) {
	switch s.sourceScheme {
	case gs.Scheme:
		return gcp.New(service), nil
	case s3.Scheme:
		return aws.New()
	}
	return nil, fmt.Errorf("unsupported scheme: %v", s.sourceScheme)
}

//Load initialises resources
func (s *service) Init(ctx context.Context, service afs.Service, resources []*config.Resource) (err error) {
	var kmsService kms.Service
	for i := range resources {
		resource := resources[i]
		if resource == nil {
			continue
		}
		if resource.Credentials == nil && resource.CustomKey == nil {
			continue
		}

		if (resource.Credentials != nil && resources[i].Credentials.Auth != nil) ||
			(resource.CustomKey != nil && resources[i].CustomKey.AES256Key != nil) {
			continue
		}

		if kmsService == nil {
			kmsService, err = s.Kms(service)
			if err != nil {
				return err
			}
		}

		if resource.Credentials != nil && resources[i].Credentials.Auth == nil {
			data, err := kmsService.Decrypt(ctx, &resource.Credentials.Secret)
			if err != nil {
				return err
			}
			resources[i].Credentials.Auth = decodeBase64IfNeeded(data)
		}

		if resource.CustomKey != nil && resources[i].CustomKey.AES256Key == nil {
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

//StorageOpts returns storage option for supplied resource
func (s service) StorageOpts(ctx context.Context, resource *config.Resource) ([]storage.Option, error) {
	var result = make([]storage.Option, 0)
	if resource == nil {
		return result, nil
	}
	if resource.CustomKey != nil && resource.CustomKey.AES256Key != nil {
		result = append(result, resource.CustomKey.AES256Key)
	}
	if resource.Proxy != nil {
		result = append(result, resource.Proxy)
	}
	if resource.Grant != nil {
		result = append(result, resource.Grant)
	}
	if resource.ACL != nil {
		result = append(result, resource.ACL)
	}
	if resource.URL == "" {
		return result, nil
	}
	var err error
	scheme := url.Scheme(resource.URL, file.Scheme)
	if resource.Credentials != nil && len(resource.Credentials.Auth) > 0 {
		if !json.Valid(resource.Credentials.Auth) {
			return nil, errors.Errorf("invalid credentials format, expected JSON but had: %s", resource.Credentials.Auth)
		}
		var authOpt interface{}
		switch scheme {
		case gs.Scheme:
			authOpt, err = gs.NewJwtConfig(resource.Credentials.Auth)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create goolge secrets: %s", resource.Credentials.Auth)
			}
			result = append(result, authOpt)
		case s3.Scheme:
			authOpt, err = s3.NewAuthConfig(resource.Credentials.Auth)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create aws secrets: %s", resource.Credentials.Auth)
			}
			result = append(result, authOpt)
			if resource.Region != "" {
				result = append(result, &option.Region{Name: resource.Region})
			}
		default:
			//do nothing init should take care of validating supported URL scheme
		}
	}
	return result, nil
}

//New creates a new secret service
func New(sourceScheme string, fs afs.Service) Service {
	return &service{
		fs:           fs,
		sourceScheme: sourceScheme,
	}
}
