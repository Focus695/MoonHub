package routing

import (
	"sync"
	"time"
)

// DecisionRecord represents a single routing decision for visualization.
// Privacy: Does NOT include the user message content.
type DecisionRecord struct {
	// ID is a unique identifier for this decision.
	ID int64 `json:"id"`

	// Timestamp when the decision was made.
	Timestamp time.Time `json:"timestamp"`

	// Tier is the selected complexity tier.
	Tier QueryTier `json:"tier"`

	// Score is the raw complexity score.
	Score float64 `json:"score"`

	// Confidence indicates how certain the classifier was.
	Confidence float64 `json:"confidence"`

	// Signals lists the factors that contributed to this decision.
	// Raw values are omitted for privacy; only signal names and contributions are kept.
	Signals []DecisionSignal `json:"signals"`

	// Model is the selected model name (if available).
	Model string `json:"model,omitempty"`

	// AgentID is the agent identifier (hashed, not raw).
	AgentID string `json:"agent_id,omitempty"`
}

// DecisionSignal is a privacy-safe representation of a classification signal.
type DecisionSignal struct {
	Name         string  `json:"name"`
	Contribution float64 `json:"contribution"`
	// RawValue is intentionally omitted for privacy
}

// DecisionRecorder stores recent routing decisions for visualization.
// It maintains a fixed-size ring buffer and is safe for concurrent use.
type DecisionRecorder struct {
	mu       sync.RWMutex
	records  []DecisionRecord
	size     int
	head     int
	count    int
	nextID   int64
	enabled  bool
	maxDepth int // Maximum signals to record per decision
}

// DecisionRecorderConfig configures the decision recorder.
type DecisionRecorderConfig struct {
	// Enabled controls whether decisions are recorded.
	Enabled bool

	// BufferSize is the maximum number of decisions to keep.
	BufferSize int

	// MaxSignalDepth limits how many signals are recorded per decision.
	MaxSignalDepth int
}

// DefaultDecisionRecorderConfig returns the default configuration.
func DefaultDecisionRecorderConfig() DecisionRecorderConfig {
	return DecisionRecorderConfig{
		Enabled:        true,
		BufferSize:     1000,
		MaxSignalDepth: 10,
	}
}

// NewDecisionRecorder creates a new decision recorder with the given config.
func NewDecisionRecorder(cfg DecisionRecorderConfig) *DecisionRecorder {
	if cfg.BufferSize <= 0 {
		cfg.BufferSize = 1000
	}
	if cfg.MaxSignalDepth <= 0 {
		cfg.MaxSignalDepth = 10
	}

	return &DecisionRecorder{
		records:  make([]DecisionRecord, cfg.BufferSize),
		size:     cfg.BufferSize,
		enabled:  cfg.Enabled,
		maxDepth: cfg.MaxSignalDepth,
	}
}

// Record stores a routing decision.
// Privacy: The message content is never stored; only statistical data is kept.
func (r *DecisionRecorder) Record(result ClassificationResult, model, agentID string) {
	if !r.enabled {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Create privacy-safe decision record
	record := DecisionRecord{
		ID:         r.nextID,
		Timestamp:  time.Now(),
		Tier:       result.Tier,
		Score:      result.Score,
		Confidence: result.Confidence,
		Model:      model,
		AgentID:    agentID,
		Signals:    r.sanitizeSignals(result.Signals),
	}
	r.nextID++

	// Add to ring buffer
	r.records[r.head] = record
	r.head = (r.head + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

// sanitizeSignals converts signals to privacy-safe format.
// Limits the number of signals and strips raw values.
func (r *DecisionRecorder) sanitizeSignals(signals []Signal) []DecisionSignal {
	n := len(signals)
	if n > r.maxDepth {
		n = r.maxDepth
	}

	result := make([]DecisionSignal, 0, n)
	for i := 0; i < n && i < len(signals); i++ {
		result = append(result, DecisionSignal{
			Name:         signals[i].Name,
			Contribution: signals[i].Contribution,
			// RawValue intentionally omitted
		})
	}
	return result
}

// GetRecent returns the most recent n decisions.
// Returns all available if n > total count.
func (r *DecisionRecorder) GetRecent(n int) []DecisionRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if n <= 0 || r.count == 0 {
		return []DecisionRecord{}
	}

	if n > r.count {
		n = r.count
	}

	result := make([]DecisionRecord, n)

	// Calculate starting position in ring buffer
	// Most recent is at (head - 1), oldest is at head
	start := r.head - n
	if start < 0 {
		start += r.size
	}

	for i := 0; i < n; i++ {
		idx := (start + i) % r.size
		result[i] = r.records[idx]
	}

	// Reverse to get most recent first
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// GetByTier returns recent decisions filtered by tier.
func (r *DecisionRecorder) GetByTier(tier QueryTier, limit int) []DecisionRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.count == 0 || limit <= 0 {
		return []DecisionRecord{}
	}

	var result []DecisionRecord
	count := 0

	// Iterate from most recent to oldest
	for i := 0; i < r.count && count < limit; i++ {
		idx := r.head - 1 - i
		if idx < 0 {
			idx += r.size
		}
		if r.records[idx].Tier == tier {
			result = append(result, r.records[idx])
			count++
		}
	}

	return result
}

// GetStats returns aggregated statistics about recent decisions.
func (r *DecisionRecorder) GetStats() DecisionStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := DecisionStats{
		TotalDecisions:  r.count,
		TierCounts:      make(map[QueryTier]int),
		SignalFrequency: make(map[string]int),
	}

	if r.count == 0 {
		return stats
	}

	var scoreSum, confSum float64

	for i := 0; i < r.count; i++ {
		record := r.records[i]
		stats.TierCounts[record.Tier]++
		scoreSum += record.Score
		confSum += record.Confidence

		// Track signal frequency
		for _, s := range record.Signals {
			stats.SignalFrequency[s.Name]++
		}
	}

	stats.AvgScore = scoreSum / float64(r.count)
	stats.AvgConfidence = confSum / float64(r.count)

	return stats
}

// DecisionStats contains aggregated statistics about routing decisions.
type DecisionStats struct {
	TotalDecisions  int               `json:"total_decisions"`
	TierCounts      map[QueryTier]int `json:"tier_counts"`
	AvgScore        float64           `json:"avg_score"`
	AvgConfidence   float64           `json:"avg_confidence"`
	SignalFrequency map[string]int    `json:"signal_frequency"`
}

// Clear removes all recorded decisions.
func (r *DecisionRecorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = make([]DecisionRecord, r.size)
	r.head = 0
	r.count = 0
}

// SetEnabled enables or disables recording.
func (r *DecisionRecorder) SetEnabled(enabled bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.enabled = enabled
}

// IsEnabled returns whether recording is enabled.
func (r *DecisionRecorder) IsEnabled() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.enabled
}

// Count returns the number of recorded decisions.
func (r *DecisionRecorder) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

// Global decision recorder instance
var globalRecorder *DecisionRecorder

// InitGlobalRecorder initializes the global decision recorder.
func InitGlobalRecorder(cfg DecisionRecorderConfig) {
	globalRecorder = NewDecisionRecorder(cfg)
}

// GetGlobalRecorder returns the global decision recorder.
func GetGlobalRecorder() *DecisionRecorder {
	if globalRecorder == nil {
		InitGlobalRecorder(DefaultDecisionRecorderConfig())
	}
	return globalRecorder
}

// RecordDecision records to the global recorder.
func RecordDecision(result ClassificationResult, model, agentID string) {
	if globalRecorder != nil {
		globalRecorder.Record(result, model, agentID)
	}
}
