package contextdemo

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestWorkWithTimeoutSuccess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := WorkWithTimeout(ctx, 5*time.Millisecond); err != nil {
		t.Fatalf("WorkWithTimeout() error = %v", err)
	}
}

func TestWorkWithTimeoutDeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	err := WorkWithTimeout(ctx, 50*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("errors.Is(err, context.DeadlineExceeded) = false, err = %v", err)
	}
}

func TestStreamNumbers(t *testing.T) {
	var got []int
	for value := range StreamNumbers(context.Background(), 4) {
		got = append(got, value)
	}

	want := []int{0, 1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("StreamNumbers() = %#v, want %#v", got, want)
	}
}

func TestRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-001")

	requestID, ok := RequestID(ctx)
	if !ok {
		t.Fatal("RequestID() ok = false, want true")
	}

	if requestID != "req-001" {
		t.Fatalf("requestID = %q, want %q", requestID, "req-001")
	}
}
