package job

import "github.com/viant/toolbox/storage"

//represents job context
type Context struct {
	Error   error
	Storage storage.Service
	Name string
	SourceURL string
}

func NewContext(err error, storage storage.Service, sourceURL string) *Context {
	return &Context{
		Error:err,
		Storage:storage,
		SourceURL:sourceURL,
	}
}