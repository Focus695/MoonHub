package routing

import (
	"testing"
)

// mockClassifierV2 is a mock classifier for testing
type mockClassifierV2 struct {
	result ClassificationResult
}

func (m *mockClassifierV2) Classify(f Features) ClassificationResult {
	return m.result
}

func TestRouterV2_4TierMode(t *testing.T) {
	tierMapping := TierModelMapping{
		TierSimple:    "model-simple",
		TierModerate:  "model-moderate",
		TierComplex:   "model-complex",
		TierReasoning: "model-reasoning",
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2(cfg)

	if !router.IsV2Enabled() {
		t.Error("expected V2 mode to be enabled")
	}
	if router.GetTierMapping() == nil {
		t.Error("expected tier mapping to be set")
	}
}

func TestRouterV2_2TierFallback(t *testing.T) {
	// No tier mapping, only light model
	cfg := RouterConfigV2{
		LightModel: "model-light",
		Threshold:  0.35,
	}

	router := NewV2(cfg)

	if router.IsV2Enabled() {
		t.Error("expected V2 mode to be disabled")
	}
	if router.LightModel() != "model-light" {
		t.Errorf("expected light model 'model-light', got %q", router.LightModel())
	}
}

func TestRouterV2_SelectModel_V2Mode(t *testing.T) {
	tierMapping := TierModelMapping{
		TierSimple:    "model-simple",
		TierModerate:  "model-moderate",
		TierComplex:   "model-complex",
		TierReasoning: "model-reasoning",
	}

	tests := []struct {
		name          string
		result        ClassificationResult
		expectedModel string
		expectedTier  QueryTier
	}{
		{
			name: "simple_tier",
			result: ClassificationResult{
				Tier:       TierSimple,
				Score:      -0.1,
				Confidence: 0.9,
			},
			expectedModel: "model-simple",
			expectedTier:  TierSimple,
		},
		{
			name: "moderate_tier",
			result: ClassificationResult{
				Tier:       TierModerate,
				Score:      0.05,
				Confidence: 0.85,
			},
			expectedModel: "model-moderate",
			expectedTier:  TierModerate,
		},
		{
			name: "complex_tier",
			result: ClassificationResult{
				Tier:       TierComplex,
				Score:      0.25,
				Confidence: 0.8,
			},
			expectedModel: "model-complex",
			expectedTier:  TierComplex,
		},
		{
			name: "reasoning_tier",
			result: ClassificationResult{
				Tier:       TierReasoning,
				Score:      0.5,
				Confidence: 0.95,
			},
			expectedModel: "model-reasoning",
			expectedTier:  TierReasoning,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := RouterConfigV2{
				TierMapping: tierMapping,
			}

			router := NewV2WithClassifier(cfg, &mockClassifierV2{result: tt.result})
			model, tier, result := router.SelectModel("test message", nil, "primary-model")

			if model != tt.expectedModel {
				t.Errorf("expected model %q, got %q", tt.expectedModel, model)
			}
			if tier != tt.expectedTier {
				t.Errorf("expected tier %s, got %s", tt.expectedTier, tier)
			}
			if result.Score != tt.result.Score {
				t.Errorf("expected score %f, got %f", tt.result.Score, result.Score)
			}
		})
	}
}

