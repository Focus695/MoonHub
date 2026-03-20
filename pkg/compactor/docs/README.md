# MoonHub Context Compactor（`pkg/compactor`）

四层上下文压缩管道：规则预压缩 → 去重 → LLM 摘要 → 分层摘要（L0/L1/L2）。与 Agent 会话历史、Token 预算及 SQLite 持久化配合工作。

## 建议阅读顺序（文档流）

1. **本文** — 职责边界与源码地图  
2. [CONFIG.md](./CONFIG.md) — `config.json` / 环境变量与默认值  
3. 实现细节与集成 — [`docs/implementation/compactor-status.md`](../../../docs/implementation/compactor-status.md)（架构表、表结构、`instance.go` / `loop.go` 钩子、测试命令）

## 源码地图（与实现对齐）

```
pkg/compactor/
├── compactor.go   # 主编排、CompactorEngine
├── config.go      # 运行时 Config（由 pkg/config 映射传入）
├── rules.go       # Layer 1：规则预压缩
├── dedup.go       # Layer 2：Shingle + Jaccard 去重
├── tiers.go       # Layer 4：分层摘要与层级选择
├── store.go       # SQLite：分层摘要与压缩状态
└── tokens.go      # Token 估算与辅助
```

测试：`go test ./pkg/compactor/...`（详见实现状态文档中的覆盖说明）。

## 相关代码（集成点）

| 区域 | 路径 | 说明 |
|------|------|------|
| 全局配置类型 | `pkg/config/config.go`（`CompactorConfig`） | 用户可见配置与 env 标签 |
| 默认值 | `pkg/config/defaults.go` | `DefaultCompactorConfig()` |
| Agent 持有与初始化 | `pkg/agent/instance.go` | `Compactor`、`CompactorTriggerTokenPercent` |
| 触发与分层注入 | `pkg/agent/loop.go` | `maybeCompact`、`getTieredSummary` 等 |

## 与 learning 文档风格的对齐

本目录与 [`pkg/learning/docs`](../learning/docs/README.md) 采用相同习惯：**包内 `docs/` 放「怎么配、从哪读代码」**，**仓库级 `docs/implementation/*-status.md` 放「设计、表结构、集成清单与测试」**，便于在后续开发与 PR 中分工维护。
