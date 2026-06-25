package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/cli/output"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

func TestPrintRun_suitePass(t *testing.T) {
	view := output.RunView{
		Version: "v0.1.0-beta.1",
		Target:  "https://openrouter.ai/api/v1",
		SuiteID: "serving-smoke",
		Cases: []output.CaseView{
			{
				ID:             "ping-gemini",
				Model:          "google/gemini-2.5-flash-lite",
				ExpectContains: "pong",
				Result: entities.CaseResult{
					CaseID:   "ping-gemini",
					Passed:   true,
					Response: entities.NewSuccessResponse("pong", 200, 479, 0, 1),
				},
			},
			{
				ID:         "json-ok",
				Model:      "openai/gpt-4o-mini",
				ExpectJSON: true,
				Result: entities.CaseResult{
					CaseID:   "json-ok",
					Passed:   true,
					Response: entities.NewSuccessResponse(`{"status":"ok"}`, 200, 630, 0, 1),
				},
			},
		},
	}

	var buf bytes.Buffer
	output.PrintRun(&buf, view, false)
	got := buf.String()

	for _, want := range []string{
		"Connor  v0.1.0-beta.1",
		"Target  https://openrouter.ai/api/v1",
		"Suite   serving-smoke (2 cases)",
		"✓  ping-gemini",
		"479ms  HTTP 200",
		"contains ✓",
		"json ✓",
		"2/2 passed",
		"GATE PASSED — safe to merge",
		"exit 0",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "body:") {
		t.Fatalf("expected no body in non-verbose pass output:\n%s", got)
	}
}

func TestPrintRun_failureShowsDetail(t *testing.T) {
	view := output.RunView{
		Version: "v0.1.0-beta.1",
		Target:  "https://openrouter.ai/api/v1",
		SuiteID: "serving-smoke",
		Cases: []output.CaseView{
			{
				ID:         "json-ok",
				Model:      "openai/gpt-4o-mini",
				ExpectJSON: true,
				Result: entities.CaseResult{
					CaseID:   "json-ok",
					Passed:   false,
					Reason:   entities.FailReasonInvalidJSON,
					Response: entities.NewSuccessResponse("not json", 200, 630, 0, 1),
				},
			},
		},
	}

	var buf bytes.Buffer
	output.PrintRun(&buf, view, false)
	got := buf.String()

	for _, want := range []string{
		"✗  json-ok",
		"gate:     expect_json",
		"reason:   invalid_json",
		"body:     not json",
		"GATE FAILED — do not merge",
		"exit 1",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}

func TestPrintRun_contentMismatchShowsDetail(t *testing.T) {
	view := output.RunView{
		Version: "v0.1.0-beta.1",
		Target:  "https://openrouter.ai/api/v1",
		SuiteID: "serving-smoke",
		Cases: []output.CaseView{
			{
				ID:             "ping-llama",
				Model:          "meta-llama/llama-3.1-8b-instruct",
				ExpectContains: "pong",
				Result: entities.CaseResult{
					CaseID:   "ping-llama",
					Passed:   false,
					Reason:   entities.FailReasonContentMismatch,
					Response: entities.NewSuccessResponse("Tabletennis", 200, 328, 0, 1),
				},
			},
		},
	}

	var buf bytes.Buffer
	output.PrintRun(&buf, view, false)
	got := buf.String()

	for _, want := range []string{
		"✗  ping-llama",
		"contains ✗",
		"gate:     expect_contains",
		"expected: pong",
		"reason:   content_mismatch",
		"body:     Tabletennis",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("output missing %q:\n%s", want, got)
		}
	}
}
