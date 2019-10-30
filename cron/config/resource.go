package config

import (
	"smirror/config"
)

//Rule represents a cron resource
type Rule struct {
	Source config.Resource
	Dest   string
	Move   bool `json:",omitempty"`
}
