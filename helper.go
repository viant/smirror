package smirror

import "strings"

const (
	notFoundMessage = "notfound"
	backendError = "backendError"
)

//IsNotFound returns true if not found error
func IsNotFound(message string) bool {
	if message == "" {
		return false
	}
	message = strings.Replace(strings.ToLower(message), " ", "", len(message))
	return strings.Contains(message, notFoundMessage)
}



//IsBackendError returns true if backend error
func IsBackendError(message string) bool {
	if message == "" {
		return false
	}
	return strings.Contains(message, backendError)
}


