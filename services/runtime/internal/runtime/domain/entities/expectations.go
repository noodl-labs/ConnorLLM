package entities

import "encoding/json"

// Expectations lists optional gates for one case (YAML → domain).
type Expectations struct {
	Contains           string // non-empty → body must contain this substring
	ContainsIgnoreCase bool   // when true, contains match is case-insensitive
	JSON               bool   // true → body must be valid JSON
	JSONSchema         json.RawMessage // non-empty → body must match this JSON Schema
}

// ExpectationsFromCase maps benchmark case fields to domain expectations.
func ExpectationsFromCase(
	expectContains string,
	expectJSON, containsIgnoreCase bool,
	jsonSchema json.RawMessage,
) Expectations {
	return Expectations{
		Contains:           expectContains,
		ContainsIgnoreCase: containsIgnoreCase,
		JSON:               expectJSON,
		JSONSchema:         jsonSchema,
	}
}

// RequiresJSONSyntax reports whether the body must be valid JSON before schema checks.
func (e Expectations) RequiresJSONSyntax() bool {
	return e.JSON || len(e.JSONSchema) > 0
}
