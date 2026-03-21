package routing

import (
	"math"
	"testing"
)

func TestRuleClassifierV2_ShortMessage(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Very short message (<20 tokens) should get negative score
	f := Features{
		TokenEstimate:     10,
		CodeBlockCount:    0,
		RecentToolCalls:   0,
		ConversationDepth: 0,
		HasAttachments:    false,
	}

	result := classifier.Classify(f)

	if result.Score >= 0 {
		t.Errorf("expected negative score for short message, got %f", result.Score)
	}
	if result.Tier != TierSimple {
		t.Errorf("expected TierSimple for short message, got %s", result.Tier)
	}
	if result.Confidence < 0.5 || result.Confidence > 1.0 {
		t.Errorf("expected confidence in [0.5, 1.0], got %f", result.Confidence)
	}

	// Check signals
	foundShortMsg := false
	for _, s := range result.Signals {
		if s.Name == "short_message" {
			foundShortMsg = true
			if s.Contribution != ShortMessageBonus {
				t.Errorf("expected contribution %f, got %f", ShortMessageBonus, s.Contribution)
			}
		}
	}
	if !foundShortMsg {
		t.Error("expected short_message signal not found")
	}
}

func TestRuleClassifierV2_AttachmentGate(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Any message with attachments should go to reasoning tier
	f := Features{
		TokenEstimate:     5,
		CodeBlockCount:    0,
		RecentToolCalls:   0,
		ConversationDepth: 0,
		HasAttachments:    true,
	}

	result := classifier.Classify(f)

	if result.Tier != TierReasoning {
		t.Errorf("expected TierReasoning for attachments, got %s", result.Tier)
	}
	if result.Score != 1.0 {
		t.Errorf("expected score 1.0 for attachments, got %f", result.Score)
	}
	if result.Confidence != 1.0 {
		t.Errorf("expected confidence 1.0 for attachments, got %f", result.Confidence)
	}
}

func TestRuleClassifierV2_CodeBlock(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Message with code block should go to complex or reasoning
	f := Features{
		TokenEstimate:     100,
		CodeBlockCount:    1,
		RecentToolCalls:   0,
		ConversationDepth: 0,
		HasAttachments:    false,
	}

	result := classifier.Classify(f)

	if result.Tier != TierComplex && result.Tier != TierReasoning {
		t.Errorf("expected TierComplex or TierReasoning for code block, got %s (score=%f)",
			result.Tier, result.Score)
	}

	// Check signals
	foundCodeBlock := false
	for _, s := range result.Signals {
		if s.Name == "code_block" {
			foundCodeBlock = true
		}
	}
	if !foundCodeBlock {
		t.Error("expected code_block signal not found")
	}
}

func TestRuleClassifierV2_LongMessage(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Long message (>200 tokens) should go to complex or reasoning
	f := Features{
		TokenEstimate:     300,
		CodeBlockCount:    0,
		RecentToolCalls:   0,
		ConversationDepth: 0,
		HasAttachments:    false,
	}

	result := classifier.Classify(f)

	if result.Tier != TierComplex && result.Tier != TierReasoning {
		t.Errorf("expected TierComplex or TierReasoning for long message, got %s (score=%f)",
			result.Tier, result.Score)
	}
}

func TestRuleClassifierV2_ToolCalls(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	tests := []struct {
		name          string
		toolCalls     int
		expectedTier  QueryTier
		minScore      float64
	}{
		{
			name:         "no_tool_calls",
			toolCalls:    0,
			expectedTier: TierSimple, // with short message
			minScore:     -0.15,
		},
		{
			name:          "few_tool_calls",
			toolCalls:     2,
			expectedTier:  TierModerate, // 0.10 from tool calls
			minScore:      0.0,
		},
		{
			name:          "many_tool_calls",
			toolCalls:     5,
			expectedTier:  TierComplex, // 0.25 from tool calls
			minScore:      0.15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Features{
				TokenEstimate:     30,
				CodeBlockCount:    0,
				RecentToolCalls:   tt.toolCalls,
				ConversationDepth: 0,
				HasAttachments:    false,
			}

			result := classifier.Classify(f)

			if result.Score < tt.minScore {
				t.Errorf("expected score >= %f, got %f", tt.minScore, result.Score)
			}
		})
	}
}

