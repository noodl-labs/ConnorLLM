package application

import (
	"context"
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/benchmark"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
)

func TestExecuteSuite_allPass(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("pong", 200, 10, 0, 1),
			entities.NewSuccessResponse(`{"ok":true}`, 200, 12, 0, 1),
		},
	}

	spec := benchmark.Spec{
		Suite: "test",
		Cases: []benchmark.CaseSpec{
			{ID: "ping", Model: "m", Prompt: "pong"},
			{ID: "json", Model: "m", Prompt: "json", ExpectJSON: true},
		},
	}

	result, err := ExecuteSuite(context.Background(), spec, fake)
	if err != nil {
		t.Fatal(err)
	}
	if !result.AllPassed() || result.PassedCount() != 2 {
		t.Fatalf("result: %+v", result)
	}
}

func TestExecuteSuite_oneFailContinues(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("pong", 200, 10, 0, 1),
			entities.NewSuccessResponse(`not json`, 200, 12, 0, 1),
		},
	}

	spec := benchmark.Spec{
		Suite: "test",
		Cases: []benchmark.CaseSpec{
			{ID: "ping", Model: "m", Prompt: "pong"},
			{ID: "json", Model: "m", Prompt: "json", ExpectJSON: true},
		},
	}

	result, err := ExecuteSuite(context.Background(), spec, fake)
	if err != nil {
		t.Fatal(err)
	}
	if result.AllPassed() || result.PassedCount() != 1 {
		t.Fatalf("result: %+v", result)
	}
	if result.Results[1].Reason != entities.FailReasonInvalidJSON {
		t.Fatalf("reason=%q", result.Results[1].Reason)
	}
}

func TestExecuteSuite_contentMismatch(t *testing.T) {
	fake := &fakeProvider{
		responses: []entities.Response{
			entities.NewSuccessResponse("Tabletennis", 200, 10, 0, 1),
		},
	}

	spec := benchmark.Spec{
		Suite: "test",
		Cases: []benchmark.CaseSpec{
			{ID: "ping-llama", Model: "m", Prompt: "pong", ExpectContains: "pong"},
		},
	}

	result, err := ExecuteSuite(context.Background(), spec, fake)
	if err != nil {
		t.Fatal(err)
	}
	if result.AllPassed() {
		t.Fatalf("result: %+v", result)
	}
	if result.Results[0].Reason != entities.FailReasonContentMismatch {
		t.Fatalf("reason=%q", result.Results[0].Reason)
	}
}

func TestResolveSuiteTimeout_defaultDeadline(t *testing.T) {
	timeout, err := resolveSuiteTimeout(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if timeout.Deadline != reliability.DefaultDeadline {
		t.Fatalf("deadline=%v", timeout.Deadline)
	}
}

func TestResolveSuiteRetries_defaults(t *testing.T) {
	if got := resolveSuiteRetries(0, 3); got != 3 {
		t.Fatalf("got %d", got)
	}
	if got := resolveSuiteRetries(0, 0); got != defaultSuiteRetries {
		t.Fatalf("got %d", got)
	}
	if got := resolveSuiteRetries(1, 3); got != 1 {
		t.Fatalf("got %d", got)
	}
}
