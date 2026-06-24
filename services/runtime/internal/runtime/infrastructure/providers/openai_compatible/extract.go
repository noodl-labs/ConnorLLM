package openai_compatible

import (
	"encoding/json"
)

// AssistantContent returns choices[0].message.content from a non-stream response body.
// Response.Body in ConnorLLM = this string (not the raw HTTP envelope).
func AssistantContent(body []byte) (string, error) {
	var resp chatCompletionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", ErrDecodeResponse
	}
	if len(resp.Choices) == 0 {
		return "", ErrEmptyChoices
	}
	return resp.Choices[0].Message.Content, nil
}
