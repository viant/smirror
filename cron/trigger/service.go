package trigger

import (
	"context"
	"github.com/viant/afs/storage"
	"smirror/cron/config"
)

//Service reresents trigger service
type Service interface {
	Trigger(ctx context.Context, resource *config.Rule, eventSource storage.Object) error
}
