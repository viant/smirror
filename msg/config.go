package msg

import (
	"fmt"
	"strings"
)

type Config struct {
	SourceFormat string
	Validate     bool
	DestURL      string
}

//IsSourceJSON returns true if source format is json
func (c *Config) IsSourceJSON() bool {
	return strings.ToUpper(c.SourceFormat) == "JSON"
}

func (c *Config) RunValidation() error {
	if c.DestURL == "" {
		return fmt.Errorf("destURL was empty")
	}
	return nil
}

//NewConfig return a config
func NewConfig(format string, validate bool, dest string) *Config {
	return &Config{
		SourceFormat: format,
		Validate:     validate,
		DestURL:      dest,
	}
}
