package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/viant/toolbox"
	"log"
	"os"
	"smirror/aws/proxy"
)

const DestEnvKey = "DEST"

func handleMessages(ctx context.Context, sqsEvent events.SQSEvent) (err error) {
	dest := os.Getenv(DestEnvKey)
	if dest == "" {
		log.Print("env.%v key was empty", DestEnvKey)
	}
	if len(sqsEvent.Records) == 0 {
		return err
	}

	proxier, err := proxy.Singleton()
	if err != nil {
		return err
	}
	toolbox.Dump(sqsEvent.Records)
	for _, record := range sqsEvent.Records {

		response := proxier.Do(ctx, dest, []byte(record.Body))
		if data, err := json.Marshal(response); err == nil {
			fmt.Printf("%v\n", string(data))
		}
	}
	return err
}

func main() {
	lambda.Start(handleMessages)
}
