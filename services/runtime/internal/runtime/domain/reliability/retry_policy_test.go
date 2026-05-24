package reliability

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewRetryPolicy_validation(t *testing.T) {
	_, err := NewRetryPolicy(-1, time.Second)
	if err == nil {
		t.Fatal("expected error for negative retries")
	}
	_, err = NewRetryPolicy(10, time.Second)
	if err == nil {
		t.Fatal("expected error above cap")
	}
}

func TestRetryPolicy_MaxAttempts(t *testing.T) {
	p, err := NewRetryPolicy(2, DefaultBackoffBase)
	if err != nil {
		t.Fatal(err)
	}
	if p.MaxAttempts() != 3 {
		t.Fatalf("max attempts: got %d want 3", p.MaxAttempts())
	}
}

func TestRetryPolicy_ShouldRetry_transientHTTP(t *testing.T) {
	p, _ := NewRetryPolicy(2, DefaultBackoffBase)
	err := TransientHTTPError(503)
	if !p.ShouldRetry(err, 1) {
		t.Fatal("503 should be retryable on attempt 1")
	}
	if p.ShouldRetry(err, 3) {
		t.Fatal("should not retry after max attempts")
	}
}

func TestRetryPolicy_ShouldRetry_noRetryOnTimeout(t *testing.T) {
	p, _ := NewRetryPolicy(2, DefaultBackoffBase)
	if p.ShouldRetry(context.DeadlineExceeded, 1) {
		t.Fatal("timeout should not be retryable")
	}
	if p.ShouldRetry(context.Canceled, 1) {
		t.Fatal("cancel should not be retryable")
	}
}

func TestRetryPolicy_ShouldRetry_nilErr(t *testing.T) {
	p, _ := NewRetryPolicy(2, DefaultBackoffBase)
	if p.ShouldRetry(nil, 1) {
		t.Fatal("nil error should not retry")
	}
}

func TestRetryPolicy_ShouldRetry_nonTransientHTTP(t *testing.T) {
	p, _ := NewRetryPolicy(2, DefaultBackoffBase)
	if p.ShouldRetry(errors.New("provider: bad request"), 1) {
		t.Fatal("generic provider error should not be retryable")
	}
}

func TestIsTransientHTTP(t *testing.T) {
	if !IsTransientHTTP(503) || !IsTransientHTTP(429) {
		t.Fatal("503/429 should be transient")
	}
	if IsTransientHTTP(400) {
		t.Fatal("400 should not be transient")
	}
}

func TestRetryPolicy_Backoff_linear(t *testing.T) {
	p, _ := NewRetryPolicy(2, 200*time.Millisecond)
	if p.Backoff(1) != 200*time.Millisecond {
		t.Fatalf("backoff 1: %v", p.Backoff(1))
	}
	if p.Backoff(2) != 400*time.Millisecond {
		t.Fatalf("backoff 2: %v", p.Backoff(2))
	}
}

func TestTransientError_message(t *testing.T) {
	var err error = TransientHTTPError(503)
	var te *TransientError
	if !errors.As(err, &te) || te.Status != 503 {
		t.Fatalf("unexpected error: %v", err)
	}
}
