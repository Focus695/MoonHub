# MoonHub 文档索引

本页是仓库文档的**入口与阅读顺序**

## 建议阅读顺序（第一次接触本仓库）

1. [仓库根 README](../README.md) — 功能概览、许可证与近期变更摘要
2. [故障排查](./troubleshooting.md)、[调试说明](./debug.md) — 运行期问题按需查阅
3. [工具与能力配置](./tools_configuration.md) — 工具侧配置

## 按子系统的文档流

### 插件架构（Channel / Provider / Tool）


| 顺序  | 文档                                                                              | 说明                      |
| --- | ------------------------------------------------------------------------------- | ----------------------- |
| 1   | `[pkg/framework/docs/README.md](../pkg/framework/docs/README.md)`               | 框架定位、包名与导入约定            |
| 2   | `[pkg/framework/docs/ARCHITECTURE.md](../pkg/framework/docs/ARCHITECTURE.md)`   | 接口、Registry、生命周期        |
| 3   | `[pkg/framework/docs/INTEGRATION.md](../pkg/framework/docs/INTEGRATION.md)`     | 与 Gateway、agent、各子系统的对接 |
| 4   | `[pkg/plugins/docs/README.md](../pkg/plugins/docs/README.md)`                   | 内置插件目录约定                |
| 5   | `[pkg/plugins/docs/PLUGIN_INDEX.md](../pkg/plugins/docs/PLUGIN_INDEX.md)`       | 插件清单与路径                 |
| 6   | `[pkg/plugins/docs/ADDING_A_PLUGIN.md](../pkg/plugins/docs/ADDING_A_PLUGIN.md)` | 新增或迁移插件的步骤              |
| 状态  | `[docs/implementation/plugin-status.md](./implementation/plugin-status.md)`     | 设计决策与实现进度               |


### 学习系统（Self-Improving）


| 顺序  | 文档                                                                                                                                      | 说明             |
| --- | --------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| 1   | `[pkg/learning/docs/README.md](../pkg/learning/docs/README.md)`                                                                         | 中文反馈与能力概述      |
| 2   | `[pkg/learning/docs/CONFIG.md](../pkg/learning/docs/CONFIG.md)`                                                                         | `learning` 配置项 |
| 3   | `[pkg/learning/docs/I18N.md](../pkg/learning/docs/I18N.md)`                                                                             | 语言与检测行为        |
| 4   | `[pkg/learning/docs/EXAMPLES.md](../pkg/learning/docs/EXAMPLES.md)`、`[SupportedPatterns.md](../pkg/learning/docs/SupportedPatterns.md)` | 示例与模式表         |
| 状态  | `[docs/implementation/learning-status.md](./implementation/learning-status.md)`                                                         | 实现状态与集成点       |


### 上下文压缩（Context Compactor）


| 顺序  | 文档                                                                                | 说明                  |
| --- | --------------------------------------------------------------------------------- | ------------------- |
| 1   | `[pkg/compactor/docs/README.md](../pkg/compactor/docs/README.md)`                 | 四层管道与代码入口           |
| 2   | `[pkg/compactor/docs/CONFIG.md](../pkg/compactor/docs/CONFIG.md)`                 | `compactor` 配置与环境变量 |
| 状态  | `[docs/implementation/compactor-status.md](./implementation/compactor-status.md)` | 架构、存储、Agent 集成与测试命令 |


### 自适应记忆（Adaptive Memory）


| 文档                                                                          | 说明                                        |
| --------------------------------------------------------------------------- | ----------------------------------------- |
| `[docs/implementation/memory-status.md](./implementation/memory-status.md)` | 包结构、能力、配置与集成说明；源码在 `pkg/adaptive_memory/` |


### 各通道（Channels）

`[docs/channels/](./channels/)` 下为各通道的说明（含部分中文 README）。

---

## `docs/implementation/` 一览


| 文件                                                            | 主题         |
| ------------------------------------------------------------- | ---------- |
| `[plugin-status.md](./implementation/plugin-status.md)`       | 插件架构实现状态   |
| `[learning-status.md](./implementation/learning-status.md)`   | 学习系统实现状态   |
| `[compactor-status.md](./implementation/compactor-status.md)` | 上下文压缩器实现状态 |
| `[memory-status.md](./implementation/memory-status.md)`       | 记忆系统实现状态   |


