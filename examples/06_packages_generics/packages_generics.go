package packagesgenerics

import "golang-improve/examples/06_packages_generics/wordutil"

type Number interface {
	~int | ~int64 | ~float64
}

func Sum[T Number](items []T) T {
	var total T
	for _, item := range items {
		total += item
	}

	return total
}

func Unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}, len(items))
	result := make([]T, 0, len(items))

	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

func NormalizeWords(words []string) []string {
	return wordutil.CleanLower(words)
}

func FirstKeyword(words ...string) string {
	return wordutil.FirstNonEmpty(words...)
}
