package config

import (
	"fmt"
	"github.com/viant/afs/matcher"
	"strings"
)

//Represents a destination
type Resource struct {
	matcher.Basic
	URL         string       `json:",omitempty"`
	Region      string       `json:",omitempty"`
	CustomKey   *CustomKey   `json:",omitempty"`
	Credentials *Credentials `json:",omitempty"`
	Topic       string       `json:",omitempty"`
	//Optional pubsub project ID, otherwise it uses default one.
	ProjectID string `json:",omitempty"`
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