func TestRouterV2_SelectModel_MissingTier(t *testing.T) {
	// Only map some tiers
	tierMapping := TierModelMapping{
		TierSimple:  "model-simple",
		TierComplex: "model-complex",
		// Missing moderate and reasoning
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2WithClassifier(cfg, &mockClassifierV2{
		result: ClassificationResult{
			Tier:       TierModerate,
			Score:      0.1,
			Confidence: 0.85,
		},
	})

	model, tier, _ := router.SelectModel("test message", nil, "primary-model")

	// Should fallback to primary model for unmapped tier
	if model != "primary-model" {
		t.Errorf("expected fallback to primary model, got %q", model)
	}
	if tier != TierModerate {
		t.Errorf("expected tier %s, got %s", TierModerate, tier)
	}
}

func TestRouterV2_SelectModel_2TierMode(t *testing.T) {
	cfg := RouterConfigV2{
		LightModel: "model-light",
		Threshold:  0.35,
	}

	router := NewV2(cfg)

	// Create a simple feature that will score low (short message)
	// The legacy router should route to light model
	model, tier, _ := router.SelectModel("hi", nil, "primary-model")

	// In 2-tier mode with short message, should use light model
	if model != "model-light" {
		t.Errorf("expected light model, got %q", model)
	}
	if tier != TierSimple {
		t.Errorf("expected TierSimple for light model, got %s", tier)
	}
}

func TestRouterV2_SelectModelLegacy(t *testing.T) {
	// Test with legacy mode (LightModel set)
	cfg := RouterConfigV2{
		LightModel: "model-light",
		Threshold:  0.35,
	}

	router := NewV2(cfg)

	model, usedLight, _ := router.SelectModelLegacy("hi", nil, "primary-model")

	// In legacy mode with short message, should use light model
	if model != "model-light" {
		t.Errorf("expected model-light, got %q", model)
	}
	if !usedLight {
		t.Error("expected usedLight to be true for short message")
	}
}

func TestRouterV2_SelectModelLegacy_ComplexTier(t *testing.T) {
	tierMapping := TierModelMapping{
		TierSimple:    "model-simple",
		TierModerate:  "model-moderate",
		TierComplex:   "model-complex",
		TierReasoning: "model-reasoning",
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2WithClassifier(cfg, &mockClassifierV2{
		result: ClassificationResult{
			Tier:       TierComplex,
			Score:      0.25,
			Confidence: 0.8,
		},
	})

	model, usedLight, _ := router.SelectModelLegacy("test message", nil, "primary-model")

	// Complex tier should not use light
	if model != "model-complex" {
		t.Errorf("expected model-complex, got %q", model)
	}
	if usedLight {
		t.Error("expected usedLight to be false for complex tier")
	}
}

func TestRouterV2_CustomBoundaries(t *testing.T) {
	customBoundaries := map[QueryTier]TierBoundary{
		TierSimple:    {Min: -1.0, Max: 0.0},
		TierModerate:  {Min: 0.0, Max: 0.3},
		TierComplex:   {Min: 0.3, Max: 0.6},
		TierReasoning: {Min: 0.6, Max: 1.0},
	}

	cfg := RouterConfigV2{
		TierBoundaries: customBoundaries,
	}

	router := NewV2(cfg)

	// Check that boundaries are used
	boundaries := router.GetTierBoundaries()
	if boundaries == nil {
		t.Error("expected custom boundaries to be set")
	}
	if boundaries[TierSimple].Max != 0.0 {
		t.Errorf("expected simple tier max 0.0, got %f", boundaries[TierSimple].Max)
	}
}

func TestRouterV2_Threshold(t *testing.T) {
	// 2-tier mode
	cfg2Tier := RouterConfigV2{
		LightModel: "model-light",
		Threshold:  0.4,
	}
	router2Tier := NewV2(cfg2Tier)
	if router2Tier.Threshold() != 0.4 {
		t.Errorf("expected threshold 0.4, got %f", router2Tier.Threshold())
	}

	// V2 mode - threshold should return complex tier boundary
	cfgV2 := RouterConfigV2{
		TierMapping: TierModelMapping{
			TierSimple:  "model-simple",
			TierComplex: "model-complex",
		},
	}
	routerV2 := NewV2(cfgV2)
	if routerV2.Threshold() != 0.15 { // Default boundary for complex
		t.Errorf("expected threshold 0.15 (complex boundary), got %f", routerV2.Threshold())
	}
}

func TestRouterV2_LightModel(t *testing.T) {
	// 2-tier mode
	cfg2Tier := RouterConfigV2{
		LightModel: "model-light",
	}
	router2Tier := NewV2(cfg2Tier)
	if router2Tier.LightModel() != "model-light" {
		t.Errorf("expected light model 'model-light', got %q", router2Tier.LightModel())
	}

	// V2 mode - light model maps to simple tier
	cfgV2 := RouterConfigV2{
		TierMapping: TierModelMapping{
			TierSimple: "model-simple",
		},
	}
	routerV2 := NewV2(cfgV2)
	if routerV2.LightModel() != "model-simple" {
		t.Errorf("expected light model 'model-simple', got %q", routerV2.LightModel())
	}
}

func TestRouterV2_NoRouting(t *testing.T) {
	// No routing configured
	cfg := RouterConfigV2{}
	router := NewV2(cfg)

	model, tier, _ := router.SelectModel("test message", nil, "primary-model")

	// Should fallback to primary model
	if model != "primary-model" {
		t.Errorf("expected fallback to primary model, got %q", model)
	}
	// Tier will be computed by classifier but won't matter
	_ = tier
}

func TestRouterV2_AllTiersMapped(t *testing.T) {
	tierMapping := TierModelMapping{
		TierSimple:    "flash",
		TierModerate:  "flash",
		TierComplex:   "plus",
		TierReasoning: "opus",
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2(cfg)

	// Verify all tiers are mapped
	mapping := router.GetTierMapping()
	for _, tier := range AllTiers() {
		if _, ok := mapping[tier]; !ok {
			t.Errorf("tier %s not mapped", tier)
		}
	}
}

// Integration test with real classifier
func TestRouterV2_Integration(t *testing.T) {
	tierMapping := TierModelMapping{
		TierSimple:    "flash",
		TierModerate:  "flash",
		TierComplex:   "plus",
		TierReasoning: "opus",
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2(cfg)

	tests := []struct {
		name          string
		msg           string
		expectedModel string
	}{
		{
			name:          "short_greeting",
			msg:           "hi",
			expectedModel: "flash", // simple or moderate
		},
		{
			name:          "code_request",
			msg:           "```go\nfmt.Println(\"hello\")\n```",
			expectedModel: "plus", // complex or reasoning (code block)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, _, _ := router.SelectModel(tt.msg, nil, "primary-model")
			if model != tt.expectedModel && model != "primary-model" {
				t.Errorf("expected model %q or primary, got %q", tt.expectedModel, model)
			}
		})
	}
}

func BenchmarkRouterV2_SelectModel(b *testing.B) {
	tierMapping := TierModelMapping{
		TierSimple:    "flash",
		TierModerate:  "flash",
		TierComplex:   "plus",
		TierReasoning: "opus",
	}

	cfg := RouterConfigV2{
		TierMapping: tierMapping,
	}

	router := NewV2(cfg)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = router.SelectModel("test message with some content", nil, "primary-model")
	}
}
