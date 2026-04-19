package packagesgenerics

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Fatalf("Sum([]int) = %d, want 6", got)
	}

	if got := Sum([]float64{1.5, 2.5}); got != 4 {
		t.Fatalf("Sum([]float64) = %f, want 4", got)
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]string{"go", "go", "gopher", "go"})
	want := []string{"go", "gopher"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Unique() = %#v, want %#v", got, want)
	}
}

func TestNormalizeWords(t *testing.T) {
	got := NormalizeWords([]string{" Go ", "", "Gopher"})
	want := []string{"go", "gopher"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeWords() = %#v, want %#v", got, want)
	}
}

func TestFirstKeyword(t *testing.T) {
	if got := FirstKeyword("", "  ", "Go"); got != "Go" {
		t.Fatalf("FirstKeyword() = %q, want %q", got, "Go")
	}
}
