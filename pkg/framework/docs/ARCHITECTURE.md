# 插件框架结构

## 核心类型（源码）

| 文件 | 内容 |
|------|------|
| `types.go` | `Plugin`、`Metadata`、`PluginType`、`RuntimeContext` |
| `channel.go` | `Channel`、`ChannelPlugin` |
| `provider.go` | `ProviderPlugin`（含 `CreateProviderFromModelConfig`、`SupportsProtocol` 等） |
| `tool.go` | `ToolPlugin` |
| `registry.go` | `Registry`、`RegisterPlugin`、`GlobalRegistry`、按类型检索 |
| `manager.go` | `Manager`：`Initialize`、`InitializeToolsOnly`、按类型 init |
| `builtin.go` | `LoadBuiltin()` no-op + 注释说明导入位置 |

## 注册模型

- 各插件包在 **`init()`** 中调用 `plugin.RegisterPlugin(...)`。
- 全局单例通过 **`plugin.GlobalRegistry()`** 访问；内部按 `PluginType` 分 map 存储。
- **顺序**：Gateway 启动时通过 import 侧保证各插件 `init` 已执行，再读 Registry。

## Manager 行为摘要

- **`InitializeToolsOnly`**：仅初始化 Tool 插件并把工具写入内部 `ToolRegistry`。Gateway 当前路径：先工具插件，再构造 AgentLoop（合并工具），Channel 仍由 `channels.Manager` 基于 Registry 启动。
- **`Initialize`**：工具 → Provider 插件实例 → Channel 插件实例（完整编排；若某入口未使用则不必调用）。

## RuntimeContext

创建 `Manager` 时注入：`Config`、`MessageBus`、可选 `MediaStore`。Tool / Provider / Channel 的 `Init(ctx *RuntimeContext)` 在此上下文中运行。
