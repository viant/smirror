package base

import (
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"os"
)

//Config represents a base config
type Config struct {
	URL string
	Region string
	ProjectID    string
	SourceScheme string
}

func (c *Config) Init() {
	var projectID string
	if c.Region== "" {
		c.Region = os.Getenv("FUNCTION_REGION")
	}
	if c.SourceScheme == "" {
		if projectID = os.Getenv("GCLOUD_PROJECT"); projectID != "" {
			c.SourceScheme = gs.Scheme
			if c.ProjectID == "" {
				c.ProjectID = projectID
			}

		} else if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
			c.SourceScheme = s3.Scheme
		}
	}
}
