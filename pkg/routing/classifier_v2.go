package routing

import (
	"math"
)

// V2 classifier weights — designed to produce scores in roughly [-0.10, 1.0] range.
// Negative scores indicate very simple messages (greetings, acknowledgments).
const (
	// ShortMessageBonus is applied when the message has fewer than 20 estimated tokens.
	// This produces a negative score contribution, routing simple messages to cheaper models.
	ShortMessageBonus = -0.10

	// TokenMediumWeight is applied for messages with 50-200 tokens (medium length).
	TokenMediumWeight = 0.15

	// TokenLongWeight is applied for messages with >200 tokens (long messages).
	TokenLongWeight = 0.35

	// TokenShortMidWeight applies to 20–49 estimated tokens (between "short" and medium-length band).
	TokenShortMidWeight = 0.06

	// CodeBlockWeight is applied when the message contains fenced code blocks.
	// Coding tasks almost always require capable models.
	CodeBlockWeight = 0.40

	// ToolCallHighWeight is applied when there are >3 recent tool calls (dense agentic workflow).
	ToolCallHighWeight = 0.25

	// ToolCallLowWeight is applied when there are 1-3 recent tool calls.
	ToolCallLowWeight = 0.10

	// ConversationDepthWeight is applied when conversation depth >10 messages.
	ConversationDepthWeight = 0.10

	// AttachmentGateScore is the score for messages with attachments (hard gate to reasoning tier).
	AttachmentGateScore = 1.0
)

// RuleClassifierV2 is the V2 implementation of the Classifier interface.
// It produces a score that can be negative (for very simple messages) and maps
// the score to one of four tiers: simple, moderate, complex, reasoning.
//
// Score calculation:
//   - Short message (<20 tokens): -0.10
//   - Short–mid (20–50 tokens): +0.06
//   - Medium message (51-200 tokens): +0.15
//   - Long message (>200 tokens): +0.35
//   - Code block present: +0.40
//   - Tool calls >3: +0.25
//   - Tool calls 1-3: +0.10
//   - Conversation depth >10: +0.10
//   - Attachments present: 1.0 (hard gate)
//
// Tier boundaries (default):
//   - Simple: [-1.0, -0.05)
//   - Moderate: [-0.05, 0.15)
//   - Complex: [0.15, 0.35)
//   - Reasoning: [0.35, 1.0]
type RuleClassifierV2 struct {
	boundaries map[QueryTier]TierBoundary
}

// NewRuleClassifierV2 creates a new V2 classifier with optional custom boundaries.
func NewRuleClassifierV2(boundaries map[QueryTier]TierBoundary) *RuleClassifierV2 {
	if boundaries == nil {
		boundaries = DefaultTierBoundaries
	}
	return &RuleClassifierV2{boundaries: boundaries}
}

// Classify evaluates the features and returns a detailed classification result.
func (c *RuleClassifierV2) Classify(f Features) ClassificationResult {
	signals := make([]Signal, 0, 6)

	// 1. Hard gate: multi-modal inputs always require the reasoning tier.
	if f.HasAttachments {
		signals = append(signals, Signal{
			Name:         "attachment_gate",
			Contribution: AttachmentGateScore,
			RawValue:     true,
		})
		return ClassificationResult{
			Tier:       TierReasoning,
			Score:      AttachmentGateScore,
			Confidence: 1.0,
			Signals:    signals,
		}
	}

	var score float64

	// 2. Token estimate — primary verbosity signal
	switch {
	case f.TokenEstimate < 20:
		// Short messages get a negative bonus to route to simple tier
		score += ShortMessageBonus
		signals = append(signals, Signal{
			Name:         "short_message",
			Contribution: ShortMessageBonus,
			RawValue:     f.TokenEstimate,
		})
	case f.TokenEstimate > 200:
		score += TokenLongWeight
		signals = append(signals, Signal{
			Name:         "token_long",
			Contribution: TokenLongWeight,
			RawValue:     f.TokenEstimate,
		})
	case f.TokenEstimate > 50:
		// Medium band: 51–200 tokens (>50 keeps token count 50 on the lighter short_mid band)
		score += TokenMediumWeight
		signals = append(signals, Signal{
			Name:         "token_medium",
			Contribution: TokenMediumWeight,
			RawValue:     f.TokenEstimate,
		})
	default:
		// 20 <= tokens < 50
		score += TokenShortMidWeight
		signals = append(signals, Signal{
			Name:         "token_short_mid",
			Contribution: TokenShortMidWeight,
			RawValue:     f.TokenEstimate,
		})
	}

	// 3. Fenced code blocks — strongest indicator of a coding/technical task
	if f.CodeBlockCount > 0 {
		score += CodeBlockWeight
		signals = append(signals, Signal{
			Name:         "code_block",
			Contribution: CodeBlockWeight,
			RawValue:     f.CodeBlockCount,
		})
	}

	// 4. Recent tool call density — indicates an ongoing agentic workflow
	switch {
	case f.RecentToolCalls > 3:
		score += ToolCallHighWeight
		signals = append(signals, Signal{
			Name:         "tool_call_high",
			Contribution: ToolCallHighWeight,
			RawValue:     f.RecentToolCalls,
		})
	case f.RecentToolCalls > 0:
		score += ToolCallLowWeight
		signals = append(signals, Signal{
			Name:         "tool_call_low",
			Contribution: ToolCallLowWeight,
			RawValue:     f.RecentToolCalls,
		})
	}

	// 5. Conversation depth — accumulated context implies compound task
	if f.ConversationDepth > 10 {
		score += ConversationDepthWeight
		signals = append(signals, Signal{
			Name:         "conversation_depth",
			Contribution: ConversationDepthWeight,
			RawValue:     f.ConversationDepth,
		})
	}

	// 6. Cap score to [-1.0, 1.0] range
	if score > 1.0 {
		score = 1.0
	}
	if score < -1.0 {
		score = -1.0
	}

	// 7. Map score to tier
	tier := ScoreToTierWithBoundaries(score, c.boundaries)

	// 7. Compute confidence using sigmoid function
	confidence := computeConfidence(score, tier, c.boundaries)

	return ClassificationResult{
		Tier:       tier,
		Score:      score,
		Confidence: confidence,
		Signals:    signals,
	}
}

// computeConfidence calculates confidence based on how far the score is
// from the tier boundaries. Uses a sigmoid function centered at the tier's midpoint.
func computeConfidence(score float64, tier QueryTier, boundaries map[QueryTier]TierBoundary) float64 {
	boundary, ok := boundaries[tier]
	if !ok {
		boundary = DefaultTierBoundaries[tier]
	}

	// Calculate the midpoint of the tier
	midpoint := (boundary.Min + boundary.Max) / 2

	// Calculate distance from midpoint, normalized by tier width
	width := boundary.Max - boundary.Min
	if width <= 0 {
		return 0.5
	}

	// Distance from midpoint as a fraction of tier width
	distance := (score - midpoint) / (width / 2)

	// Sigmoid: 1 / (1 + e^(-k*x)) where k controls steepness
	// Higher k = sharper transition at boundaries
	k := 3.0
	confidence := 1.0 / (1.0 + math.Exp(-k*math.Abs(distance)))

	// Scale to [0.5, 1.0] range (center of tier = 1.0, boundary = 0.5)
	confidence = 0.5 + 0.5*confidence

	return math.Round(confidence*1000) / 1000 // Round to 3 decimal places
}

// GetBoundaries returns the tier boundaries used by this classifier.
func (c *RuleClassifierV2) GetBoundaries() map[QueryTier]TierBoundary {
	return c.boundaries
}
