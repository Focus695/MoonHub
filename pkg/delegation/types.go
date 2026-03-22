// Package delegation implements the sub-agent delegation system with intercom.
// It enables non-blocking task delegation, role template reuse, and adaptive timeouts.
package delegation

import (
	"time"
)

// SubAgentStatus represents the current state of a sub-agent
type SubAgentStatus string

const (
	StatusActive    SubAgentStatus = "active"
	StatusSuspended SubAgentStatus = "suspended"
	StatusDismissed SubAgentStatus = "dismissed"
)

// TaskCategory represents the type of task for timeout estimation
type TaskCategory string

const (
	CategoryResearch TaskCategory = "research"
	CategoryCode     TaskCategory = "code"
	CategoryAnalysis TaskCategory = "analysis"
	CategoryWriting  TaskCategory = "writing"
	CategoryGeneral  TaskCategory = "general"
	CategoryCustom   TaskCategory = "custom"
)

// SubAgentRecord represents a persistent sub-agent instance
type SubAgentRecord struct {
	ID             string         `json:"id"`
	UserID         string         `json:"user_id"`
	Label          string         `json:"label"`
	RolePrompt     string         `json:"role_prompt"`
	TaskKeywords   []string       `json:"task_keywords"`
	Tools          []string       `json:"tools"`
	Status         SubAgentStatus `json:"status"`
	CreatedAt      int64          `json:"created_at"`
	LastActiveAt   int64          `json:"last_active_at"`
	CompletedTasks int            `json:"completed_tasks"`
	SuccessRate    float64        `json:"success_rate"`
	SessionKey     string         `json:"session_key"`
}

// BackgroundTaskRecord represents a background task execution
type BackgroundTaskRecord struct {
	ID            string       `json:"id"`
	SubAgentID    string       `json:"sub_agent_id"`
	UserID        string       `json:"user_id"`
	Task          string       `json:"task"`
	Category      TaskCategory `json:"category"`
	Status        string       `json:"status"` // pending, running, completed, failed
	Result        string       `json:"result"`
	CreatedAt     int64        `json:"created_at"`
	CompletedAt   int64        `json:"completed_at"`
	Delivered     bool         `json:"delivered"`
	OriginChannel string       `json:"origin_channel"`
	OriginChatID  string       `json:"origin_chat_id"`
}

// RoleTemplate represents a reusable role configuration
type RoleTemplate struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	Name          string       `json:"name"`
	RolePrompt    string       `json:"role_prompt"`
	TaskKeywords  []string     `json:"task_keywords"`
	Tools         []string     `json:"tools"`
	Category      TaskCategory `json:"category"`
	SuccessCount  int          `json:"success_count"`
	FailureCount  int          `json:"failure_count"`
	AvgDurationMs int64        `json:"avg_duration_ms"`
	CreatedAt     int64        `json:"created_at"`
	LastUsedAt    int64        `json:"last_used_at"`
}

// TaskMetricRecord represents timing metrics for adaptive timeouts
type TaskMetricRecord struct {
	ID           string       `json:"id"`
	UserID       string       `json:"user_id"`
	TaskCategory TaskCategory `json:"task_category"`
	TaskHint     string       `json:"task_hint"` // First 100 chars of task
	DurationMs   int64        `json:"duration_ms"`
	Success      bool         `json:"success"`
	CreatedAt    int64        `json:"created_at"`
}

// SubAgentMessage represents a message in a sub-agent's history
type SubAgentMessage struct {
	ID         string `json:"id"`
	SubAgentID string `json:"sub_agent_id"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}

// OrientationContext provides context for sub-agent orientation
type OrientationContext struct {
	SessionSummary  string   `json:"session_summary"`
	RecentTasks     []string `json:"recent_tasks"`
	UserPreferences []string `json:"user_preferences"`
}

// IntercomEvent represents an event in the intercom system
type IntercomEvent struct {
	Topic     string         `json:"topic"`
	UserID    string         `json:"user_id"`
	Data      map[string]any `json:"data"`
	Timestamp int64          `json:"timestamp"`
}

// DelegationConfig holds configuration for the delegation system
type DelegationConfig struct {
	Enabled              bool     `json:"enabled"`
	DefaultSubAgentTools []string `json:"default_sub_agent_tools"`
	MaxActivePerUser     int      `json:"max_active_per_user"`
	MaxConcurrentTasks   int      `json:"max_concurrent_tasks"`
	RetentionDays        int      `json:"retention_days"`
	ReuseThreshold       float64  `json:"reuse_threshold"` // Keyword overlap threshold for reuse
}

// DefaultDelegationConfig returns the default configuration
func DefaultDelegationConfig() DelegationConfig {
	return DelegationConfig{
		Enabled:              true,
		DefaultSubAgentTools: []string{"read_file", "list_dir", "web_search", "web_fetch", "memory_recall"},
		MaxActivePerUser:     10,
		MaxConcurrentTasks:   3,
		RetentionDays:        14,
		ReuseThreshold:       0.6,
	}
}

// DelegateTaskRequest represents a task delegation request
type DelegateTaskRequest struct {
	Task          string       `json:"task"`
	Label         string       `json:"label"`
	Category      TaskCategory `json:"category"`
	Tools         []string     `json:"tools"`
	RolePrompt    string       `json:"role_prompt"`
	Background    bool         `json:"background"`
	TimeoutHint   int          `json:"timeout_hint"` // seconds, 0 for auto
	ReuseExisting bool         `json:"reuse_existing"`
}

// DelegateTaskResult represents the result of task delegation
type DelegateTaskResult struct {
	SubAgentID   string `json:"sub_agent_id"`
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	Result       string `json:"result,omitempty"`
	IsBackground bool   `json:"is_background"`
}

// SubAgentSummary represents a summary of a sub-agent for listing
type SubAgentSummary struct {
	ID             string         `json:"id"`
	Label          string         `json:"label"`
	Status         SubAgentStatus `json:"status"`
	CompletedTasks int            `json:"completed_tasks"`
	SuccessRate    float64        `json:"success_rate"`
	LastActiveAt   int64          `json:"last_active_at"`
	CreatedAt      int64          `json:"created_at"`
}

// TimeoutEstimate represents an estimated timeout for a task
type TimeoutEstimate struct {
	Duration   time.Duration `json:"duration"`
	Category   TaskCategory  `json:"category"`
	Confidence float64       `json:"confidence"` // 0-1, based on sample size
}

// BlackboardProposal represents a proposal in the blackboard collaboration mode
type BlackboardProposal struct {
	SubAgentID string    `json:"sub_agent_id"`
	Proposal   string    `json:"proposal"`
	Timestamp  time.Time `json:"timestamp"`
}
