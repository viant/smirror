package job

import (
	"context"
)

//Context represents job context
type Context struct {
	context.Context
	Error     error
	Name      string
	SourceURL string
}

//NewContext creates a context
func NewContext(ctx context.Context, err error, sourceURL string) *Context {
	return &Context{
		Context:   ctx,
		Error:     err,
		SourceURL: sourceURL,
	}
}
