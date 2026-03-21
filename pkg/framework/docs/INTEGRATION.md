# Integration with MoonHub Main Program

## Entry Point: `cmd/moonhub/internal/gateway/helpers.go`

1. **Provider**
   - Early startup calls `providers.SetPluginProviderResolver(...)`.
   - Resolver iterates `plugin.GlobalRegistry().GetProviderPlugins()`, tries creation via `SupportsProtocol` + `CreateProviderFromModelConfig`; `providers.ErrSkipProvider` means pass to next plugin or fallback to built-in factory.
   - See `pkg/providers/factory_provider.go`: `CreateProviderFromConfig` falls back to built-in implementation **only** when resolver returns `ErrSkipProvider`.

2. **Tool**
   - `plugin.NewManager(cfg, msgBus, nil)` → `InitializeToolsOnly(ctx)`.
   - `agent.NewAgentLoopWithPluginTools(..., pluginMgr.ToolRegistry())`: first `MergeFrom` plugin tools into each Agent, then skip `web_search` / `web_fetch` / `message` in `registerSharedTools` if already provided by plugin.

3. **Channel**
   - `channels.Manager` retrieves plugins from `plugin.GlobalRegistry().GetChannelPlugins()`, calls `CreateChannel` for those where `IsEnabled` is true, asserts result as `channels.Channel`, then injects `MediaStore` / `PlaceholderRecorder` / `Owner`, etc.

## Handling Circular Dependencies

Plugin implementations depend on `pkg/channels/*` etc.; if `pkg/framework` imports plugins, it would form a cycle. Therefore:

- `**builtin.go` does NOT import any `pkg/plugins/...`**.
- **Gateway** triggers each package's `init()` registration via `_ ".../pkg/plugins/..."`.

## Related Code Index

| Concern | Main Location |
|---------|---------------|
| Plugin Registry | `pkg/framework/registry.go` |
| Tool Merge | `pkg/tools/registry.go` — `MergeFrom` |
| Agent Side | `pkg/agent/loop.go` — `NewAgentLoopWithPluginTools`, `registerSharedTools` |
| Channel Side | `pkg/channels/manager.go` — plugin iteration initialization |
