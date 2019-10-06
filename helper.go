package smirror

import "strings"

const (
	notFoundMessage = "notfound"
)

//IsNotFound returns true if not found error
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	message := strings.Replace(strings.ToLower(err.Error()), " ", "", len(err.Error()))
	return strings.Contains(message, notFoundMessage)
}
