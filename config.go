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
}


//Init initialises routes
func (c *Config) Init() error {
	if c.SourceScheme == "" {
		if os.Getenv("GCLOUD_PROJECT") != "" {
			c.SourceScheme = gs.Scheme
		} else if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
			c.SourceScheme = s3.Scheme
		}
	}
	return nil
}


//Resources returns
func (c *Config) Resources() []*config.Resource {
	var result = make([]*config.Resource,0)
	for _, resource := range c.Routes {
		if resource.Source != nil {
			result = append(result, resource.Source)
		}
		if resource.Dest.Credentials != nil || resource.Dest.CustomKey != nil {
			result = append(result, resource.Source)
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
	config := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(config)
	if err == nil {
		err = config.Init()
	}
	return config, err
}

//NewConfigFromURL creates a new config from env
func NewConfigFromURL(ctx context.Context, URL string) (*Config, error) {
	service := afs.New()
	reader, err := service.DownloadWithURL(ctx, URL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download: %v", URL)
	}
	config := &Config{}
	err = json.NewDecoder(reader).Decode(config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode: %v ", URL)
	}
	return config, config.Init()
}
