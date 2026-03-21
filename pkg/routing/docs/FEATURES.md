# Feature Extraction and Scoring

**Package**: `pkg/routing/` · **Main Doc**: [README.md](./README.md)

## Overview

The Smart Router uses a rule-based feature extraction system to classify messages. This approach is:

- **Fast**: Sub-microsecond classification
- **Language-agnostic**: Works with any language
- **Deterministic**: Same input always produces same output
- **Transparent**: Every decision can be explained via signals

## Feature Extraction

Features are extracted in `features.go`:

```go
type Features struct {
    TokenEstimate      int  // Estimated token count
    CodeBlockCount     int  // Number of fenced code blocks
    RecentToolCalls    int  // Tool calls in recent history
    ConversationDepth  int  // Total messages in conversation
    HasAttachments     bool // Multi-modal content present
}
```

### Token Estimation

Tokens are estimated using a simple heuristic (4 characters ≈ 1 token):

```go
func EstimateTokens(text string) int {
    return len(text) / 4
}
```

This is intentionally simple for performance. The estimate is used for relative comparison, not exact counts.

### Code Block Detection

Detects fenced code blocks with optional language specifiers:

```
```python
def hello():
    print("world")
```
```

Counts the number of code blocks, not lines of code.

### Tool Call Counting

Counts tool calls from the conversation history. Only recent tool calls (within the last N messages) are considered.

### Attachment Detection

Checks for media references in the message:

```
media://uuid-here
```

## Scoring Weights

Each feature contributes to the final score:

| Feature | Condition | Weight | Notes |
|---------|-----------|--------|-------|
| Short message | < 20 tokens | **-0.10** | Negative! Routes to simple |
| Short-mid | 20-50 tokens | +0.06 | Between greeting and medium |
| Medium | 51-200 tokens | +0.15 | Standard length |
| Long | > 200 tokens | +0.35 | Needs stronger model |
| Code block | count > 0 | +0.40 | Strong signal for complex |
| Tool calls | > 3 | +0.25 | Dense agentic workflow |
| Tool calls | 1-3 | +0.10 | Active workflow |
| Conversation depth | > 10 | +0.10 | Accumulated context |
| Attachments | present | **1.0** | Hard gate to reasoning |

### Score Range

The final score is clamped to `[-1.0, 1.0]`:

- **Negative scores** indicate very simple messages (greetings)
- **Scores near 0** are moderate complexity
- **Scores near 1** require the most capable models

## Tier Mapping

Scores are mapped to tiers using boundaries:

```
Score:  -1.0 ─────── -0.05 ─────── 0.15 ─────── 0.35 ─────── 1.0
          │   simple   │   moderate  │   complex   │  reasoning │
          └────────────┴─────────────┴─────────────┴────────────┘
```

## Scoring Examples

### Example 1: Greeting

```
Input: "Hi"
Features:
  - TokenEstimate: 1 (< 20)
  - CodeBlockCount: 0
  - RecentToolCalls: 0
  - ConversationDepth: 1
  - HasAttachments: false

Signals:
  - short_message: -0.10

Final Score: -0.10
Tier: simple
```

### Example 2: Simple Question

```
Input: "What's the weather today?"
Features:
  - TokenEstimate: 12 (< 20)
  - CodeBlockCount: 0
  - RecentToolCalls: 0
  - ConversationDepth: 1
  - HasAttachments: false

Signals:
  - short_message: -0.10

Final Score: -0.10
Tier: simple
```

### Example 3: Medium Question

```
Input: "Can you explain the difference between TCP and UDP protocols and when to use each?"
Features:
  - TokenEstimate: 24 (20-50)
  - CodeBlockCount: 0
  - RecentToolCalls: 0
  - ConversationDepth: 1
  - HasAttachments: false

Signals:
  - token_short_mid: +0.06

Final Score: 0.06
Tier: moderate
```

### Example 4: Code Review

```
Input: "Please review this code:\n```python\ndef foo():\n    return 42\n```\nIs it correct?"
Features:
  - TokenEstimate: 25
  - CodeBlockCount: 1
  - RecentToolCalls: 0
  - ConversationDepth: 1
  - HasAttachments: false

Signals:
  - token_short_mid: +0.06
  - code_block: +0.40

Final Score: 0.46
Tier: reasoning (crosses 0.35 boundary)
```

### Example 5: Multi-Modal

```
Input: "What's in this image?"
Attachments: [image.png]
Features:
  - TokenEstimate: 6
  - CodeBlockCount: 0
  - RecentToolCalls: 0
  - ConversationDepth: 1
  - HasAttachments: true

Signals:
  - attachment_gate: 1.0

Final Score: 1.0
Tier: reasoning (hard gate)
```

### Example 6: Long Conversation

```
Input: "Continue from where we left off"
Features:
  - TokenEstimate: 10
  - CodeBlockCount: 0
  - RecentToolCalls: 5
  - ConversationDepth: 25
  - HasAttachments: false

Signals:
  - short_message: -0.10
  - tool_call_high: +0.25
  - conversation_depth: +0.10

Final Score: 0.25
Tier: complex
```

## Confidence Calculation

Confidence is computed using a sigmoid function based on how far the score is from tier boundaries:

```go
func computeConfidence(score float64, tier QueryTier, boundaries map[QueryTier]TierBoundary) float64 {
    midpoint := (boundary.Min + boundary.Max) / 2
    distance := (score - midpoint) / (width / 2)
    confidence := 1.0 / (1.0 + math.Exp(-k*math.Abs(distance)))
    return 0.5 + 0.5*confidence  // Scale to [0.5, 1.0]
}
```

**Interpretation**:
- Score at tier center → confidence ≈ 1.0
- Score near boundary → confidence ≈ 0.5
- Attachment hard gate → confidence = 1.0

## Signal Tracing

Every classification includes a signal trace for debugging:

```go
type Signal struct {
    Name         string  // e.g., "code_block", "short_message"
    Contribution float64 // Weight applied
    RawValue     any     // Original value (e.g., token count)
}
```

Example output:
```json
{
  "tier": "complex",
  "score": 0.25,
  "confidence": 0.78,
  "signals": [
    {"name": "token_medium", "contribution": 0.15, "raw_value": 75},
    {"name": "tool_call_low", "contribution": 0.10, "raw_value": 2}
  ]
}
```

## Customization

### Adjusting Weights

Weights are defined as constants in `classifier_v2.go`:

```go
const (
    ShortMessageBonus    = -0.10
    TokenShortMidWeight  = 0.06
    TokenMediumWeight    = 0.15
    TokenLongWeight      = 0.35
    CodeBlockWeight      = 0.40
    ToolCallHighWeight   = 0.25
    ToolCallLowWeight    = 0.10
    ConversationDepthWeight = 0.10
    AttachmentGateScore  = 1.0
)
```

To customize, modify these values and rebuild.

### Adjusting Boundaries

See [CONFIG.md](./CONFIG.md) for customizing tier boundaries via configuration.
