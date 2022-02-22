package base

import (
	"cloud.google.com/go/compute/metadata"
	"context"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"golang.org/x/oauth2/google"
	"os"
)

//Config represents a base config
type Config struct {
	URL          string
	Region       string
	ProjectID    string
	SourceScheme string
}

func (c *Config) Init() {
	if c.Region == "" {
		c.Region = os.Getenv("FUNCTION_REGION")
	}
	if c.SourceScheme == "" {
		projectID, err := metadata.ProjectID()
		if err == nil && projectID != "" {
			c.SourceScheme = gs.Scheme
			if c.ProjectID == "" {
				c.ProjectID = projectID
			}
		} else if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
			c.SourceScheme = s3.Scheme
		}
	}
	if c.ProjectID == "" {
		if credentials, err := google.FindDefaultCredentials(context.Background()); err == nil {
			c.ProjectID = credentials.ProjectID
		}
	}
}
