package proxy

import (
	"github.com/pkg/errors"
	"smirror/config"
)

//Request represents proxy
type Request struct {
	Source *config.Resource
	Dest   *config.Resource
	Move   bool
	Stream bool
}

//Validate checks if request is valid
func (r *Request) Validate() error {
	if r.Source == nil {
		return errors.Errorf("source was empty")
	}
	if r.Source.URL == "" {
		return errors.Errorf("source.url was empty")
	}
	if r.Dest.URL == "" {
		return errors.Errorf("dest.url was empty")
	}
	return nil
}
