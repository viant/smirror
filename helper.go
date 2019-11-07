package smirror

import "strings"

const (
	notFound     = "404"
	backendError = "backendError"
)

//IsNotFound returns true if not found error
func IsNotFound(message string) bool {
	if message == "" {
		return false
	}
	return strings.Contains(message, notFound)
}

//IsBackendError returns true if backend error
func IsBackendError(message string) bool {
	if message == "" {
		return false
	}
	return strings.Contains(message, backendError)
}
