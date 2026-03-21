package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMetricsCollector_Basic(t *testing.T) {
	collector := NewDefaultMetricsCollector()

	// Record some classifications
	result1 := ClassificationResult{
		Tier:       TierSimple,
		Score:      -0.08,
		Confidence: 0.95,
		Signals: []Signal{
			{Name: "short_message", Contribution: -0.10},
		},
	}
	collector.RecordClassification(result1)

	result2 := ClassificationResult{
		Tier:       TierComplex,
		Score:      0.25,
		Confidence: 0.75,
		Signals: []Signal{
			{Name: "code_block", Contribution: 0.40},
			{Name: "token_medium", Contribution: 0.15},
		},
	}
	collector.RecordClassification(result2)

	result3 := ClassificationResult{
		Tier:       TierReasoning,
		Score:      1.0,
		Confidence: 1.0,
		Signals: []Signal{
			{Name: "attachment_gate", Contribution: 1.0},
		},
	}
	collector.RecordClassification(result3)

	// Get snapshot
	snapshot := collector.GetSnapshot()

	assert.Equal(t, int64(3), snapshot.TotalClassifications)
	assert.Equal(t, int64(1), snapshot.TierCounts[TierSimple])
	assert.Equal(t, int64(1), snapshot.TierCounts[TierComplex])
	assert.Equal(t, int64(1), snapshot.TierCounts[TierReasoning])

	// Check averages
	expectedAvgScore := (-0.08 + 0.25 + 1.0) / 3.0
	assert.InDelta(t, expectedAvgScore, snapshot.AvgScore, 0.001)

	expectedAvgConf := (0.95 + 0.75 + 1.0) / 3.0
	assert.InDelta(t, expectedAvgConf, snapshot.AvgConfidence, 0.001)

	// Check signal counts
	assert.Equal(t, int64(1), snapshot.SignalCounts["short_message"])
	assert.Equal(t, int64(1), snapshot.SignalCounts["code_block"])
	assert.Equal(t, int64(1), snapshot.SignalCounts["attachment_gate"])
}

func TestDefaultMetricsCollector_Reset(t *testing.T) {
	collector := NewDefaultMetricsCollector()

	// Record some classifications
	for i := 0; i < 10; i++ {
		result := ClassificationResult{
			Tier:       TierModerate,
			Score:      0.1,
			Confidence: 0.8,
			Signals:    []Signal{{Name: "test_signal", Contribution: 0.1}},
		}
		collector.RecordClassification(result)
	}

	// Verify data was recorded
	snapshot := collector.GetSnapshot()
	assert.Equal(t, int64(10), snapshot.TotalClassifications)

	// Reset
	collector.Reset()

	// Verify data was cleared
	snapshot = collector.GetSnapshot()
	assert.Equal(t, int64(0), snapshot.TotalClassifications)
	assert.Equal(t, int64(0), snapshot.TierCounts[TierModerate])
	assert.Equal(t, int64(0), snapshot.SignalCounts["test_signal"])
}

func TestDefaultMetricsCollector_Concurrent(t *testing.T) {
	collector := NewDefaultMetricsCollector()

	// Record classifications concurrently
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(idx int) {
			result := ClassificationResult{
				Tier:       AllTiers()[idx%4],
				Score:      float64(idx) / 100.0,
				Confidence: 0.5 + float64(idx%50)/100.0,
				Signals:    []Signal{{Name: "concurrent_test", Contribution: 0.1}},
			}
			collector.RecordClassification(result)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify all classifications were recorded
	snapshot := collector.GetSnapshot()
	assert.Equal(t, int64(100), snapshot.TotalClassifications)
	assert.Equal(t, int64(100), snapshot.SignalCounts["concurrent_test"])
}

func TestNoOpMetricsCollector(t *testing.T) {
	collector := &NoOpMetricsCollector{}

	result := ClassificationResult{
		Tier:       TierSimple,
		Score:      0.5,
		Confidence: 0.9,
	}

	// Should not panic or record anything
	collector.RecordClassification(result)
	snapshot := collector.GetSnapshot()
	// NoOp collector returns a snapshot with current timestamp but empty data
	assert.Equal(t, int64(0), snapshot.TotalClassifications)
	assert.Empty(t, snapshot.TierCounts)

	collector.Reset() // Should not panic
}

func TestGlobalMetrics(t *testing.T) {
	// Save original
	original := GetGlobalMetrics()
	defer SetGlobalMetrics(original)

	// Set new collector
	collector := NewDefaultMetricsCollector()
	SetGlobalMetrics(collector)

	// Verify it's the same
	assert.Equal(t, MetricsCollector(collector), GetGlobalMetrics())

	// Set nil should use no-op
	SetGlobalMetrics(nil)
	assert.IsType(t, &NoOpMetricsCollector{}, GetGlobalMetrics())
}

func TestMetricsCollector_Histograms(t *testing.T) {
	collector := NewDefaultMetricsCollector()

	// Record scores in different buckets
	scores := []float64{-0.8, -0.1, 0.0, 0.1, 0.2, 0.4, 0.6, 0.9}
	for _, score := range scores {
		result := ClassificationResult{
			Tier:       ScoreToTier(score),
			Score:      score,
			Confidence: 0.75,
		}
		collector.RecordClassification(result)
	}

	snapshot := collector.GetSnapshot()

	// Verify histogram has data
	totalInBuckets := int64(0)
	for _, bucket := range snapshot.ScoreDistribution {
		totalInBuckets += bucket.Count
	}
	assert.Equal(t, int64(len(scores)), totalInBuckets)

	// Verify confidence histogram
	totalConfInBuckets := int64(0)
	for _, bucket := range snapshot.ConfidenceDistribution {
		totalConfInBuckets += bucket.Count
	}
	assert.Equal(t, int64(len(scores)), totalConfInBuckets)
}
