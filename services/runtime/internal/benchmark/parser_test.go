package benchmark

import "testing"

func TestParse_valid(t *testing.T) {
	data := []byte(`
suite: test
cases:
  - id: ping
    model: gpt-4o-mini
    prompt: pong
`)
	spec, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Suite != "test" || len(spec.Cases) != 1 {
		t.Fatalf("%+v", spec)
	}
}

func TestParse_duplicateID(t *testing.T) {
	data := []byte(`
suite: test
cases:
  - id: a
    model: m
    prompt: p
  - id: a
    model: m
    prompt: p2
`)
	_, err := Parse(data)
	if err == nil {
		t.Fatal("expected error")
	}
}
