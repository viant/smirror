package job

import (
	"context"
)

//Context represents job context
type Context struct {
	context.Context
	Error        error
	Name         string
	RelativePath string
	SourceURL    string
}

//NewContext creates a context
func NewContext(ctx context.Context, err error, sourceURL, relativePath string) *Context {
	return &Context{
		Context:      ctx,
		Error:        err,
		RelativePath: relativePath,
		SourceURL:    sourceURL,
	}
}
