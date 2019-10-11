package smirror

import "strings"

const (
	notFoundMessage = "notfound"
)

//IsNotFound returns true if not found error
func IsNotFound(message string) bool {
	if message == "" {
		return false
	}
	message = strings.Replace(strings.ToLower(message), " ", "", len(message))
	return strings.Contains(message, notFoundMessage)
}

