package smirror

import (
	"os"
	"strings"
)

//LoggingEnvKey logging key
const LoggingEnvKey = "LOGGING"

//IsFnLoggingEnabled returns true if logging is enabled
func IsFnLoggingEnabled(key string) bool {
	return strings.ToLower(os.Getenv(key)) == "true"
}
