package kms

import (
	"context"
	"smirror/config"
)

type Service interface {
	Decrypt(ctx context.Context, secret *config.Secret) ([]byte, error)
}
