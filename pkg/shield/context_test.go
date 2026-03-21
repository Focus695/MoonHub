package shield

import (
	"context"
	"testing"
)

func TestApprovedToolExecution(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	if ApprovedToolExecution(ctx) {
		t.Fatal("expected false on plain context")
	}
	ctx2 := ContextWithApprovedToolExecution(ctx)
	if !ApprovedToolExecution(ctx2) {
		t.Fatal("expected true after ContextWithApprovedToolExecution")
	}
}
