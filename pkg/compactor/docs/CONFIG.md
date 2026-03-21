# Context Compactor Configuration

Configuration structure is defined in **`pkg/config/config.go`** under `CompactorConfig` and `TierBudgetsConfig`; defaults are in **`DefaultCompactorConfig()`** (`pkg/config/defaults.go`). Field names below are consistent with JSON / environment variables and source code.

## config.json Example

```json
{
  "compactor": {
    "enabled": true,
    "trigger_token_percent": 70,
    "keep_recent": 10,
    "tier_budgets": {
      "l0": 200,
      "l1": 1000,
      "l2": 3000
    },
    "dedup_enabled": true,
    "dedup_similarity_threshold": 0.6,
    "strip_emoji": true,
    "remove_duplicate_lines": true,
    "normalize_cjk": true,
    "smart_rule_selection": true,
    "parallel_processing": true,
    "incremental_compaction": true,
    "summarization_model": ""
  }
}
```

When `summarization_model` is empty, the default conversation model is used; to specify a separate summarization model, provide a model name available to the provider.

## Environment Variables (Consistent with struct tags)

| Environment Variable | Description |
|---------------------|-------------|
| `MOONHUB_COMPACTOR_ENABLED` | Enable/disable |
| `MOONHUB_COMPACTOR_TRIGGER_TOKEN_PERCENT` | Trigger percentage relative to context window |
| `MOONHUB_COMPACTOR_KEEP_RECENT` | Number of recent messages to keep after compaction |
| `MOONHUB_COMPACTOR_DEDUP_ENABLED` | Enable deduplication |
| `MOONHUB_COMPACTOR_DEDUP_SIMILARITY_THRESHOLD` | Jaccard similarity threshold |
| `MOONHUB_COMPACTOR_STRIP_EMOJI` | Strip emojis |
| `MOONHUB_COMPACTOR_REMOVE_DUPLICATE_LINES` | Remove duplicate lines (rule layer) |
| `MOONHUB_COMPACTOR_NORMALIZE_CJK` | CJK punctuation normalization |
| `MOONHUB_COMPACTOR_SMART_RULE_SELECTION` | Enable rule subset by content |
| `MOONHUB_COMPACTOR_PARALLEL_PROCESSING` | Parallel processing |
| `MOONHUB_COMPACTOR_INCREMENTAL_COMPACTION` | Incremental compaction |
| `MOONHUB_COMPACTOR_SUMMARIZATION_MODEL` | Dedicated summarization model (optional) |

## Default Values Summary

| Item | Default |
|------|---------|
| `enabled` | `true` |
| `trigger_token_percent` | `70` |
| `keep_recent` | `10` |
| `tier_budgets.l0/l1/l2` | `200` / `1000` / `3000` |
| `dedup_similarity_threshold` | `0.6` |
| Other boolean switches | Mostly `true` (see `DefaultCompactorConfig`) |

Persistence path, trigger logic, and loop integration are documented in [`docs/implementation/compactor-status.md`](../../../docs/implementation/compactor-status.md).
