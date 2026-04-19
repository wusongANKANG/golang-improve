package testingdemo

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	if got := Add(2, 3); got != 5 {
		t.Fatalf("Add() = %d, want 5", got)
	}
}

func TestDivide(t *testing.T) {
	testCases := []struct {
		name    string
		a       float64
		b       float64
		want    float64
		wantErr error
	}{
		{name: "normal", a: 10, b: 2, want: 5},
		{name: "zero divisor", a: 10, b: 0, wantErr: ErrDivideByZero},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := Divide(testCase.a, testCase.b)

			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("Divide() error = %v, want %v", err, testCase.wantErr)
			}

			if got != testCase.want {
				t.Fatalf("Divide() = %f, want %f", got, testCase.want)
			}
		})
	}
}
