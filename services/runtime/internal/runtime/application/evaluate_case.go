package application

import (
	"context"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/validation"
)

// EvaluateCase runs one benchmark case: ExecuteCase then validation gates.
//
// Gates (via Expectations): expect_contains, expect_json, expect_json_schema.
// HTTP failure → Passed=false, Reason=call_failed.
// HTTP 2xx → validation.Evaluate applies active gates.
func EvaluateCase(
	ctx context.Context,
	caseID string,
	exp entities.Expectations,
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

	passed, reason := validation.Evaluate(resp.Body, exp)
	result.Passed = passed
	result.Reason = reason
	return result, nil
}
