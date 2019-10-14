package slack

import (
	"context"
	"github.com/viant/afs"
	"smirror/auth"
	"smirror/job"
	"smirror/secret"
)

//Slack represents storage service
type Slack interface {
	Notify(ctx context.Context, request *job.NotifyRequest) error
}

type service struct {
	projectID   string
	Region      string
	Secret      secret.Service
	Storage     afs.Service
	Credentials *auth.Credentials
}

//NewSlack creates slack service
func NewSlack(region, projectID string, storageService afs.Service, secretService secret.Service, credentials *auth.Credentials) Slack {
	return &service{
		Region:      region,
		projectID:   projectID,
		Secret:      secretService,
		Storage:     storageService,
		Credentials: credentials,
	}
}
