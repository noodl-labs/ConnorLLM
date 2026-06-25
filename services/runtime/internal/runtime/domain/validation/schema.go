package validation

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

const schemaResourceURL = "connor://expect_json_schema"

// CheckSchema reports whether body matches schemaJSON (JSON Schema).
// body must already be valid JSON syntax.
func CheckSchema(body string, schemaJSON json.RawMessage) bool {
	if len(schemaJSON) == 0 {
		return true
	}

	sch, err := CompileSchema(schemaJSON)
	if err != nil {
		return false
	}

	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(body))
	if err != nil {
		return false
	}
	return sch.Validate(doc) == nil
}

// CompileSchema parses and compiles a JSON Schema document for validation.
func CompileSchema(schemaJSON json.RawMessage) (*jsonschema.Schema, error) {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		return nil, err
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(schemaResourceURL, doc); err != nil {
		return nil, err
	}
	return compiler.Compile(schemaResourceURL)
}
