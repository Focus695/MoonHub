# 与 MoonHub 主程序的集成

## 入口：`cmd/moonhub/internal/gateway/helpers.go`

1. **Provider**
  - 启动早期调用 `providers.SetPluginProviderResolver(...)`。  
  - Resolver 遍历 `plugin.GlobalRegistry().GetProviderPlugins()`，按 `SupportsProtocol` + `CreateProviderFromModelConfig` 尝试创建；`providers.ErrSkipProvider` 表示交给下一个插件或回退内置工厂。  
  - 详见 `pkg/providers/factory_provider.go`：`CreateProviderFromConfig` 在 resolver 报错时**仅**在 `ErrSkipProvider` 时回退内置实现。
2. **Tool**
  - `plugin.NewManager(cfg, msgBus, nil)` → `InitializeToolsOnly(ctx)`。  
  - `agent.NewAgentLoopWithPluginTools(..., pluginMgr.ToolRegistry())`：先把插件工具 `MergeFrom` 进各 Agent，再在 `registerSharedTools` 里跳过已由插件提供的 `web_search` / `web_fetch` / `message`。
3. **Channel**
  - `channels.Manager` 从 `plugin.GlobalRegistry().GetChannelPlugins()` 取插件，对 `IsEnabled` 为真的调用 `CreateChannel`，结果断言为 `channels.Channel` 后注入 `MediaStore` / `PlaceholderRecorder` / `Owner` 等。

## 循环依赖的处理

插件实现会依赖 `pkg/channels/*` 等包，若由 `pkg/framework` 再 import 插件会形成环。因此：

- `**builtin.go` 不 import 任何 `pkg/plugins/...`**。  
- **Gateway** 通过 `_ "…/pkg/plugins/..."` 触发各包 `init()` 注册。

## 相关代码索引


| Concern   | 主要位置                                                                      |
| --------- | ------------------------------------------------------------------------- |
| 插件注册表     | `pkg/framework/registry.go`                                               |
| Tool 合并   | `pkg/tools/registry.go` — `MergeFrom`                                     |
| Agent 侧   | `pkg/agent/loop.go` — `NewAgentLoopWithPluginTools`、`registerSharedTools` |
| Channel 侧 | `pkg/channels/manager.go` — 插件迭代初始化                                       |


