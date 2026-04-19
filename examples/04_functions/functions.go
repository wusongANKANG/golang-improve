package functions

import (
	"errors"
	"fmt"
)

var ErrDivideByZero = errors.New("divide by zero")

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}

	return a / b, nil
}

func Accumulator(start int) func(delta int) int {
	sum := start

	return func(delta int) int {
		sum += delta
		return sum
	}
}

func SafeExecute(fn func() string) (result string, recovered any) {
	defer func() {
		if r := recover(); r != nil {
			recovered = r
			result = fmt.Sprintf("recovered: %v", r)
		}
	}()

	return fn(), nil
}
