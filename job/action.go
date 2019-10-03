package job

import (
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"smirror/auth"
	"strings"
)

const (
	//ActionDelete delete action
	ActionDelete = "delete"
	//ActionMove move action
	ActionMove = "move"
)

//Action represents an action
type Action struct {
	Action string //empty Delete,Move
	URL    string
	Message interface{}
	Credentials *auth.Credentials
}

//DestURL returns destination URL
func (a Action) DestURL(relativePath string) string {
	return url.Join(a.URL, relativePath)
}

//WriteError writes an error file if context has error
func (a Action) WriteError(context *Context, service afs.Service) error {
	moveURL := a.DestURL(context.SourceURL) + "-error"
	return service.Upload(context.Context, moveURL, file.DefaultFileOsMode, strings.NewReader(context.Error.Error()))
}

//Do perform an action
func (a Action) Do(context *Context, service afs.Service) error {
	URL := context.SourceURL
	switch a.Action {
	case ActionDelete:
		return service.Delete(context.Context, URL)
	case ActionMove:
		targetURL := a.DestURL(context.RelativePath)
		return service.Move(context.Context, URL, targetURL)
	}
	return nil
}
