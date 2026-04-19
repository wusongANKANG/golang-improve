package performancedemo

import "strings"

func JoinWithPlus(parts []string) string {
	joined := ""
	for _, part := range parts {
		joined += part
	}

	return joined
}

func JoinWithBuilder(parts []string) string {
	var builder strings.Builder

	totalLength := 0
	for _, part := range parts {
		totalLength += len(part)
	}

	builder.Grow(totalLength)
	for _, part := range parts {
		builder.WriteString(part)
	}

	return builder.String()
}

func BuildNumbersNoPrealloc(n int) []int {
	var numbers []int
	for i := 0; i < n; i++ {
		numbers = append(numbers, i)
	}

	return numbers
}

func BuildNumbersPrealloc(n int) []int {
	numbers := make([]int, 0, n)
	for i := 0; i < n; i++ {
		numbers = append(numbers, i)
	}

	return numbers
}

func SafeSubset(items []int, n int) []int {
	if n > len(items) {
		n = len(items)
	}

	return append([]int(nil), items[:n]...)
}
