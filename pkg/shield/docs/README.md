# MoonHub SHIELD (`pkg/shield`)

Runtime threat evaluation: Parse YAML threat definitions from workspace `SHIELD.md`, perform matching and disposition (`block` / `require_approval` / `log`) on tool calls, skill installations, network egress events, etc., working with the approval manager.

**Repository Documentation Index** (reading order for subsystems): [docs/README.md](../../../docs/README.md).

## Recommended Reading Order (Documentation Flow)

1. **This document** — Responsibility boundaries and source map
2. [SHIELD_MD.md](./SHIELD_MD.md) — `SHIELD.md` format, conditional DSL, confidence and `action` semantics
3. [EXAMPLES.md](./EXAMPLES.md) — Minimal examples and runtime commands
4. Implementation details and integration checklist — [`docs/implementation/shield-status.md`](../../../docs/implementation/shield-status.md) (integration points, `instance.go` / `loop.go`, built-in threat tables, test commands)

## Source Map (Aligned with Implementation)

```
pkg/shield/
├── types.go            # ShieldAction, ThreatEntry, ShieldEvent, ShieldDecision, etc.
├── parser.go           # SHIELD.md → YAML block parsing, ParseDirectives
├── matcher.go          # Condition and event matching
├── engine.go           # ShieldEngine, Evaluate, confidence threshold, Reload/AddThreat
├── approval.go         # ApprovalManager, WaitForApproval
├── default_threats.go  # Built-in default threat library (string form)
├── context.go          # Approved tool execution marker (context.WithValue, avoid repeated interception)
└── *_test.go
```

Testing: `go test ./pkg/shield/...` (see implementation status document for details).

## Related Code (Integration Points)

| Area | Path | Description |
|------|------|-------------|
| Engine & Approval Holder | `pkg/agent/instance.go` | `Shield`, `ApprovalManager`; prioritize `workspace/SHIELD.md`, otherwise default library |
| Pre-tool Evaluation & Approval Flow | `pkg/agent/loop.go` | `Evaluate` before tool call; `require_approval` with user interaction |
| User Commands | `pkg/commands/cmd_approve.go`, `builtin.go` | `/approve`, `/reject` |
| Skill Installation | `pkg/tools/skills_install.go` | `ScopeSkillInstall` |
| Web Fetch | `pkg/tools/web.go` | `ScopeNetworkEgress` |

## Alignment with learning / compactor Documentation Style

This directory follows the same convention as [`pkg/learning/docs`](../learning/docs/README.md) and [`pkg/compactor/docs`](../compactor/docs/README.md): **package-level `docs/` contains "how to write policies, where to read code"**, **repository-level `docs/implementation/shield-status.md` contains "design, integration checklist, testing and troubleshooting"**, facilitating future development and PR division of labor.
