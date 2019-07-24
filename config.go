package mirror

import (
	"smirror/secret"
)



//Config represents routes
type Config struct {
	Routes Routes
	Secrets []*secret.Config
}


