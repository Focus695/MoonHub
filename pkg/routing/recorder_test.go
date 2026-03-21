package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecisionRecorder_Basic(t *testing.T) {
	cfg := DefaultDecisionRecorderConfig()
	cfg.BufferSize = 10
	recorder := NewDecisionRecorder(cfg)

	// Record some decisions
	for i := 0; i < 5; i++ {
		result := ClassificationResult{
			Tier:       AllTiers()[i%4],
			Score:      float64(i) / 10.0,
			Confidence: 0.5 + float64(i)/10.0,
			Signals: []Signal{
				{Name: "test_signal", Contribution: 0.1, RawValue: "sensitive_data"},
			},
		}
		recorder.Record(result, "test-model", "agent-123")
	}

	// Get recent decisions
	decisions := recorder.GetRecent(10)
	assert.Len(t, decisions, 5)

	// Verify most recent is first (ID 4)
	assert.Equal(t, int64(4), decisions[0].ID)
	// ID 4 corresponds to i=4, so tier is AllTiers()[4%4] = AllTiers()[0] = TierSimple
	assert.Equal(t, TierSimple, decisions[0].Tier)

	// Verify privacy: raw values should NOT be present
	for _, d := range decisions {
		for _, s := range d.Signals {
			// DecisionSignal should only have Name and Contribution
			assert.NotEmpty(t, s.Name)
			assert.NotEmpty(t, s.Contribution)
		}
	}
}

func TestDecisionRecorder_RingBuffer(t *testing.T) {
	cfg := DecisionRecorderConfig{
		Enabled:        true,
		BufferSize:     5,
		MaxSignalDepth: 10,
	}
	recorder := NewDecisionRecorder(cfg)

	// Record more than buffer size
	for i := 0; i < 10; i++ {
		result := ClassificationResult{
			Tier:       TierSimple,
			Score:      float64(i) / 10.0,
			Confidence: 0.9,
		}
		recorder.Record(result, "model", "agent")
	}

	// Should only have last 5
	decisions := recorder.GetRecent(10)
	assert.Len(t, decisions, 5)

	// Most recent should be ID 9
	assert.Equal(t, int64(9), decisions[0].ID)

	// Oldest should be ID 5
	assert.Equal(t, int64(5), decisions[4].ID)
}

func TestDecisionRecorder_GetByTier(t *testing.T) {
	cfg := DefaultDecisionRecorderConfig()
	recorder := NewDecisionRecorder(cfg)

	// Record decisions for different tiers
	tiers := []QueryTier{TierSimple, TierModerate, TierComplex, TierReasoning, TierSimple}
	for i, tier := range tiers {
		result := ClassificationResult{
			Tier:       tier,
			Score:      float64(i) / 10.0,
			Confidence: 0.8,
		}
		recorder.Record(result, "model", "agent")
	}

	// Get only simple tier
	simpleDecisions := recorder.GetByTier(TierSimple, 10)
	assert.Len(t, simpleDecisions, 2)

	// All should be simple tier
	for _, d := range simpleDecisions {
		assert.Equal(t, TierSimple, d.Tier)
	}

	// Get complex tier
	complexDecisions := recorder.GetByTier(TierComplex, 10)
	assert.Len(t, complexDecisions, 1)
	assert.Equal(t, TierComplex, complexDecisions[0].Tier)
}

func TestDecisionRecorder_GetStats(t *testing.T) {
	cfg := DefaultDecisionRecorderConfig()
	recorder := NewDecisionRecorder(cfg)

	// Record some decisions
	for i := 0; i < 20; i++ {
		result := ClassificationResult{
			Tier:       AllTiers()[i%4],
			Score:      float64(i) / 20.0,
			Confidence: 0.5 + float64(i%10)/20.0,
			Signals: []Signal{
				{Name: "signal_a", Contribution: 0.1},
				{Name: "signal_b", Contribution: 0.2},
			},
		}
		recorder.Record(result, "model", "agent")
	}

	stats := recorder.GetStats()

	assert.Equal(t, 20, stats.TotalDecisions)
	assert.Equal(t, 5, stats.TierCounts[TierSimple])
	assert.Equal(t, 5, stats.TierCounts[TierModerate])
	assert.Equal(t, 5, stats.TierCounts[TierComplex])
	assert.Equal(t, 5, stats.TierCounts[TierReasoning])

	// Check signal frequency
	assert.Equal(t, 20, stats.SignalFrequency["signal_a"])
	assert.Equal(t, 20, stats.SignalFrequency["signal_b"])
}

