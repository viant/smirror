package config

import (
	"github.com/viant/smirror/config"
)

//Rule represents a cron resource
type Rule struct {
	Source   config.Resource
	Dest     config.Resource
	Move     bool `json:",omitempty"`
}
