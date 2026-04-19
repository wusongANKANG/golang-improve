package performancedemo

import "testing"

func BenchmarkJoinWithPlus(b *testing.B) {
	parts := []string{"go", "-", "lang", "-", "builder"}
	for i := 0; i < b.N; i++ {
		_ = JoinWithPlus(parts)
	}
}

func BenchmarkJoinWithBuilder(b *testing.B) {
	parts := []string{"go", "-", "lang", "-", "builder"}
	for i := 0; i < b.N; i++ {
		_ = JoinWithBuilder(parts)
	}
}

func BenchmarkBuildNumbersNoPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BuildNumbersNoPrealloc(1024)
	}
}

func BenchmarkBuildNumbersPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BuildNumbersPrealloc(1024)
	}
}
