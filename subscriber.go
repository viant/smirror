package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"smirror/event"
	"smirror/proxy"
)

var proxyConfug *proxy.Config

//StorageMirrorSubscriber cloud function entry point
func StorageMirrorSubscriber(ctx context.Context, event event.PubsubBucketNotification) (err error) {
	storageEvent := event.StorageEvent()
	if storageEvent == nil {
		JSON, _ := json.Marshal(event)
		log.Printf("storage event was empty: %s", JSON)
		return nil
	}
	if proxyConfug == nil {
		proxyConfug, err = proxy.NewConfig(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to create config")
		}
	}
	proxier := proxy.Singleton(proxyConfug)
	response := proxier.Proxy(ctx, &proxy.Request{
		Source: proxyConfug.Source.CloneWithURL(storageEvent.URL()),
		Dest:   &proxyConfug.Dest,
		Move:   proxyConfug.Move,
	})
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%v\n", string(data))
	}
	return nil
}
