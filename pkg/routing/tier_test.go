package routing

import "testing"

func TestTryParseTier(t *testing.T) {
	tests := []struct {
		in       string
		wantTier QueryTier
		wantOK   bool
	}{
		{"simple", TierSimple, true},
		{" moderate ", TierModerate, true},
		{"complex", TierComplex, true},
		{"reasoning", TierReasoning, true},
		{"simpe", "", false},
		{"", "", false},
		{"SIMPLE", "", false},
	}
	for _, tt := range tests {
		got, ok := TryParseTier(tt.in)
		if ok != tt.wantOK || got != tt.wantTier {
			t.Errorf("TryParseTier(%q) = (%q, %v), want (%q, %v)", tt.in, got, ok, tt.wantTier, tt.wantOK)
		}
	}
}

func TestParseTier_LegacyUnknown(t *testing.T) {
	if got := ParseTier("nope"); got != TierModerate {
		t.Errorf("ParseTier unknown: got %s, want moderate", got)
	}
}

func TestValidateTierCutpoints(t *testing.T) {
	if err := ValidateTierCutpoints(-0.05, 0.15, 0.35); err != nil {
		t.Errorf("default cutpoints: %v", err)
	}
	if err := ValidateTierCutpoints(0, 0.2, 0.4); err != nil {
		t.Errorf("explicit zero sm: %v", err)
	}
	if err := ValidateTierCutpoints(-0.05, 0.15, 0.10); err == nil {
		t.Error("expected error when cr <= mc")
	}
	if err := ValidateTierCutpoints(0.2, 0.15, 0.35); err == nil {
		t.Error("expected error when sm >= mc")
	}
	if err := ValidateTierCutpoints(-1.1, 0.15, 0.35); err == nil {
		t.Error("expected error when sm < -1")
	}
	if err := ValidateTierCutpoints(-0.05, 0.15, 1.1); err == nil {
		t.Error("expected error when cr > 1")
	}
}

func TestRuleClassifierV2_TokenShortMidBand(t *testing.T) {
	c := NewRuleClassifierV2(nil)
	// 25 tokens: only token_short_mid (+0.06) → moderate
	r := c.Classify(Features{TokenEstimate: 25})
	found := false
	for _, s := range r.Signals {
		if s.Name == "token_short_mid" && s.Contribution == TokenShortMidWeight {
			found = true
		}
	}
	if !found {
		t.Errorf("expected token_short_mid signal, signals=%v", r.Signals)
	}
	if r.Score != TokenShortMidWeight {
		t.Errorf("score: got %v want %v", r.Score, TokenShortMidWeight)
	}
}
