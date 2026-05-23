package entities

/*
This file contains the entities for the runtime. What is sent to the LLM provider.
Messages are what the runtime should ask the LLM provider to execute.
It is a domain value object: no HTTP, no API keys, no retry/timeout policies.
*/
import "errors"

// Message is one turn in a chat/completions call (OpenAI-compatible roles).
type Message struct {
	// Role must be system, user, or assistant for providers that follow the OpenAI schema.
	Role string
	// Content is the text sent for this turn (MVP: plain string only).
	Content string
}

// Correlation IDs (run_id, case_id) belong on ExecuteCase input or logging context,
// not on Request — keeps the payload stable and easy to replay.
type Request struct {
	// Model is the provider model id (e.g. "meta-llama-3-8b" on vLLM).
	Model string
	// Messages is the full prompt thread for this benchmark case.
	Messages []Message
	// Stream enables token streaming; when true, infrastructure must measure TTFT.
	Stream bool
}

var (
	ErrEmptyModel    = errors.New("request: model is required")
	ErrEmptyMessages = errors.New("request: at least one message is required")
	ErrInvalidRole   = errors.New("request: invalid message role")
)

// NewRequest builds a validated Request. Call this from application layer
// when mapping YAML test cases (or tests) into a provider call.
func NewRequest(model string, messages []Message, stream bool) (Request, error) {
	if model == "" {
		return Request{}, ErrEmptyModel
	}
	if len(messages) == 0 {
		return Request{}, ErrEmptyMessages
	}

	// Defensive copy so callers cannot mutate internal state after construction.
	out := make([]Message, len(messages))
	for i, m := range messages {
		if !isValidRole(m.Role) {
			return Request{}, ErrInvalidRole
		}
		if m.Content == "" {
			return Request{}, errors.New("request: message content is required")
		}
		out[i] = m
	}

	return Request{
		Model:    model,
		Messages: out,
		Stream:   stream,
	}, nil
}

func isValidRole(role string) bool {
	switch role {
	case "system", "user", "assistant":
		return true
	default:
		return false
	}
}
