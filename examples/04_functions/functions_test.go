package functions

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil {
		t.Fatalf("Divide() error = %v", err)
	}

	if got != 5 {
		t.Fatalf("Divide() = %f, want 5", got)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(10, 0)
	if !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide() error = %v, want %v", err, ErrDivideByZero)
	}
}

func TestAccumulator(t *testing.T) {
	acc := Accumulator(10)

	if got := acc(5); got != 15 {
		t.Fatalf("first acc() = %d, want 15", got)
	}

	if got := acc(3); got != 18 {
		t.Fatalf("second acc() = %d, want 18", got)
	}
}

func TestSafeExecuteRecover(t *testing.T) {
	result, recovered := SafeExecute(func() string {
		panic("boom")
	})

	if recovered == nil {
		t.Fatal("expected recovered value, got nil")
	}

	if result != "recovered: boom" {
		t.Fatalf("result = %q, want %q", result, "recovered: boom")
	}
}
