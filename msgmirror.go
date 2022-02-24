package smirror

import (
	"cloud.google.com/go/functions/metadata"
	"context"
	"encoding/json"
	"fmt"
	"github.com/viant/smirror/msg"
	"github.com/viant/toolbox"
	"os"
)

//MessageMirror represents a message mirror
func MessageMirror(ctx context.Context, event *msg.Request) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	event.EventID = meta.EventID
	format := os.Getenv("FORMAT")
	validate := toolbox.AsBoolean(os.Getenv("VALIDATE"))
	destURL := os.Getenv("DEST_URL")
	config := msg.NewConfig(format, validate, destURL)
	if err := config.RunValidation(); err != nil {
		return err
	}
	service := msg.Singleton(config)
	response := service.Proxy(ctx, event)
	JSON, _ := json.Marshal(response)
	fmt.Printf("%s\n", JSON)
	return nil
}
