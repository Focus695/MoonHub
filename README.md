# MoonHub

> [!NOTE]
> **基于声明 (Based On)**
>
> 本项目基于 [TinyClaw](https://github.com/wgtechlabs/tinyclaw) 的诸多特性与 [PicoClaw](https://github.com/sipeed/picoclaw) 的轻量化设计融合更改，并将在本项目上发展独有功能。本项目遵循 GPL-3.0 许可证。

**文档索引与阅读顺序**（插件 / 学习 / 压缩器 / 记忆等）：[`docs/README.md`](docs/README.md)。

## Features

### Implemented

- **Adaptive Memory** — 3-layer memory system (episodic, semantic FTS5, temporal decay) that learns what to remember and forget over time.
- **Self-Improving** — Behavioral pattern detection system that learns from user feedback, tracks tool usage preferences, and evolves patterns over time.
- **Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.
- **Context Compactor** — 4-layer context compaction pipeline (rules, dedup, LLM summary, L0/L1/L2 tiers) integrated in the agent loop; see [`docs/implementation/compactor-status.md`](docs/implementation/compactor-status.md) and [`pkg/compactor/docs/`](pkg/compactor/docs/README.md).

### Planned

- ~~**Self-Improving** — Behavioral pattern detection that makes the agent better with every interaction. It grows with you.~~ → **Implemented**
- ~~**Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.~~ → **Implemented**
- ~~**Context Compactor** — 4-layer context compaction pipeline with rule-based pre-compression, deduplication, LLM summarization, and tiered summaries.~~ → **Implemented**
- **Smart Routing** — 8-dimension query classifier that routes simple queries to cheap models and complex ones to powerful ones, cutting LLM costs.
- **SHIELD.md Anti-Malware** — Runtime SHIELD.md enforcement engine with threat parsing, pattern matching, and built-in anti-malware protection.
- **Delegation System** — Autonomous sub-agent orchestration with self-improving role templates, blackboard collaboration, and adaptive timeouts.
- **Inter-Agent Comms** — Lightweight pub/sub event bus for real-time inter-agent communication with wildcard subscriptions and bounded history.

## Changelog

<details>
<summary><strong>2026-03-20 — Context Compactor &amp; documentation flow</strong></summary>

#### Features

- **Context Compactor** (`pkg/compactor/`) — 四层管道（规则预压缩、去重、LLM 摘要、L0/L1/L2 分层）已接入 Agent；配置见 `compactor` 与 [`pkg/compactor/docs/CONFIG.md`](pkg/compactor/docs/CONFIG.md)。

#### Documentation

- 新增仓库文档入口 [`docs/README.md`](docs/README.md)，与 `pkg/learning/docs` 一致区分「包内 docs」与 `docs/implementation/*-status.md`。
- 新增 [`pkg/compactor/docs/`](pkg/compactor/docs/README.md)（README + CONFIG）。
- 修正指向 `plugin-architecture-status.md` 的链接为实际文件 [`docs/implementation/plugin-status.md`](docs/implementation/plugin-status.md)。

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
