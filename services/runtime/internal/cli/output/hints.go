package output

import (
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

func failHint(reason entities.FailReason, resp entities.Response) string {
	switch reason {
	case entities.FailReasonContentMismatch:
		return "response does not include expected text (case-sensitive)"
	case entities.FailReasonInvalidJSON:
		return "response is not valid JSON (markdown wrapper?)"
	case entities.FailReasonCallFailed:
		if resp.IsTimeout() {
			return "increase timeout_ms or check endpoint latency"
		}
		if resp.HTTPStatus >= 400 {
			return "check model slug and API key"
		}
		if resp.IsRetryExhausted() {
			return "transient errors exhausted retries; check endpoint health"
		}
	}
	return ""
}
