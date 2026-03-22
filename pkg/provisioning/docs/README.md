# MoonHub Device Provisioning (`pkg/provisioning`)

WiFi provisioning and network management for edge devices: scan/connect WiFi, hotspot (AP) control, connectivity diagnostics, automatic recovery with cooldown, factory reset, six-digit pairing **auth code** and device ID, and **SSE** event broadcast for the web UI. The package is **inert** until the web launcher enables it via environment variables (see [CONFIG.md](./CONFIG.md)).

**Repository documentation index** (subsystem reading order): [docs/README.md](../../../docs/README.md).

## Recommended Reading Order (Documentation Flow)

1. **This document** — Responsibilities, source map, integration table
2. [CONFIG.md](./CONFIG.md) — Environment variables, persisted keys, `provisioning.json`, HTTP security
3. Implementation details — [`docs/implementation/provisioning-status.md`](../../../docs/implementation/provisioning-status.md) (Chinese; API endpoints, frontend file tree, curl examples, tests)

## Source Map

```
pkg/provisioning/
├── config.go       # Modes, runtime phases, recovery types, config key constants
├── events.go       # SSE-oriented event types and broadcaster
├── network.go      # nmcli / NetworkManager command runner and parsing
├── network_test.go   # Parser and helper tests
├── manager.go        # DeviceManager state machine and operations
├── manager_test.go   # Manager and flow tests (mock runner)
└── docs/
    ├── README.md     # This file
    └── CONFIG.md     # Configuration and security reference
```

## Related Code (Integration Points)

| Area | Path | Description |
| --- | --- | --- |
| HTTP API | `web/backend/api/provisioning.go` | REST + SSE handlers under `/api/provisioning/*` |
| Route registration | `web/backend/api/router.go` | `SetProvisioningHandler` before `RegisterRoutes` |
| Launcher wiring | `web/backend/main.go` | `MOONHUB_PROVISIONING_ENABLED`, `NewJSONConfigStore`, `Start` goroutine, middleware order |
| Auth middleware | `web/backend/middleware/provisioning_auth.go` | Optional Bearer / header token for all provisioning API routes |
| Provisioning UI | `web/frontend/src/routes/provisioning/` | TanStack Router pages; API client in `src/api/provisioning.ts` |
| systemd (example) | `deployment/systemd/moonhub.service` | Optional unit for edge deployment |

Testing: `go test ./pkg/provisioning/... -v` (see implementation status for scenarios).

## Alignment with Other Subsystems

This directory follows the same convention as [`pkg/delegation/docs`](../../delegation/docs/README.md), [`pkg/routing/docs`](../../routing/docs/README.md), and [`pkg/shield/docs`](../../shield/docs/README.md): **package-level `docs/`** covers configuration and where to read code; **repository-level `docs/implementation/provisioning-status.md`** holds the full feature list, API tables, and frontend structure.
