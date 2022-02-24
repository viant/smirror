package config

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/viant/afs/file"
	"github.com/viant/afs/url"
	"github.com/viant/smirror/event"
)

//Overflow defines overflow controls
type Overflow struct {
	SizeMB  int64
	DestURL string
	Queue   string
	ProjectID string
	Topic   string
}

func (o Overflow) MessageEvent(URL string) interface{} {
	baseURL, URLPath := url.Base(URL, file.Scheme)
	bucket := url.Host(baseURL)
	if o.Queue != "" {
		s3event := event.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: bucket,
						},
						Object: events.S3Object{
							Key: URLPath,
						},
					},
				},
			},
		}
		return s3event
	}
	gsEvent := event.StorageEvent{
		Bucket: bucket,
		Name: URLPath,
	}
	return gsEvent
}

func (o Overflow) MessageDest() string {
	if o.Topic != "" {
		return o.Topic
	}
	return o.Queue
}

func (o Overflow) Size() int64 {
	return o.SizeMB * 1024 * 1024
}
