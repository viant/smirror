package kms

import (
	"context"
	"github.com/viant/smirror/auth"
)

type Service interface {
	Decrypt(ctx context.Context, secret *auth.Secret) ([]byte, error)
}
