package application

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/infrastructure/providers/openai_compatible"
)

func TestEvaluateCase_withOpenAIClient_invalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{ \"a\":"}}]}`))
	}))
	defer srv.Close()

	client, err := openai_compatible.NewClient(srv.URL+"/v1", "test-key", srv.Client())
	if err != nil {
		t.Fatal(err)
	}

	timeout, _ := reliability.NewTimeoutPolicy(5 * time.Second)
	retry, _ := reliability.NewRetryPolicy(1, 10*time.Millisecond)
	req, _ := entities.NewRequest("gpt-4o-mini", []entities.Message{
		{Role: "user", Content: "json"},
	}, false)

	result, err := EvaluateCase(context.Background(), "json-case", entities.Expectations{JSON: true}, req, timeout, retry, client)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Fatal("expected fail")
	}
	if result.Reason != entities.FailReasonInvalidJSON {
		t.Fatalf("reason=%q", result.Reason)
	}
	if !result.Response.Succeeded() {
		t.Fatal("HTTP should succeed")
	}
}
