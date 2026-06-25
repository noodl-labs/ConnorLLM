package benchmark

type Spec struct {
	Suite    string       `yaml:"suite"`
	Defaults CaseDefaults `yaml:"defaults"`
	Cases    []CaseSpec   `yaml:"cases"`
}

type CaseDefaults struct {
	TimeoutMS                int64 `yaml:"timeout_ms"`
	Retries                  int   `yaml:"retries"`
	ExpectContainsIgnoreCase bool  `yaml:"expect_contains_ignore_case"`
}

type CaseSpec struct {
	ID                       string             `yaml:"id"`
	Model                    string             `yaml:"model"`
	Prompt                   string             `yaml:"prompt"`
	ExpectContains           string             `yaml:"expect_contains"`
	ExpectContainsIgnoreCase bool               `yaml:"expect_contains_ignore_case"`
	ExpectJSON               bool               `yaml:"expect_json"`
	ExpectJSONSchema         JSONSchemaDocument `yaml:"expect_json_schema"`
	TimeoutMS                int64              `yaml:"timeout_ms"`
	Retries                  int                `yaml:"retries"`
}
