package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"log"
	"os"
	"smirror/base"
	"smirror/contract"
	"smirror/event"
)

//StorageMirror cloud function entry point
func StorageMirror(ctx context.Context, event event.StorageEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	_, err = storageMirror(ctx, event)
	if err != nil {
		err = errors.Wrap(err, "failed to mirror "+event.URL())
		return err
	}
	return err
}

func outputResponse(response *contract.Response) {
	if response == nil {
		return
	}
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%v\n", string(data))
	}
}

//StorageMirrorSubscriber cloud function entry point
func StorageMirrorSubscriber(ctx context.Context, event event.PubsubBucketNotification) (err error) {
	storageEvent := event.StorageEvent()
	if storageEvent == nil {
		JSON, _ := json.Marshal(event)
		log.Printf("storage event was empty: %s", JSON)
		return nil
	}
	return StorageMirror(ctx, *storageEvent)
}

var fs afs.Service

func proxy(ctx context.Context, destination string, evnt event.StorageEvent) (*contract.Response, error) {
	if fs == nil {
		fs = afs.New()
	}
	destBucket := url.Host(destination)
	response := contract.NewResponse(evnt.URL())
	response.DestURLs = []string{evnt.ProxyDestURL(destBucket)}
	response.Status = base.StatusProxy
	isMove := os.Getenv(base.ProxyMethod) == base.ProxyMethodMove
	err := base.Trigger(ctx, fs, isMove, evnt.URL(), evnt.ProxyDestURL(destBucket), response.Triggered)
	if err != nil {
		if exists, e := fs.Exists(ctx, evnt.URL()); e == nil && ! exists {
			response.Status = base.StatusNoFound
			response.Error = ""
			response.NotFoundError = err.Error()
			return response, nil
		}
		err = errors.Wrapf(err, "failed to copy: %v to %v", evnt.URL(), evnt.ProxyDestURL(destBucket))
	}
	return response, err
}

func storageMirror(ctx context.Context, event event.StorageEvent) (response *contract.Response, err error) {
	destination := os.Getenv(base.DestEnvKey)
	if base.IsURL(destination) {
		if response, err = proxy(ctx, destination, event); err != nil {
			return response, err
		}
	} else {
		service, err := NewFromEnv(ctx, base.ConfigEnvKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage mirror: %v", err)
		}
		response = service.Mirror(ctx, contract.NewRequest(event.URL()))
	}
	outputResponse(response)
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
