package smirror

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afsc/gs"
	"github.com/viant/afsc/s3"
	"github.com/viant/toolbox"
	"os"
	"smirror/config"
	"strings"
)

//ConfigEnvKey config eng key
const ConfigEnvKey = "CONFIG"

//Config represents routes
type Config struct {
	Routes config.Routes
	//SourceScheme, currently gs or s3
	SourceScheme string
	ProjectID    string
	Region       string
}

//Init initialises routes
func (c *Config) Init() error {
	var projectID string
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

	for i := range c.Routes {
		c.Routes[i].Dest.Init(projectID)
	}

	return nil
}

//UseMessageDest returns true if any routes uses message bus
func (c *Config) UseMessageDest() bool {
	for _, resource := range c.Routes {
		if resource.Dest.Topic != "" {
			return true
		}
	}
	return false
}

//Resources returns
func (c *Config) Resources() []*config.Resource {
	var result = make([]*config.Resource, 0)

	for _, resource := range c.Routes {
		if resource.Source != nil {
			result = append(result, resource.Source)
		}
		if resource.Dest.Credentials != nil || resource.Dest.CustomKey != nil {
			result = append(result, &resource.Dest)
		}
	}
	return result
}

//NewConfigFromEnv returns new config from env
func NewConfigFromEnv(ctx context.Context, key string) (*Config, error) {
	JSONOrURL := strings.TrimSpace(os.Getenv(key))
	if toolbox.IsStructuredJSON(JSONOrURL) {
		return NewConfigFromJSON(JSONOrURL)
	}
	return NewConfigFromURL(ctx, JSONOrURL)
}

//NewConfigFromJSON creates a new config from env
func NewConfigFromJSON(payload string) (*Config, error) {
	cfg := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(cfg)
	if err == nil {
		err = cfg.Init()
	}
	return cfg, err
}

//NewConfigFromURL creates a new config from env
func NewConfigFromURL(ctx context.Context, URL string) (*Config, error) {
	service := afs.New()
	reader, err := service.DownloadWithURL(ctx, URL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download: %v", URL)
	}
	cfg := &Config{}
	err = json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode: %v ", URL)
	}
	return cfg, cfg.Init()
}