func TestRuleClassifierV2_ConversationDepth(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Deep conversation should add to score
	f := Features{
		TokenEstimate:     50,
		CodeBlockCount:    0,
		RecentToolCalls:   0,
		ConversationDepth: 15,
		HasAttachments:    false,
	}

	result := classifier.Classify(f)

	// 0.15 (medium) + 0.10 (depth) = 0.25 → moderate
	if result.Score < 0.10 {
		t.Errorf("expected score >= 0.10 for deep conversation, got %f", result.Score)
	}

	// Check signals
	foundDepth := false
	for _, s := range result.Signals {
		if s.Name == "conversation_depth" {
			foundDepth = true
		}
	}
	if !foundDepth {
		t.Error("expected conversation_depth signal not found")
	}
}

func TestRuleClassifierV2_CustomBoundaries(t *testing.T) {
	customBoundaries := map[QueryTier]TierBoundary{
		TierSimple:    {Min: -1.0, Max: 0.0},
		TierModerate:  {Min: 0.0, Max: 0.3},
		TierComplex:   {Min: 0.3, Max: 0.6},
		TierReasoning: {Min: 0.6, Max: 1.0},
	}
	classifier := NewRuleClassifierV2(customBoundaries)

	tests := []struct {
		name         string
		features     Features
		expectedTier QueryTier
	}{
		{
			name: "short_message_custom_boundary",
			features: Features{
				TokenEstimate:  10,
				HasAttachments: false,
			},
			expectedTier: TierSimple,
		},
		{
			name: "code_block_custom_boundary",
			features: Features{
				TokenEstimate:  100,
				CodeBlockCount: 1,
				HasAttachments: false,
			},
			expectedTier: TierComplex, // 0.40 + 0.15 = 0.55 → complex with custom boundary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.Classify(tt.features)
			if result.Tier != tt.expectedTier {
				t.Errorf("expected tier %s, got %s (score=%f)",
					tt.expectedTier, result.Tier, result.Score)
			}
		})
	}
}

func TestRuleClassifierV2_ConfidenceCalculation(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Test that confidence is calculated correctly
	// High confidence when score is clearly within a tier

	// Very short message should have high confidence for simple tier
	f1 := Features{TokenEstimate: 5}
	result1 := classifier.Classify(f1)
	if result1.Confidence < 0.7 {
		t.Errorf("expected high confidence for clearly simple message, got %f", result1.Confidence)
	}

	// Attachment gate should have perfect confidence
	f2 := Features{HasAttachments: true}
	result2 := classifier.Classify(f2)
	if result2.Confidence != 1.0 {
		t.Errorf("expected confidence 1.0 for attachment gate, got %f", result2.Confidence)
	}
}

func TestRuleClassifierV2_SignalTracing(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	f := Features{
		TokenEstimate:     300,
		CodeBlockCount:    2,
		RecentToolCalls:   5,
		ConversationDepth: 15,
		HasAttachments:    false,
	}

	result := classifier.Classify(f)

	// Should have signals for all contributing factors
	expectedSignals := map[string]bool{
		"token_long":         false,
		"code_block":         false,
		"tool_call_high":     false,
		"conversation_depth": false,
	}

	for _, s := range result.Signals {
		if _, ok := expectedSignals[s.Name]; ok {
			expectedSignals[s.Name] = true
		}
	}

	for signal, found := range expectedSignals {
		if !found {
			t.Errorf("expected signal %s not found", signal)
		}
	}
}

func TestScoreToTier(t *testing.T) {
	tests := []struct {
		score        float64
		expectedTier QueryTier
	}{
		{-0.5, TierSimple},
		{-0.1, TierSimple},
		{-0.05, TierModerate}, // boundary
		{0.0, TierModerate},
		{0.1, TierModerate},
		{0.15, TierComplex}, // boundary
		{0.2, TierComplex},
		{0.3, TierComplex},
		{0.35, TierReasoning}, // boundary
		{0.5, TierReasoning},
		{1.0, TierReasoning},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tier := ScoreToTier(tt.score)
			if tier != tt.expectedTier {
				t.Errorf("score %f: expected tier %s, got %s", tt.score, tt.expectedTier, tier)
			}
		})
	}
}

