package entities

import (
	"errors"
	"testing"
)

func TestResponse_Succeeded(t *testing.T) {
	r := NewSuccessResponse(`{"ok":true}`, 200, 100, 0, 1)
	if !r.Succeeded() || r.Failed() {
		t.Fatalf("expected success, got %+v", r)
	}
}

func TestResponse_Succeeded_requires2xx(t *testing.T) {
	r := NewSuccessResponse("body", 503, 100, 0, 1)
	if r.Succeeded() {
		t.Fatal("503 should not succeed")
	}
}

func TestResponse_Succeeded_withRetries(t *testing.T) {
	r := NewSuccessResponse("ok", 200, 1200, 0, 3)
	if !r.Succeeded() || r.Attempts != 3 {
		t.Fatalf("expected success after retries, got %+v", r)
	}
}

func TestResponse_IsTimeout_byKind(t *testing.T) {
	r := NewFailedResponse(FailureTimeout, 0, 1, 2000, nil)
	if !r.IsTimeout() || r.Succeeded() {
		t.Fatalf("expected timeout failure, got %+v", r)
	}
	if !errors.Is(r.Err, ErrCallTimeout) {
		t.Fatalf("expected ErrCallTimeout, got %v", r.Err)
	}
}

func TestResponse_IsTimeout_byErr(t *testing.T) {
	r := Response{
		HTTPStatus: 0,
		Kind:       FailureNone,
		Err:        ErrCallTimeout,
	}
	if !r.IsTimeout() {
		t.Fatal("expected IsTimeout via Err")
	}
}

func TestResponse_IsRetryExhausted(t *testing.T) {
	r := NewFailedResponse(FailureRetryExhausted, 503, 3, 4000, nil)
	if !r.IsRetryExhausted() || r.Succeeded() {
		t.Fatalf("expected retry exhausted, got %+v", r)
	}
	if r.Attempts != 3 {
		t.Fatalf("attempts: %d", r.Attempts)
	}
	if !errors.Is(r.Err, ErrCallRetryExhausted) {
		t.Fatalf("err: %v", r.Err)
	}
}

func TestResponse_HasTTFT(t *testing.T) {
	r := NewSuccessResponse("hello", 200, 500, 120, 1)
	if !r.HasTTFT() {
		t.Fatal("expected TTFT")
	}
	nonStream := NewSuccessResponse("hello", 200, 500, 0, 1)
	if nonStream.HasTTFT() {
		t.Fatal("non-stream should have TTFTMs 0")
	}
}

func TestResponse_BodyPreview(t *testing.T) {
	tests := []struct {
		name string
		body string
		max  int
		want string
	}{
		{"empty body", "", 10, ""},
		{"zero max", "hello", 0, ""},
		{"short", "hi", 10, "hi"},
		{"truncate ascii", "hello world", 5, "hello…"},
		{"truncate utf8", "héllo 🌍", 3, "hél…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewSuccessResponse(tt.body, 200, 1, 0, 1)
			got := r.BodyPreview(tt.max)
			if got != tt.want {
				t.Fatalf("BodyPreview(%d) = %q, want %q", tt.max, got, tt.want)
			}
		})
	}
}

func TestNewSuccessResponse_defaults(t *testing.T) {
	r := NewSuccessResponse("x", 0, 50, 0, 0)
	if r.HTTPStatus != 200 || r.Attempts != 1 {
		t.Fatalf("defaults: status=%d attempts=%d", r.HTTPStatus, r.Attempts)
	}
}

func TestNewFailedResponse_providerErr(t *testing.T) {
	r := NewFailedResponse(FailureProvider, 400, 1, 10, nil)
	if r.Succeeded() || r.Err == nil {
		t.Fatalf("expected provider failure, got %+v", r)
	}
}

func TestNewFailedResponse_cancelled(t *testing.T) {
	r := NewFailedResponse(FailureCancelled, 0, 1, 0, nil)
	if !errors.Is(r.Err, ErrCallCancelled) {
		t.Fatalf("err: %v", r.Err)
	}
}
