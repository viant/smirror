package cron

import (
	"github.com/viant/smirror/cron/config"
	"github.com/viant/smirror/proxy"
	"github.com/viant/afs/storage"
)

//Response represents schedule response
type Response struct {
	*proxy.Response
	Matched []*Matched `json:",omitempty"`
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
func NewResponse(baseResponse *proxy.Response) *Response {
	return &Response{
		Response: baseResponse,
		Matched:  make([]*Matched, 0),
	}
}
