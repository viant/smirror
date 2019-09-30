package smirror

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/afs/matcher"
	"github.com/viant/afs/storage"
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
	Routes        config.Routes
	RoutesBaseURL string
	//SourceScheme, currently gs or s3
	SourceScheme string
	ProjectID    string
	Region       string
}

func (c *Config) loadRoutes(ctx context.Context) error {
	if c.RoutesBaseURL == "" {
		return nil
	}
	fs := afs.New()

	suffixMatcher, _ := matcher.NewBasic("", ".json", "", nil)
	routesObject, err := fs.List(ctx, c.RoutesBaseURL, suffixMatcher)
	if err != nil {
		return err
	}
	for _, object := range routesObject {
		if err = c.loadRoute(ctx, fs, object); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) loadRoute(ctx context.Context, storage afs.Service, object storage.Object) error {
	reader, err := storage.Download(ctx, object)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	routes := config.Routes{}
	if err = json.NewDecoder(reader).Decode(&routes); err == nil {
		c.Routes = append(c.Routes, routes...)
	}
	return err
}

//Init initialises routes
func (c *Config) Init(ctx context.Context) error {
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
	if len(c.Routes) == 0 {
		c.Routes = make([]*config.Route, 0)
	}
	if err := c.loadRoutes(ctx); err != nil {
		return err
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
		return NewConfigFromJSON(ctx, JSONOrURL)
	}
	return NewConfigFromURL(ctx, JSONOrURL)
}

//NewConfigFromJSON creates a new config from env
func NewConfigFromJSON(ctx context.Context, payload string) (*Config, error) {
	cfg := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(cfg)
	if err == nil {
		err = cfg.Init(ctx)
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
	return cfg, cfg.Init(ctx)
}
