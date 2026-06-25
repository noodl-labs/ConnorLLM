package validation_test

import (
	"encoding/json"
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/validation"
)

var flightIntentSchema = json.RawMessage(`{
  "type": "object",
  "required": ["intent", "confidence"],
  "properties": {
    "intent": { "type": "string" },
    "confidence": { "type": "number" }
  }
}`)

func TestCheckSchema_valid(t *testing.T) {
	body := `{"intent":"book_flight","confidence":0.9}`
	if !validation.CheckSchema(body, flightIntentSchema) {
		t.Fatal("expected schema match")
	}
}

func TestCheckSchema_missingField(t *testing.T) {
	body := `{"intent":"book_flight"}`
	if validation.CheckSchema(body, flightIntentSchema) {
		t.Fatal("expected schema mismatch")
	}
}

func TestCheckSchema_wrongType(t *testing.T) {
	body := `{"intent":42,"confidence":0.9}`
	if validation.CheckSchema(body, flightIntentSchema) {
		t.Fatal("expected schema mismatch")
	}
}

func TestCheckSchema_emptySchema(t *testing.T) {
	if !validation.CheckSchema(`{"any":true}`, nil) {
		t.Fatal("empty schema should pass")
	}
}

func TestEvaluate_schemaPass(t *testing.T) {
	exp := entities.Expectations{JSONSchema: flightIntentSchema}
	passed, reason := validation.Evaluate(`{"intent":"book_flight","confidence":0.9}`, exp)
	if !passed || reason != entities.FailReasonNone {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestEvaluate_schemaMismatch(t *testing.T) {
	exp := entities.Expectations{JSONSchema: flightIntentSchema}
	passed, reason := validation.Evaluate(`{"status":"ok"}`, exp)
	if passed || reason != entities.FailReasonSchemaMismatch {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}

func TestEvaluate_invalidJSONBeforeSchema(t *testing.T) {
	exp := entities.Expectations{JSONSchema: flightIntentSchema}
	passed, reason := validation.Evaluate("not json", exp)
	if passed || reason != entities.FailReasonInvalidJSON {
		t.Fatalf("passed=%v reason=%q", passed, reason)
	}
}
