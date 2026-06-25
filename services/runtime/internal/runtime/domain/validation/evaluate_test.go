package validation_test

import (
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/validation"
)

func TestEvaluate_containsMismatch(t *testing.T) {
	exp := entities.Expectations{Contains: "pong"}
	passed, reason := validation.Evaluate("Tabletennis", exp)
	if passed || reason != entities.FailReasonContentMismatch {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestEvaluate_containsPass(t *testing.T) {
	exp := entities.Expectations{Contains: "pong"}
	passed, reason := validation.Evaluate("pong", exp)
	if !passed || reason != entities.FailReasonNone {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestEvaluate_containsIgnoreCasePass(t *testing.T) {
	exp := entities.Expectations{Contains: "pong", ContainsIgnoreCase: true}
	passed, reason := validation.Evaluate("Pong", exp)
	if !passed || reason != entities.FailReasonNone {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestEvaluate_containsDisabled(t *testing.T) {
	passed, reason := validation.Evaluate("anything", entities.Expectations{})
	if !passed || reason != entities.FailReasonNone {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestContains_caseSensitive(t *testing.T) {
	if validation.Contains("Ping", "pong", false) {
		t.Fatal("expected case-sensitive mismatch")
	}
}

func TestContains_ignoreCase(t *testing.T) {
	if !validation.Contains("Pong", "pong", true) {
		t.Fatal("expected case-insensitive match")
	}
}
