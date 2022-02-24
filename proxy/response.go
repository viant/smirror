package proxy

import (
	"github.com/viant/smirror/base"
	"sync"
)

type Response struct {
	SourceURL string  `json:",omitempty"`
	Copied  map[string]string `json:",omitempty"`
	Moved   map[string]string `json:",omitempty"`
	Invoked map[string]string `json:",omitempty"`
	Status  string            `json:",omitempty"`
	Error   string            `json:",omitempty"`
	mux     *sync.Mutex
}

//AddCopied adds to copied map
func (r *Response) AddCopied(key, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Copied[key] = value
}

//AddCopied adds to move map
func (r *Response) AddMoved(key, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Moved[key] = value
}

//AddCopied adds to move map
func (r *Response) AddInvoked(key, value string) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.Invoked[key] = value
}

//NewResponse create a response
func NewResponse() *Response {
	return &Response{
		Status:  base.StatusOK,
		Copied:  make(map[string]string),
		Moved:   make(map[string]string),
		Invoked: make(map[string]string),
		mux:     &sync.Mutex{},
	}
}