func TestComputeConfidence(t *testing.T) {
	// Test confidence at various points
	boundaries := DefaultTierBoundaries

	// At the midpoint of a tier, confidence should be high
	confidence := computeConfidence(0.0, TierModerate, boundaries) // midpoint of moderate
	if confidence < 0.7 {
		t.Errorf("expected high confidence at midpoint, got %f", confidence)
	}

	// Confidence should be in [0.5, 1.0] range
	confidence = computeConfidence(-0.06, TierSimple, boundaries)
	if confidence < 0.5 || confidence > 1.0 {
		t.Errorf("confidence %f out of range [0.5, 1.0]", confidence)
	}
}

func TestQueryTier_StringAndParse(t *testing.T) {
	tests := []struct {
		s    string
		tier QueryTier
	}{
		{"simple", TierSimple},
		{"moderate", TierModerate},
		{"complex", TierComplex},
		{"reasoning", TierReasoning},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			parsed, ok := TryParseTier(tt.s)
			if !ok || parsed != tt.tier {
				t.Errorf("TryParseTier(%q): expected %s, got %s ok=%v", tt.s, tt.tier, parsed, ok)
			}
			legacy := ParseTier(tt.s)
			if legacy != tt.tier {
				t.Errorf("ParseTier(%q): expected %s, got %s", tt.s, tt.tier, legacy)
			}
			if parsed.String() != tt.s {
				t.Errorf("String(): expected %q, got %q", tt.s, parsed.String())
			}
		})
	}

	t.Run("unknown_ParseTier_defaults_moderate", func(t *testing.T) {
		if _, ok := TryParseTier("unknown"); ok {
			t.Fatal("TryParseTier should reject unknown")
		}
		if ParseTier("unknown") != TierModerate {
			t.Error("ParseTier(unknown) should default to moderate for legacy callers")
		}
	})
}

func TestQueryTier_IsValid(t *testing.T) {
	if !TierSimple.IsValid() || !TierModerate.IsValid() || !TierComplex.IsValid() || !TierReasoning.IsValid() {
		t.Error("valid tiers should return true")
	}

	invalidTier := QueryTier("invalid")
	if invalidTier.IsValid() {
		t.Error("invalid tier should return false")
	}
}

func BenchmarkRuleClassifierV2_Classify(b *testing.B) {
	classifier := NewRuleClassifierV2(nil)
	f := Features{
		TokenEstimate:     100,
		CodeBlockCount:    1,
		RecentToolCalls:   2,
		ConversationDepth: 5,
		HasAttachments:    false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = classifier.Classify(f)
	}
}

func TestRuleClassifierV2_ScoreRange(t *testing.T) {
	classifier := NewRuleClassifierV2(nil)

	// Test that scores stay within expected bounds
	tests := []Features{
		{TokenEstimate: 5, CodeBlockCount: 0, RecentToolCalls: 0, ConversationDepth: 0, HasAttachments: false},
		{TokenEstimate: 100, CodeBlockCount: 0, RecentToolCalls: 0, ConversationDepth: 0, HasAttachments: false},
		{TokenEstimate: 500, CodeBlockCount: 3, RecentToolCalls: 10, ConversationDepth: 20, HasAttachments: false},
		{TokenEstimate: 50, CodeBlockCount: 1, RecentToolCalls: 1, ConversationDepth: 5, HasAttachments: false},
	}

	for _, f := range tests {
		result := classifier.Classify(f)
		if result.Score < -1.0 || result.Score > 1.0 {
			t.Errorf("score %f out of expected range [-1, 1]", result.Score)
		}
		if math.IsNaN(result.Score) || math.IsInf(result.Score, 0) {
			t.Errorf("score is invalid: %f", result.Score)
		}
	}
}
