package proxy

import (
	"github.com/viant/afs"
	"smirror/secret"
)

var singleton Service


//Singleton returns proxy singleton service
func Singleton(config *Config) Service {
	if singleton != nil {
		return singleton
	}
	fs := afs.New()
	singleton = New(fs, secret.New(config.SourceScheme, fs))
	return singleton
}