func TestDecisionRecorder_Disabled(t *testing.T) {
	cfg := DecisionRecorderConfig{
		Enabled:    false,
		BufferSize: 10,
	}
	recorder := NewDecisionRecorder(cfg)

	// Record a decision
	result := ClassificationResult{
		Tier:       TierSimple,
		Score:      0.1,
		Confidence: 0.9,
	}
	recorder.Record(result, "model", "agent")

	// Should not have recorded anything
	assert.Equal(t, 0, recorder.Count())
	decisions := recorder.GetRecent(10)
	assert.Empty(t, decisions)

	// Enable and record again
	recorder.SetEnabled(true)
	recorder.Record(result, "model", "agent")
	assert.Equal(t, 1, recorder.Count())
}

func TestDecisionRecorder_Clear(t *testing.T) {
	cfg := DefaultDecisionRecorderConfig()
	recorder := NewDecisionRecorder(cfg)

	// Record some decisions
	for i := 0; i < 5; i++ {
		result := ClassificationResult{
			Tier:       TierSimple,
			Score:      0.1,
			Confidence: 0.9,
		}
		recorder.Record(result, "model", "agent")
	}

	assert.Equal(t, 5, recorder.Count())

	// Clear
	recorder.Clear()
	assert.Equal(t, 0, recorder.Count())

	decisions := recorder.GetRecent(10)
	assert.Empty(t, decisions)
}

func TestDecisionRecorder_MaxSignalDepth(t *testing.T) {
	cfg := DecisionRecorderConfig{
		Enabled:        true,
		BufferSize:     10,
		MaxSignalDepth: 3,
	}
	recorder := NewDecisionRecorder(cfg)

	// Record with many signals
	signals := make([]Signal, 10)
	for i := range signals {
		signals[i] = Signal{
			Name:         "signal",
			Contribution: float64(i) / 10.0,
			RawValue:     "sensitive",
		}
	}
	result := ClassificationResult{
		Tier:       TierComplex,
		Score:      0.5,
		Confidence: 0.8,
		Signals:    signals,
	}
	recorder.Record(result, "model", "agent")

	decisions := recorder.GetRecent(1)
	assert.Len(t, decisions, 1)

	// Should only have MaxSignalDepth signals
	assert.Len(t, decisions[0].Signals, 3)
}

func TestDecisionRecorder_Privacy(t *testing.T) {
	cfg := DefaultDecisionRecorderConfig()
	recorder := NewDecisionRecorder(cfg)

	// Record with sensitive raw values
	result := ClassificationResult{
		Tier:       TierComplex,
		Score:      0.5,
		Confidence: 0.8,
		Signals: []Signal{
			{Name: "user_input", Contribution: 0.1, RawValue: "user's secret message"},
			{Name: "code_block", Contribution: 0.4, RawValue: "sensitive code"},
		},
	}
	recorder.Record(result, "model", "agent")

	decisions := recorder.GetRecent(1)
	assert.Len(t, decisions, 1)

	// Verify raw values are NOT in the decision record
	for _, s := range decisions[0].Signals {
		// DecisionSignal struct only has Name and Contribution
		// RawValue should not be accessible
		assert.NotEmpty(t, s.Name)
		assert.NotEmpty(t, s.Contribution)
	}
}

func TestGlobalRecorder(t *testing.T) {
	// Initialize global recorder
	InitGlobalRecorder(DefaultDecisionRecorderConfig())

	recorder := GetGlobalRecorder()
	assert.NotNil(t, recorder)
	assert.True(t, recorder.IsEnabled())

	// Clear any existing data
	recorder.Clear()

	// Record using global function
	result := ClassificationResult{
		Tier:       TierSimple,
		Score:      0.1,
		Confidence: 0.9,
	}
	RecordDecision(result, "test-model", "test-agent")

	assert.Equal(t, 1, recorder.Count())
}
