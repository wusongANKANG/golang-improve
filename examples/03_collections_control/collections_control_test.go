package collectionscontrol

import "testing"

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3, 4}); got != 10 {
		t.Fatalf("Sum() = %d, want 10", got)
	}
}

func TestWordFrequency(t *testing.T) {
	freq := WordFrequency([]string{"Go", "go", "  Gopher ", "", "go"})

	if freq["go"] != 3 {
		t.Fatalf("freq[go] = %d, want 3", freq["go"])
	}

	if freq["gopher"] != 1 {
		t.Fatalf("freq[gopher] = %d, want 1", freq["gopher"])
	}
}

func TestGrade(t *testing.T) {
	if got := Grade(88); got != "B" {
		t.Fatalf("Grade(88) = %q, want %q", got, "B")
	}
}

func TestSliceSharingDemo(t *testing.T) {
	base, shared, safeCopy := SliceSharingDemo()

	if base[2] != 9 {
		t.Fatalf("base[2] = %d, want 9", base[2])
	}

	if shared[2] != 9 {
		t.Fatalf("shared[2] = %d, want 9", shared[2])
	}

	if safeCopy[0] != 100 {
		t.Fatalf("safeCopy[0] = %d, want 100", safeCopy[0])
	}

	if base[0] != 1 {
		t.Fatalf("base[0] = %d, want 1", base[0])
	}
}
