package benchmark

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// JSONSchemaDocument holds an inline YAML/JSON Schema object from expect_json_schema.
type JSONSchemaDocument []byte

// UnmarshalYAML decodes a YAML mapping into canonical JSON bytes for the domain layer.
func (s *JSONSchemaDocument) UnmarshalYAML(value *yaml.Node) error {
	var raw any
	if err := value.Decode(&raw); err != nil {
		return err
	}
	if raw == nil {
		*s = nil
		return nil
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	*s = b
	return nil
}

// Raw returns schema bytes for domain expectations (nil when unset).
func (s JSONSchemaDocument) Raw() json.RawMessage {
	if len(s) == 0 {
		return nil
	}
	return json.RawMessage(s)
}

// IsSet reports whether expect_json_schema was provided in YAML.
func (s JSONSchemaDocument) IsSet() bool {
	return len(s) > 0
}
