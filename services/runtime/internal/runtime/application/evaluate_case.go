package application

import (
	"context"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/validation"
)

// EvaluateCase runs one benchmark case: provider call (ExecuteCase) then optional JSON check (L1).
//
// expectJSON false: Passed iff Response.Succeeded().
// expectJSON true:  Passed iff Succeeded() and validation.Check(Response.Body).
func EvaluateCase(
	ctx context.Context,
	caseID string,
	expectJSON bool,
	req entities.Request,
	timeout reliability.TimeoutPolicy,
	retry reliability.RetryPolicy,
	executor domain.ProviderExecutor,
) (entities.CaseResult, error) {
	resp, err := ExecuteCase(ctx, req, timeout, retry, executor)
	if err != nil {
		return entities.CaseResult{}, err
	}

	result := entities.CaseResult{
		CaseID:   caseID,
		Response: resp,
		Passed:   false,
		Reason:   entities.FailReasonCallFailed,
	}

	if !resp.Succeeded() {
		return result, nil
	}

	if !expectJSON {
		result.Passed = true
		result.Reason = entities.FailReasonNone
		return result, nil
	}

	if validation.Check(resp.Body) {
		result.Passed = true
		result.Reason = entities.FailReasonNone
	} else {
		result.Reason = entities.FailReasonInvalidJSON
	}
	return result, nil
}
