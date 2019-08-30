package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"runtime/debug"
	"smirror/gs"
	"time"
)

//Fn cloud function entry point
func Fn(ctx context.Context, event gs.Event) (err error) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = fmt.Errorf("%v", r)
		}
	}()
	response, err := fn(ctx, event)
	elapsed := time.Since(start)
	if err != nil {
		err = errors.Wrap(err, "failed to mirror "+event.URL())
		return err
	}

	if IsFnLoggingEnabled(LoggingEnvKey) {
		fmt.Printf("mirrored %v: %v in %v", response.Status, event.URL(), elapsed)
	}
	return err
}

func fn(ctx context.Context, event gs.Event) (*Response, error) {
	service, err := NewFromEnv(ctx, ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	if IsFnLoggingEnabled(LoggingEnvKey) {
		fmt.Printf("uses service %T(%p), err: %v\n", service, service, err)
	}
	response := service.Mirror(ctx, NewRequest(event.URL()))
	if IsFnLoggingEnabled(LoggingEnvKey) {
		if data, err := json.Marshal(response); err == nil {
			fmt.Print(string(data))
		}
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
