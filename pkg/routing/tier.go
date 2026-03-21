package routing

import (
	"fmt"
	"strings"
)

// QueryTier represents the complexity tier of an incoming message.
// The 4-tier system allows fine-grained model selection:
//   - Simple: greetings, trivial Q&A (routed to fastest/cheapest models)
//   - Moderate: short questions, simple tasks (routed to capable but efficient models)
//   - Complex: coding, long context, multi-step reasoning (routed to powerful models)
//   - Reasoning: deep analysis, multi-modal, heavy computation (routed to most capable models)
type QueryTier string

const (
	// TierSimple is for very short messages like greetings, acknowledgments.
	// Score range: [-1.0, -0.05)
	TierSimple QueryTier = "simple"

	// TierModerate is for short questions and simple tasks.
	// Score range: [-0.05, 0.15)
	TierModerate QueryTier = "moderate"

	// TierComplex is for coding tasks, long messages, tool usage.
	// Score range: [0.15, 0.35)
	TierComplex QueryTier = "complex"

	// TierReasoning is for multi-modal, deep analysis, or very complex tasks.
	// Score range: [0.35, 1.0]
	TierReasoning QueryTier = "reasoning"
)

// TierBoundary defines the score range for a single tier.
type TierBoundary struct {
	Min float64
	Max float64
}

// DefaultTierBoundaries defines the default score ranges for each tier.
// These boundaries are calibrated so that:
//   - Pure greetings score around -0.10 → simple
//   - Short questions score around 0.00-0.10 → moderate
//   - Code blocks score 0.40+ → complex or reasoning
//   - Attachments score 1.0 → reasoning (hard gate)
var DefaultTierBoundaries = map[QueryTier]TierBoundary{
	TierSimple:    {Min: -1.0, Max: -0.05},
	TierModerate:  {Min: -0.05, Max: 0.15},
	TierComplex:   {Min: 0.15, Max: 0.35},
	TierReasoning: {Min: 0.35, Max: 1.0},
}

// AllTiers returns all defined tiers in order from simplest to most complex.
func AllTiers() []QueryTier {
	return []QueryTier{TierSimple, TierModerate, TierComplex, TierReasoning}
}

// ScoreToTier maps a complexity score to a tier using the default boundaries.
func ScoreToTier(score float64) QueryTier {
	for _, tier := range AllTiers() {
		boundary := DefaultTierBoundaries[tier]
		if score >= boundary.Min && score < boundary.Max {
			return tier
		}
	}
	// If score >= max of reasoning tier, return reasoning
	if score >= DefaultTierBoundaries[TierReasoning].Min {
		return TierReasoning
	}
	// Default to simple for very low scores
	return TierSimple
}

// ScoreToTierWithBoundaries maps a complexity score to a tier using custom boundaries.
func ScoreToTierWithBoundaries(score float64, boundaries map[QueryTier]TierBoundary) QueryTier {
	if boundaries == nil {
		return ScoreToTier(score)
	}

	for _, tier := range AllTiers() {
		if boundary, ok := boundaries[tier]; ok {
			if score >= boundary.Min && score < boundary.Max {
				return tier
			}
		}
	}

	// Check if score is in reasoning range
	if rb, ok := boundaries[TierReasoning]; ok && score >= rb.Min {
		return TierReasoning
	}

	return TierSimple
}

// String returns the string representation of the tier.
func (t QueryTier) String() string {
	return string(t)
}

// IsValid checks if the tier is one of the defined values.
func (t QueryTier) IsValid() bool {
	switch t {
	case TierSimple, TierModerate, TierComplex, TierReasoning:
		return true
	default:
		return false
	}
}

// ParseTier converts a string to a QueryTier.
// Returns TierModerate as default for unrecognized values (legacy callers,
// e.g. query filters). For config keys, use TryParseTier and reject unknowns.
func ParseTier(s string) QueryTier {
	t, _ := TryParseTier(s)
	if t == "" {
		return TierModerate
	}
	return t
}

// TryParseTier parses a tier name. Returns false for unknown or empty strings.
func TryParseTier(s string) (QueryTier, bool) {
	switch QueryTier(strings.TrimSpace(s)) {
	case TierSimple:
		return TierSimple, true
	case TierModerate:
		return TierModerate, true
	case TierComplex:
		return TierComplex, true
	case TierReasoning:
		return TierReasoning, true
	default:
		return "", false
	}
}

const tierScoreMin = -1.0
const tierScoreMax = 1.0

// ValidateTierCutpoints checks simple/moderate/complex/reasoning cutpoints sm, mc, cr
// where tiers are [-1,sm), [sm,mc), [mc,cr), [cr,1]. Requires strict sm < mc < cr
// and all cutpoints within [tierScoreMin, tierScoreMax].
func ValidateTierCutpoints(sm, mc, cr float64) error {
	if sm < tierScoreMin || cr > tierScoreMax {
		return fmt.Errorf("cutpoints must lie within [%.1f, %.1f] (got sm=%g cr=%g)", tierScoreMin, tierScoreMax, sm, cr)
	}
	if mc < tierScoreMin || mc > tierScoreMax {
		return fmt.Errorf("moderate/complex cutpoint mc=%g out of [%.1f, %.1f]", mc, tierScoreMin, tierScoreMax)
	}
	if !(sm < mc && mc < cr) {
		return fmt.Errorf("need sm < mc < cr, got sm=%g mc=%g cr=%g", sm, mc, cr)
	}
	return nil
}
