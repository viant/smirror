package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

func storageMirror(ctx context.Context, event event.StorageEvent) (response *contract.Response, err error) {
	service, err := NewFromEnv(ctx, base.ConfigEnvKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage mirror: %v", err)
	}
	response = service.Mirror(ctx, contract.NewRequest(event.URL()))
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%v\n", string(data))
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
