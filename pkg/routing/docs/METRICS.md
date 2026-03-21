# Metrics and Observability

**Package**: `pkg/routing/` · **Main Doc**: [README.md](./README.md)

## Overview

The Smart Router includes a privacy-safe metrics and observability system that collects statistics without storing user message content.

## Privacy by Design

**What is collected:**
- Aggregate statistics (counts, averages, distributions)
- Tier assignments
- Signal triggers
- Score distributions

**What is NOT collected:**
- Raw message content
- User identifiers (use hashed agent IDs)
- Session keys
- Any PII

## Metrics Collector

### Interface

```go
type MetricsCollector interface {
    RecordClassification(result ClassificationResult)
    GetSnapshot() MetricsSnapshot
    Reset()
}
```

### Implementation

```go
type DefaultMetricsCollector struct {
    mu                    sync.RWMutex
    totalClassifications  int64
    tierCounts            map[QueryTier]int64
    scoreSum              float64
    confidenceSum         float64
    signalCounts          map[string]int64
    scoreDistribution     [10]int64      // Histogram buckets
    confidenceDistribution[10]int64      // Histogram buckets
}
```

### Usage

```go
// Create and set global collector
collector := routing.NewDefaultMetricsCollector()
routing.SetGlobalMetrics(collector)

// Router automatically records metrics
model, tier, result := router.SelectModel(msg, history, primaryModel)

// Get snapshot
snapshot := collector.GetSnapshot()
fmt.Printf("Total classifications: %d\n", snapshot.TotalClassifications)
fmt.Printf("Average score: %.3f\n", snapshot.AvgScore)
```

### Metrics Snapshot

```go
type MetricsSnapshot struct {
    Timestamp            time.Time
    TotalClassifications int64
    AvgScore             float64
    AvgConfidence        float64
    TierCounts           map[string]int64
    SignalCounts         map[string]int64
    ScoreDistribution    []int64    // 10 buckets: [-1,-0.8), [-0.8,-0.6), ...
    ConfidenceDistribution []int64  // 10 buckets: [0,0.1), [0.1,0.2), ...
}
```

### Score Distribution Buckets

| Bucket | Score Range |
|--------|-------------|
| 0 | [-1.0, -0.8) |
| 1 | [-0.8, -0.6) |
| 2 | [-0.6, -0.4) |
| 3 | [-0.4, -0.2) |
| 4 | [-0.2, 0.0) |
| 5 | [0.0, 0.2) |
| 6 | [0.2, 0.4) |
| 7 | [0.4, 0.6) |
| 8 | [0.6, 0.8) |
| 9 | [0.8, 1.0] |

## Decision Recorder

For debugging and visualization, the decision recorder stores recent classification decisions.

### Interface

```go
type DecisionRecorder struct {
    // Configurable settings
    BufferSize    int     // Max decisions to store (default: 1000)
    MaxSignalDepth int    // Max signals per decision (default: 10)
    Enabled       bool    // Enable/disable recording
}
```

### Recorded Decision

```go
type RecordedDecision struct {
    ID         int64
    Timestamp  time.Time
    Tier       QueryTier
    Score      float64
    Confidence float64
    Signals    []SignalSummary  // RawValue stripped for privacy
    Model      string
    AgentID    string           // Should be hashed, not raw ID
}
```

### Usage

```go
// Create recorder
cfg := routing.DefaultDecisionRecorderConfig()
cfg.BufferSize = 500
recorder := routing.NewDecisionRecorder(cfg)

// Set as global
routing.SetGlobalRecorder(recorder)

// Router automatically records decisions
model, tier, result := router.SelectModelWithRecord(msg, history, primaryModel, hashedAgentID)

// Query decisions
recent := recorder.GetRecent(50)
simpleDecisions := recorder.GetByTier(routing.TierSimple, 10)
stats := recorder.GetStats()
```

### Privacy Protections

1. **RawValue Stripping**: Signal RawValue fields are removed
2. **AgentID Hashing**: Use `routing.HashAgentID()` to hash agent IDs
3. **Ring Buffer**: Old decisions automatically discarded
4. **No Content**: Message content is never stored

