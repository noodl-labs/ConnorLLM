package validation

import "testing"

func TestCheck(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"object ok", `{"a":1,"b":2}`, true},
		{"broken object", `{ "name": "John"`, false},
		{"empty", "", false},
		{"plain text", "not json", false},
		{"number alone", "42", true},
		{"string alone", `"hello"`, true},
		{"empty array", "[]", true},
		{"empty object", "{}", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.in); got != tt.want {
				t.Fatalf("Check(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
