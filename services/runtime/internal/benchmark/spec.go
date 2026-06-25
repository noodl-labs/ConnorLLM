package benchmark

type Spec struct {
	Suite    string       `yaml:"suite"`
	Defaults CaseDefaults `yaml:"defaults"`
	Cases    []CaseSpec   `yaml:"cases"`
}

type CaseDefaults struct {
	TimeoutMS int64 `yaml:"timeout_ms"`
	Retries   int   `yaml:"retries"`
}

type CaseSpec struct {
	ID         string `yaml:"id"`
	Model      string `yaml:"model"`
	Prompt     string `yaml:"prompt"`
	ExpectJSON bool   `yaml:"expect_json"`
	TimeoutMS  int64  `yaml:"timeout_ms"`
	Retries    int    `yaml:"retries"`
}
