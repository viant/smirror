package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"runtime/debug"
	"smirror/replay"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Println("Recovered in f", r)
		}
	}()
	lambda.Start(handleRequest)
}



func handleRequest(ctx context.Context, request *replay.Request) (*replay.Response, error) {
	service := replay.Singleton()
	return service.Replay(ctx, request), nil
}

