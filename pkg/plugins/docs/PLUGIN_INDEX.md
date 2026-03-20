# 内置插件索引

路径均为 `pkg/plugins/.../plugin.go`（包名一般为目录名，如 `telegram`、`openai_compat`）。

## Channel（16）

| 目录 | 用途 |
|------|------|
| `channels/telegram` | Telegram |
| `channels/discord` | Discord |
| `channels/slack` | Slack |
| `channels/matrix` | Matrix |
| `channels/feishu` | 飞书 |
| `channels/qq` | QQ |
| `channels/dingtalk` | 钉钉 |
| `channels/line` | LINE |
| `channels/onebot` | OneBot |
| `channels/wecom` | 企业微信 |
| `channels/wecom_app` | 企业微信应用 |
| `channels/wecom_aibot` | 企业微信 AI Bot |
| `channels/pico` | Pico |
| `channels/irc` | IRC |
| `channels/maixcam` | MaixCam |
| `channels/whatsapp` | WhatsApp（桥接） |
| `channels/whatsapp_native` | WhatsApp Native |

## Provider（8）

| 目录 | 说明 |
|------|------|
| `providers/openai_compat` | OpenAI 兼容 HTTP 等多协议 |
| `providers/openai_oauth` | OpenAI OAuth 等场景 |
| `providers/anthropic` | Anthropic |
| `providers/anthropic_messages` | Anthropic Messages |
| `providers/antigravity` | Antigravity |
| `providers/claude_cli` | Claude CLI |
| `providers/codex_cli` | Codex CLI |
| `providers/github_copilot` | GitHub Copilot |

## Tool（2）

| 目录 | 说明 |
|------|------|
| `tools/web` | Web 搜索 / 抓取等 |
| `tools/message` | Message 工具 |

> 若本文件与仓库不同步，以 `pkg/plugins` 下实际目录为准，并同步更新 `helpers.go` 中的 blank import。
