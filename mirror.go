package smirror

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"smirror/base"
	"smirror/contract"
	"smirror/event"
	"smirror/shared"
)

//StorageMirror cloud function entry point
func StorageMirror(ctx context.Context, event event.StorageEvent) (err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		err = fmt.Errorf("%v", r)
	//	}
	//}()
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
	shared.LogLn(response)
	//Schema error
	if response.Error != "" && response.SchemaError == ""{
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
