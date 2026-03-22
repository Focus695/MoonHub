// Package routing provides intelligent model routing based on message complexity.
// This file implements optional metrics collection for observability.
package routing

import (
	"sync"
	"time"
)

// MetricsCollector defines the interface for collecting routing metrics.
// Implementations must be safe for concurrent use.
type MetricsCollector interface {
	// RecordClassification records a single classification decision.
	// Does NOT record the message content - only aggregated statistics.
	RecordClassification(result ClassificationResult)

	// GetSnapshot returns a point-in-time snapshot of all metrics.
	GetSnapshot() MetricsSnapshot

	// Reset clears all collected metrics.
	Reset()
}

// MetricsSnapshot contains a point-in-time view of routing metrics.
// All fields are safe to read without locking.
type MetricsSnapshot struct {
	// Timestamp when this snapshot was taken.
	Timestamp time.Time

	// TotalClassifications is the total number of routing decisions made.
	TotalClassifications int64

	// TierCounts shows how many messages were routed to each tier.
	TierCounts map[QueryTier]int64

	// ScoreDistribution contains histogram buckets for score ranges.
	ScoreDistribution []ScoreBucket

	// ConfidenceDistribution contains histogram buckets for confidence ranges.
	ConfidenceDistribution []ConfidenceBucket

	// SignalCounts shows how often each signal type was triggered.
	SignalCounts map[string]int64

	// AvgScore is the running average of all classification scores.
	AvgScore float64

	// AvgConfidence is the running average of all confidence values.
	AvgConfidence float64
}

// ScoreBucket represents a histogram bucket for score distribution.
type ScoreBucket struct {
	Min   float64
	Max   float64
	Count int64
}

// ConfidenceBucket represents a histogram bucket for confidence distribution.
type ConfidenceBucket struct {
	Min   float64
	Max   float64
	Count int64
}

// DefaultMetricsCollector is the default implementation of MetricsCollector.
// It collects aggregated statistics without storing any user content.
type DefaultMetricsCollector struct {
	mu sync.RWMutex

	totalClassifications int64
	tierCounts           map[QueryTier]int64
	scoreSum             float64
	confidenceSum        float64
	signalCounts         map[string]int64

	// Histogram buckets for score distribution
	scoreBuckets []scoreBucketInternal

	// Histogram buckets for confidence distribution
	confidenceBuckets []confidenceBucketInternal
}

type scoreBucketInternal struct {
	min   float64
	max   float64
	count int64
}

type confidenceBucketInternal struct {
	min   float64
	max   float64
	count int64
}

// NewDefaultMetricsCollector creates a new metrics collector with default buckets.
func NewDefaultMetricsCollector() *DefaultMetricsCollector {
	c := &DefaultMetricsCollector{
		tierCounts:   make(map[QueryTier]int64),
		signalCounts: make(map[string]int64),
		scoreBuckets: []scoreBucketInternal{
			{min: -1.0, max: -0.5},
			{min: -0.5, max: -0.05},
			{min: -0.05, max: 0.15},
			{min: 0.15, max: 0.35},
			{min: 0.35, max: 0.5},
			{min: 0.5, max: 0.75},
			{min: 0.75, max: 1.0},
		},
		confidenceBuckets: []confidenceBucketInternal{
			{min: 0.0, max: 0.5},
			{min: 0.5, max: 0.6},
			{min: 0.6, max: 0.7},
			{min: 0.7, max: 0.8},
			{min: 0.8, max: 0.9},
			{min: 0.9, max: 1.0},
		},
	}

	// Initialize tier counters
	for _, tier := range AllTiers() {
		c.tierCounts[tier] = 0
	}

	return c
}

