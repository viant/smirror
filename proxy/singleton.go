package proxy

import (
	"github.com/viant/smirror/secret"
	"github.com/viant/afs"
)

var singleton Service

//Singleton returns proxy singleton service
func Singleton(config *Config) Service {
	if singleton != nil {
		return singleton
	}
	fs := afs.New()
	singleton = New(fs, config, secret.New(config.SourceScheme, fs))
	return singleton
}
