package collectionscontrol

import "strings"

func Sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}

	return total
}

func WordFrequency(words []string) map[string]int {
	result := make(map[string]int, len(words))
	for _, word := range words {
		normalized := strings.ToLower(strings.TrimSpace(word))
		if normalized == "" {
			continue
		}

		result[normalized]++
	}

	return result
}

func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	default:
		return "D"
	}
}

func SliceSharingDemo() (base []int, shared []int, safeCopy []int) {
	base = make([]int, 3, 4)
	copy(base, []int{1, 2, 3})

	shared = base[:2]
	shared = append(shared, 9)

	safeCopy = append([]int(nil), base...)
	safeCopy[0] = 100

	return base, shared, safeCopy
}
