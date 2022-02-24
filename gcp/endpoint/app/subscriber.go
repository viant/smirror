package main

import (
	"context"
	"github.com/viant/afs"
	"log"
	"github.com/viant/smirror/gcp/endpoint"
)

const envConfig = "APP_CONFIG"

func main() {
	fs := afs.New()
	config, err := endpoint.NewConfigFromEnv(envConfig)
	if err != nil {
		log.Fatalf("failed to load conifg: %v %v", envConfig, err)
	}
	err = config.Init(context.Background(), fs)
	if err != nil {
		log.Fatalf("failed to init conifg: %v %v", envConfig, err)
	}
	err = config.Validate()
	if err != nil {
		log.Fatalf("failed to validate conifg: %v %v", envConfig, err)
	}
	srv, err := endpoint.New(config,  fs)
	if err != nil {
		log.Fatalf("failed to create subscriber service: %v", err)
	}
	err = srv.Consume(context.Background())
	if err != nil {
		log.Fatalf("failed to run service: %v ", err)
	}
}
