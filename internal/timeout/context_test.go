package timeout

import (
	"context"
	"testing"
	"time"
)

func TestContextZeroTimeout(t *testing.T) {
	ctx, cancel := Context(0)
	defer cancel()

	if _, ok := ctx.Deadline(); ok {
		t.Errorf("expected no deadline when timeout is 0")
	}
	if err := ctx.Err(); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestContextWithTimeout(t *testing.T) {
	ctx, cancel := Context(50 * time.Millisecond)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatalf("expected a deadline when timeout > 0")
	}
	if time.Until(deadline) <= 0 {
		t.Errorf("expected deadline in the future, got %v", deadline)
	}
}

func TestContextDeadlineExceeded(t *testing.T) {
	ctx, cancel := Context(10 * time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != context.DeadlineExceeded {
			t.Errorf("expected context.DeadlineExceeded, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("context did not time out within expected window")
	}
}

func TestContextCancelable(t *testing.T) {
	ctx, cancel := Context(time.Hour)
	cancel()

	select {
	case <-ctx.Done():
		if err := ctx.Err(); err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Error("context did not cancel after cancel()")
	}
}
