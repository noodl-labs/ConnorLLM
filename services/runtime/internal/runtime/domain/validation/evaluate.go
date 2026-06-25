package validation

import "github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"

// Evaluate runs gates in order after a successful HTTP response.
// Returns (passed, failReason).
func Evaluate(body string, exp entities.Expectations) (bool, entities.FailReason) {
	if exp.Contains != "" && !Contains(body, exp.Contains, exp.ContainsIgnoreCase) {
		return false, entities.FailReasonContentMismatch
	}
	if exp.RequiresJSONSyntax() && !Check(body) {
		return false, entities.FailReasonInvalidJSON
	}
	if len(exp.JSONSchema) > 0 && !CheckSchema(body, exp.JSONSchema) {
		return false, entities.FailReasonSchemaMismatch
	}
	return true, entities.FailReasonNone
}
