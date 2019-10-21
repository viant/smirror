package smirror

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/toolbox"
	"os"
	"smirror/auth"
	"smirror/base"
	"smirror/config"
	"strings"
)

//Config represents routes
type Config struct {
	base.Config
	SlackCredentials *auth.Credentials
	Mirrors          config.Ruleset
}

//Init initialises routes
func (c *Config) Init(ctx context.Context, fs afs.Service) (err error) {
	c.Config.Init()
	if err = c.Mirrors.Init(ctx, fs, c.ProjectID); err != nil {
		return err
	}
	for i := range c.Mirrors.Rules {
		c.Mirrors.Rules[i].Dest.Init(c.ProjectID)
	}
	return nil
}

//UseMessageDest returns true if any routes uses message bus
func (c *Config) UseMessageDest() bool {
	for _, resource := range c.Mirrors.Rules {
		if resource.Dest.Topic != "" || resource.Dest.Queue != "" {
			return true
		}
	}
	return false
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
		err = cfg.Init(ctx, afs.New())
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
	return cfg, cfg.Init(ctx, afs.New())
}
