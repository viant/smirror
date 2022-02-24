package proxy

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/afs"
	"github.com/viant/toolbox"
	"os"
	"github.com/viant/smirror/base"
	"github.com/viant/smirror/config"
	"strings"
)

//Config represents proxy config
type Config struct {
	base.Config
	Dest   config.Resource
	Source config.Resource
	Move   bool
}

//Validate checks if config is valid
func (c *Config) Validate() error {
	c.Init()
	if c.Dest.URL == "" {
		return errors.Errorf("dest.url was empty")
	}
	return nil
}

//NewConfig creates a config
func NewConfig(ctx context.Context) (*Config, error) {
	if os.Getenv(base.ConfigEnvKey) != "" {
		return NewConfigFromEnv(ctx, base.ConfigEnvKey)
	}
	cfg := &Config{
		Dest: config.Resource{
			URL: os.Getenv(base.DestEnvKey),
		},
	}
	return cfg, cfg.Validate()
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
		err = cfg.Validate()
	}
	return cfg, err
}

//NewConfigFromURL creates a new config from env
func NewConfigFromURL(ctx context.Context, URL string) (*Config, error) {
	service := afs.New()
	reader, err := service.OpenURL(ctx, URL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download: %v", URL)
	}
	cfg := &Config{}
	err = json.NewDecoder(reader).Decode(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode: %v ", URL)
	}
	err = cfg.Validate()
	return cfg, err
}
