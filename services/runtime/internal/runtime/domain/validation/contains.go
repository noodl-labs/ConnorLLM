package validation

import "strings"

// Contains reports whether body includes want (case-sensitive, beta.2).
func Contains(body, want string) bool {
	if want == "" {
		return true // gate disabled
	}
	return strings.Contains(body, want)
}
