package errorsdemo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidAge  = errors.New("invalid age")
	ErrNonPositive = errors.New("value must be positive")
)

type FieldError struct {
	Field  string
	Reason string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func ValidateAge(age int) error {
	if age < 0 || age > 150 {
		return ErrInvalidAge
	}

	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return &FieldError{
			Field:  "name",
			Reason: "cannot be blank",
		}
	}

	return nil
}

func RegisterUser(name string, age int) error {
	if err := ValidateName(name); err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	if err := ValidateAge(age); err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return nil
}

func ParsePositiveInt(raw string) (int, error) {
	number, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("parse positive int: %w", err)
	}

	if number <= 0 {
		return 0, fmt.Errorf("parse positive int: %w", ErrNonPositive)
	}

	return number, nil
}
