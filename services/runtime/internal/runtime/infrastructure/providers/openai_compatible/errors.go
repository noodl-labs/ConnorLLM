package openai_compatible

import "errors"

var (
	ErrEmptyBaseURL   = errors.New("openai_compatible: base URL is required")
	ErrStreamingBeta1 = errors.New("openai_compatible: streaming not supported in Beta.1")
	ErrEnvNotWired    = errors.New("openai_compatible: CONNOR_BASE_URL and CONNOR_API_KEY not wired")
	ErrDecodeResponse = errors.New("openai_compatible: decode response")
	ErrEmptyChoices   = errors.New("openai_compatible: empty choices")
)
