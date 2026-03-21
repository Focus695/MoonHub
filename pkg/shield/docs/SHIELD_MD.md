# SHIELD.md Format and Conditional Syntax

MoonHub looks for `SHIELD.md` in the workspace root directory (consistent with Agent initialization logic). The file is Markdown, with **each threat** placed in a fenced code block with language `yaml` or `yml`.

## YAML Fields (ThreatEntry)

| Field | Description |
|-------|-------------|
| `id` | Unique ID, must contain `THREAT-` substring (used by parser to identify threat blocks) |
| `fingerprint` | Short identifier for logging and deduplication |
| `category` | `tool` / `skill` / `network` / `file` / `prompt` / `memory` / `message` / `other` |
| `severity` | `critical` / `high` / `medium` / `low` |
| `confidence` | 0.0–1.0, participates in engine threshold logic (see below) |
| `action` | **Actual disposition on match**: `block`, `require_approval`, `log` |
| `title` / `description` | Display and troubleshooting text |
| `recommendation_agent` | Multi-line string: one condition per line, see below |
| `expires_at` / `revoked` / `revoked_at` | Optional; expired or revoked entries are ignored |

### Disposition Actions and Confidence (Engine Behavior)

- Constant `ConfidenceThreshold = 0.85` (see `engine.go`).
- `confidence >= 0.85`: Use the entry's `action`.
- `confidence < 0.85`: Non-`log` actions are downgraded to `require_approval`; **exception**: when `severity == critical` and `action == block`, it remains `block`.
- When multiple threats match simultaneously, the strongest action is selected: `block` > `require_approval` > `log`.

### `recommendation_agent` Line Prefixes

Each non-empty line should be one of the following (case-insensitive), with a condition string after the colon (used for matching only, see `ParseDirectives` for parsing):

- `BLOCK:`
- `APPROVE:` (parsed as internal `require_approval` semantics)
- `LOG:`

**Note**: The current `ShieldEngine.Evaluate` uses the YAML top-level **`action` field** (adjusted by confidence rules) for decision-making, not individual line prefixes. It's recommended that **`action` matches the disposition you expect in the text** to avoid reading ambiguity.

## Conditional DSL (Aligned with `matcher.go`)

Condition strings are `TrimSpace`d before matching; prefix matching uses case-insensitive detection in some implementations. Common patterns:

| Condition Pattern | Meaning |
|-------------------|---------|
| `tool.call with arguments containing (...)` | Tool arguments (serialized) match comma-separated substrings in parentheses |
| `tool.call <name>` | Tool name equals or contains (case-insensitive) |
| `skill name equals <value>` | Skill name exact match |
| `skill name contains <value>` | Skill name substring |
| `outbound request to <domain>` | Outbound domain pattern |
| `file path equals <path>` | Path exact match |
| `file path contains <path>` | Path substring |
| `secrets read path equals <path>` | Secrets path (supports wildcards and other matching logic, see source) |
| `incoming message contains <text>` | User input text |
| `message contains <text>` | Alias for the above |

**Order-sensitive**: `tool.call with arguments containing` must be recognized before the generalized `tool.call`; implementation follows this order.

## Event Scopes (ShieldScope)

Consistent with `ShieldEvent.Scope` passed by integration code, e.g.: `tool.call`, `skill.install`, `skill.execute`, `network.egress`, `secrets.read`, `prompt`. The matcher filters threats by `category` and `scope` compatibility (see `isScopeCompatible`).
