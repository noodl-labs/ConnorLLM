package reliability

import (
	"context"
	"fmt"
	"time"
)

// DefaultDeadline is used when a benchmark suite does not set timeout_ms.
// It bounds a single provider attempt so runs never hang indefinitely (CI-safe).
// Suites should override with explicit timeout_ms for strict smoke (e.g. 5000).
const DefaultDeadline = 30 * time.Second

// TimeoutPolicy defines the maximum duration allowed for one provider attempt.
// Each retry gets a fresh deadline via Apply (per-attempt timeout, MVP design).
// It does not perform HTTP calls or build Response — ExecuteCase owns that.
type TimeoutPolicy struct {
	// Deadline is the per-attempt time budget (from YAML timeout_ms or DefaultDeadline).
	Deadline time.Duration
}

// NewTimeoutPolicy creates a policy with a positive duration.
// Prefer NewTimeoutPolicyFromMS when loading suite YAML (timeout_ms).
func NewTimeoutPolicy(deadline time.Duration) (TimeoutPolicy, error) {
	if deadline <= 0 {
		return TimeoutPolicy{}, fmt.Errorf("reliability: deadline must be positive")
	}
	return TimeoutPolicy{Deadline: deadline}, nil
}

// NewTimeoutPolicyFromMS converts suite timeout_ms into a TimeoutPolicy.
func NewTimeoutPolicyFromMS(ms int64) (TimeoutPolicy, error) {
	if ms <= 0 {
		return TimeoutPolicy{}, fmt.Errorf("reliability: timeout_ms must be positive")
	}
	return NewTimeoutPolicy(time.Duration(ms) * time.Millisecond)
}

// Apply returns a child context cancelled after p.Deadline.
// The caller must call cancel (typically defer cancel()) to release the timer.
// Pass the returned context to ProviderExecutor.Execute so the call stops on timeout.
// On deadline exceeded, ExecuteCase maps context.DeadlineExceeded to Response timeout failure.
func (p TimeoutPolicy) Apply(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, p.Deadline)
}
