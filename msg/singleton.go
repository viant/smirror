package msg

import (
	"github.com/viant/afs"
)

var singleton Service

//Singleton returns new service for env key
func Singleton(config *Config) Service {
	if singleton != nil {
		return singleton
	}
	singleton := New(config, afs.New())
	return singleton
}
