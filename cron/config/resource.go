package config

import (
	"smirror/config"
)

//Rule represents a cron resource
type Rule struct {
	config.Resource
	Dest string
}


