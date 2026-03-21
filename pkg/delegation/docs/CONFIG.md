# Delegation Configuration

The `delegation` object in **`config.json`** is defined as **`DelegationConfig`** in `pkg/config/config.go`; defaults are in **`DefaultDelegationConfig()`** (`pkg/config/defaults.go`). Field names below match JSON keys and code.

## config.json Example

```json
{
  "delegation": {
    "enabled": true,
    "default_sub_agent_tools": [
      "read_file",
      "list_dir",
      "web_search",
      "web_fetch",
      "memory_recall"
    ],
    "max_active_per_user": 10,
    "max_concurrent_tasks": 3,
    "retention_days": 14,
    "reuse_threshold": 0.6
  }
}
```

When `enabled` is `false` (default), delegation tools are not registered and no `delegation.db` activity occurs for the agent path described in the implementation status doc.

## Field Reference

| JSON field | Type | Default | Description |
|------------|------|---------|-------------|
| `enabled` | bool | `false` | Master switch for sub-agent delegation |
| `default_sub_agent_tools` | string[] | see defaults | Tool names allowed on spawned sub-agents unless overridden per task |
| `max_active_per_user` | int | `10` | Cap on active sub-agents per delegation partition / user |
| `max_concurrent_tasks` | int | `3` | Concurrent task limit for delegation scheduling |
| `retention_days` | int | `14` | Retention window for persisted delegation data (see store / cleanup in code) |
| `reuse_threshold` | float | `0.6` | Jaccard-style template reuse threshold (0–1) |

Default tool list if omitted: `read_file`, `list_dir`, `web_search`, `web_fetch`, `memory_recall`.

## Environment Variables (Struct Tags)

| Environment Variable | Maps To |
|---------------------|---------|
| `MOONHUB_DELEGATION_ENABLED` | `enabled` |
| `MOONHUB_DELEGATION_MAX_ACTIVE_PER_USER` | `max_active_per_user` |
| `MOONHUB_DELEGATION_MAX_CONCURRENT_TASKS` | `max_concurrent_tasks` |
| `MOONHUB_DELEGATION_RETENTION_DAYS` | `retention_days` |
| `MOONHUB_DELEGATION_REUSE_THRESHOLD` | `reuse_threshold` |

`default_sub_agent_tools` is JSON-only (no env tag on the struct).

## Further Reading

- Tool names, schema, and loop integration: [`docs/implementation/delegation-status.md`](../../../docs/implementation/delegation-status.md)
- Package overview: [README.md](./README.md)
