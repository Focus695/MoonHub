# MoonHub Learning System 配置
本文档描述 MoonHub 学习系统的配置选项。
## 如述
学习系统支持多种配置选项来自定义行为。
## 配置结构
在 `~/.moonhub/config.json` 中配置学习系统:
```json
{
  "learning": {
    "language": "auto",
    "db_path": "memory/learning.db",
    "min_confidence": 0.7,
    "max_examples": 10,
    "enable_semantic_detection": true,
    "enable_implicit_signals": true,
    "enable_behavioral_scoring": true,
    "enable_pattern_evolution": false,
    "enable_contradiction_detection": false,
    "enable_suggestions": false,
    "decay_older_than_days": 7,
    "prune_older_than_days": 30,
    "merge_similarity_threshold": 0.8
  }
}
```
## 配置选项说明
### language
- `auto` - 自动检测语言（默认)
- `en` - 仅使用英文模式
- `zh` - 仅使用中文模式
### db_path
学习数据库存储路径，默认为 `memory/learning.db`
### min_confidence
模式置信度阈值，低于此值的模式将被忽略, 默认 0.7
### max_examples
每个模式保留的最大示例数量, 默认 10
### enable_semantic_detection
是否启用语义关键词检测, 默认 true
### enable_implicit_signals
是否从工具使用中推断隐式信号, 默认 true
### enable_behavioral_scoring
是否启用行为评分系统, 默认 true
### enable_pattern_evolution
是否启用模式演化(合并/修剪/衰减), 默认 false
### enable_contradiction_detection
是否检测矛盾的模式, 默认 false
### enable_suggestions
是否生成主动建议, 默认 false
### decay_older_than_days
模式衰减天数, 默认 7 天
### prune_older_than_days
模式修剪天数, 默认 30 天
### merge_similarity_threshold
模式合并相似度阈值, 默认 0.8
## 使用示例
### 中文模式配置
```json
{
  "learning": {
    "language": "zh"
  }
}
```
### 英文模式配置
```json
{
  "learning": {
    "language": "en"
  }
}
```
### 宷式化配置
```json
{
  "learning": {
    "language": "auto",
    "enable_suggestions": true,
    "enable_pattern_evolution": true
  }
}
```
## 注意事项
- 配置更改需要重启服务才能生效
- `db_path` 目录需要写入权限
- 高级功能(演化、矛盾检测)可能影响性能
