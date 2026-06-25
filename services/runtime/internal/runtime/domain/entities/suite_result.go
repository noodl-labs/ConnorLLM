package entities

type SuiteResult struct {
	SuiteID string
	Results []CaseResult
}

func (s SuiteResult) AllPassed() bool {
	for _, r := range s.Results {
		if !r.Passed {
			return false
		}
	}
	return len(s.Results) > 0
}

func (s SuiteResult) PassedCount() int {
	n := 0
	for _, r := range s.Results {
		if r.Passed {
			n++
		}
	}
	return n
}
