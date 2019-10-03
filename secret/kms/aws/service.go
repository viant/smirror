package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	akms "github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
	"smirror/auth"
	"smirror/secret/kms"
	"strings"
)

type service struct {
	*ssm.SSM
	*akms.KMS
}

func (s *service) Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error) {
	if secret.Parameter == "" {
		return nil, errors.New("parameter was empty")
	}
	if secret.Key == "" {
		return nil, errors.New("key was empty")
	}
	parameter, err := s.getParameters(secret.Parameter, true)
	if err != nil {
		return nil, err
	}
	return []byte(*parameter.Value), nil
}

func (s *service) getKeyByAlias(keyOrAlias string) (string, error) {
	if strings.Count(keyOrAlias, ":") > 0 {
		return keyOrAlias, nil
	}
	var nextMarker *string
	for {
		output, err := s.ListAliases(&akms.ListAliasesInput{
			Marker: nextMarker,
		})
		if err != nil {
			return "", err
		}
		if len(output.Aliases) == 0 {
			break
		}
		for _, candidate := range output.Aliases {
			if *candidate.AliasName == keyOrAlias {
				return *candidate.TargetKeyId, nil
			}
		}
		nextMarker = output.NextMarker
		if nextMarker == nil {
			break
		}
	}
	return "", fmt.Errorf("key for alias %v no found", keyOrAlias)
}

func (s *service) getParameters(name string, withDecryption bool) (*ssm.Parameter, error) {
	output, err := s.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return nil, err
	}
	return output.Parameter, nil
}

//New create AWS kms service
func New() (kms.Service, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &service{
		SSM: ssm.New(sess),
		KMS: akms.New(sess),
	}, nil
}
