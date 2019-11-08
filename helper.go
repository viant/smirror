package smirror

import "strings"

const (
	notFoundCode = "404"
	notFound     = "not found"
	backendError = "backendError"
	connectionReset = "connection reset by peer"
)

//IsNotFound returns true if not found error
func IsNotFound(message string) bool {
	if message == "" {
		return false
	}
	return strings.Contains(message, notFound) || strings.Contains(message, notFoundCode)
}

//IsRetryError returns true if backend error
func IsRetryError(message string) bool {
	if message == "" {
		return false
	}
	return strings.Contains(message, backendError) || strings.Contains(message, connectionReset)
}
