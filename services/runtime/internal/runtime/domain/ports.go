package domain

import (
	"context"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

// ProviderExecutor runs a single provider call (one attempt). Retry and timeout are applied by ExecuteCase.
type ProviderExecutor interface {
	Execute(ctx context.Context, req entities.Request) (entities.Response, error)
}
