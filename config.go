package mirror

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/storage"
	"os"
	"smirror/secret"
	"strings"
)

//Config represents routes
type Config struct {
	Routes  Routes
	Secrets []*secret.Config
}

//NewConfigFromEnv returns new config from env
func NewConfigFromEnv(key string) (*Config, error) {
	JSONOrURL := strings.TrimSpace(os.Getenv(key))
	if toolbox.IsStructuredJSON(JSONOrURL) {
		return NewConfigFromJSON(JSONOrURL)
	}
	return NewConfigFromURL(JSONOrURL)
}

//NewConfigFromJSON creates a new config from env
func NewConfigFromJSON(payload string) (*Config, error) {
	config := &Config{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(config)
	return config, err
}

//NewConfigFromURL creates a new config from env
func NewConfigFromURL(URL string) (*Config, error) {
	storageService, err := storage.NewServiceForURL(URL, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get storage service "+URL)
	}
	reader, err := storageService.DownloadWithURL(URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download "+URL)
	}
	config := &Config{}
	err = json.NewDecoder(reader).Decode(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode "+URL)
	}
	return config, err
}
