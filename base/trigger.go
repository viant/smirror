package base

import (
	"context"
	"github.com/viant/afs"
)

//Trigger triggers storage event
func Trigger(ctx context.Context, fs afs.Service, isMove bool, sourceURL, destURL string, triggered map[string]string) error {
	triggerFunc := fs.Copy
	if isMove {
		triggerFunc = fs.Move
	}
	triggered[sourceURL] = destURL
	err := triggerFunc(ctx, sourceURL, destURL)
	if exists, e := fs.Exists(ctx, sourceURL); e == nil && ! exists {
		err = nil
		triggered[sourceURL] = StatusNoFound
	}
	return err
}
