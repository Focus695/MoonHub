package shield

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ApprovalStatus represents the current status of an approval request.
type ApprovalStatus string

const (
	// ApprovalPending indicates the request is waiting for a response.
	ApprovalPending ApprovalStatus = "pending"
	// ApprovalApproved indicates the request was approved.
	ApprovalApproved ApprovalStatus = "approved"
	// ApprovalRejected indicates the request was rejected.
	ApprovalRejected ApprovalStatus = "rejected"
	// ApprovalExpired indicates the request expired without a response.
	ApprovalExpired ApprovalStatus = "expired"
)

// ApprovalRequest represents a pending approval request for a shield action.
type ApprovalRequest struct {
	// ID is the unique identifier for this approval request.
	ID string `json:"id"`
	// ThreatID is the ID of the threat that triggered the approval.
	ThreatID string `json:"threat_id"`
	// ToolName is the name of the tool that requires approval.
	ToolName string `json:"tool_name"`
	// ToolArgs are the arguments passed to the tool.
	ToolArgs map[string]any `json:"tool_args"`
	// Reason explains why approval is required.
	Reason string `json:"reason"`
	// CreatedAt is when the request was created.
	CreatedAt time.Time `json:"created_at"`
	// ExpiresAt is when the request will expire.
	ExpiresAt time.Time `json:"expires_at"`
	// Status is the current status of the request.
	Status ApprovalStatus `json:"status"`
	// Scope is the event scope that triggered the approval.
	Scope ShieldScope `json:"scope"`
	// SkillName is the skill name if this is a skill-related approval.
	SkillName string `json:"skill_name,omitempty"`
	// Domain is the target domain if this is a network-related approval.
	Domain string `json:"domain,omitempty"`
	// URL is the target URL if this is a network-related approval.
	URL string `json:"url,omitempty"`
}

// ApprovalManager manages approval requests and their lifecycle.
type ApprovalManager struct {
	mu        sync.RWMutex
	requests  map[string]*ApprovalRequest
	responses map[string]chan bool
	timeout   time.Duration

	// pendingCh notifies when a new pending request is created.
	pendingCh chan *ApprovalRequest
}

// NewApprovalManager creates a new ApprovalManager with the specified timeout.
func NewApprovalManager(timeout time.Duration) *ApprovalManager {
	if timeout <= 0 {
		timeout = 5 * time.Minute
	}
	m := &ApprovalManager{
		requests:  make(map[string]*ApprovalRequest),
		responses: make(map[string]chan bool),
		timeout:   timeout,
		pendingCh: make(chan *ApprovalRequest, 100),
	}
	go m.cleanupExpired()
	return m
}

// CreateRequest creates a new approval request from a shield event and decision.
func (m *ApprovalManager) CreateRequest(event ShieldEvent, decision ShieldDecision) *ApprovalRequest {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	req := &ApprovalRequest{
		ID:        uuid.New().String(),
		ThreatID:  decision.ThreatID,
		ToolName:  event.ToolName,
		ToolArgs:  event.ToolArgs,
		Reason:    decision.Reason,
		CreatedAt: now,
		ExpiresAt: now.Add(m.timeout),
		Status:    ApprovalPending,
		Scope:     event.Scope,
		SkillName: event.SkillName,
		Domain:    event.Domain,
		URL:       event.URL,
	}

	m.requests[req.ID] = req
	m.responses[req.ID] = make(chan bool, 1)

	// Notify about new pending request (non-blocking)
	select {
	case m.pendingCh <- req:
	default:
	}

	return req
}

// WaitForApproval waits for a response to the approval request.
// Returns true if approved, false if rejected or expired.
func (m *ApprovalManager) WaitForApproval(ctx context.Context, id string) (bool, error) {
	m.mu.RLock()
	responseCh, exists := m.responses[id]
	m.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("approval request %s not found", id)
	}

	// Create a timer for expiration
	m.mu.RLock()
	req, reqExists := m.requests[id]
	m.mu.RUnlock()

	if !reqExists {
		return false, fmt.Errorf("approval request %s not found", id)
	}

	expireTimer := time.NewTimer(time.Until(req.ExpiresAt))
	defer expireTimer.Stop()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case approved, ok := <-responseCh:
		if !ok {
			return false, fmt.Errorf("approval channel closed")
		}
		return approved, nil
	case <-expireTimer.C:
		m.mu.Lock()
		if r, ok := m.requests[id]; ok && r.Status == ApprovalPending {
			r.Status = ApprovalExpired
		}
		delete(m.responses, id)
		m.mu.Unlock()
		return false, fmt.Errorf("approval request expired")
	}
}

// Approve approves an approval request by ID.
// Returns true if the request was successfully approved.
func (m *ApprovalManager) Approve(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	req, exists := m.requests[id]
	if !exists {
		return false
	}

	if req.Status != ApprovalPending {
		return false
	}

	req.Status = ApprovalApproved

	if responseCh, ok := m.responses[id]; ok {
		select {
		case responseCh <- true:
		default:
		}
	}

	return true
}

// Reject rejects an approval request by ID.
// Returns true if the request was successfully rejected.
func (m *ApprovalManager) Reject(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	req, exists := m.requests[id]
	if !exists {
		return false
	}

	if req.Status != ApprovalPending {
		return false
	}

	req.Status = ApprovalRejected

	if responseCh, ok := m.responses[id]; ok {
		select {
		case responseCh <- false:
		default:
		}
	}

	return true
}

// GetPending returns all pending approval requests.
func (m *ApprovalManager) GetPending() []*ApprovalRequest {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var pending []*ApprovalRequest
	for _, req := range m.requests {
		if req.Status == ApprovalPending {
			pending = append(pending, req)
		}
	}
	return pending
}

// GetRequest returns an approval request by ID.
func (m *ApprovalManager) GetRequest(id string) (*ApprovalRequest, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	req, exists := m.requests[id]
	if !exists {
		return nil, false
	}
	return req, true
}

// PendingChannel returns a channel that receives new pending approval requests.
func (m *ApprovalManager) PendingChannel() <-chan *ApprovalRequest {
	return m.pendingCh
}

// cleanupExpired periodically removes expired requests.
func (m *ApprovalManager) cleanupExpired() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for id, req := range m.requests {
			if req.Status == ApprovalPending && now.After(req.ExpiresAt) {
				req.Status = ApprovalExpired
				if ch, ok := m.responses[id]; ok {
					close(ch)
					delete(m.responses, id)
				}
			}
			// Remove old completed/expired requests after 5 minutes
			if req.Status != ApprovalPending && now.Sub(req.ExpiresAt) > 5*time.Minute {
				delete(m.requests, id)
			}
		}
		m.mu.Unlock()
	}
}

// Close releases resources held by the approval manager.
func (m *ApprovalManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, ch := range m.responses {
		close(ch)
		delete(m.responses, id)
	}
	m.requests = make(map[string]*ApprovalRequest)
}
