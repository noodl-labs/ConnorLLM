package entities

// Expectations lists optional gates for one case (YAML → domain).
type Expectations struct {
	Contains           string // non-empty → body must contain this substring
	ContainsIgnoreCase bool   // when true, contains match is case-insensitive
	JSON               bool   // true → body must be valid JSON
}

// ExpectationsFromCase maps benchmark case fields to domain expectations.
func ExpectationsFromCase(expectContains string, expectJSON, containsIgnoreCase bool) Expectations {
	return Expectations{
		Contains:           expectContains,
		ContainsIgnoreCase: containsIgnoreCase,
		JSON:               expectJSON,
	}
}
