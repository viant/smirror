package config

import (
	"fmt"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/option"
	"smirror/auth"
	"strings"
)

//Represents a destination
type Resource struct {
	matcher.Basic
	Bucket      string            `json:",omitempty"`
	URL         string            `json:",omitempty"`
	Region      string            `json:",omitempty"`
	CustomKey   *CustomKey        `json:",omitempty"`
	Credentials *auth.Credentials `json:",omitempty"`
	Proxy       *option.Proxy
	Topic       string `json:",omitempty"`
	Queue       string `json:",omitempty"`
	//Optional pubsub project ID, otherwise it uses default one.
	ProjectID string `json:",omitempty"`
}

//CloneWithURL clone resource with URL
func (r Resource) CloneWithURL(URL string) *Resource {
	return &Resource{
		Basic:       r.Basic,
		URL:         URL,
		Region:      r.Region,
		CustomKey:   r.CustomKey,
		Credentials: r.Credentials,
		Topic:       r.Topic,
		Queue:       r.Queue,
		ProjectID:   r.ProjectID,
	}
}

func (r *Resource) Init(projectID string) {
	if r.Topic == "" {
		return
	}
	if r.ProjectID == "" {
		r.ProjectID = projectID
	}
	if r.Topic != "" {
		if !strings.Contains(r.Topic, "/") && r.ProjectID != "" {
			r.Topic = fmt.Sprintf("projects/%s/topics/%s", r.ProjectID, r.Topic)
		}
	}
	if r.ProjectID == "" {
		if elements := strings.Split(r.Topic, ""); len(elements) == 4 {
			r.ProjectID = elements[1]
		}
	}
}
