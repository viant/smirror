package job

import (
	"context"
	"errors"
	"fmt"
	"smirror/auth"
	"strings"
)

const DefaultRegion = "us-central1"

//Notify represents notify function
type Notify func(ctx context.Context, request *NotifyRequest) error

//NotifyRequest represents a notify request
type NotifyRequest struct {
	Channels []string
	From     string
	Title    string
	Message  string
	Filename string
	Body     interface{}
	*auth.Credentials
	BodyType string
}

//Init initializes request
func (r *NotifyRequest) Init(location, projectID string) error {
	if location == "" {
		location = DefaultRegion
	}
	if r.Credentials != nil {
		if strings.Count(r.Secret.Key, "/") == 1 {
			pair := strings.Split(r.Secret.Key, "/")
			ring := strings.TrimSpace(pair[0])
			key := strings.TrimSpace(pair[1])
			r.Secret.Key = fmt.Sprintf("projects/%v/locations/%v/keyRings/%v/cryptoKeys/%v", projectID, location, ring, key)
		}
	}
	return nil
}

//Validate checks if request is valid
func (r *NotifyRequest) Validate() error {
	if r.Credentials == nil || (r.Credentials.Token == "" && r.Credentials.Secret.Key == "") {
		return errors.New("notify.secret was empty")
	}
	if len(r.Channels) == 0 {
		return errors.New("channels was empty")
	}
	if r.Title == "" {
		return errors.New("title was empty")
	}
	return nil
}
