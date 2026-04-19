package performancedemo

import "testing"

func TestJoinFunctions(t *testing.T) {
	parts := []string{"go", "-", "pher"}
	want := "go-pher"

	if got := JoinWithPlus(parts); got != want {
		t.Fatalf("JoinWithPlus() = %q, want %q", got, want)
	}

	if got := JoinWithBuilder(parts); got != want {
		t.Fatalf("JoinWithBuilder() = %q, want %q", got, want)
	}
}

func TestBuildNumbersPrealloc(t *testing.T) {
	numbers := BuildNumbersPrealloc(5)

	if len(numbers) != 5 {
		t.Fatalf("len(numbers) = %d, want 5", len(numbers))
	}

	if cap(numbers) != 5 {
		t.Fatalf("cap(numbers) = %d, want 5", cap(numbers))
	}
}

func TestSafeSubset(t *testing.T) {
	source := []int{1, 2, 3, 4}
	subset := SafeSubset(source, 2)
	subset[0] = 100

	if source[0] != 1 {
		t.Fatalf("source[0] = %d, want 1", source[0])
	}
}