// RecordClassification records a single classification decision.
// Privacy: Only records aggregated statistics, never the message content.
func (c *DefaultMetricsCollector) RecordClassification(result ClassificationResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalClassifications++

	// Update tier count
	c.tierCounts[result.Tier]++

	// Update running averages
	c.scoreSum += result.Score
	c.confidenceSum += result.Confidence

	// Update score histogram
	for i := range c.scoreBuckets {
		b := &c.scoreBuckets[i]
		if result.Score >= b.min && result.Score < b.max {
			b.count++
			break
		}
		// Handle edge case for max bucket
		if i == len(c.scoreBuckets)-1 && result.Score >= b.min && result.Score <= b.max {
			b.count++
		}
	}

	// Update confidence histogram
	for i := range c.confidenceBuckets {
		b := &c.confidenceBuckets[i]
		if result.Confidence >= b.min && result.Confidence < b.max {
			b.count++
			break
		}
		// Handle edge case for max bucket
		if i == len(c.confidenceBuckets)-1 && result.Confidence >= b.min && result.Confidence <= b.max {
			b.count++
		}
	}

	// Update signal counts (only signal names, not raw values)
	for _, signal := range result.Signals {
		c.signalCounts[signal.Name]++
	}
}

// GetSnapshot returns a point-in-time snapshot of all metrics.
func (c *DefaultMetricsCollector) GetSnapshot() MetricsSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := MetricsSnapshot{
		Timestamp:              time.Now(),
		TierCounts:             make(map[QueryTier]int64),
		SignalCounts:           make(map[string]int64),
		ScoreDistribution:      make([]ScoreBucket, len(c.scoreBuckets)),
		ConfidenceDistribution: make([]ConfidenceBucket, len(c.confidenceBuckets)),
		TotalClassifications:   c.totalClassifications,
	}

	if c.totalClassifications > 0 {
		snapshot.AvgScore = c.scoreSum / float64(c.totalClassifications)
		snapshot.AvgConfidence = c.confidenceSum / float64(c.totalClassifications)
	}

	// Copy tier counts
	for k, v := range c.tierCounts {
		snapshot.TierCounts[k] = v
	}

	// Copy signal counts
	for k, v := range c.signalCounts {
		snapshot.SignalCounts[k] = v
	}

	// Copy score distribution
	for i, b := range c.scoreBuckets {
		snapshot.ScoreDistribution[i] = ScoreBucket{
			Min:   b.min,
			Max:   b.max,
			Count: b.count,
		}
	}

	// Copy confidence distribution
	for i, b := range c.confidenceBuckets {
		snapshot.ConfidenceDistribution[i] = ConfidenceBucket{
			Min:   b.min,
			Max:   b.max,
			Count: b.count,
		}
	}

	return snapshot
}

// Reset clears all collected metrics.
func (c *DefaultMetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalClassifications = 0
	c.scoreSum = 0
	c.confidenceSum = 0

	// Reset tier counts
	for k := range c.tierCounts {
		c.tierCounts[k] = 0
	}

	// Reset signal counts
	for k := range c.signalCounts {
		c.signalCounts[k] = 0
	}

	// Reset score buckets
	for i := range c.scoreBuckets {
		c.scoreBuckets[i].count = 0
	}

	// Reset confidence buckets
	for i := range c.confidenceBuckets {
		c.confidenceBuckets[i].count = 0
	}
}

// NoOpMetricsCollector is a no-op implementation that discards all metrics.
// Use this when metrics collection is disabled.
type NoOpMetricsCollector struct{}

// RecordClassification does nothing.
func (n *NoOpMetricsCollector) RecordClassification(result ClassificationResult) {}

// GetSnapshot returns an empty snapshot.
func (n *NoOpMetricsCollector) GetSnapshot() MetricsSnapshot {
	return MetricsSnapshot{
		Timestamp:            time.Now(),
		TierCounts:           make(map[QueryTier]int64),
		SignalCounts:         make(map[string]int64),
		TotalClassifications: 0,
	}
}

// Reset does nothing.
func (n *NoOpMetricsCollector) Reset() {}

// Global metrics collector instance
var globalMetrics MetricsCollector = &NoOpMetricsCollector{}

// SetGlobalMetrics sets the global metrics collector.
func SetGlobalMetrics(collector MetricsCollector) {
	if collector == nil {
		collector = &NoOpMetricsCollector{}
	}
	globalMetrics = collector
}

// GetGlobalMetrics returns the global metrics collector.
func GetGlobalMetrics() MetricsCollector {
	return globalMetrics
}

// RecordClassification records to the global metrics collector.
func RecordClassification(result ClassificationResult) {
	globalMetrics.RecordClassification(result)
}
