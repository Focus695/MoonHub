// API client for device provisioning

export const PROVISIONING_TOKEN_STORAGE_KEY = "moonhub.provisioningToken"

export const PROVISIONING_AUTH_REQUIRED_EVENT = "moonhub-provisioning-auth-required"

export function getProvisioningAuthHeaders(): Record<string, string> {
  if (typeof sessionStorage === "undefined") {
    return {}
  }
  const t = sessionStorage.getItem(PROVISIONING_TOKEN_STORAGE_KEY)?.trim()
  if (!t) {
    return {}
  }
  return { Authorization: `Bearer ${t}` }
}

export function hasProvisioningToken(): boolean {
  if (typeof sessionStorage === "undefined") {
    return false
  }
  return Boolean(sessionStorage.getItem(PROVISIONING_TOKEN_STORAGE_KEY)?.trim())
}

export function notifyProvisioningAuthRequired(): void {
  if (typeof window === "undefined") {
    return
  }
  window.dispatchEvent(new CustomEvent(PROVISIONING_AUTH_REQUIRED_EVENT))
}

export type DeviceMode =
  | "provisioning"
  | "connecting"
  | "onboarding"
  | "ready"
  | "maintenance"
  | "error"

export type DeviceRuntimePhase =
  | "provisioning"
  | "connecting"
  | "onboarding"
  | "ready"
  | "degraded"
  | "maintenance"
  | "restarting"
  | "error"

export type HealthStatus = "ok" | "warning" | "error"

export interface DeviceNetworkSummary {
  provisioned: boolean
  connected: boolean
  lastSsid?: string
  hotspotSsid: string
  hotspotProfile: string
  interfaceName: string
  apEnabled: boolean
  activeConnection?: string
  ipv4Address?: string
  gateway?: string
  internetReachable?: boolean
  lastError?: string
}

export interface DeviceRuntimeStatus {
  phase: DeviceRuntimePhase
  health: HealthStatus
  summary: string
  warnings: string[]
}

export type DeviceRecoveryResult =
  | "idle"
  | "pending"
  | "cooldown"
  | "hotspot_restored"
  | "rejoin_succeeded"
  | "rejoin_failed"
  | "failed"

export type DeviceRecoveryLock = "none" | "cooldown" | "recovery"

export interface DeviceRecoveryStatus {
  enabled: boolean
  failureThreshold: number
  cooldownMs: number
  degradedNetworkStreak: number
  lastAttemptAt?: number
  lastRecoveredAt?: number
  lastResult?: DeviceRecoveryResult
  cooldownUntil?: number
  activeLock: DeviceRecoveryLock
  lastReason?: string
  inCooldown: boolean
}

export interface DeviceServicesStatus {
  nmcliAvailable: boolean
  networkManagerAvailable: boolean
}

export interface DeviceStatus {
  mode: DeviceMode
  hostname: string
  network: DeviceNetworkSummary
  runtime: DeviceRuntimeStatus
  recovery: DeviceRecoveryStatus
  restartPending: boolean
  lastRestartReason?: string
  services: DeviceServicesStatus
}

export interface WifiNetwork {
  ssid: string
  signal: number
  security: string
  inUse: boolean
}

export interface SavedNetworkProfile {
  name: string
  uuid?: string
  type: string
  autoconnect: boolean
  active: boolean
}

export interface NetworkDiagnostics {
  connected: boolean
  internetReachable: boolean
  dnsResolved: boolean
  nmcliState: string
  activeConnection?: string
  ipv4Address?: string
  gateway?: string
  details: string[]
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const headers = new Headers(options?.headers)
  const auth = getProvisioningAuthHeaders()
  if (auth.Authorization) {
    headers.set("Authorization", auth.Authorization)
  }
  const res = await fetch(path, { ...options, headers })
  if (!res.ok) {
    if (res.status === 401 && path.startsWith("/api/provisioning")) {
      notifyProvisioningAuthRequired()
    }
    let message = `API error: ${res.status} ${res.statusText}`
    try {
      const body = (await res.json()) as { error?: string; errors?: string[] }
      if (Array.isArray(body.errors) && body.errors.length > 0) {
        message = body.errors.join("; ")
      } else if (typeof body.error === "string" && body.error.trim() !== "") {
        message = body.error
      }
    } catch {
      // Keep fallback error message when response body is not JSON.
    }
    throw new Error(message)
  }
  return res.json() as Promise<T>
}

export async function getStatus(): Promise<{ status: DeviceStatus }> {
  return request<{ status: DeviceStatus }>("/api/provisioning/status")
}

export async function scanNetworks(): Promise<{ networks: WifiNetwork[] }> {
  return request<{ networks: WifiNetwork[] }>("/api/provisioning/networks")
}

export async function getSavedNetworks(): Promise<{
  networks: SavedNetworkProfile[]
}> {
  return request<{ networks: SavedNetworkProfile[] }>(
    "/api/provisioning/networks/saved"
  )
}

export async function connectWifi(
  ssid: string,
  password: string,
  hidden?: boolean
): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/network/connect", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ ssid, password, hidden }),
  })
}

export async function forgetNetwork(name: string): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/network/forget", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name }),
  })
}

export async function getDiagnostics(): Promise<{
  diagnostics: NetworkDiagnostics
}> {
  return request<{ diagnostics: NetworkDiagnostics }>(
    "/api/provisioning/network/diagnostics"
  )
}

export async function testInternet(): Promise<{
  connected: boolean
  diagnostics: NetworkDiagnostics
}> {
  return request<{ connected: boolean; diagnostics: NetworkDiagnostics }>(
    "/api/provisioning/network/test",
    { method: "POST" }
  )
}

export async function enableHotspot(): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/hotspot/enable", {
    method: "POST",
  })
}

export async function disableHotspot(): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/hotspot/disable", {
    method: "POST",
  })
}

export async function getAuthCode(): Promise<{ code: string; deviceId: string }> {
  return request<{ code: string; deviceId: string }>(
    "/api/provisioning/auth-code"
  )
}

export async function regenerateAuthCode(): Promise<{
  code: string
  deviceId: string
}> {
  return request<{ code: string; deviceId: string }>(
    "/api/provisioning/auth-code/regenerate",
    { method: "POST" }
  )
}

export async function triggerRecovery(): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/recovery/trigger", {
    method: "POST",
  })
}

export async function factoryReset(): Promise<{ ok: boolean }> {
  return request<{ ok: boolean }>("/api/provisioning/factory-reset", {
    method: "POST",
  })
}
