# MoonHub

> [!NOTE]
> **Based On**
>
> This project is based on features from [TinyClaw](https://github.com/wgtechlabs/tinyclaw) and the lightweight design of [PicoClaw](https://github.com/sipeed/picoclaw), with unique features developed for this project. Licensed under GPL-3.0.

**Documentation Index** (plugins, learning, compactor, SHIELD, memory, delegation, etc.): [`docs/README.md`](docs/README.md).

**[English](README.md)**

## 简介

> **A ready-to-use agent that's perfect for everyday users.**

MoonHub 是一款面向边缘计算的**本地优先 AI 助手**——简单、快速、安全。

我们相信，AI 不应该是少数技术专家的专属工具。MoonHub 专为**日常用户**打造，无需任何技术背景，插电即用。目前适配 Linux 平台，完美运行于树莓派、工业网关等各类嵌入式设备。

### 设计理念

| 原则 | 描述 |
|------|------|
| **Simple** | 零学习曲线。开箱即用，像使用家用电器一样简单。 |
| **Fast** | 极致轻量。<10MB 内存，1 秒冷启动，毫秒级响应。 |
| **Secure** | 本地优先。数据不出设备，隐私完全由你掌控。 |

### 核心特性

**🤖 Agent 协作引擎**

你是指挥官，Agent 团队为你服务。每个 Agent 各司其职，它们可以彼此通信、主动协作，在关键决策点向你请示。这不是简单的问答机器人，而是一个真正理解上下文、能自主推进任务的智能团队。

**🎨 动态 UI 生成**

告别"只能输出文字"的传统 Agent。MoonHub 能够根据你的需求，实时生成可视化交互界面——财务仪表盘、任务管理器、数据看板，一切随需而变。你描述想法，Agent 为你构建。

**📱 专属应用**

用户通过专属应用与设备交互。目前以 **PWA** 形式提供，支持快速安装、离线使用；后续将推出原生 **APP**，覆盖更多平台与使用场景。

### 无限可能

MoonHub 可运行于任何边缘设备，为你的场景注入智能。基于本项目二次开发，打造专属解决方案：

| 场景 | 描述 |
|------|------|
| **智能灌溉** | 接入土壤湿度、气象传感器，Agent 根据实时数据动态调整灌溉策略，实现精准节水农业。 |
| **工业监控** | 部署于生产车间，实时采集设备状态，预测性维护告警，生成可视化运维看板。 |
| **智慧门店** | 连接客流统计、库存传感器，自动生成补货建议、销售分析报告，辅助经营决策。 |
| **能源管理** | 对接智能电表、光伏逆变器，实时优化用电策略，生成能耗报告与节能建议。 |
| **智能 CRM** | 集成客户数据、沟通记录，AI 分析客户画像，自动生成跟进提醒与销售机会洞察。 |
| **智能运维** | 接入服务器、应用监控数据，AI 识别异常模式，自动告警并生成故障诊断报告。 |
| **智能安防** | 对接摄像头、门窗传感器，AI 识别异常行为，实时推送告警并生成安全日志。 |
| **智慧养殖** | 接入水质、投喂设备，实时监测养殖环境，自动调节投喂量并生成生长分析报告。 |
| **智慧教室** | 接入考勤设备、互动大屏，自动记录出勤，辅助教师生成个性化学习报告与教学建议。 |
| **智能电商** | 对接订单、库存、物流系统，AI 分析销售趋势，自动生成补货建议与营销策略。 |

你的想象力，就是 MoonHub 的边界。

### 快速开始

1. **插电启动** — 设备开机后自动创建 WiFi 热点（`MoonHub-XXXX`）
2. **手机配网** — 连接热点，访问配网页面，配置 WiFi 并设置授权码
3. **安装 PWA** — 配网完成后引导安装 PWA 应用
4. **开始使用** — PWA 自动扫描本地设备，输入授权码即可开始对话、管理和配置

## Features

### Implemented

- **Adaptive Memory** — 3-layer memory system (episodic, semantic FTS5, temporal decay) that learns what to remember and forget over time.
- **Self-Improving** — Behavioral pattern detection system that learns from user feedback, tracks tool usage preferences, and evolves patterns over time.
- **Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.
- **Context Compactor** — 4-layer context compaction pipeline (rules, dedup, LLM summary, L0/L1/L2 tiers) integrated in the agent loop; see [`docs/implementation/compactor-status.md`](docs/implementation/compactor-status.md) and [`pkg/compactor/docs/`](pkg/compactor/docs/README.md).
- **SHIELD.md Anti-Malware** — Runtime threat evaluation engine with YAML threat parsing, pattern matching, approval workflow, and 8 built-in threats; see [`docs/implementation/shield-status.md`](docs/implementation/shield-status.md) and [`pkg/shield/docs/`](pkg/shield/docs/README.md).
- **Delegation System** — Sub-agent orchestration (non-blocking and background tasks, template reuse, adaptive timeouts, SQLite persistence, Intercom pub/sub). Opt-in via `delegation.enabled` in `config.json` (default off); see [`pkg/delegation/docs/README.md`](pkg/delegation/docs/README.md), [`pkg/delegation/docs/CONFIG.md`](pkg/delegation/docs/CONFIG.md), and [`docs/implementation/delegation-status.md`](docs/implementation/delegation-status.md).
- **Inter-Agent Comms (Intercom)** — In-process pub/sub for delegation-time signals: subscribe per topic (`On`), catch-all via `OnAny`, bounded per-topic retention with `Recent` / `RecentAll`. TinyClaw-compatible topic constants; see [`pkg/delegation/intercom.go`](pkg/delegation/intercom.go) and the Intercom section in [`docs/implementation/delegation-status.md`](docs/implementation/delegation-status.md).
- **Smart Router V2** — 4-tier model routing system (simple/moderate/complex/reasoning) with rule-based scoring, feature extraction, and privacy-safe metrics. Routes simple queries to cheap models and complex ones to powerful models; see [`pkg/routing/docs/README.md`](pkg/routing/docs/README.md) and [`docs/implementation/routing-status.md`](docs/implementation/routing-status.md).

### Planned

- **Device Provisioning** — Zero-config device WiFi setup with hotspot, network management, SSE events, and recovery mechanisms (backend API layer; see [`docs/implementation/provisioning-status.md`](docs/implementation/provisioning-status.md))
- **PWA Frontend** — User-facing PWA application for device discovery, pairing, and interaction with the MoonHub device
- **Dynamic UI Generation** — Real-time visual component generation based on user needs (dashboards, task managers, data visualizations)
- **Native APP** — Native mobile applications for iOS and Android platforms
- **Factory Reset** — Device factory reset functionality (interface already reserved)
- **Enhanced Network Recovery** — Optional automatic network recovery mechanism optimization

<details>
<summary><strong>Completed (click to expand)</strong></summary>

- ~~**Self-Improving** — Behavioral pattern detection that makes the agent better with every interaction. It grows with you.~~ → **Implemented**
- ~~**Plugin Architecture** — Channels, providers, and tools are all plugins. The core stays tiny — everything else is extensible.~~ → **Implemented**
- ~~**Context Compactor** — 4-layer context compaction pipeline with rule-based pre-compression, deduplication, LLM summarization, and tiered summaries.~~ → **Implemented**
- ~~**SHIELD.md Anti-Malware** — Runtime SHIELD.md enforcement engine with threat parsing, pattern matching, and built-in anti-malware protection.~~ → **Implemented**
- ~~**Delegation System** — Autonomous sub-agent orchestration with self-improving role templates, blackboard collaboration, and adaptive timeouts.~~ → **Implemented** (opt-in; see [`pkg/delegation/docs/`](pkg/delegation/docs/README.md) and [`docs/implementation/delegation-status.md`](docs/implementation/delegation-status.md))
- ~~**Smart Routing** — 4-tier query classifier that routes simple queries to cheap models and complex ones to powerful ones, cutting LLM costs.~~ → **Implemented** (see [`pkg/routing/docs/`](pkg/routing/docs/README.md) and [`docs/implementation/routing-status.md`](docs/implementation/routing-status.md))
- ~~**Inter-Agent Comms** — Lightweight pub/sub event bus for real-time inter-agent communication with wildcard subscriptions and bounded history.~~ → **Implemented** (delegation **Intercom** in [`pkg/delegation/intercom.go`](pkg/delegation/intercom.go); enabled with delegation)

</details>

## Changelog

<details>
<summary><strong>2026-03-21 — Smart Router V2 (4-Tier Model Routing)</strong></summary>

#### Summary

4-tier model routing system (TinyClaw-style) with rule-based classification, replacing the original 2-tier (light/heavy) system. Automatically selects the appropriate LLM based on message complexity.

#### New Features

**Smart Router V2** (`pkg/routing/`)
- **4-Tier Classification** — simple, moderate, complex, reasoning tiers with configurable boundaries
- **Rule-Based Scoring** — Sub-microsecond classification using structural features (no API calls)
- **Feature Extraction** — Token estimate, code blocks, tool calls, conversation depth, attachments
- **Attachment Hard Gate** — Multi-modal inputs automatically route to reasoning tier
- **Confidence Scoring** — Sigmoid-based confidence calculation for each classification
- **Signal Tracing** — Every decision includes explainable signals for debugging
- **Privacy-Safe Metrics** — Aggregate statistics without storing message content
- **Decision Recorder** — Ring buffer for recent decisions with tier filtering
- **HTTP Endpoints** — `/metrics`, `/routing/decisions`, `/routing/stats`
- **Backward Compatible** — Falls back to 2-tier mode when `light_model` is configured

#### Documentation

- [`pkg/routing/docs/README.md`](pkg/routing/docs/README.md) — Overview, architecture, quick start
- [`pkg/routing/docs/CONFIG.md`](pkg/routing/docs/CONFIG.md) — Configuration options, tier mapping, custom boundaries
- [`pkg/routing/docs/FEATURES.md`](pkg/routing/docs/FEATURES.md) — Feature extraction, scoring weights, examples
- [`pkg/routing/docs/METRICS.md`](pkg/routing/docs/METRICS.md) — Metrics collection, decision recorder, HTTP endpoints
- [`docs/implementation/routing-status.md`](docs/implementation/routing-status.md) — Full implementation status
- [`docs/README.md`](docs/README.md) — Repository documentation index (updated)

#### Technical Details

- Score range: [-1.0, 1.0] with negative scores for simple messages
- Default boundaries: simple [-1, -0.05), moderate [-0.05, 0.15), complex [0.15, 0.35), reasoning [0.35, 1.0]
- Weights: short message (-0.10), code block (+0.40), long message (+0.35), attachment (1.0 hard gate)
- 42 unit tests, all passing

#### Files Changed

- `pkg/routing/` — Core implementation (tier, classifier, router, features, metrics, recorder)
- `pkg/routing/docs/` — New documentation directory (README, CONFIG, FEATURES, METRICS)
- `pkg/config/config.go` — RoutingConfig with TierMapping and TierBoundariesConfig
- `pkg/agent/instance.go` — RouterV2, TierCandidates fields and initialization
- `pkg/agent/loop.go` — Updated selectCandidates for 4-tier routing
- `pkg/health/server.go` — HTTP endpoints for metrics and decisions
- `docs/README.md` — Updated Smart Router section with full documentation links

</details>

<details>
<summary><strong>2026-03-21 — Delegation System (sub-agent orchestration)</strong></summary>

#### Summary

Sub-agent delegation aligned with TinyClaw-style workflows: eight tools, SQLite store, session queue, background task injection in the agent loop, and `DelegationUserID` on tool execution context.

#### Documentation

- [`pkg/delegation/docs/README.md`](pkg/delegation/docs/README.md), [`pkg/delegation/docs/CONFIG.md`](pkg/delegation/docs/CONFIG.md) — Package-level docs (flow + config)
- [`docs/implementation/delegation-status.md`](docs/implementation/delegation-status.md) — Deep implementation reference
- [`docs/README.md`](docs/README.md) — Repository index (delegation subsystem table)

#### Code (high level)

- `pkg/delegation/` — Core implementation
- `pkg/agent/delegation_integration.go`, `pkg/agent/instance.go`, `pkg/agent/loop.go` — Runtime wiring
- `pkg/config/config.go`, `pkg/config/defaults.go` — `DelegationConfig`

</details>

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
