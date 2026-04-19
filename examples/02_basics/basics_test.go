package basics

import "testing"

func TestZeroValueProfile(t *testing.T) {
	profile := ZeroValueProfile()

	if profile.Name != "" || profile.Age != 0 || profile.Active || profile.Status != "" {
		t.Fatalf("unexpected zero value profile: %+v", profile)
	}
}

func TestNewProfile(t *testing.T) {
	profile := NewProfile("alice", 18)

	if profile.Name != "alice" || profile.Age != 18 || !profile.Active || profile.Status != StatusActive {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestSwap(t *testing.T) {
	left, right := Swap(1, 2)

	if left != 2 || right != 1 {
		t.Fatalf("Swap(1, 2) = (%d, %d), want (2, 1)", left, right)
	}
}

func TestTypedAndUntypedConstants(t *testing.T) {
	retries, pi := TypedAndUntypedConstants()

	if retries != 3 {
		t.Fatalf("retries = %d, want 3", retries)
	}

	if pi != 3.14 {
		t.Fatalf("pi = %f, want 3.14", pi)
	}
}
