# Device Provisioning Configuration

Reference for enabling and securing **device provisioning** from the **web launcher** (`web/backend`). Gateway (`moonhub gateway`) does not load this subsystem unless you integrate it elsewhere.

## Environment Variables

| Variable | When set | Effect |
| --- | --- | --- |
| `MOONHUB_PROVISIONING_ENABLED` | `1` | Creates `DeviceManager`, starts background loop, registers `/api/provisioning/*` |
| `MOONHUB_ALLOW_SYSTEM_CONTROL` | `1` | Allows API-triggered **real** system control (e.g. `systemctl` restart paths used by manager) |
| `MOONHUB_PROVISIONING_TOKEN` | non-empty | All `/api/provisioning/*` requests must present the same secret via `Authorization: Bearer <token>` or `X-MoonHub-Provisioning-Token: <token>` |

**Security:** If the launcher listens on all interfaces (`-public` or launcher config) with provisioning enabled, set `MOONHUB_PROVISIONING_TOKEN`. The launcher logs a warning when provisioning is on, public listen is used, and the token is unset.

**Middleware order:** IP allowlist (if `allowedCIDRs` is configured) is applied **before** provisioning auth — see `web/backend/main.go`.

## Sidecar Store: `provisioning.json`

Device state is persisted next to the main config file:

- Path: `<directory of config.json>/provisioning.json`
- Implementation: `JSONConfigStore` in `manager.go` (file mode `0600`)

Sensitive values (hotspot PSK, auth code, etc.) must not be world-readable; the store uses restrictive permissions.

## Persisted Keys (`device.*`, `onboarding.*`)

Defined in `pkg/provisioning/config.go` as `Key*` constants. Examples:

| Key | Purpose |
| --- | --- |
| `device.mode` | High-level mode (`provisioning`, `connecting`, `onboarding`, `ready`, …) |
| `device.network.*` | Provisioned flag, last SSID, hotspot SSID/password, interface, AP enabled, errors |
| `device.restart.*` | Pending restart flags |
| `device.recovery.*` | Recovery enabled, failure threshold, cooldown, streaks, locks |
| `device.auth.code` | Six-digit pairing code |
| `device.auth.deviceId` | Display device id (e.g. `YSHU-YYYY-XXX`) |
| `onboarding.completed` | Onboarding completion flag |

Defaults for recovery and interface names are in `config.go` (`Default*` constants).

## HTTP Surface

All routes are registered by `ProvisioningHandler.RegisterRoutes` under prefix `/api/provisioning/`. See [`docs/implementation/provisioning-status.md`](../../../docs/implementation/provisioning-status.md) for the full method/path table.

SSE: `GET /api/provisioning/events` — when a token is configured, clients cannot use plain `EventSource` with custom headers; the frontend uses `fetch` streaming (see `web/frontend/src/hooks/use-provisioning-sse.ts`).

## Frontend Token (browser)

- Do **not** bake `MOONHUB_PROVISIONING_TOKEN` into Vite env or the repository.
- On `401`, the provisioning layout prompts for the token; it is stored in **`sessionStorage`** under `moonhub.provisioningToken`, then the page reloads.

## Related Documentation

- [README.md](./README.md) — Package overview and source map
- [`docs/implementation/provisioning-status.md`](../../../docs/implementation/provisioning-status.md) — End-to-end status, UI paths, curl examples
