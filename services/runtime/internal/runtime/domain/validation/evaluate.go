package validation

import "github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"

// Evaluate runs gates in order after a successful HTTP response.
// Returns (passed, failReason).
func Evaluate(body string, exp entities.Expectations) (bool, entities.FailReason) {
	if exp.Contains != "" && !Contains(body, exp.Contains) {
		return false, entities.FailReasonContentMismatch
	}
	if exp.JSON && !Check(body) {
		return false, entities.FailReasonInvalidJSON
	}
	return true, entities.FailReasonNone
}
