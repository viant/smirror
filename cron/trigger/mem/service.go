package mem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/afs"
	"github.com/viant/afs/mem"
	"github.com/viant/afs/storage"
	"smirror/cron/config"
	"smirror/cron/trigger"
)

//Scheme represents mem scheme
const Scheme = mem.Scheme

type service struct {
	afs.Service
}

//Trigger triggers lambda execution
func (s *service) Trigger(ctx context.Context, resource *config.Rule, eventSource storage.Object) error {
	URL := fmt.Sprintf("%v://localhost/%v", mem.Scheme, resource.Dest)
	event := Event{URL: eventSource.URL(), Size: eventSource.Size()}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.Upload(ctx, URL, 0644, bytes.NewReader(payload))
}

//New creates a new memory trigger destination
func New(fs afs.Service) trigger.Service {
	return &service{Service:fs}
}
