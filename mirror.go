package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"smirror/base"
	"smirror/event"
	"smirror/gs"
)

//StorageMirror cloud function entry point
func StorageMirror(ctx context.Context, event event.StorageEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	_, err = storageMirror(ctx, event)
	//elapsed := time.Since(start)
	if err != nil {
		err = errors.Wrap(err, "failed to mirror "+event.URL())
		return err
	}
	//if base.IsLoggingEnabled() {
	//	fmt.Printf("mirrored %v: %v in %v", response.Status, event.URL(), elapsed)
	//}
	return err
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



func storageMirror(ctx context.Context, event event.StorageEvent) (*Response, error) {
	service, err := NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	response := service.Mirror(ctx, NewRequest(event.URL()))
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%v\n", string(data))
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}


