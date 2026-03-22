# 设备配网系统 - 实现状态

**源码包**：`pkg/provisioning/` · **仓库文档索引**：[docs/README.md](../README.md)

## 实现概述

MoonHub 设备配网系统已成功从 TinyClaw 移植，实现了"插电即用"的零配置设备体验。用户后续可用 PWA 实现 UI 层，本次只实现后端 API 层。

## 已完成的工作

### 1. 核心包结构

```
pkg/provisioning/
├── config.go          # 类型定义和配置 Key
├── events.go          # SSE 事件广播
├── network.go         # nmcli 网络操作封装
├── manager.go         # DeviceManager 核心实现
└── manager_test.go    # 单元测试
```

### 2. 核心功能

#### 2.1 设备模式状态机
- 六种设备模式：
  - `provisioning`: 等待初始 WiFi 配置（热点激活）
  - `connecting`: 正在连接配置的 WiFi
  - `onboarding`: 网络已连接，等待用户完成引导
  - `ready`: 设备完全就绪
  - `maintenance`: 维护模式
  - `error`: 遇到不可恢复的错误
- 运行时阶段：provisioning, connecting, onboarding, ready, degraded, maintenance, restarting, error

#### 2.2 WiFi 管理
- WiFi 网络扫描（nmcli）
- WiFi 连接（支持密码、隐藏网络）
- 已保存网络列表
- 网络遗忘功能
- 自动重连机制

#### 2.3 热点管理
- 自动创建热点配置文件
- 热点启用/禁用
- 默认热点 SSID：`MoonHub-XXXX`（XXXX 为设备名后4位）

#### 2.4 网络诊断
- 连接状态检测
- 互联网可达性测试
- DNS 解析测试
- 网络详细信息获取

#### 2.5 恢复机制
- 可配置的故障阈值
- 冷却时间机制
- 网络退化检测
- 自动回滚到配网模式

#### 2.6 SSE 实时事件推送
- 事件类型：status, action, system
- 事件级别：info, success, warning, error
- 动作阶段：started, completed, failed
- 最近事件缓存（25条）

### 3. 集成点

#### 3.1 web/backend/api/provisioning.go
- HTTP API 处理器
- SSE 事件流处理
- 请求验证和响应格式化

#### 3.2 web/backend/api/router.go
- 添加 `provisioning` 字段
- 添加 `SetProvisioningHandler()` 方法
- 自动注册配网路由

#### 3.3 web/backend/main.go
- 环境变量控制启用/禁用
- 设备管理器初始化
- 后台监控循环启动

### 4. API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/provisioning/status` | 获取设备状态 |
| GET | `/api/provisioning/events` | SSE 事件流 |
| GET | `/api/provisioning/networks` | 扫描 WiFi 网络 |
| GET | `/api/provisioning/networks/saved` | 已保存网络列表 |
| POST | `/api/provisioning/network/connect` | 连接 WiFi |
| POST | `/api/provisioning/network/forget` | 遗忘网络 |
| GET | `/api/provisioning/network/diagnostics` | 网络诊断 |
| POST | `/api/provisioning/network/test` | 测试互联网连接 |
| POST | `/api/provisioning/hotspot/enable` | 启用热点 |
| POST | `/api/provisioning/hotspot/disable` | 禁用热点 |
| POST | `/api/provisioning/restart` | 请求重启 |
| POST | `/api/provisioning/rollback` | 回滚到配网模式 |

### 5. 配置选项

#### 环境变量
```bash
MOONHUB_PROVISIONING_ENABLED=1    # 启用配网系统
MOONHUB_ALLOW_SYSTEM_CONTROL=1    # 允许系统控制（重启等）
```

#### 配置存储（provisioning.json）
```json
{
  "device.mode": "provisioning",
  "device.network.provisioned": false,
  "device.network.lastSsid": "",
  "device.network.hotspotSsid": "MoonHub-XXXX",
  "device.network.hotspotProfile": "MoonHub Hotspot",
  "device.network.hotspotPassword": "MoonHub-XXXX-wifi",
  "device.network.interfaceName": "wlan0",
  "device.network.apEnabled": true,
  "device.recovery.enabled": true,
  "device.recovery.failureThreshold": 3,
  "device.recovery.cooldownMs": 120000
}
```

### 6. systemd 服务

```ini
[Unit]
Description=MoonHub AI Assistant
After=network-online.target NetworkManager.service
Wants=network-online.target

[Service]
Type=simple
Environment=MOONHUB_PROVISIONING_ENABLED=1
Environment=MOONHUB_ALLOW_SYSTEM_CONTROL=1
ExecStart=/usr/local/bin/moonhub-web -public /etc/moonhub/config.json
Restart=always

[Install]
WantedBy=multi-user.target
```

## 使用示例

### 获取设备状态

```bash
curl http://localhost:18800/api/provisioning/status
```

响应：
```json
{
  "status": {
    "mode": "provisioning",
    "hostname": "moonhub-device",
    "network": {
      "provisioned": false,
      "connected": false,
      "hotspotSsid": "MoonHub-VICE",
      "apEnabled": true
    },
    "runtime": {
      "phase": "provisioning",
      "health": "ok",
      "summary": "Waiting for initial WiFi provisioning."
    }
  }
}
```

### 扫描 WiFi 网络

```bash
curl http://localhost:18800/api/provisioning/networks
```

### 连接 WiFi

```bash
curl -X POST http://localhost:18800/api/provisioning/network/connect \
  -H "Content-Type: application/json" \
  -d '{"ssid": "MyNetwork", "password": "mypassword"}'
```

### SSE 事件流

```javascript
const eventSource = new EventSource('/api/provisioning/events');
eventSource.addEventListener('device', (event) => {
  const data = JSON.parse(event.data);
  console.log('Event:', data);
});
```

## 测试覆盖

### 单元测试
- `TestNewDeviceManager` - 设备管理器创建
- `TestNormalizeMode` - 模式规范化
- `TestSanitizeHotspotSuffix` - 热点 SSID 后缀生成
- `TestDefaultHotspotSSID` - 默认热点 SSID 生成
- `TestGetStatus` - 状态获取
- `TestListWifiNetworks` - WiFi 网络列表
- `TestEventBroadcaster` - 事件广播
- `TestConfigStore` - 配置存储
- `TestDeriveRuntimeState` - 运行时状态派生

### 运行测试

```bash
go test ./pkg/provisioning/... -v
```

## 文件清单

### 新建文件
- `pkg/provisioning/config.go`
- `pkg/provisioning/events.go`
- `pkg/provisioning/network.go`
- `pkg/provisioning/manager.go`
- `pkg/provisioning/manager_test.go`
- `web/backend/api/provisioning.go`
- `deployment/systemd/moonhub.service`

### 修改文件
- `web/backend/api/router.go` - 添加 provisioning 字段和 SetProvisioningHandler 方法
- `web/backend/main.go` - 添加设备管理器初始化代码

## 参考实现

TinyClaw TypeScript 实现：
- `TinyClaw/tinyclaw/packages/device/src/index.ts` - DeviceManager 完整实现
- `TinyClaw/tinyclaw/packages/device/src/http.ts` - HTTP API 定义

## 后续工作

1. **PWA 前端** - 用户计划使用 PWA 实现 UI 层
2. **网络恢复增强** - 可选的自动恢复机制优化
3. **Factory Reset** - 恢复出厂设置功能（已预留接口）
