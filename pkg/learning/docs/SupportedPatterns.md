# MoonHub Learning System - Supported Chinese Patterns

This document lists all Chinese patterns supported by the learning system.

## Overview

## Positive Feedback Patterns

### Patterns

| Type | Regex | Examples |
|------|-------|----------|
| Simple Agreement | `^(好的?|行|可以|没问题|OK|ok)$` | `好的`, `行`, `没问题` |
| Gratitude | `^(谢谢|感谢|多谢|谢了|辛苦了)$` | `谢谢`, `感谢`, `多谢`, `谢了` |
| Colloquial Praise | `^(很好?|棒|赞|牛|给力|厉害|牛逼|666)$` | `很好`, `棒`, `赞`, `牛`, `给力`, `厉害`, `牛逼`, `666` |
| Strong Affirmation | `^(完美|太棒了?|太好了?|太赞了?|绝了)$` | `完美`, `太棒了`, `太好了`, `太赞了`, `绝了` |
| Brief Affirmation | `^(嗯|嗯嗯|嗯呐|昂|噢了)$` | `嗯`, `嗯嗯`, `嗯呐`, `昂`, `噢了` |
| Correct Confirmation | `^(对|是的|没错|正确|正是|就是)$` | `对`, `是`, `没错`, `正确`, `正是` |
| Satisfaction | `^(靠谱|稳|到位|舒服|顺滑)$` | `靠谱`, `稳`, `到位`, `舒服`, `顺滑` |
| Choice Confirmation | `^(就这个|这个对|这个好|要这个)$` | `就这个`, `这个对`, `这个好`, `要这个` |
| Convenience | `^(省事|方便|快捷|高效)$` | `省事`, `方便`, `快捷`, `高效` |
| Professionalism | `^(靠谱|稳|专业|到位)$` | `靠谱`, `稳`, `专业`, `到位` |

## Negative Feedback Patterns

| Type | Regex | Examples |
|------|-------|----------|
| Polite Rejection | `^(不太对|不太行|差点意思|差点火候)$` | `不太对`, `不太行`, `差点意思`, `差点火候` |
| Neutral-to-Negative | `^(一般|还行吧|马马虎虎|凑合)$` | `一般`, `还凑合`, `马马虎虎`, `凑合` |
| Ineffective | `^(没帮上忙?|没解决|还是有问题)$` | `没帮上忙`, `没解决`, `还是有问题` |
| Frustration | `^(浪费(时间|没意思|无语)$` | `浪费时间`, `没意思`, `无语` |
| Strong Negation | `^(不好|太差|差劲|垃圾|废)$` | `不好`, `太差`, `差劲`, `垃圾`, `废` |
| Retry Request | `^(再试(一次)?|重来|重新来|换个)$` | `再试一次`, `重来`, `重新来`, `换个` |
| Error Identification | `^(错误|不正确|有误|有问题)$` | `错误`, `不正确`, `有误`, `有问题` |
| Misunderstanding | `^(不是这个|不是我要的|理解错了?|会错意了?)$` | `不是这个`, `不是我要的`, `理解错了`, `会错意了` |

## Correction/Preference Patterns

| Type | Regex | Examples |
|------|-------|----------|
| Clarification | `^(其实|实际上|说真的|老实说)$` | `其实`, `实际上`, `说真的`, `老实说` |
| Intent Correction | `^(我的意思(是|不是说)|我是说)$` | `我的意思是`, `我的意思不是说`, `我是说` |
| Preference Expression | `^(我更(喜欢|希望|想要|倾向于))$` | `我更喜欢`, `我希望`, `我想要`, `我倾向于` |
| Future Instruction | `^(下次|以后|之后|以后记得)$` | `下次`, `以后`, `之后`, `以后记得` |
| Prohibition | `^(不要|别|不用|不需要)$` | `不要`, `别`, `不用`, `不需要` |
| Reverse Instruction | `^(相反|反过来|倒过来|反过来)$` | `相反`, `反过来`, `倒过来` |
| Memory Request | `^(记住(我|我喜欢|我想要)|记得)$` | `记住我`, `记住我喜欢`, `记住我想要`, `记得` |
| Expectation | `^(最好是|最好是能|希望)$` | `最好是`, `最好是能`, `希望` |
| Request | `^(能不能|可以不可以|麻烦)$` | `能不能`, `可以不可以`, `麻烦` |
| Modification Request | `^(改(成|为)|换成|修改成)$` | `改成`, `换成`, `修改成` |
| Prompt/Suggestion | `^(你(应该|最好|能不能))$` | `你应该`, `你最好`, `你能不能` |

## Workflow Preference Patterns

| Type | Regex | Examples |
|------|-------|----------|
| Frequency | `(每次|总是|一直|老)` | Check before each code modification, Always check config before running tests, Always prefer Python, Keeps forgetting to check config |
| Negative Frequency | `(从不|从没|从来不)` | Never deploy to production, Never done code review |
| Timing | `(当.+时|在.+之前|在.+之后)` | Run tests when deploying, Check before committing, Clean up after modification |
| Habit | `(喜欢用|习惯用|常用|爱用)` | Like using Python for scripts, Used to debugging with IDE, Often use Docker for development |
| Sequence | `(先.+再|先.+然后|先.+接着)` | Analyze requirements before coding, Test before deploying |
| Mandatory | `(必须|一定要|务必)` | Must pass tests before commit, Must backup code, Must follow coding standards |

## Tool Preference Patterns

| Type | Regex | Examples |
|------|-------|----------|
| Replacement | `(用.+代替|换成|改用)` | Use pytest instead of unittest, Switch to VSCode from Vim, Switch to Docker deployment |
| Prohibition | `(不要用|别用|不用)` | Don't use bash scripts, Don't use XML config, No manual deployment |
| Mandatory | `(总是用|一直用|每次都用)` | Always use TypeScript, Always use Python for development, Use Git every time |
| Exclusion | `(从不用|从不使用|从来不用)` | Never use manual testing, Never use jQuery, Never use FTP |
| Priority | `(优先(使用|选|用))` | Prefer Docker, Prefer Go language, Prefer cloud services |
