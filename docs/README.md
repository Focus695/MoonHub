# MoonHub Documentation Index

This page is the **entry point and reading guide** for the repository documentation.

## Recommended Reading Order (First Time)

1. [Repository Root README](../README.md) — Feature overview, license, and recent changelog summary
2. [Troubleshooting](./troubleshooting.md), [Debug Guide](./debug.md) — Runtime issues reference
3. [Tools & Capabilities Configuration](./tools_configuration.md) — Tool-side configuration (`tools.*` in `config.json`)
4. If you enable **sub-agent delegation** (`delegation.enabled`): [pkg/delegation/docs/README.md](../pkg/delegation/docs/README.md) → [CONFIG.md](../pkg/delegation/docs/CONFIG.md) → [implementation status](./implementation/delegation-status.md)

## Documentation Flow by Subsystem

### Plugin Architecture (Channel / Provider / Tool)

| Order | Document | Description |
| --- | --- | --- |
| 1 | [`pkg/framework/docs/README.md`](../pkg/framework/docs/README.md) | Framework positioning, package names, and import conventions |
| 2 | [`pkg/framework/docs/ARCHITECTURE.md`](../pkg/framework/docs/ARCHITECTURE.md) | Interfaces, Registry, lifecycle |
| 3 | [`pkg/framework/docs/INTEGRATION.md`](../pkg/framework/docs/INTEGRATION.md) | Integration with Gateway, agent, and subsystems |
| 4 | [`pkg/plugins/docs/README.md`](../pkg/plugins/docs/README.md) | Built-in plugin directory conventions |
| 5 | [`pkg/plugins/docs/PLUGIN_INDEX.md`](../pkg/plugins/docs/PLUGIN_INDEX.md) | Plugin list and paths |
| 6 | [`pkg/plugins/docs/ADDING_A_PLUGIN.md`](../pkg/plugins/docs/ADDING_A_PLUGIN.md) | Steps to add or migrate plugins |
| Status | [`docs/implementation/plugin-status.md`](./implementation/plugin-status.md) | Design decisions and implementation progress |

### Self-Improving System

| Order | Document | Description |
| --- | --- | --- |
| 1 | [`pkg/learning/docs/README.md`](../pkg/learning/docs/README.md) | Feedback and capability overview |
| 2 | [`pkg/learning/docs/CONFIG.md`](../pkg/learning/docs/CONFIG.md) | `learning` configuration options |
| 3 | [`pkg/learning/docs/I18N.md`](../pkg/learning/docs/I18N.md) | Language and detection behavior |
| 4 | [`pkg/learning/docs/EXAMPLES.md`](../pkg/learning/docs/EXAMPLES.md), [`SupportedPatterns.md`](../pkg/learning/docs/SupportedPatterns.md) | Examples and pattern tables |
| Status | [`docs/implementation/learning-status.md`](./implementation/learning-status.md) | Implementation status and integration points |

### Context Compactor

| Order | Document | Description |
| --- | --- | --- |
| 1 | [`pkg/compactor/docs/README.md`](../pkg/compactor/docs/README.md) | Four-layer pipeline and code entry points |
| 2 | [`pkg/compactor/docs/CONFIG.md`](../pkg/compactor/docs/CONFIG.md) | `compactor` configuration and environment variables |
| Status | [`docs/implementation/compactor-status.md`](./implementation/compactor-status.md) | Architecture, storage, agent integration, and test commands |

### SHIELD Runtime Security (`pkg/shield`)

| Order | Document | Description |
| --- | --- | --- |
| 1 | [`pkg/shield/docs/README.md`](../pkg/shield/docs/README.md) | Responsibilities, source map, and integration entry |
| 2 | [`pkg/shield/docs/SHIELD_MD.md`](../pkg/shield/docs/SHIELD_MD.md) | `SHIELD.md` YAML, conditional DSL, confidence semantics |
| 3 | [`pkg/shield/docs/EXAMPLES.md`](../pkg/shield/docs/EXAMPLES.md) | Workspace policy examples and approval commands |
| Status | [`docs/implementation/shield-status.md`](./implementation/shield-status.md) | Design, built-in threat tables, integration checklist, and tests |

### Adaptive Memory

| Document | Description |
| --- | --- |
| [`docs/implementation/memory-status.md`](./implementation/memory-status.md) | Package structure, capabilities, configuration, and integration notes; source in `pkg/adaptive_memory/` |

### Delegation System (sub-agent orchestration)

| Order | Document | Description |
| --- | --- | --- |
| 1 | [`pkg/delegation/docs/README.md`](../pkg/delegation/docs/README.md) | Responsibilities, source map, integration table |
| 2 | [`pkg/delegation/docs/CONFIG.md`](../pkg/delegation/docs/CONFIG.md) | `delegation` in `config.json` and environment variables |
| Status | [`docs/implementation/delegation-status.md`](./implementation/delegation-status.md) | Schema, eight tools, Intercom topics, examples, tests |

### Channels

Channel architecture, migration, and how to implement a channel: [`pkg/channels/README.md`](../pkg/channels/README.md). Per-channel behavior also lives with each plugin under [`pkg/plugins/channels/`](../pkg/plugins/channels/) (see [`pkg/plugins/docs/PLUGIN_INDEX.md`](../pkg/plugins/docs/PLUGIN_INDEX.md)).

---

## `docs/implementation/` Overview

| File | Topic |
| --- | --- |
| [`plugin-status.md`](./implementation/plugin-status.md) | Plugin architecture implementation status |
| [`learning-status.md`](./implementation/learning-status.md) | Self-improving system implementation status |
| [`compactor-status.md`](./implementation/compactor-status.md) | Context compactor implementation status |
| [`shield-status.md`](./implementation/shield-status.md) | SHIELD runtime security implementation status |
| [`memory-status.md`](./implementation/memory-status.md) | Memory system implementation status |
| [`delegation-status.md`](./implementation/delegation-status.md) | Sub-agent delegation orchestration implementation status |
