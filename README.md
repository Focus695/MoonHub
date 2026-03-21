# MoonHub

> [!NOTE]
> **Based On**
>
> This project is based on features from [TinyClaw](https://github.com/wgtechlabs/tinyclaw) and the lightweight design of [PicoClaw](https://github.com/sipeed/picoclaw), with unique features developed for this project. Licensed under GPL-3.0.

**Documentation Index** (plugins / learning / compactor / memory, etc.): [`docs/README.md`](docs/README.md).

## Features

### Implemented

- **Adaptive Memory** — 3-layer memory system (episodic, semantic FTS5, temporal decay) that learns what to remember and forget over time.
- **Self-Improving** — Behavioral pattern detection system that learns from user feedback, tracks tool usage preferences, and evolves patterns over time.
- **Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.
- **Context Compactor** — 4-layer context compaction pipeline (rules, dedup, LLM summary, L0/L1/L2 tiers) integrated in the agent loop; see [`docs/implementation/compactor-status.md`](docs/implementation/compactor-status.md) and [`pkg/compactor/docs/`](pkg/compactor/docs/README.md).
- **SHIELD.md Anti-Malware** — Runtime threat evaluation engine with YAML threat parsing, pattern matching, approval workflow, and 8 built-in threats; see [`docs/implementation/shield-status.md`](docs/implementation/shield-status.md) and [`pkg/shield/docs/`](pkg/shield/docs/README.md).

### Planned

- ~~**Self-Improving** — Behavioral pattern detection that makes the agent better with every interaction. It grows with you.~~ → **Implemented**
- ~~**Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.~~ → **Implemented**
- ~~**Context Compactor** — 4-layer context compaction pipeline with rule-based pre-compression, deduplication, LLM summarization, and tiered summaries.~~ → **Implemented**
- ~~**SHIELD.md Anti-Malware** — Runtime SHIELD.md enforcement engine with threat parsing, pattern matching, and built-in anti-malware protection.~~ → **Implemented**
- **Smart Routing** — 8-dimension query classifier that routes simple queries to cheap models and complex ones to powerful ones, cutting LLM costs.
- **Delegation System** — Autonomous sub-agent orchestration with self-improving role templates, blackboard collaboration, and adaptive timeouts.
- **Inter-Agent Comms** — Lightweight pub/sub event bus for real-time inter-agent communication with wildcard subscriptions and bounded history.

## Changelog

<details>
<summary><strong>2026-03-21 — SHIELD.md Anti-Malware Implementation</strong></summary>

#### New Features

**SHIELD.md Anti-Malware System** (`pkg/shield/`)

A runtime threat evaluation engine inspired by TinyClaw design:

- **Threat Parser** — YAML-formatted SHIELD.md parser with support for threat definitions, directives, and metadata
- **Pattern Matcher** — Condition syntax support for tool calls, file paths, network egress, skill operations
- **Enforcement Actions** — Three action types: `block`, `require_approval`, `log` with priority-based resolution
- **Approval Workflow** — `/approve` and `/reject` commands for user-confirmed actions with 5-minute timeout
- **Tool Integration** — Shield evaluation integrated into `web_fetch` and `install_skill` tools
- **Default Threats** — 8 built-in threats covering SQL injection, command injection, path traversal, credential access, etc.

#### Technical Details

- Confidence threshold (0.85) with severity override for critical threats
- Action priority: `block` > `require_approval` > `log`
- Context-based approval bypass to prevent double-evaluation
- Comprehensive unit tests (31 tests, 100% pass rate)

#### Files Changed

- `pkg/shield/` — New package (11 core files + 5 test files)
- `pkg/agent/instance.go` — Shield and ApprovalManager initialization
- `pkg/agent/loop.go` — Shield evaluation in tool execution flow
- `pkg/commands/cmd_approve.go` — Approve/reject command handlers
- `pkg/tools/web.go` — Shield integration for network egress
- `pkg/tools/skills_install.go` — Shield integration for skill installation
- `docs/implementation/shield-status.md` — Implementation status

</details>

<details>
<summary><strong>2026-03-20 — Context Compactor &amp; documentation flow</strong></summary>

#### Features

- **Context Compactor** (`pkg/compactor/`) — Four-layer pipeline (rule-based pre-compression, deduplication, LLM summary, L0/L1/L2 tiers) integrated into Agent; see `compactor` config and [`pkg/compactor/docs/CONFIG.md`](pkg/compactor/docs/CONFIG.md).

#### Documentation

- Added repository documentation entry [`docs/README.md`](docs/README.md), distinguishing "package docs" from `docs/implementation/*-status.md` consistent with `pkg/learning/docs`.
- Added [`pkg/compactor/docs/`](pkg/compactor/docs/README.md) (README + CONFIG).
- Fixed link to `plugin-architecture-status.md` pointing to actual file [`docs/implementation/plugin-status.md`](docs/implementation/plugin-status.md).

</details>

<details>
<summary><strong>2025-03-20 — Plugin Architecture Implementation</strong></summary>

#### New Features

**Plugin Architecture System** (`pkg/framework/`, `pkg/plugins/`)

A comprehensive plugin system that makes channels, providers, and tools all extensible plugins:

- **Core Framework** — Plugin types, interfaces (Channel/Provider/Tool), registry system, and lifecycle manager
- **Channel Plugins** — 16 channel plugins migrated (Telegram, Discord, Slack, Matrix, Feishu, QQ, DingTalk, LINE, OneBot, WeCom, WeCom App, WeCom AIBot, Pico, IRC, MaixCam, WhatsApp)
- **Provider Plugins** — 8 provider plugins migrated (OpenAI Compat, OpenAI OAuth, Anthropic, Anthropic Messages, Antigravity, Claude CLI, Codex CLI, GitHub Copilot)
- **Tool Plugins** — Web tools (web_search, web_fetch) and message tool migrated to plugin system
- **Plugin Resolver** — Provider factory now supports plugin-first resolution with built-in fallback

#### Technical Improvements

- Added `SetPluginProviderResolver` for provider plugin integration
- Added `NewAgentLoopWithPluginTools` for tool plugin merging
- Added `MergeFrom` method to ToolRegistry for combining plugin tools
- Added `InitializeToolsOnly` to plugin manager for early tool initialization
- Updated channel manager to use plugin system for initialization
- Removed legacy factory pattern code from channel registry

#### Files Changed

- `pkg/framework/` — New package (7 core files)
- `pkg/plugins/channels/` — 16 channel plugins
- `pkg/plugins/providers/` — 8 provider plugins
- `pkg/plugins/tools/` — 2 tool plugins
- `pkg/plugins/docs/` — Plugin documentation
- `cmd/moonhub/internal/gateway/helpers.go` — Plugin imports and initialization
- `pkg/agent/loop.go` — Plugin tool integration
- `pkg/channels/manager.go` — Plugin-based channel initialization
- `pkg/providers/factory_provider.go` — Plugin resolver support
- `pkg/tools/registry.go` — MergeFrom method
- `docs/implementation/plugin-status.md` — Implementation status

</details>

<details>
<summary><strong>2025-03-19 — Self-Improving System Implementation</strong></summary>

#### New Features

**Self-Improving Behavioral Pattern Detection System** (`pkg/learning/`)

A comprehensive learning system that makes the agent better with every interaction:

- **Pattern Detector** — Multi-layer signal detection using regex patterns, semantic keyword analysis, and conversation flow analysis. Supports both English and Chinese feedback detection.
- **Tool Tracker** — Tracks tool usage statistics including success rates, user acceptance/rejection, duration metrics (avg, P50, P95), and preference scoring.
- **Behavioral Scorer** — Multi-dimensional scoring system measuring response quality, tool efficiency, context relevance, correction rate, and adaptation speed.
- **Pattern Evolution** — Ebbinghaus-inspired decay algorithm, pattern merging, pruning of stale patterns, and contradiction detection.
- **Proactive Suggestions** — Generates optimization suggestions based on detected patterns and behavioral trends.
- **Agent Integration** — Seamless integration with the agent loop for automatic learning during conversations.

#### Technical Improvements

- Added FTS5 full-text search for pattern queries with proper special character escaping
- Implemented proper JSON error handling throughout the persistence layer
- Fixed SQL parameter mismatches in tool usage tracking
- Added comprehensive unit tests (43 tests, 100% pass rate)
- Updated golangci-lint configuration to v2 format

#### Bug Fixes

- Fixed `GetAllPatterns()`, `GetPatternsByCategory()`, `GetToolUsagePatterns()` returning nil instead of empty slices
- Fixed `deduplicateSignals()` returning nil instead of empty slice
- Fixed `RecordUserAcceptance()` not incrementing `UserAcceptedCalls` counter
- Fixed `calculatePreference()` using wrong metric (success rate instead of acceptance rate)
- Removed premature `break` statements in pattern detection loops to capture all matching signals

#### Files Changed

- `pkg/learning/` — New package (14 files, ~3000 lines)
- `pkg/agent/loop.go` — Learning integration in agent loop
- `pkg/agent/context.go` — Learning context injection
- `pkg/agent/memory.go` — Memory-learning bridge
- `pkg/config/config.go` — Learning configuration options
- `.golangci.yaml` — Updated to v2 format

</details>
