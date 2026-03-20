# Context Compactor 配置

配置结构定义见 **`pkg/config/config.go`** 中的 `CompactorConfig`、`TierBudgetsConfig`；默认值见 **`DefaultCompactorConfig()`**（`pkg/config/defaults.go`）。下列字段名与 JSON / 环境变量与源码一致。

## config.json 示例

```json
{
  "compactor": {
    "enabled": true,
    "trigger_token_percent": 70,
    "keep_recent": 10,
    "tier_budgets": {
      "l0": 200,
      "l1": 1000,
      "l2": 3000
    },
    "dedup_enabled": true,
    "dedup_similarity_threshold": 0.6,
    "strip_emoji": true,
    "remove_duplicate_lines": true,
    "normalize_cjk": true,
    "smart_rule_selection": true,
    "parallel_processing": true,
    "incremental_compaction": true,
    "summarization_model": ""
  }
}
```

`summarization_model` 为空时由上层使用默认对话模型；若需单独指定摘要模型，填入对应 provider 可用的模型名。

## 环境变量（与 struct 标签一致）

| 环境变量 | 说明 |
|----------|------|
| `MOONHUB_COMPACTOR_ENABLED` | 是否启用 |
| `MOONHUB_COMPACTOR_TRIGGER_TOKEN_PERCENT` | 相对上下文窗口的触发百分比 |
| `MOONHUB_COMPACTOR_KEEP_RECENT` | 压缩后保留的最近消息条数 |
| `MOONHUB_COMPACTOR_DEDUP_ENABLED` | 是否启用去重 |
| `MOONHUB_COMPACTOR_DEDUP_SIMILARITY_THRESHOLD` | Jaccard 相似度阈值 |
| `MOONHUB_COMPACTOR_STRIP_EMOJI` | 是否剥离 emoji |
| `MOONHUB_COMPACTOR_REMOVE_DUPLICATE_LINES` | 是否删除重复行（规则层） |
| `MOONHUB_COMPACTOR_NORMALIZE_CJK` | CJK 标点规范化 |
| `MOONHUB_COMPACTOR_SMART_RULE_SELECTION` | 按内容启用规则子集 |
| `MOONHUB_COMPACTOR_PARALLEL_PROCESSING` | 并行处理 |
| `MOONHUB_COMPACTOR_INCREMENTAL_COMPACTION` | 增量压缩 |
| `MOONHUB_COMPACTOR_SUMMARIZATION_MODEL` | 摘要专用模型（可选） |

## 默认值摘要

| 项 | 默认 |
|----|------|
| `enabled` | `true` |
| `trigger_token_percent` | `70` |
| `keep_recent` | `10` |
| `tier_budgets.l0/l1/l2` | `200` / `1000` / `3000` |
| `dedup_similarity_threshold` | `0.6` |
| 其余布尔开关 | 多为 `true`（见 `DefaultCompactorConfig`） |

持久化路径、触发逻辑与循环内调用方式见 [`docs/implementation/compactor-status.md`](../../../docs/implementation/compactor-status.md)。
