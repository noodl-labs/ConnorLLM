package validation

import "strings"

// Contains reports whether body includes want. When ignoreCase is false, match is case-sensitive.

func Contains(body, want string, ignoreCase bool) bool {
	if want == "" {
		return true
	}
	if ignoreCase {
		return strings.Contains(
			strings.ToLower(body),
			strings.ToLower(want),
		)
	}
	return strings.Contains(body, want)
}
