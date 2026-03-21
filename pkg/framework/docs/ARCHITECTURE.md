# Plugin Framework Architecture

## Core Types (Source)

| File | Content |
|------|---------|
| `types.go` | `Plugin`, `Metadata`, `PluginType`, `RuntimeContext` |
| `channel.go` | `Channel`, `ChannelPlugin` |
| `provider.go` | `ProviderPlugin` (including `CreateProviderFromModelConfig`, `SupportsProtocol`, etc.) |
| `tool.go` | `ToolPlugin` |
| `registry.go` | `Registry`, `RegisterPlugin`, `GlobalRegistry`, type-based retrieval |
| `manager.go` | `Manager`: `Initialize`, `InitializeToolsOnly`, type-based init |
| `builtin.go` | `LoadBuiltin()` no-op + comments explaining import location |

## Registration Model

- Each plugin package calls `plugin.RegisterPlugin(...)` in **`init()`**.
- Global singleton accessed via **`plugin.GlobalRegistry()`**; internally stored by `PluginType` in separate maps.
- **Order**: Gateway startup ensures all plugin `init` functions have executed via imports before reading Registry.

## Manager Behavior Summary

- **`InitializeToolsOnly`**: Only initializes Tool plugins and writes tools to internal `ToolRegistry`. Current Gateway path: tool plugins first, then construct AgentLoop (merge tools), Channels still started by `channels.Manager` based on Registry.
- **`Initialize`**: Tools → Provider plugin instances → Channel plugin instances (full orchestration; call only if needed by entry point).

## RuntimeContext

Injected when creating `Manager`: `Config`, `MessageBus`, optional `MediaStore`. Tool / Provider / Channel `Init(ctx *RuntimeContext)` runs within this context.
