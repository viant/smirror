package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/url"
	"log"
	"smirror/base"
)

const (
	ProxyTypeLambda  = "lambda"
	ProxyTypeStorage = "storage"
)

type Proxy interface {
	Do(ctx context.Context, destination string, payload []byte) *Response
}

type proxy struct {
	*lambda.Lambda
	fs afs.Service
}

//Do proxy payload to destination
func (p *proxy) Do(ctx context.Context, destination string, payload []byte) *Response {
	response := &Response{Status: base.StatusOK, Copy: make(map[string]string)}
	err := p.do(ctx, destination, payload, response)
	if err != nil {
		response.Status = base.StatusError
		response.Error = err.Error()
	}
	return response
}

//Do proxy payload to destination
func (p *proxy) do(ctx context.Context, destination string, payload []byte, response *Response) (err error) {
	if destination == "" {
		log.Println("dest is empty, ignoring event: %s\n", payload)
		return nil
	}
	response.Source = string(payload)
	response.Destination = destination
	if !base.IsURL(destination) {
		response.ProxyType = ProxyTypeLambda
		_, err = p.Invoke(&lambda.InvokeInput{
			FunctionName:   &destination,
			Payload:        payload,
			InvocationType: aws.String(lambda.InvocationTypeEvent),
		})
		return err
	}
	response.ProxyType = ProxyTypeStorage
	s3Event := &events.S3Event{}
	if err = json.Unmarshal(payload, s3Event); err != nil {
		return errors.Wrapf(err, "failed to decode %T, from %s", s3Event, payload)
	}
	destBucket := url.Host(destination)
	for _, record := range s3Event.Records {
		sourceURL := sourceURL(record)
		destURL := destinationURL(record, destBucket)
		response.Copy[sourceURL] = destURL
		if err := p.fs.Copy(ctx, sourceURL, destURL); err != nil {
			if exists, e := p.fs.Exists(ctx, sourceURL); e == nil && ! exists {
				response.Status = base.StatusNoFound
				return nil
			}
			return err
		}
	}
	return err
}

//sourceURL returns resource URL
func sourceURL(resource events.S3EventRecord) string {
	return fmt.Sprintf("s3://%s/%s", resource.S3.Bucket.Name, resource.S3.Object.Key)
}

//destinationURL returns resource URL
func destinationURL(resource events.S3EventRecord, destBucket string) string {
	return fmt.Sprintf("s3://%s/%s/%s", destBucket, resource.S3.Bucket.Name, resource.S3.Object.Key)
}

//New returns new proxy
func New() (Proxy, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &proxy{
		Lambda: lambda.New(sess),
		fs:     afs.New(),
	}, nil
}
