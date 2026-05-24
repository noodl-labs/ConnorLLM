package reliability

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewTimeoutPolicy_invalid(t *testing.T) {
	_, err := NewTimeoutPolicy(0)
	if err == nil {
		t.Fatal("expected error for zero deadline")
	}
	_, err = NewTimeoutPolicy(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative deadline")
	}
}

func TestNewTimeoutPolicyFromMS(t *testing.T) {
	p, err := NewTimeoutPolicyFromMS(5000)
	if err != nil {
		t.Fatal(err)
	}
	if p.Deadline != 5*time.Second {
		t.Fatalf("deadline: got %v want 5s", p.Deadline)
	}
	_, err = NewTimeoutPolicyFromMS(0)
	if err == nil {
		t.Fatal("expected error for zero ms")
	}
}

func TestTimeoutPolicy_Apply_deadlineExceeded(t *testing.T) {
	p, err := NewTimeoutPolicy(50 * time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := p.Apply(context.Background())
	defer cancel()

	select {
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", ctx.Err())
	}
}

func TestTimeoutPolicy_Apply_parentCancel(t *testing.T) {
	p, err := NewTimeoutPolicy(time.Second)
	if err != nil {
		t.Fatal(err)
	}
	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := p.Apply(parent)
	defer cancel()

	parentCancel()
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Fatalf("expected Canceled, got %v", ctx.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("child context should be cancelled when parent is cancelled")
	}
}

func TestDefaultDeadline_positive(t *testing.T) {
	if DefaultDeadline <= 0 {
		t.Fatal("DefaultDeadline must be positive")
	}
}
