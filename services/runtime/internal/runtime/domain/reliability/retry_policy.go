package reliability

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	DefaultBackoffBase = 200 * time.Millisecond
	maxRetriesCap      = 5
)

// TransientError marks an HTTP or provider error that may be retried.
type TransientError struct {
	Status int
}

func (e *TransientError) Error() string {
	return fmt.Sprintf("reliability: transient provider error (status %d)", e.Status)
}

// RetryPolicy defines how many times to retry and which errors are retryable.
// YAML retries: N means N retries after the first attempt (MaxAttempts = N + 1).
type RetryPolicy struct {
	MaxRetries  int
	BackoffBase time.Duration
}

// NewRetryPolicy creates a retry policy. maxRetries is the number of retries (not total attempts).
func NewRetryPolicy(maxRetries int, backoffBase time.Duration) (RetryPolicy, error) {
	if maxRetries < 0 {
		return RetryPolicy{}, fmt.Errorf("reliability: max retries must be non-negative")
	}
	if maxRetries > maxRetriesCap {
		return RetryPolicy{}, fmt.Errorf("reliability: max retries must be <= %d", maxRetriesCap)
	}
	if backoffBase <= 0 {
		backoffBase = DefaultBackoffBase
	}
	return RetryPolicy{MaxRetries: maxRetries, BackoffBase: backoffBase}, nil
}

// MaxAttempts returns total tries including the first attempt.
func (p RetryPolicy) MaxAttempts() int {
	return p.MaxRetries + 1
}

// ShouldRetry reports whether another attempt is allowed after attempt (1-based).
func (p RetryPolicy) ShouldRetry(err error, attempt int) bool {
	if attempt >= p.MaxAttempts() {
		return false
	}
	return p.IsRetryable(err)
}

// IsRetryable classifies whether an error represents a transient failure worth retrying.
func (p RetryPolicy) IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	var transient *TransientError
	if errors.As(err, &transient) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	return false
}

// Backoff returns wait duration before the next attempt (linear: base * attempt).
func (p RetryPolicy) Backoff(attempt int) time.Duration {
	if attempt < 1 {
		return p.BackoffBase
	}
	return p.BackoffBase * time.Duration(attempt)
}

// TransientHTTPError wraps an HTTP status as a retryable error.
func TransientHTTPError(status int) error {
	return &TransientError{Status: status}
}

// IsTransientHTTP reports whether an HTTP status should trigger a retry.
func IsTransientHTTP(status int) bool {
	switch status {
	case 429, 502, 503, 504:
		return true
	default:
		return false
	}
}
