package application

import (
	"context"
	"testing"
	"time"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
)

func TestEvaluateCase_successNoJSON(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("hello", 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	result, err := EvaluateCase(context.Background(), "ping", entities.Expectations{}, testRequest(t), timeout, retry, fake)
	if err != nil || !result.Passed || result.Reason != entities.FailReasonNone {
		t.Fatalf("result: %+v err=%v", result, err)
	}
}

func TestEvaluateCase_contentMismatch(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("Tabletennis", 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	exp := entities.ExpectationsFromCase("pong", false, false, nil)
	result, err := EvaluateCase(context.Background(), "ping-llama", exp, testRequest(t), timeout, retry, fake)
	if err != nil {
		t.Fatal(err)
	}
	if result.Passed {
		t.Fatal("expected fail")
	}
	if result.Reason != entities.FailReasonContentMismatch {
		t.Fatalf("reason=%q", result.Reason)
	}
	if !result.Response.Succeeded() {
		t.Fatal("HTTP-level call should succeed")
	}
}

func TestEvaluateCase_containsPass(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("pong", 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	exp := entities.ExpectationsFromCase("pong", false, false, nil)
	result, err := EvaluateCase(context.Background(), "ping-gemini", exp, testRequest(t), timeout, retry, fake)
	if err != nil || !result.Passed {
		t.Fatalf("result: %+v err=%v", result, err)
	}
}

func TestEvaluateCase_containsIgnoreCasePass(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("Pong", 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	exp := entities.ExpectationsFromCase("pong", false, true, nil)
	result, err := EvaluateCase(context.Background(), "ping-gemini", exp, testRequest(t), timeout, retry, fake)
	if err != nil || !result.Passed {
		t.Fatalf("result: %+v err=%v", result, err)
	}
}

func TestEvaluateCase_invalidJSON(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse(`{ "name": "John"`, 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	result, err := EvaluateCase(context.Background(), "json-case", entities.Expectations{JSON: true}, testRequest(t), timeout, retry, fake)
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
		t.Fatal("HTTP-level call should succeed")
	}
}

func TestEvaluateCase_validJSON(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse(`{"name":"John"}`, 200, 10, 0, 1),
		},
	}
	timeout, _ := reliability.NewTimeoutPolicy(time.Second)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	result, err := EvaluateCase(context.Background(), "json-ok", entities.Expectations{JSON: true}, testRequest(t), timeout, retry, fake)
	if err != nil || !result.Passed {
		t.Fatalf("result: %+v", result)
	}
}

func TestEvaluateCase_callFailedTimeout(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("slow", 200, 0, 0, 1),
		},
		sleep: 200 * time.Millisecond,
	}
	timeout, _ := reliability.NewTimeoutPolicy(50 * time.Millisecond)
	retry, _ := reliability.NewRetryPolicy(2, 5*time.Millisecond)

	result, err := EvaluateCase(context.Background(), "slow", entities.Expectations{JSON: true}, testRequest(t), timeout, retry, fake)
	if err != nil || result.Passed || result.Reason != entities.FailReasonCallFailed {
		t.Fatalf("result: %+v", result)
	}
}
