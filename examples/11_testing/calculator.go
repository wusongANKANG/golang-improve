package testingdemo

import "errors"

var ErrDivideByZero = errors.New("divide by zero")

func Add(a, b int) int {
	return a + b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}

	return a / b, nil
}

func Fibonacci(n int) int {
	if n < 2 {
		return n
	}

	prev, current := 0, 1
	for i := 2; i <= n; i++ {
		prev, current = current, prev+current
	}

	return current
}
