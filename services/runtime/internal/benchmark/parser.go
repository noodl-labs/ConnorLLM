package benchmark

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/validation"
	"gopkg.in/yaml.v3"
)

func ParseFile(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, err
	}
	return Parse(data)
}

func Parse(data []byte) (Spec, error) {
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return Spec{}, fmt.Errorf("benchmark: invalid yaml: %w", err)
	}
	if err := validate(spec); err != nil {
		return Spec{}, err
	}
	return spec, nil
}

func validate(spec Spec) error {
	if spec.Suite == "" {
		return fmt.Errorf("benchmark: suite is required")
	}
	if len(spec.Cases) == 0 {
		return fmt.Errorf("benchmark: at least one case is required")
	}
	seen := make(map[string]struct{}, len(spec.Cases))
	for i, c := range spec.Cases {
		if c.ID == "" {
			return fmt.Errorf("benchmark: cases[%d]: id is required", i)
		}
		if _, ok := seen[c.ID]; ok {
			return fmt.Errorf("benchmark: duplicate case id %q", c.ID)
		}
		seen[c.ID] = struct{}{}
		if c.Model == "" {
			return fmt.Errorf("benchmark: case %q: model is required", c.ID)
		}
		if c.Prompt == "" {
			return fmt.Errorf("benchmark: case %q: prompt is required", c.ID)
		}
		if c.TimeoutMS < 0 {
			return fmt.Errorf("benchmark: case %q: timeout_ms must be >= 0", c.ID)
		}
		if c.Retries < 0 {
			return fmt.Errorf("benchmark: case %q: retries must be >= 0", c.ID)
		}
		if err := validateJSONSchema(c.ID, c.ExpectJSONSchema); err != nil {
			return err
		}
	}
	if spec.Defaults.TimeoutMS < 0 {
		return fmt.Errorf("benchmark: defaults.timeout_ms must be >= 0")
	}
	if spec.Defaults.Retries < 0 {
		return fmt.Errorf("benchmark: defaults.retries must be >= 0")
	}
	return nil
}

func validateJSONSchema(caseID string, schema JSONSchemaDocument) error {
	if !schema.IsSet() {
		return nil
	}
	raw := schema.Raw()
	if !json.Valid(raw) {
		return fmt.Errorf("benchmark: case %q: expect_json_schema must be valid JSON", caseID)
	}
	if _, err := validation.CompileSchema(raw); err != nil {
		return fmt.Errorf("benchmark: case %q: expect_json_schema: %w", caseID, err)
	}
	return nil
}
