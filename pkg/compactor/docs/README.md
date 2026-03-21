# MoonHub Context Compactor (`pkg/compactor`)

Four-layer context compression pipeline: rule-based pre-compression → deduplication → LLM summary → tiered summaries (L0/L1/L2). Works with Agent conversation history, token budget, and SQLite persistence.

## Recommended Reading Order (Documentation Flow)

1. **This document** — Responsibility boundaries and source map
2. [CONFIG.md](./CONFIG.md) — `config.json` / environment variables and defaults
3. Implementation details and integration — [`docs/implementation/compactor-status.md`](../../../docs/implementation/compactor-status.md) (architecture table, table structure, `instance.go` / `loop.go` hooks, test commands)

## Source Map (Aligned with Implementation)

```
pkg/compactor/
├── compactor.go   # Main orchestration, CompactorEngine
├── config.go      # Runtime Config (mapped from pkg/config)
├── rules.go       # Layer 1: Rule-based pre-compression
├── dedup.go       # Layer 2: Shingle + Jaccard deduplication
├── tiers.go       # Layer 4: Tiered summaries and tier selection
├── store.go       # SQLite: tiered summaries and compression state
└── tokens.go      # Token estimation and helpers
```

Testing: `go test ./pkg/compactor/...` (see coverage details in implementation status document).

## Related Code (Integration Points)

| Area | Path | Description |
|------|------|-------------|
| Global Config Type | `pkg/config/config.go` (`CompactorConfig`) | User-visible config and env tags |
| Defaults | `pkg/config/defaults.go` | `DefaultCompactorConfig()` |
| Agent Holder & Initialization | `pkg/agent/instance.go` | `Compactor`, `CompactorTriggerTokenPercent` |
| Trigger & Tiered Injection | `pkg/agent/loop.go` | `maybeCompact`, `getTieredSummary`, etc. |

## Alignment with learning Documentation Style

This directory follows the same convention as [`pkg/learning/docs`](../learning/docs/README.md): **package-level `docs/` contains "how to configure, where to read code"**, **repository-level `docs/implementation/*-status.md` contains "design, table structure, integration checklist and testing"**, facilitating future development and PR division of labor.
