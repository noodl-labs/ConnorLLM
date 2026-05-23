package entities

import (
	"errors"
	"testing"
)

func TestNewRequest_ok(t *testing.T) {
	msgs := []Message{{Role: "user", Content: "ping"}}
	req, err := NewRequest("llama-3", msgs, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "llama-3" || !req.Stream {
		t.Fatalf("unexpected request: %+v", req)
	}
	if len(req.Messages) != 1 || req.Messages[0].Content != "ping" {
		t.Fatalf("messages: %+v", req.Messages)
	}
}

func TestNewRequest_emptyModel(t *testing.T) {
	_, err := NewRequest("", []Message{{Role: "user", Content: "x"}}, false)
	if !errors.Is(err, ErrEmptyModel) {
		t.Fatalf("want ErrEmptyModel, got %v", err)
	}
}

func TestNewRequest_emptyMessages(t *testing.T) {
	_, err := NewRequest("m", nil, false)
	if !errors.Is(err, ErrEmptyMessages) {
		t.Fatalf("want ErrEmptyMessages, got %v", err)
	}
}

func TestNewRequest_invalidRole(t *testing.T) {
	_, err := NewRequest("m", []Message{{Role: "tool", Content: "x"}}, false)
	if !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("want ErrInvalidRole, got %v", err)
	}
}

func TestNewRequest_emptyContent(t *testing.T) {
	_, err := NewRequest("m", []Message{{Role: "user", Content: ""}}, false)
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestNewRequest_defensiveCopy(t *testing.T) {
	msgs := []Message{{Role: "user", Content: "original"}}
	req, err := NewRequest("m", msgs, false)
	if err != nil {
		t.Fatal(err)
	}
	msgs[0].Content = "mutated"
	if req.Messages[0].Content != "original" {
		t.Fatalf("caller mutation affected request: %q", req.Messages[0].Content)
	}
}

func TestNewRequest_allValidRoles(t *testing.T) {
	roles := []string{"system", "user", "assistant"}
	for _, role := range roles {
		_, err := NewRequest("m", []Message{{Role: role, Content: "ok"}}, false)
		if err != nil {
			t.Fatalf("role %q: %v", role, err)
		}
	}
}
