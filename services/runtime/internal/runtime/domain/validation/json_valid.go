package validation

import "encoding/json"

// Check reports whether s is valid JSON (syntax only, not schema).
func Check(s string) bool {
	if s == "" {
		return false
	}
	return json.Valid([]byte(s))
}
