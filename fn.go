package mirror

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
	"os"

	"runtime/debug"
	"smirror/gs"
	"time"
)

const ConfigEnvKey = "CONFIG"
const LoggingKey = "LOGGING"

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

	if isFnLoggingEnabled(LoggingKey) {
		fmt.Printf("Mirrorred %v: %v in %v", response.Status, event.URL(), elapsed)
	}
	return err
}

func isFnLoggingEnabled(key string) bool {
	return toolbox.AsBoolean(os.Getenv(key))
}

func fn(ctx context.Context, event gs.Event) (*Response, error) {
	fmt.Printf("Start file %v\n", event.URL())
	service, err := NewFromEnv(ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	response := service.Mirror(NewRequest(event.URL()))
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
