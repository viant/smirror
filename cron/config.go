package cron

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
	"smirror/cron/config"
	"smirror/cron/trigger/mem"
	"strings"
)

//Config represents cron config
type Config struct {
	ProjectID string
	MetaURL string
	TimeWindow *config.TimeWindow
	Resources []*config.Resource
	ResourcesBaseURL string
	SourceScheme string
}

func (c *Config) loadAllResources(ctx context.Context) error {
	if c.ResourcesBaseURL == "" {
		return nil
	}
	fs := afs.New()

	suffixMatcher, _ := matcher.NewBasic("", ".json", "")
	routesObject, err := fs.List(ctx, c.ResourcesBaseURL, suffixMatcher)
	if err != nil {
		return err
	}
	for _, object := range routesObject {
		if err = c.loadResources(ctx, fs, object);err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) loadResources(ctx context.Context, storage afs.Service, object storage.Object) error {
	reader, err := storage.Download(ctx, object)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	resources := make([]*config.Resource, 0)
	if err = json.NewDecoder(reader).Decode(&resources);err == nil {
		c.Resources = append(c.Resources, resources...)
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
		} else {
			c.SourceScheme = mem.Scheme
		}
	}
	if len(c.Resources) == 0 {
		c.Resources = make([]*config.Resource, 0)
	}
	if err := c.loadAllResources(ctx);err != nil {
		return err
	}
	for i := range c.Resources {
		c.Resources[i].Init(projectID)
	}
	return nil
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
