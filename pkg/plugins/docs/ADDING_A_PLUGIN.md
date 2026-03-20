# 新增内置插件（Checklist）

## 通用

1. 在 `pkg/plugins/<channels|providers|tools>/<newname>/` 新增 `plugin.go`。  
2. **`init()`** 中调用 `plugin.RegisterPlugin(&YourPlugin{})`（导入 `github.com/sipeed/moonhub/pkg/framework`，代码前缀使用包名 **`plugin`**）。  
3. 实现 **`Metadata()`**、`Init(*plugin.RuntimeContext)`、`Validate(*config.Config)` 及对应子接口的全部方法。  
4. 在 **`cmd/moonhub/internal/gateway/helpers.go`** 增加一行 `_ "github.com/sipeed/moonhub/pkg/plugins/.../newname"`，与同类插件放在一起。  
5. 运行 `go build ./...`，必要时补测试。

## Channel 插件

- 实现 `plugin.ChannelPlugin`：`ChannelPrefix`、`CreateChannel`、`IsEnabled`。  
- `CreateChannel` 通常委托 `pkg/channels/<impl>` 已有构造函数。  
- 确认 `channels.Manager` 对返回类型与 `channels.Channel`、可选的 `SetMediaStore` / `SetPlaceholderRecorder` / `SetOwner` 仍兼容。

## Provider 插件

- 实现 `ProviderPlugin`，特别注意：  
  - **`SupportsProtocol`** 与模型配置里的协议前缀一致。  
  - **`CreateProviderFromModelConfig`**：不适用当前配置时返回 **`providers.ErrSkipProvider`**，便于 resolver 尝试下一个插件或回退内置工厂。  
- 可复用 `pkg/providers` 中已导出的工厂函数（如 OAuth 相关）。

## Tool 插件

- 实现 `ToolPlugin`：`CreateTools`、`IsCore`（决定是否对 LLM **可见**注册）。  
- 若工具名与 Agent 内置 `registerSharedTools` 重复（如 `web_search`、`web_fetch`、`message`），应确保 Gateway 已用插件 Registry 合并，避免重复注册；参考 `pkg/agent/loop.go` 中的跳过逻辑。

## 文档

- 更新 [`PLUGIN_INDEX.md`](./PLUGIN_INDEX.md) 与本目录或仓库级 `docs/implementation/plugin-status.md` 中的清单（若项目要求保持同步）。
