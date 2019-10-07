package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	flambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
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
	for _, message := range sqsEvent.Records {
		if err = notify(dest, []byte(message.Body)); err != nil {
			return err
		}
	}
	return err
}

func main() {
	lambda.Start(handleMessages)
}

func notify(destination string, payload []byte) error {
	if destination == "" {
		log.Println("deast is empty, ingoring event: %s\n", payload)
		return nil
	}
	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	service := flambda.New(sess)
	_, err = service.Invoke(&flambda.InvokeInput{
		FunctionName:   &destination,
		Payload:        payload,
		InvocationType: aws.String(flambda.InvocationTypeEvent),
	})
	return err
}
