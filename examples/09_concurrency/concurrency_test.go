package concurrencydemo

import (
	"reflect"
	"testing"
	"time"
)

func TestSquareAll(t *testing.T) {
	got := SquareAll([]int{1, 2, 3, 4})
	want := []int{1, 4, 9, 16}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SquareAll() = %#v, want %#v", got, want)
	}
}

func TestRaceMessages(t *testing.T) {
	got := RaceMessages(map[string]time.Duration{
		"slow": 20 * time.Millisecond,
		"fast": 2 * time.Millisecond,
	})

	if got != "fast" {
		t.Fatalf("RaceMessages() = %q, want %q", got, "fast")
	}
}

func TestWorkerPool(t *testing.T) {
	got := WorkerPool(2, []int{2, 3, 4})
	want := []int{4, 9, 16}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("WorkerPool() = %#v, want %#v", got, want)
	}
}
