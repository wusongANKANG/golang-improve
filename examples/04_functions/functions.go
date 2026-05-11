package functions

import (
	"errors"
	"fmt"
)

var ErrDivideByZero = errors.New("divide by zero")
var ErrInvalidRetryTimes = errors.New("retry times must be positive")

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

func Retry(times int, fn func() error) error {
	if times <= 0 {
		return ErrInvalidRetryTimes
	}

	var err error
	for i := 0; i < times; i++ {
		if err = fn(); err == nil {
			return nil
		}
	}

	return err
}

func SafeCall(fn func() int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case error:
				err = fmt.Errorf("safe call panic: %w", v)
			default:
				err = fmt.Errorf("safe call panic: %v", v)
			}
		}
	}()

	return fn(), nil
}
