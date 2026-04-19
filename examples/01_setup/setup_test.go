package setup

import "testing"

func TestGreeting(t *testing.T) {
	if got := Greeting("alice"); got != "hello, alice" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, alice")
	}

	if got := Greeting(""); got != "hello, gopher" {
		t.Fatalf("Greeting() = %q, want %q", got, "hello, gopher")
	}
}

func TestProgramShape(t *testing.T) {
	shape := ProgramShape()
	if len(shape) != 3 {
		t.Fatalf("len(ProgramShape()) = %d, want 3", len(shape))
	}
}
