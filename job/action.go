package job

import (
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"smirror/auth"
	"smirror/base"
	"strings"
)

const (
	//ActionDelete delete action
	ActionDelete = "delete"
	//ActionMove move action
	ActionMove = "move"
	//Action notify
	ActionNotify = "notify"
)

//Action represents an action
type Action struct {
	Action          string //empty Delete,Move
	URL             string `json:",omitempty"`
	Message         string `json:",omitempty"`
	Title           string `json:",omitempty"`
	Body			interface{} `json:",omitempty"`
	Channels        []string `json:",omitempty"`
	Credentials     *auth.Credentials `json:",omitempty"`
}

//DestURL returns destination URL
func (a Action) DestURL(relativePath string) string {
	if strings.Contains(relativePath, "://") {
		_, relativePath = url.Base(relativePath, "file")
	}
	return url.Join(a.URL, relativePath)
}

//WriteError writes an error file if context has error
func (a Action) WriteError(context *Context, service afs.Service) error {
	moveURL := a.DestURL(context.SourceURL) + "-error"
	return service.Upload(context.Context, moveURL, file.DefaultFileOsMode, strings.NewReader(context.Error.Error()))
}

//Do perform an action
func (a Action) Do(context *Context, service afs.Service, notify Notify, info *base.Info, response interface{}) (err error) {
	URL := context.SourceURL
	switch a.Action {
	case ActionDelete:
		err = service.Delete(context.Context, URL)
	case ActionNotify:
		body := a.Body
		if textBody, ok := a.Body.(string);ok && textBody == "$Response" {
			body =  response
		}
		title := strings.Replace(a.Title, "$SourceURL", context.SourceURL, 1)
		message := strings.Replace(a.Message, "$SourceURL", context.SourceURL, 1)
		if len(a.Channels) == 0 && info.SlackChannel != "" {
			a.Channels = []string{info.SlackChannel}
		}
		err = notify(context.Context, &NotifyRequest{
			From:        base.App,
			Title:       title,
			Channels:    a.Channels,
			Credentials: a.Credentials,
			Message:     message,
			Body:        body,
		})
	case ActionMove:
		targetURL := a.DestURL(context.RelativePath)
		return service.Move(context.Context, URL, targetURL)
	default:
		err = fmt.Errorf("unsupported action: %v", a.Action)
	}
	return err
}
