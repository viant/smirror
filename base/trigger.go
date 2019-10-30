package base

import (
	"context"
	"github.com/viant/afs"
	"github.com/viant/afs/storage"
)

//Trigger triggers storage event
func Trigger(ctx context.Context, fs afs.Service, isMove bool, sourceURL, destURL string, triggered map[string]string, options ...storage.Option) error {
	triggerFunc := fs.Copy
	if isMove {
		triggerFunc = fs.Move
	}
	triggered[sourceURL] = destURL
	err := triggerFunc(ctx, sourceURL, destURL, options...)
	if exists, e := fs.Exists(ctx, sourceURL); e == nil && !exists {
		err = nil
		triggered[sourceURL] = StatusNoFound
	}
	return err
}
