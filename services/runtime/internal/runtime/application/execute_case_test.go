package application

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
)

type fakeProvider struct {
	mu       sync.Mutex
	responses []entities.Response
	calls    int
	sleep    time.Duration
}

func (f *fakeProvider) Execute(ctx context.Context, req entities.Request) (entities.Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.sleep > 0 {
		select {
		case <-time.After(f.sleep):
		case <-ctx.Done():
			return entities.Response{}, ctx.Err()
		}
	}
	idx := f.calls
	f.calls++
	if idx >= len(f.responses) {
		return entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil), nil
	}
	return f.responses[idx], nil
}

func testRequest(t *testing.T) entities.Request {
	t.Helper()
	req, err := entities.NewRequest("test-model", []entities.Message{{Role: "user", Content: "ping"}}, false)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func TestExecuteCase_success(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("ok", 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 10*time.Millisecond)

	resp, err := ExecuteCase(context.Background(), testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Succeeded() || resp.Attempts != 1 {
		t.Fatalf("response: %+v", resp)
	}
}

func TestExecuteCase_retryThenSuccess(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil),
			entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil),
			entities.NewSuccessResponse("ok", 200, 5, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	resp, err := ExecuteCase(context.Background(), testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Succeeded() || resp.Attempts != 3 {
		t.Fatalf("response: %+v calls=%d", resp, fake.calls)
	}
}

func TestExecuteCase_timeoutNoRetry(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("slow", 200, 0, 0, 1),
		},
		sleep: 200 * time.Millisecond,
	}
	timeout, _ := reliability.NewTimeoutPolicy(50 * time.Millisecond)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	resp, err := ExecuteCase(context.Background(), testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsTimeout() {
		t.Fatalf("expected timeout, got %+v", resp)
	}
	if resp.Attempts != 1 {
		t.Fatalf("timeout should not retry, attempts=%d", resp.Attempts)
	}
}

func TestExecuteCase_noRetryOn400(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewFailedResponse(entities.FailureProvider, 400, 1, 0, nil),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	resp, err := ExecuteCase(context.Background(), testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Succeeded() || resp.Attempts != 1 {
		t.Fatalf("400 should fail once: %+v", resp)
	}
	if resp.HTTPStatus != 400 {
		t.Fatalf("status: %d", resp.HTTPStatus)
	}
}

func TestExecuteCase_retryExhausted(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil),
			entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil),
			entities.NewFailedResponse(entities.FailureProvider, 503, 1, 0, nil),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	resp, err := ExecuteCase(context.Background(), testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsRetryExhausted() || resp.Attempts != 3 {
		t.Fatalf("response: %+v", resp)
	}
}
