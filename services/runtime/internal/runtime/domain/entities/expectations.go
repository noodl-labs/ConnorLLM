package entities

// Expectations lists optional gates for one case (YAML → domain).
type Expectations struct {
	Contains string // non-empty → body must contain this substring
	JSON     bool   // true → body must be valid JSON
}

// ExpectationsFromCase maps benchmark case fields to domain expectations.
func ExpectationsFromCase(expectContains string, expectJSON bool) Expectations {
	return Expectations{
		Contains: expectContains,
		JSON:     expectJSON,
	}
}
