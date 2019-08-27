package job

import (
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/toolbox"
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

//WriteError writes an error file if context has error
func (a Action) WriteError(context *Context, service afs.Service) error {
	_, name := toolbox.URLSplit(context.SourceURL)
	moveURL := toolbox.URLPathJoin(a.URL, name) + "-error"
	return service.Upload(context.Context, moveURL, file.DefaultFileOsMode, strings.NewReader(context.Error.Error()))
}

//Do perform an action
func (a Action) Do(context *Context, service afs.Service) error {
	URL := context.SourceURL
	switch a.Action {
	case ActionDelete:
		return service.Delete(context.Context, URL)
	case ActionMove:
		_, name := toolbox.URLSplit(context.SourceURL)
		targetURL := toolbox.URLPathJoin(a.URL, name)

		fmt.Printf("move: %v -> %v\n", URL, targetURL)

		return service.Move(context.Context, URL, targetURL)
	}
	return nil
}