```go
// Hash agent ID before passing to recorder
hashedID := routing.HashAgentID("my-agent-name")
model, tier, result := router.SelectModelWithRecord(msg, history, primaryModel, hashedID)
```

## HTTP Endpoints

When the health server is running, these endpoints are available:

### GET /metrics

Returns aggregate routing metrics:

```json
{
  "timestamp": "2026-03-21T13:00:00Z",
  "total_classifications": 1234,
  "avg_score": 0.15,
  "avg_confidence": 0.82,
  "tier_counts": {
    "simple": 400,
    "moderate": 300,
    "complex": 350,
    "reasoning": 184
  },
  "signal_counts": {
    "short_message": 400,
    "code_block": 350,
    "tool_call_high": 50,
    "attachment_gate": 25
  },
  "score_distribution": [50, 40, 30, 80, 200, 300, 250, 150, 100, 29],
  "confidence_distribution": [10, 20, 30, 40, 50, 100, 200, 300, 350, 134]
}
```

### GET /routing/decisions

Returns recent decision records:

**Query Parameters:**
- `limit` - Number of decisions to return (default: 50, max: 1000)
- `tier` - Filter by tier: simple, moderate, complex, reasoning

```json
{
  "total": 1234,
  "limit": 5,
  "decisions": [
    {
      "id": 1233,
      "timestamp": "2026-03-21T12:59:59Z",
      "tier": "complex",
      "score": 0.25,
      "confidence": 0.78,
      "signals": [
        {"name": "code_block", "contribution": 0.40},
        {"name": "token_medium", "contribution": 0.15}
      ],
      "model": "glm-4-plus"
    }
  ]
}
```

### GET /routing/stats

Returns aggregated statistics:

```json
{
  "total_decisions": 1234,
  "tier_percentages": {
    "simple": 32.4,
    "moderate": 24.3,
    "complex": 28.4,
    "reasoning": 14.9
  },
  "avg_score": 0.15,
  "avg_confidence": 0.82,
  "most_common_signals": [
    {"name": "short_message", "count": 400},
    {"name": "code_block", "count": 350}
  ]
}
```

## Integration with RouterV2

The RouterV2 automatically integrates with global metrics and recorder:

```go
// Global setup (once at startup)
collector := routing.NewDefaultMetricsCollector()
routing.SetGlobalMetrics(collector)

recorder := routing.NewDecisionRecorder(routing.DefaultDecisionRecorderConfig())
routing.SetGlobalRecorder(recorder)

// Router automatically uses globals
router := routing.NewV2(cfg)

// Or set custom instances
router.SetMetrics(customCollector)
router.SetRecorder(customRecorder)
```

## NoOp Implementation

For testing or when metrics are not needed:

```go
collector := &routing.NoOpMetricsCollector{}
routing.SetGlobalMetrics(collector)
```

## Monitoring Recommendations

### Key Metrics to Watch

| Metric | Alert Threshold | Meaning |
|--------|-----------------|---------|
| High `reasoning` % | > 30% | Too many complex tasks, review tier boundaries |
| Low `simple` % | < 10% | May be over-routing to expensive models |
| Low avg confidence | < 0.6 | Boundary tuning may be needed |
| High `attachment_gate` | - | Normal if multi-modal is common |

### Grafana Dashboard Example

```promql
# Tier distribution
rate(routing_tier_total[5m])

# Average score
rate(routing_score_sum[5m]) / rate(routing_classifications_total[5m])

# Signal triggers
rate(routing_signals_total[5m])
```

## Testing

### Unit Tests

```bash
go test ./pkg/routing/... -v

# Specific test files
go test ./pkg/routing/ -run TestMetrics -v
go test ./pkg/routing/ -run TestRecorder -v
```

### Test Coverage

| File | Tests | Coverage |
|------|-------|----------|
| `metrics.go` | 6 tests | Collector, snapshot, reset, concurrent, noop |
| `recorder.go` | 9 tests | Record, ring buffer, filter, stats, privacy |
