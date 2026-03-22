package provisioning

import (
	"crypto/rand"
	"encoding/json"
	"sync"
	"time"
)

// EventType categorizes device events.
type EventType string

const (
	EventTypeStatus EventType = "status"
	EventTypeAction EventType = "action"
	EventTypeSystem EventType = "system"
)

// EventLevel indicates severity of the event.
type EventLevel string

const (
	EventLevelInfo    EventLevel = "info"
	EventLevelSuccess EventLevel = "success"
	EventLevelWarning EventLevel = "warning"
	EventLevelError   EventLevel = "error"
)

// ActionPhase indicates the phase of an action.
type ActionPhase string

const (
	ActionPhaseStarted   ActionPhase = "started"
	ActionPhaseCompleted ActionPhase = "completed"
	ActionPhaseFailed    ActionPhase = "failed"
)

// DeviceManagerEvent represents an event emitted by the device manager.
type DeviceManagerEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Level     EventLevel             `json:"level"`
	Message   string                 `json:"message"`
	Timestamp int64                  `json:"timestamp"`
	Action    string                 `json:"action,omitempty"`
	Phase     ActionPhase            `json:"phase,omitempty"`
	Status    *DeviceStatus          `json:"status,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// EventBroadcaster manages SSE client connections and event distribution.
type EventBroadcaster struct {
	mu           sync.RWMutex
	clients      map[chan []byte]struct{}
	recentEvents []DeviceManagerEvent
	maxRecent    int
}

// NewEventBroadcaster creates a new event broadcaster.
func NewEventBroadcaster(maxRecent int) *EventBroadcaster {
	if maxRecent <= 0 {
		maxRecent = 25
	}
	return &EventBroadcaster{
		clients:      make(map[chan []byte]struct{}),
		recentEvents: make([]DeviceManagerEvent, 0, maxRecent),
		maxRecent:    maxRecent,
	}
}

// Subscribe creates a new client channel for receiving events.
func (b *EventBroadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 16)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes a client channel.
func (b *EventBroadcaster) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
}

// Emit creates and broadcasts a new event.
func (b *EventBroadcaster) Emit(eventType EventType, level EventLevel, message string, opts ...EventOption) *DeviceManagerEvent {
	event := &DeviceManagerEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UnixMilli(),
		Details:   make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(event)
	}

	b.publish(event)
	return event
}

// EmitAction emits an action event with the specified phase.
func (b *EventBroadcaster) EmitAction(action string, phase ActionPhase, level EventLevel, message string, opts ...EventOption) *DeviceManagerEvent {
	allOpts := append([]EventOption{WithAction(action), WithPhase(phase)}, opts...)
	return b.Emit(EventTypeAction, level, message, allOpts...)
}

// EmitWithStatus emits an event and includes current device status.
func (b *EventBroadcaster) EmitWithStatus(eventType EventType, level EventLevel, message string, status DeviceStatus, opts ...EventOption) *DeviceManagerEvent {
	allOpts := append([]EventOption{WithStatus(&status)}, opts...)
	return b.Emit(eventType, level, message, allOpts...)
}

func (b *EventBroadcaster) publish(event *DeviceManagerEvent) {
	// Store in recent events
	b.mu.Lock()
	b.recentEvents = append(b.recentEvents, *event)
	if len(b.recentEvents) > b.maxRecent {
		b.recentEvents = b.recentEvents[len(b.recentEvents)-b.maxRecent:]
	}
	clients := make([]chan []byte, 0, len(b.clients))
	for ch := range b.clients {
		clients = append(clients, ch)
	}
	b.mu.Unlock()

	// Encode and broadcast
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	for _, ch := range clients {
		select {
		case ch <- data:
		default:
			// Client buffer full, skip
		}
	}
}

// GetRecentEvents returns a copy of recent events for new SSE clients.
func (b *EventBroadcaster) GetRecentEvents() []DeviceManagerEvent {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]DeviceManagerEvent, len(b.recentEvents))
	copy(result, b.recentEvents)
	return result
}

// ClientCount returns the number of connected SSE clients.
func (b *EventBroadcaster) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}

// EventOption is a functional option for configuring events.
type EventOption func(*DeviceManagerEvent)

// WithAction sets the action name for the event.
func WithAction(action string) EventOption {
	return func(e *DeviceManagerEvent) {
		e.Action = action
	}
}

// WithPhase sets the action phase for the event.
func WithPhase(phase ActionPhase) EventOption {
	return func(e *DeviceManagerEvent) {
		e.Phase = phase
	}
}

// WithStatus sets the device status for the event.
func WithStatus(status *DeviceStatus) EventOption {
	return func(e *DeviceManagerEvent) {
		e.Status = status
	}
}

// WithDetails sets additional details for the event.
func WithDetails(details map[string]interface{}) EventOption {
	return func(e *DeviceManagerEvent) {
		if e.Details == nil {
			e.Details = make(map[string]interface{})
		}
		for k, v := range details {
			e.Details[k] = v
		}
	}
}

func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomSuffix()
}

func randomSuffix() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 6)
	randBytes := make([]byte, 6)
	if _, err := rand.Read(randBytes); err != nil {
		now := time.Now().UnixNano()
		for i := range result {
			result[i] = chars[(now+int64(i))%int64(len(chars))]
		}
		return string(result)
	}
	for i := range result {
		result[i] = chars[int(randBytes[i])%len(chars)]
	}
	return string(result)
}
