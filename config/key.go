package config

import (
	"github.com/viant/afs/option"
	"github.com/viant/smirror/auth"
)

//CustomKey represents custom key
type CustomKey struct {
	auth.Secret `json:",omitempty"`
	AES256Key   *option.AES256Key `json:"-"`
}
