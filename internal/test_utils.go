package internal

import (
	"context"
	"testing"
	"time"
)

// TryUntil tries to execute fn, but terminates early with a test failure if the context completes before the function does.
// Used to test that a function will terminate before the context is done.
func TryUntil(t *testing.T, deadline time.Duration, fn func() any) {
	t.Helper()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(deadline))
	defer cancel()

	success := make(chan bool)
	go func() {
		fn()
		success <- true
	}()
	go func() {
		<-ctx.Done()
		success <- false
	}()

	res := <-success
	if !res {
		t.Fatalf("failed to complete execution within %v", deadline)
	}
}
