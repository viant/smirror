package history

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/file"
	"github.com/viant/afs/option"
	"io/ioutil"
	"smirror/base"
)

//Events represents history events
type Events struct {
	URL    string
	Events []*Source
	index  map[string]*Source
}

//Put adds events
func (e *Events) Put(event *Source) bool {
	if prev, ok := e.index[event.URL]; ok && e.index[event.URL].Status == base.StatusOK {
		if prev.Time.Equal(event.Time) {
			return false
		}
	}
	e.Events = append(e.Events, event)
	e.index[event.URL] = event
	return true
}

//New creates events
func New(URL string) *Events {
	return &Events{
		URL:    URL,
		Events: make([]*Source, 0),
		index:  make(map[string]*Source),
	}
}

func (e *Events) indexEvents() {
	e.index = make(map[string]*Source)
	if len(e.Events) == 0 {
		e.Events = make([]*Source, 0)
	}
	for _, source := range e.Events {
		e.index[source.URL] = source
	}
}

//Persist persist history events
func (e *Events) Persist(ctx context.Context, fs afs.Service) error {
	data, err := json.Marshal(e)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal: %v", e.URL)
	}
	return fs.Upload(ctx, e.URL, file.DefaultFileOsMode, bytes.NewReader(data))
}

//FromURL creates events from URL
func FromURL(ctx context.Context, URL string, fs afs.Service) (*Events, error) {
	exists, _ := fs.Exists(ctx, URL, option.NewObjectKind(true))
	if !exists {
		return New(URL), nil
	}
	events := &Events{}
	reader, err := fs.OpenURL(ctx, URL)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read: %v", URL)
	}
	err = json.Unmarshal(data, events)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal: %v, %s", URL, data)
	}
	events.indexEvents()
	return events, nil
}
