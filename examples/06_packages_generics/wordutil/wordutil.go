package wordutil

import "strings"

func CleanLower(words []string) []string {
	result := make([]string, 0, len(words))
	for _, word := range words {
		normalized := strings.ToLower(strings.TrimSpace(word))
		if normalized == "" {
			continue
		}

		result = append(result, normalized)
	}

	return result
}

func FirstNonEmpty(words ...string) string {
	for _, word := range words {
		if strings.TrimSpace(word) != "" {
			return word
		}
	}

	return ""
}
