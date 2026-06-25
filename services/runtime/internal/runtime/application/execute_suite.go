package application

import (
	"context"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/benchmark"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
)

const defaultSuiteRetries = 2

func ExecuteSuite(
	ctx context.Context,
	spec benchmark.Spec,
	executor domain.ProviderExecutor,
) (entities.SuiteResult, error) {
	out := entities.SuiteResult{SuiteID: spec.Suite}

	for _, c := range spec.Cases {
		timeout, err := resolveSuiteTimeout(c.TimeoutMS, spec.Defaults.TimeoutMS)
		if err != nil {
			return entities.SuiteResult{}, err
		}

		retries := resolveSuiteRetries(c.Retries, spec.Defaults.Retries)
		retry, err := reliability.NewRetryPolicy(retries, reliability.DefaultBackoffBase)
		if err != nil {
			return entities.SuiteResult{}, err
		}

		req, err := entities.NewRequest(c.Model, []entities.Message{
			{Role: "user", Content: c.Prompt},
		}, false)
		if err != nil {
			return entities.SuiteResult{}, err
		}

		exp := entities.ExpectationsFromCase(c.ExpectContains, c.ExpectJSON)
		result, err := EvaluateCase(ctx, c.ID, exp, req, timeout, retry, executor)
		if err != nil {
			return entities.SuiteResult{}, err
		}
		out.Results = append(out.Results, result)
	}
	return out, nil
}

// resolveSuiteTimeout merges per-case and suite defaults; 0 everywhere uses DefaultDeadline.
func resolveSuiteTimeout(caseMS, defaultsMS int64) (reliability.TimeoutPolicy, error) {
	ms := caseMS
	if ms == 0 {
		ms = defaultsMS
	}
	if ms == 0 {
		return reliability.NewTimeoutPolicy(reliability.DefaultDeadline)
	}
	return reliability.NewTimeoutPolicyFromMS(ms)
}

// resolveSuiteRetries merges per-case and suite defaults; 0 everywhere uses defaultSuiteRetries.
// YAML int 0 means unset (inherit); explicit zero retries is not supported in beta.1.
func resolveSuiteRetries(caseRetries, defaultsRetries int) int {
	if caseRetries != 0 {
		return caseRetries
	}
	if defaultsRetries != 0 {
		return defaultsRetries
	}
	return defaultSuiteRetries
}
