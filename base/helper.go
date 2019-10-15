package base

import "strings"

//IntPtr returns int pointer
func IntPtr(i int) *int {
	return &i
}

//IsURL returns true if candidate is URL
func IsURL(candidate string) bool {
	return strings.Contains(candidate, "://")
}
