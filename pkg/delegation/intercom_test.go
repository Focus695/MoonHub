package delegation

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestIntercom_UnsubscribeFirstDoesNotBreakSecond(t *testing.T) {
	ic := NewIntercom(10)
	var second atomic.Bool
	u1 := ic.On("t", func(context.Context, IntercomEvent) {})
	u2 := ic.On("t", func(context.Context, IntercomEvent) { second.Store(true) })
	u1()
	ic.Emit(context.Background(), "t", "user", nil)
	if !second.Load() {
		t.Fatal("expected second handler to run after first unsubscribed")
	}
	_ = u2
}

func TestIntercom_OnAnyUnsubscribeByID(t *testing.T) {
	ic := NewIntercom(10)
	var hit atomic.Bool
	u := ic.OnAny(func(context.Context, IntercomEvent) { hit.Store(true) })
	u()
	ic.Emit(context.Background(), "any", "u", nil)
	if hit.Load() {
		t.Fatal("handler should not run after unsubscribe")
	}
}
