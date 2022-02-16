package base

import (
	"context"
	"fmt"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
)

var cfMetaURL = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
var cfEnvironmentVariableKeys = []string{"GOOGLE_CLOUD_PROJECT", "GCLOUD_PROJECT", "GCP_PROJECT"}

//Config represents a base config
type Config struct {
	URL          string
	Region       string
	ProjectID    string
	SourceScheme string
}

func (c *Config) Init() {
	var projectID string
	if c.Region == "" {
		c.Region = os.Getenv("FUNCTION_REGION")
	}
	if c.SourceScheme == "" {
		if projectID = gcpProjectId(); projectID != "" {
			os.Setenv("GOOGLE_CLOUD_PROJECT", projectID)
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

func gcpProjectId() string {
	if projectId := gcpProjectIdFromEnvKeys(); projectId != "" {
		return projectId
	}
	return gcpProjectIdFromURL()
}

func gcpProjectIdFromEnvKeys() string {
	var projectID = ""
	for _, key := range cfEnvironmentVariableKeys {
		if projectID = os.Getenv(key); projectID != "" {
			break
		}
	}
	return projectID
}

func gcpProjectIdFromURL() string {
	req, err := http.NewRequest("GET", cfMetaURL, nil)
	if err != nil {
		fmt.Printf("couldn't make http request to gcp metaurl: %v due to error: %v", cfMetaURL, err)
		return ""
	}
	req.Header.Set("Metadata-Flavor", "Google")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("error getting response from gcp meta request due to error: %v", err)
		return ""
	}
	data, _ := ioutil.ReadAll(res.Body)
	return string(data)
}
