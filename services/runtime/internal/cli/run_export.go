package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

// writeRunArtifact writes RunArtifact to path as indented JSON (RFC 0001 §4).
func writeRunArtifact(path string, artifact entities.RunArtifact) error {
	data, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return fmt.Errorf("connor: encode run artifact: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("connor: write %s: %w", path, err)
	}
	return nil
}
