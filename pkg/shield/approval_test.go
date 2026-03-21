package shield

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewApprovalManager(t *testing.T) {
	// Test with default timeout
	m := NewApprovalManager(0)
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if m.timeout != 5*time.Minute {
		t.Errorf("expected default timeout of 5m, got %v", m.timeout)
	}
	m.Close()

	// Test with custom timeout
	m = NewApprovalManager(10 * time.Minute)
	if m.timeout != 10*time.Minute {
		t.Errorf("expected timeout of 10m, got %v", m.timeout)
	}
	m.Close()
}

func TestCreateRequest(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
		ToolArgs: map[string]any{"command": "rm -rf /"},
	}

	decision := ShieldDecision{
		Action:    ActionRequireApproval,
		ThreatID:  "THREAT-001",
		Reason:    "Dangerous command detected",
		Scope:     ScopeToolCall,
		MatchedOn: "tool_name",
	}

	req := m.CreateRequest(event, decision)

	if req.ID == "" {
		t.Error("expected non-empty request ID")
	}
	if req.ThreatID != "THREAT-001" {
		t.Errorf("expected threat ID THREAT-001, got %s", req.ThreatID)
	}
	if req.ToolName != "exec" {
		t.Errorf("expected tool name exec, got %s", req.ToolName)
	}
	if req.Status != ApprovalPending {
		t.Errorf("expected status pending, got %s", req.Status)
	}
	if req.ExpiresAt.Before(req.CreatedAt) {
		t.Error("expected expires_at to be after created_at")
	}
}

func TestApproveReject(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	req := m.CreateRequest(event, decision)

	// Test approve
	if !m.Approve(req.ID) {
		t.Error("expected approve to succeed")
	}

	stored, exists := m.GetRequest(req.ID)
	if !exists {
		t.Fatal("expected request to exist")
	}
	if stored.Status != ApprovalApproved {
		t.Errorf("expected status approved, got %s", stored.Status)
	}

	// Test approve already approved request
	if m.Approve(req.ID) {
		t.Error("expected approve to fail for already approved request")
	}

	// Create new request for reject test
	req2 := m.CreateRequest(event, decision)

	// Test reject
	if !m.Reject(req2.ID) {
		t.Error("expected reject to succeed")
	}

	stored2, exists := m.GetRequest(req2.ID)
	if !exists {
		t.Fatal("expected request to exist")
	}
	if stored2.Status != ApprovalRejected {
		t.Errorf("expected status rejected, got %s", stored2.Status)
	}

	// Test reject already rejected request
	if m.Reject(req2.ID) {
		t.Error("expected reject to fail for already rejected request")
	}

	// Test non-existent request
	if m.Approve("non-existent") {
		t.Error("expected approve to fail for non-existent request")
	}
	if m.Reject("non-existent") {
		t.Error("expected reject to fail for non-existent request")
	}
}

func TestWaitForApproval(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	t.Run("approved", func(t *testing.T) {
		req := m.CreateRequest(event, decision)

		var approved bool
		var err error
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			approved, err = m.WaitForApproval(context.Background(), req.ID)
		}()

		// Small delay to ensure WaitForApproval is waiting
		time.Sleep(10 * time.Millisecond)
		m.Approve(req.ID)
		wg.Wait()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !approved {
			t.Error("expected approval to be true")
		}
	})

	t.Run("rejected", func(t *testing.T) {
		req := m.CreateRequest(event, decision)

		var approved bool
		var err error
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			approved, err = m.WaitForApproval(context.Background(), req.ID)
		}()

		time.Sleep(10 * time.Millisecond)
		m.Reject(req.ID)
		wg.Wait()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if approved {
			t.Error("expected approval to be false")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := m.WaitForApproval(context.Background(), "non-existent")
		if err == nil {
			t.Error("expected error for non-existent request")
		}
	})
}

func TestExpiration(t *testing.T) {
	// Use very short timeout for testing
	m := NewApprovalManager(100 * time.Millisecond)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	req := m.CreateRequest(event, decision)

	// Wait for expiration
	approved, err := m.WaitForApproval(context.Background(), req.ID)

	if err == nil {
		t.Error("expected expiration error")
	}
	if approved {
		t.Error("expected approval to be false on expiration")
	}

	// Verify status is expired
	stored, exists := m.GetRequest(req.ID)
	if !exists {
		t.Fatal("expected request to exist")
	}
	if stored.Status != ApprovalExpired {
		t.Errorf("expected status expired, got %s", stored.Status)
	}
}

func TestConcurrentAccess(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	var wg sync.WaitGroup
	numRequests := 100

	// Concurrent create
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.CreateRequest(event, decision)
		}()
	}
	wg.Wait()

	pending := m.GetPending()
	if len(pending) != numRequests {
		t.Errorf("expected %d pending requests, got %d", numRequests, len(pending))
	}

	// Concurrent approve/reject
	for i, req := range pending {
		wg.Add(1)
		go func(idx int, reqID string) {
			defer wg.Done()
			if idx%2 == 0 {
				m.Approve(reqID)
			} else {
				m.Reject(reqID)
			}
		}(i, req.ID)
	}
	wg.Wait()

	// Verify all requests are processed
	approvedCount := 0
	rejectedCount := 0
	for _, req := range pending {
		stored, _ := m.GetRequest(req.ID)
		switch stored.Status {
		case ApprovalApproved:
			approvedCount++
		case ApprovalRejected:
			rejectedCount++
		}
	}

	if approvedCount+rejectedCount != numRequests {
		t.Errorf("expected all requests to be approved or rejected, got %d approved, %d rejected",
			approvedCount, rejectedCount)
	}
}

func TestGetPending(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	// Create multiple requests
	req1 := m.CreateRequest(event, decision)
	req2 := m.CreateRequest(event, decision)
	req3 := m.CreateRequest(event, decision)

	// Approve one
	m.Approve(req1.ID)

	pending := m.GetPending()
	if len(pending) != 2 {
		t.Errorf("expected 2 pending requests, got %d", len(pending))
	}

	// Verify pending contains only pending requests
	for _, p := range pending {
		if p.ID == req1.ID {
			t.Error("expected req1 to not be in pending list")
		}
		if p.ID != req2.ID && p.ID != req3.ID {
			t.Errorf("unexpected request ID in pending: %s", p.ID)
		}
	}
}

func TestPendingChannel(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	// Create request and verify it's sent to pending channel
	req := m.CreateRequest(event, decision)

	select {
	case pending := <-m.PendingChannel():
		if pending.ID != req.ID {
			t.Errorf("expected request ID %s, got %s", req.ID, pending.ID)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected request on pending channel")
	}
}

func TestContextCancellation(t *testing.T) {
	m := NewApprovalManager(5 * time.Minute)
	defer m.Close()

	event := ShieldEvent{
		Scope:    ScopeToolCall,
		ToolName: "exec",
	}

	decision := ShieldDecision{
		Action:   ActionRequireApproval,
		ThreatID: "THREAT-001",
		Reason:   "Test",
	}

	req := m.CreateRequest(event, decision)

	ctx, cancel := context.WithCancel(context.Background())

	var approved bool
	var err error
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		approved, err = m.WaitForApproval(ctx, req.ID)
	}()

	// Cancel context
	cancel()
	wg.Wait()

	if err == nil {
		t.Error("expected context cancellation error")
	}
	if approved {
		t.Error("expected approval to be false on context cancellation")
	}
}
