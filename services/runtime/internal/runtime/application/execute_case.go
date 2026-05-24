package application

import (
	"context"
	"errors"
	"time"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
)

// ExecuteCase runs one benchmark case: per-attempt timeout, retries, and provider execution.
func ExecuteCase(
	ctx context.Context,
	req entities.Request,
	timeout reliability.TimeoutPolicy,
	retry reliability.RetryPolicy,
	executor domain.ProviderExecutor,
) (entities.Response, error) {
	var lastResp entities.Response
	var lastErr error

	for attempt := 1; attempt <= retry.MaxAttempts(); attempt++ {
		if ctx.Err() != nil {
			return failureFromContext(ctx, attempt), nil
		}

		ctxAttempt, cancel := timeout.Apply(ctx)
		resp, err := executor.Execute(ctxAttempt, req)
		cancel()

		lastResp = resp
		lastErr = err

		if err == nil && resp.Succeeded() {
			resp.Attempts = attempt
			return resp, nil
		}

		retryErr := classifyRetryError(err, resp)
		if !retry.ShouldRetry(retryErr, attempt) {
			if attempt >= retry.MaxAttempts() && retry.IsRetryable(retryErr) {
				return mapRetryExhausted(resp, attempt, retryErr), nil
			}
			return mapFailure(resp, attempt, retryErr), nil
		}

		backoff := retry.Backoff(attempt)
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return failureFromContext(ctx, attempt), nil
		}
	}

	return mapRetryExhausted(lastResp, retry.MaxAttempts(), lastErr), nil
}

func classifyRetryError(err error, resp entities.Response) error {
	if err != nil {
		return err
	}
	if reliability.IsTransientHTTP(resp.HTTPStatus) {
		return reliability.TransientHTTPError(resp.HTTPStatus)
	}
	if resp.Failed() {
		return errors.New("entities: provider call failed")
	}
	return nil
}

func mapFailure(resp entities.Response, attempt int, err error) entities.Response {
	if errors.Is(err, context.DeadlineExceeded) {
		return entities.NewFailedResponse(entities.FailureTimeout, resp.HTTPStatus, attempt, resp.LatencyMs, entities.ErrCallTimeout)
	}
	if errors.Is(err, context.Canceled) {
		return entities.NewFailedResponse(entities.FailureCancelled, resp.HTTPStatus, attempt, resp.LatencyMs, entities.ErrCallCancelled)
	}
	var transient *reliability.TransientError
	if errors.As(err, &transient) {
		return entities.NewFailedResponse(entities.FailureProvider, transient.Status, attempt, resp.LatencyMs, err)
	}
	return entities.NewFailedResponse(entities.FailureProvider, resp.HTTPStatus, attempt, resp.LatencyMs, err)
}

func mapRetryExhausted(last entities.Response, maxAttempts int, lastErr error) entities.Response {
	status := last.HTTPStatus
	if status == 0 {
		status = 503
	}
	_ = lastErr
	return entities.NewFailedResponse(entities.FailureRetryExhausted, status, maxAttempts, last.LatencyMs, entities.ErrCallRetryExhausted)
}

func failureFromContext(ctx context.Context, attempt int) entities.Response {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return entities.NewFailedResponse(entities.FailureTimeout, 0, attempt, 0, entities.ErrCallTimeout)
	}
	return entities.NewFailedResponse(entities.FailureCancelled, 0, attempt, 0, entities.ErrCallCancelled)
}
