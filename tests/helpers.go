package tests

import (
	"context"
	"errors"
	"testing"
)

var ErrTestTimeout = errors.New("tests timeout exceeded")

func GetContext(t *testing.T) context.Context {
	t.Helper()

	ctx := context.Background()

	deadline, ok := t.Deadline()
	if ok {
		ctx, cancel := context.WithDeadlineCause(ctx, deadline, ErrTestTimeout)
		t.Cleanup(cancel)

		return ctx
	}

	return ctx
}
