package event

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/viant/afs/url"
	"strings"
	"time"
)


//S3Event represents s3 event
type S3Event events.S3Event


//Each iterates each record calling supplied handler with the record URL
func (e S3Event) Each(handler func(URL string) error) error {
	for i:= range e.Records {
		URL := resourceURL(e.Records[i])
		if err := handler(URL);err != nil {
			return err
		}
	}
	return nil
}



//resourceURL returns resource URL
func resourceURL(resource events.S3EventRecord) string {
	return fmt.Sprintf("s3://%s/%s", resource.S3.Bucket.Name, resource.S3.Object.Key)
}


//NewS3EventFromJSON creates a new s3 events
func NewS3EventFromJSON(data []byte) (*S3Event, error) {
	s3Event := &S3Event{}
	return s3Event, json.Unmarshal(data, s3Event)
}


//NewS3EventForURL creates s3 events for supplied URL
func NewS3EventForURL(URL string) *S3Event {
	bucket := url.Host(URL)
	URLPath := url.Path(URL)
	s3Event := &S3Event{Records: make([]events.S3EventRecord, 0)}
	s3Event.Records = append(s3Event.Records, events.S3EventRecord{
		EventTime:   time.Now(),
		EventSource: "s3",
		S3: events.S3Entity{
			Bucket: events.S3Bucket{
				Name: bucket,
			},
			Object: events.S3Object{
				Key: strings.Trim(URLPath, "/"),
			},
		},
	})
	return s3Event
}
