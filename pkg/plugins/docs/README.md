# MoonHub 内置插件集（`pkg/plugins`）

本树 **只放具体插件实现**（Channel / Provider / Tool），依赖 [`pkg/framework`](../framework/docs/README.md) 的接口与注册 API。

## 建议阅读顺序（文档流）

1. **本文** — 目录约定与和 framework 的关系  
2. [PLUGIN_INDEX.md](./PLUGIN_INDEX.md) — 当前内置插件清单与路径  
3. [ADDING_A_PLUGIN.md](./ADDING_A_PLUGIN.md) — 新增或复制一类插件时的步骤  

设计决策、阶段完成情况：[`docs/implementation/plugin-status.md`](../../../docs/implementation/plugin-status.md)。

## 目录布局

```
pkg/plugins/
├── channels/<name>/plugin.go   # 各即时通讯 / 通道
├── providers/<name>/plugin.go  # LLM / 协议后端
└── tools/<name>/plugin.go      # 共享工具（如 web、message）
```

每个子包通常包含：

- `init()` 内 `plugin.RegisterPlugin(&XxxPlugin{})`
- 实现对应 `ChannelPlugin` / `ProviderPlugin` / `ToolPlugin`
- 对现有 `pkg/channels/*`、`pkg/providers/*`、`pkg/tools/*` 的薄封装

## 必须被 Gateway  import 才会生效

插件仅依赖 `init()` 副作用注册到全局 Registry。若新增插件包，请在 **`cmd/moonhub/internal/gateway/helpers.go`** 中增加对应的 blank import（与现有 channel / provider / tool 块保持一致）。
