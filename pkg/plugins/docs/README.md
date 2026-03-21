# MoonHub Built-in Plugins (`pkg/plugins`)

This tree contains **only specific plugin implementations** (Channel / Provider / Tool), depending on interfaces and registration APIs from [`pkg/framework`](../framework/docs/README.md).

## Recommended Reading Order (Documentation Flow)

1. **This document** — Directory conventions and relationship with framework
2. [PLUGIN_INDEX.md](./PLUGIN_INDEX.md) — Current built-in plugin list and paths
3. [ADDING_A_PLUGIN.md](./ADDING_A_PLUGIN.md) — Steps when adding or copying a type of plugin

Design decisions and phase completion: [`docs/implementation/plugin-status.md`](../../../docs/implementation/plugin-status.md).

## Directory Layout

```
pkg/plugins/
├── channels/<name>/plugin.go   # IM / channel plugins
├── providers/<name>/plugin.go  # LLM / protocol backends
└── tools/<name>/plugin.go      # Shared tools (e.g., web, message)
```

Each subpackage typically contains:

- `init()` with `plugin.RegisterPlugin(&XxxPlugin{})`
- Implements corresponding `ChannelPlugin` / `ProviderPlugin` / `ToolPlugin`
- Thin wrapper around existing `pkg/channels/*`, `pkg/providers/*`, `pkg/tools/*`

## Must Be Imported by Gateway to Take Effect

Plugins only register to global Registry via `init()` side effects. When adding new plugin packages, add corresponding blank imports in **`cmd/moonhub/internal/gateway/helpers.go`** (consistent with existing channel / provider / tool blocks).
