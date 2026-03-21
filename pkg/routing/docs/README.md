# Smart Router V2 - 4-Tier Model Routing

**Package**: `pkg/routing/` · **Status Doc**: [`docs/implementation/routing-status.md`](../../../docs/implementation/routing-status.md)

## Overview

MoonHub's Smart Router V2 is a 4-tier model routing system that automatically selects the most appropriate LLM based on message complexity. It evolved from the original 2-tier (light/heavy) system to provide finer-grained control over model selection and cost optimization.

### Key Features

- **4-Tier Classification**: simple, moderate, complex, reasoning
- **Rule-Based Scoring**: Sub-microsecond classification without API calls
- **Language-Agnostic**: Uses structural features, not language-specific patterns
- **Privacy-Safe Metrics**: Collects statistics without storing message content
- **Backward Compatible**: Falls back to 2-tier mode when configured

## Architecture

```
┌─────────────────┐     ExtractFeatures()     ┌──────────────────┐
│  InboundMessage │ ───────────────────────▶ │    Features      │
│  + History      │                           │                  │
└─────────────────┘                           └────────┬─────────┘
                                                       │
                                                       ▼
┌─────────────────┐     Classify()           ┌──────────────────┐
│ RouterV2        │ ◀─────────────────────── │ RuleClassifierV2 │
│                 │                          │                  │
│ tier_mapping:   │     Tier + Score         │ - Token estimate │
│   simple → A    │ ◀─────────────────────── │ - Code blocks    │
│   moderate → B  │     + Confidence         │ - Tool calls     │
│   complex → C   │     + Signals            │ - Conv depth     │
│   reasoning → D │                          │ - Attachments    │
└────────┬────────┘                          └──────────────────┘
         │
         │ SelectModel()
         ▼
┌─────────────────┐
│   Model Name    │ ──▶ Provider Candidates ──▶ LLM API
└─────────────────┘
```

## Quick Start

### Basic 4-Tier Configuration

```json
{
  "agents": {
    "defaults": {
      "routing": {
        "enabled": true,
        "tier_mapping": {
          "simple": "glm-4-flash",
          "moderate": "glm-4-flash",
          "complex": "glm-4-plus",
          "reasoning": "claude-opus-4-6"
        }
      }
    }
  }
}
```

### Legacy 2-Tier Configuration

```json
{
  "agents": {
    "defaults": {
      "routing": {
        "enabled": true,
        "light_model": "glm-4-flash",
        "threshold": 0.35
      }
    }
  }
}
```

## Package Structure

```
pkg/routing/
├── tier.go                # QueryTier types and boundaries
├── classifier.go          # Classifier interface + legacy RuleClassifier
├── classifier_v2.go       # RuleClassifierV2 (4-tier implementation)
├── features.go            # Feature extraction (language-agnostic)
├── router.go              # Router (2-tier) + RouterV2 (4-tier)
├── metrics.go             # MetricsCollector for statistics
├── recorder.go            # DecisionRecorder for visualization
├── agent_id.go            # Agent ID hashing for privacy
├── session_key.go         # Session key generation
└── docs/
    ├── README.md          # This file
    ├── CONFIG.md          # Configuration options
    ├── FEATURES.md        # Feature weights and scoring
    └── METRICS.md         # Observability and metrics
```

## 4-Tier System

| Tier | Score Range | Description | Typical Use Case |
|------|-------------|-------------|------------------|
| `simple` | [-1.0, -0.05) | Greetings, acknowledgments | "Hi", "Thanks", "OK" |
| `moderate` | [-0.05, 0.15) | Short questions, simple tasks | "What time is it?" |
| `complex` | [0.15, 0.35) | Coding, long context, tools | Code review, file edits |
| `reasoning` | [0.35, 1.0] | Deep analysis, multi-modal | Image analysis, research |

## Related Documentation

| Document | Description |
|----------|-------------|
| [CONFIG.md](./CONFIG.md) | Configuration options and examples |
| [FEATURES.md](./FEATURES.md) | Feature extraction and weight details |
| [METRICS.md](./METRICS.md) | Metrics collection and HTTP endpoints |
| [Implementation Status](../../../docs/implementation/routing-status.md) | Full implementation details |

## Integration Points

### Agent Integration (`pkg/agent/`)

The router is integrated into the agent loop:

```go
// pkg/agent/instance.go
type AgentInstance struct {
    Router         *routing.Router     // 2-tier
    RouterV2       *routing.RouterV2   // 4-tier
    TierCandidates map[routing.QueryTier][]providers.FallbackCandidate
}

// pkg/agent/loop.go
func (al *AgentLoop) selectCandidates(agent *AgentInstance, msg string, history []providers.Message) {
    if agent.RouterV2 != nil && agent.RouterV2.IsV2Enabled() {
        model, tier, result := agent.RouterV2.SelectModel(msg, history, agent.Model)
        // Use tier-specific candidates
    }
}
```

### Config Integration (`pkg/config/`)

```go
type RoutingConfig struct {
    Enabled    bool    `json:"enabled"`

    // 2-tier mode
    LightModel string  `json:"light_model"`
    Threshold  float64 `json:"threshold"`

    // 4-tier mode (V2)
    TierMapping    map[string]string     `json:"tier_mapping,omitempty"`
    TierBoundaries *TierBoundariesConfig `json:"tier_boundaries,omitempty"`
}
```

### Health Server Integration (`pkg/health/`)

```go
// HTTP endpoints for observability
GET /metrics              // Aggregate routing metrics
GET /routing/decisions    // Recent decision records
GET /routing/stats        // Statistical summary
```

## Performance

- **Classification latency**: < 1µs per message (pure rules, no API calls)
- **Memory overhead**: Zero extra allocations (except return structs)
- **Concurrency**: RouterV2 is safe for concurrent use from multiple goroutines
