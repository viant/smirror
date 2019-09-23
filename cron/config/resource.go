package config

import "smirror/config"

//Resource represents a cron resource
type Resource struct {
	config.Resource
	DestFunction string
}




