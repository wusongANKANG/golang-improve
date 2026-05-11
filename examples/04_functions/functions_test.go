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

func TestSafeExecuteNoPanic(t *testing.T) {
	result, recovered := SafeExecute(func() string {
		return "hello recover"
	})

	if recovered != nil {
		t.Fatalf("recovered = %v, want nil", recovered)
	}

	if result != "hello recover" {
		t.Fatalf("result = %q, want %q", result, "hello recover")
	}
}

func TestRetryEventuallySuccess(t *testing.T) {
	attempts := 0
	err := Retry(3, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary failure")
		}

		return nil
	})

	if err != nil {
		t.Fatalf("Retry() error = %v, want nil", err)
	}

	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestRetryExhausted(t *testing.T) {
	wantErr := errors.New("still failing")
	attempts := 0

	err := Retry(2, func() error {
		attempts++
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Retry() error = %v, want %v", err, wantErr)
	}

	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestRetryInvalidTimes(t *testing.T) {
	called := false
	err := Retry(0, func() error {
		called = true
		return nil
	})

	if !errors.Is(err, ErrInvalidRetryTimes) {
		t.Fatalf("Retry() error = %v, want %v", err, ErrInvalidRetryTimes)
	}

	if called {
		t.Fatal("expected fn not to be called")
	}
}

func TestSafeCall(t *testing.T) {
	got, err := SafeCall(func() int {
		return 42
	})
	if err != nil {
		t.Fatalf("SafeCall() error = %v, want nil", err)
	}

	if got != 42 {
		t.Fatalf("SafeCall() = %d, want 42", got)
	}
}

func TestSafeCallRecover(t *testing.T) {
	_, err := SafeCall(func() int {
		panic("boom")
	})
	if err == nil {
		t.Fatal("SafeCall() error = nil, want non-nil")
	}

	if err.Error() != "safe call panic: boom" {
		t.Fatalf("SafeCall() error = %q, want %q", err.Error(), "safe call panic: boom")
	}
}
