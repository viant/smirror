package msg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"smirror/base"
	"strings"
)

//Service service to proxy pubsub topic payload to google storage with or without validation
type Service interface {
	Proxy(context.Context, *Request) *Response
}

type service struct {
	config *Config
	fs     afs.Service
}

func (p *service) Proxy(ctx context.Context, request *Request) *Response {
	response := NewResponse(request.EventID)
	err := p.proxy(ctx, request, response)
	if err != nil {
		response.Error = err.Error()
		response.Status = base.StatusError
	}
	return response
}

func (p *service) proxy(ctx context.Context, request *Request, response *Response) error {

	if p.config.Validate {
		switch strings.ToUpper(p.config.SourceFormat) {
		case "JSON":
			if !json.Valid(request.Data) {
				return fmt.Errorf("invaid JSON: %s", request.Data)
			}
		default:
			return fmt.Errorf("unsupported sourceFornat: %s", p.config.SourceFormat)
		}
	}
	ext := ""
	if p.config.IsSourceJSON() {
		ext = ".json"
	}
	URL := url.Join(p.config.DestURL, request.EventID+ext)
	response.URL = URL
	response.Size = len(request.Data)
	return p.fs.Upload(ctx, URL, 0666, bytes.NewReader(request.Data))
}

//New create a service service
func New(config *Config, fs afs.Service) Service {
	return &service{
		fs:     fs,
		config: config,
	}
}
