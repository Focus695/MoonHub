# MoonHub Learning System - Chinese Language Support

The MoonHub learning system now supports full Chinese feedback recognition and learning.

**Repository Documentation Index** (reading order for subsystems): [docs/README.md](../../../docs/README.md).

## Overview

This document describes the Chinese language adaptation features of the MoonHub self-improving learning system.

## Core Components

The system implements Chinese language support through the following components:
- **i18n module**: Internationalization manager
- **patterns module**: Pattern manager
- **detector**: Pattern detector
- **suggester**: Suggestion generator
- **engine**: Context builder

## Configuration

Configure the language in `~/.moonhub/config.json`:
```json
{
  "learning": {
    "language": "auto"  // "en", "zh", "auto"
  }
}
```

## Supported Expressions

### Positive Feedback (Chinese)
- "太棒了", "给力", "完美", "靠谱"
- "牛", "666", "绝了", "稳"
- "好的", "行", "可以", "没问题"
- "谢谢", "感谢", "辛苦了"
- "对", "是的", "没错", "正确"

### Negative Feedback (Chinese)
- "不对", "错了", "不行", "不可以"
- "不太对", "不太行", "差点意思"
- "没帮助", "没用", "无效"
- "浪费", "没意思", "无语"

### Corrections/Preferences (Chinese)
- "其实...", "我的意思是...", "我更喜欢..."
- "下次记得...", "不要用..."
- "用...代替...", "优先使用..."

### Workflow (Chinese)
- "每次...", "总是...", "从不..."
- "先...再...", "当...时"
- "必须...", "务必..."

## Examples

User says "太棒了，正是我想要的" → System identifies as positive feedback
User says "不对,理解错了,我的意思是..." → System identifies as correction signal
User says "下次修改代码前先运行测试" → System records workflow preference
User says "不要用 bash,用 Python" → System records tool disable preference
