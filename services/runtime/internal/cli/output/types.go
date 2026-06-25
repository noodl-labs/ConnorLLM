package output

import "github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"

type CaseView struct {
	ID                       string
	Model                    string
	ExpectContains           string
	ExpectContainsIgnoreCase bool
	ExpectJSON               bool
	ExpectJSONSchema         bool
	Result                   entities.CaseResult
}

type RunView struct {
	Version string
	Target  string
	SuiteID string
	Cases   []CaseView
}
