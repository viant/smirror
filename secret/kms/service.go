package kms

import (
	"context"
	"smirror/auth"
)

type Service interface {
	Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error)
}
