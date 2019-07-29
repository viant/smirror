package smirror

import (
	"github.com/viant/toolbox"
	"os"
)

//LoggingEnvKey logging key
const LoggingEnvKey = "LOGGING"

//IsFnLoggingEnabled returns true if logging is enabled
func IsFnLoggingEnabled(key string) bool {
	return toolbox.AsBoolean(os.Getenv(key))
}
