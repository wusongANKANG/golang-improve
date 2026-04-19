package errorsdemo

import (
	"errors"
	"strconv"
	"testing"
)

func TestRegisterUserInvalidAge(t *testing.T) {
	err := RegisterUser("alice", -1)
	if !errors.Is(err, ErrInvalidAge) {
		t.Fatalf("errors.Is(err, ErrInvalidAge) = false, err = %v", err)
	}
}

func TestRegisterUserInvalidName(t *testing.T) {
	err := RegisterUser("", 18)

	var fieldErr *FieldError
	if !errors.As(err, &fieldErr) {
		t.Fatalf("errors.As(err, *FieldError) = false, err = %v", err)
	}

	if fieldErr.Field != "name" {
		t.Fatalf("fieldErr.Field = %q, want %q", fieldErr.Field, "name")
	}
}

func TestParsePositiveIntSyntax(t *testing.T) {
	_, err := ParsePositiveInt("abc")

	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Fatalf("errors.As(err, *strconv.NumError) = false, err = %v", err)
	}
}

func TestParsePositiveIntNonPositive(t *testing.T) {
	_, err := ParsePositiveInt("0")
	if !errors.Is(err, ErrNonPositive) {
		t.Fatalf("errors.Is(err, ErrNonPositive) = false, err = %v", err)
	}
}
