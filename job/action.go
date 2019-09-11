package job

import (
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"path"
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
}

//DestURL returns destination URL
func (a Action) DestURL(sourceURL string) string {
	_, URLPath := url.Base(sourceURL, file.Scheme)
	_, name := path.Split(URLPath)
	return url.Join(a.URL, name)
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
		targetURL := a.DestURL(context.SourceURL)
		return service.Move(context.Context, URL, targetURL)
	}
	return nil
}
