# MoonHub Learning System - Chinese Examples

This document demonstrates usage examples for the MoonHub learning system's Chinese language support.

## Positive Feedback Examples

### Simple Agreement
User: 嗯好的
AI: Identified as positive feedback ✓
User: 行
AI: Identified as positive feedback ✓
User: 可以
AI: Identified as positive feedback ✓
User: 没问题
AI: Identified as positive feedback ✓

### Gratitude Expressions
User: 谢谢
AI: Identified as positive feedback ✓
User: 感谢
AI: Identified as positive feedback ✓
User: 辛苦了
AI: Identified as positive feedback ✓

### Confirmation
User: 对
AI: Identified as positive feedback ✓
User: 是的
AI: Identified as positive feedback ✓
User: 没错
AI: Identified as positive feedback ✓
User: 正是
AI: Identified as positive feedback ✓

### Satisfaction
User: 饭谱
AI: Identified as positive feedback ✓
User: 稳
AI: Identified as positive feedback ✓
User: 到位
AI: Identified as positive feedback ✓
User: 舒服
AI: Identified as positive feedback ✓
User: 顺滑
AI: Identified as positive feedback ✓

### Colloquial
User: 太棒了
AI: Identified as positive feedback ✓
User: 太牛了
AI: Identified as positive feedback ✓
User: 给力
AI: Identified as positive feedback ✓
User: 666
AI: Identified as positive feedback ✓
User: 绝了
AI: Identified as positive feedback ✓

### Brief Affirmation
User: 嗯
AI: Identified as positive feedback ✓
User: 嗯嗯
AI: Identified as positive feedback ✓
User: 嗯呐
AI: Identified as positive feedback ✓
User: 昂
AI: Identified as positive feedback ✓

## Negative Feedback Examples

### Polite Rejection
User: 不太对
AI: Identified as negative feedback ⚠
User: 不太行
AI: Identified as negative feedback ⚠
User: 差点意思
AI: Identified as negative feedback ⚠
User: 差点火候
AI: Identified as negative feedback ⚠

### Neutral-to-Negative
User: 一般
AI: Identified as negative feedback ⚠
User: 还行吧
AI: Identified as negative feedback ⚠
User: 马马虎虎
AI: Identified as negative feedback ⚠
User: 凑合
AI: Identified as negative feedback ⚠

### Ineffective Expression
User: 没帮上忙
AI: Identified as negative feedback ⚠
User: 没解决
AI: Identified as negative feedback ⚠
User: 还是有问题
AI: Identified as negative feedback ⚠

### Strong Negation
User: 不好
AI: Identified as negative feedback ⚠
User: 太差
AI: Identified as negative feedback ⚠
User: 垃圾
AI: Identified as negative feedback ⚠

## Correction/Preference Examples

### Clarifying Intent
User: 其实...
AI: Identified as correction signal 🔄
User: 实际上...
AI: Identified as correction signal 🔄
User: 说真的...
AI: Identified as correction signal 🔄

### Correcting Intent
User: 我的意思是...
AI: Identified as correction signal 🔄
User: 我不是说...
AI: Identified as correction signal 🔄
User: 我是说...
AI: Identified as correction signal 🔄

### Expressing Preference
User: 我更喜欢...
AI: Identified as preference signal 💡
User: 我希望...
AI: Identified as preference signal 💡
User: 我想要...
AI: Identified as preference signal 💡
User: 我倾向于...
AI: Identified as preference signal 💡

### Future Instruction
User: 下次记得...
AI: Identified as preference signal 💡
User: 以后记得...
AI: Identified as preference signal 💡
User: 之后...
AI: Identified as preference signal 💡

### Prohibition
User: 不要...
AI: Identified as disable signal 🚫
User: 别...
AI: Identified as disable signal 🚫
User: 不用...
AI: Identified as disable signal 🚫

## Workflow Preference Examples

### Frequency
User: 每次修改代码前都先运行测试
AI: Identified as workflow preference 💾
User: 总是先检查语法
AI: Identified as workflow preference 💾

### Timing
User: 当部署时，自动运行测试
AI: Identified as workflow preference 💾
User: 在提交之前先拉取代码
AI: Identified as workflow preference 💾

### Sequence
User: 先分析需求，再写代码
AI: Identified as workflow preference 💾
User: 先运行测试,再提交
AI: Identified as workflow preference 💯

## Tool Preference Examples

### Replacement
User: 用 pytest 代替 unittest
AI: Identified as tool replacement preference 💡
User: 换成 ESLint
AI: Identified as tool replacement preference 💡

### Disable
User: 不要用 bash
AI: Identified as tool disable preference 🚫
User: 别用 Python 2
AI: Identified as tool disable preference 🚫

### Priority
User: 优先使用 TypeScript
AI: Identified as tool priority preference 💯

## Mixed Chinese-English Examples

User: 太棒了, perfect!
AI: Identified as positive feedback ✅
User: 不对, please fix it
AI: Identified as negative feedback + needs fix ⚠️
User: 好的, 继续用 Python
AI: Identified as positive feedback + preference signal ✅
User: 下次记得 check 边界情况
AI: Identified as workflow preference 💾
