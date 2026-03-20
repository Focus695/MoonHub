# MoonHub 插件框架（`pkg/framework`）

本目录是 **插件运行时核心**：类型定义、全局注册表、生命周期管理。实现代码在上一级目录（`../*.go`），包名在 Go 源码中为 `**plugin`**，导入路径为 `**github.com/sipeed/moonhub/pkg/framework**`。

## 建议阅读顺序

1. **本文** — 定位与术语
2. [ARCHITECTURE.md](./ARCHITECTURE.md) — 接口、Registry、Manager、初始化顺序
3. [INTEGRATION.md](./INTEGRATION.md) — 与 Gateway、channels、providers、agent 如何对接

实现状态与阶段说明见仓库级文档：`[docs/implementation/plugin-status.md](../../../docs/implementation/plugin-status.md)`。

## 包名与导入（常见困惑）


| 你看到的                       | 含义                                                         |
| -------------------------- | ---------------------------------------------------------- |
| 目录 `pkg/framework/`        | 物理路径                                                       |
| `package plugin`           | Go 包名，因此代码里写 `plugin.GlobalRegistry()`、`plugin.Metadata{}` |
| `import "…/pkg/framework"` | 标准库式导入，**不提供** `framework.` 前缀                             |


## 本包不负责的事

- **不**在 `init()` 里列举各业务插件；具体 Channel / Provider / Tool 实现在 `[pkg/plugins](../../plugins/docs/README.md)`，由 **Gateway** 用 blank import 拉齐注册。
- `**LoadBuiltin()`** 为兼容保留的 no-op；见 `builtin.go` 注释。

