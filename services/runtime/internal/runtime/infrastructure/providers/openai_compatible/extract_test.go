package openai_compatible

import "testing"

func TestAssistantContent_ok(t *testing.T) {
	body := []byte(`{
		"choices":[{"message":{"role":"assistant","content":"{\"ok\":true}"}}]
	}`)
	got, err := AssistantContent(body)
	if err != nil || got != `{"ok":true}` {
		t.Fatalf("got %q err=%v", got, err)
	}
}

func TestAssistantContent_emptyChoices(t *testing.T) {
	_, err := AssistantContent([]byte(`{"choices":[]}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAssistantContent_invalidJSON(t *testing.T) {
	_, err := AssistantContent([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
}
