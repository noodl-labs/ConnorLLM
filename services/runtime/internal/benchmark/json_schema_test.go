package benchmark_test

import (
	"encoding/json"
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/benchmark"
)

func TestJSONSchemaDocument_UnmarshalYAML(t *testing.T) {
	data := []byte(`
suite: test
cases:
  - id: flight
    model: m
    prompt: p
    expect_json_schema:
      type: object
      required: [intent]
      properties:
        intent: { type: string }
`)
	spec, err := benchmark.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !spec.Cases[0].ExpectJSONSchema.IsSet() {
		t.Fatal("expected schema to be set")
	}
	var doc map[string]any
	if err := json.Unmarshal(spec.Cases[0].ExpectJSONSchema.Raw(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc["type"] != "object" {
		t.Fatalf("type=%v", doc["type"])
	}
}

func TestParse_invalidJSONSchema(t *testing.T) {
	data := []byte(`
suite: test
cases:
  - id: bad
    model: m
    prompt: p
    expect_json_schema:
      type: not-a-valid-type
`)
	_, err := benchmark.Parse(data)
	if err == nil {
		t.Fatal("expected parse error for invalid schema")
	}
}
