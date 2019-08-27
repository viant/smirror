package config

import (
	"github.com/viant/afs/option"
)

//CustomKey represents custom key
type CustomKey struct {
	Secret
	AES256Key *option.AES256Key
}
