package job

import (
	"github.com/viant/toolbox"
	"strings"
)

const (
	ActionDelete = "delete"
	ActionMove   = "move"
)

type Action struct {
	Action string //empty Delete,Move
	URL    string
}

func (a Action) WriteError(context *Context) error {
	_, name := toolbox.URLSplit(context.SourceURL)
	moveURL := toolbox.URLPathJoin(a.URL, name) + "-error"
	return context.Storage.Upload(moveURL, strings.NewReader(context.Error.Error()))
}

func (a Action) Do(context *Context) error {
	URL := context.SourceURL
	switch a.Action {
	case ActionDelete:
		if object, err := context.Storage.StorageObject(URL); err == nil {
			return context.Storage.Delete(object)
		}
	case ActionMove:
		//TODO optimize it
		rawData, err := context.Storage.DownloadWithURL(URL)
		if err != nil {
			return err
		}
		_, name := toolbox.URLSplit(context.SourceURL)
		moveURL := toolbox.URLPathJoin(a.URL, name)
		if err := context.Storage.Upload(moveURL, rawData); err != nil {
			return err
		}

		if object, err := context.Storage.StorageObject(URL); err == nil {
			return context.Storage.Delete(object)
		}

	}
	return nil
}
