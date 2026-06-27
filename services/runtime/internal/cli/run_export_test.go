package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

func TestWriteRunArtifact(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "run.json")

	artifact := entities.RunArtifact{
		Version: entities.RunArtifactVersion,
		SuiteID: "test-suite",
		Target:  "https://example/v1",
		Cases: []entities.RunCase{
			{ID: "c1", Model: "gpt-4o-mini", Passed: true, LatencyMs: 42, Attempts: 1},
		},
		Summary: entities.RunSummary{
			Total: 1, Passed: 1, Failed: 0, PassRate: 100, P50Ms: 42, P95Ms: 42,
		},
	}

	if err := writeRunArtifact(path, artifact); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var decoded entities.RunArtifact
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.SuiteID != "test-suite" || decoded.Cases[0].Model != "gpt-4o-mini" {
		t.Fatalf("decoded: %+v", decoded)
	}
}
