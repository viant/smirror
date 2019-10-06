package config

import (
	"github.com/viant/afs/option"
	"smirror/auth"
)

//CustomKey represents custom key
type CustomKey struct {
	auth.Secret `json:",omitempty"`
	AES256Key   *option.AES256Key `json:",omitempty"`
}
