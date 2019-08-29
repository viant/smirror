package config

import (
	"fmt"
	"strings"
)

//Represents a destination
type Resource struct {
	URL         string
	CustomKey   *CustomKey
	Credentials *Credentials

	Topic string
	//Optional pubsub project ID, otherwise it uses default one.
	ProjectID string
}

func (r *Resource) Init(projectID string) {
	if r.Topic == "" {
		return
	}
	if r.ProjectID == "" {
		r.ProjectID = projectID
	}
	if r.Topic != "" {
		if ! strings.Contains(r.Topic, "/") && r.ProjectID != "" {
			r.Topic = fmt.Sprintf("projects/%s/topics/%s", r.ProjectID, r.Topic)
		}
	}

	if r.ProjectID == "" {
		if elements := strings.Split(r.Topic, ""); len(elements) == 4 {
			r.ProjectID = elements[1]
		}
	}
}
