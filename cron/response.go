package cron

import (
	"github.com/viant/afs/storage"
	"smirror/base"
	"smirror/cron/config"
)

//Response represents schedule response
type Response struct {
	Matched []*Matched        `json:",omitempty"`
	Copied  map[string]string `json:",omitempty"`
	Moved   map[string]string `json:",omitempty"`
	Status  string            `json:",omitempty"`
	Error   string            `json:",omitempty"`
}

type Matched struct {
	Resource *config.Rule `json:",omitempty"`
	URLs     []string     `json:",omitempty"`
}

func (m *Matched) Add(objects ...storage.Object) {
	for _, object := range objects {
		m.URLs = append(m.URLs, object.URL())
	}
}

//NewResponse create a response
func NewResponse() *Response {
	return &Response{
		Status:  base.StatusOK,
		Matched: make([]*Matched, 0),
		Copied:  make(map[string]string),
		Moved:   make(map[string]string),
	}
}
