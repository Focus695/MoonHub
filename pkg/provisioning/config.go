// Package provisioning provides device provisioning and network management
// for MoonHub devices, enabling "plug-and-play" zero-configuration setup.
package provisioning

// DeviceMode represents the current operational state of the device.
type DeviceMode string

const (
	// ModeProvisioning: Device is waiting for initial WiFi setup (hotspot active)
	ModeProvisioning DeviceMode = "provisioning"
	// ModeConnecting: Device is attempting to connect to configured WiFi
	ModeConnecting DeviceMode = "connecting"
	// ModeOnboarding: Network connected, waiting for user setup to complete
	ModeOnboarding DeviceMode = "onboarding"
	// ModeReady: Device is fully operational
	ModeReady DeviceMode = "ready"
	// ModeMaintenance: Device is in maintenance mode
	ModeMaintenance DeviceMode = "maintenance"
	// ModeError: Device encountered an unrecoverable error
	ModeError DeviceMode = "error"
)

// DeviceRuntimePhase represents the runtime phase derived from mode and network state.
type DeviceRuntimePhase string

const (
	PhaseProvisioning DeviceRuntimePhase = "provisioning"
	PhaseConnecting   DeviceRuntimePhase = "connecting"
	PhaseOnboarding   DeviceRuntimePhase = "onboarding"
	PhaseReady        DeviceRuntimePhase = "ready"
	PhaseDegraded     DeviceRuntimePhase = "degraded"
	PhaseMaintenance  DeviceRuntimePhase = "maintenance"
	PhaseRestarting   DeviceRuntimePhase = "restarting"
	PhaseError        DeviceRuntimePhase = "error"
)

// DeviceRecoveryResult represents the outcome of a recovery attempt.
type DeviceRecoveryResult string

const (
	RecoveryIdle            DeviceRecoveryResult = "idle"
	RecoveryPending         DeviceRecoveryResult = "pending"
	RecoveryCooldown        DeviceRecoveryResult = "cooldown"
	RecoveryHotspotRestored DeviceRecoveryResult = "hotspot_restored"
	RecoveryRejoinSucceeded DeviceRecoveryResult = "rejoin_succeeded"
	RecoveryRejoinFailed    DeviceRecoveryResult = "rejoin_failed"
	RecoveryFailed          DeviceRecoveryResult = "failed"
)

// DeviceRecoveryLock represents the current recovery lock state.
type DeviceRecoveryLock string

const (
	LockNone     DeviceRecoveryLock = "none"
	LockCooldown DeviceRecoveryLock = "cooldown"
	LockRecovery DeviceRecoveryLock = "recovery"
)

// HealthStatus represents the health level of the device.
type HealthStatus string

const (
	HealthOK      HealthStatus = "ok"
	HealthWarning HealthStatus = "warning"
	HealthError   HealthStatus = "error"
)

// DeviceNetworkSummary contains network state information.
type DeviceNetworkSummary struct {
	Provisioned       bool   `json:"provisioned"`
	Connected         bool   `json:"connected"`
	LastSSID          string `json:"lastSsid,omitempty"`
	HotspotSSID       string `json:"hotspotSsid"`
	HotspotProfile    string `json:"hotspotProfile"`
	InterfaceName     string `json:"interfaceName"`
	APEnabled         bool   `json:"apEnabled"`
	ActiveConnection  string `json:"activeConnection,omitempty"`
	IPv4Address       string `json:"ipv4Address,omitempty"`
	Gateway           string `json:"gateway,omitempty"`
	InternetReachable bool   `json:"internetReachable,omitempty"`
	LastError         string `json:"lastError,omitempty"`
}

// DeviceRuntimeStatus contains derived runtime state.
type DeviceRuntimeStatus struct {
	Phase    DeviceRuntimePhase `json:"phase"`
	Health   HealthStatus       `json:"health"`
	Summary  string             `json:"summary"`
	Warnings []string           `json:"warnings"`
}

// DeviceRecoveryStatus contains recovery mechanism state.
type DeviceRecoveryStatus struct {
	Enabled               bool                 `json:"enabled"`
	FailureThreshold      int                  `json:"failureThreshold"`
	CooldownMs            int                  `json:"cooldownMs"`
	DegradedNetworkStreak int                  `json:"degradedNetworkStreak"`
	LastAttemptAt         int64                `json:"lastAttemptAt,omitempty"`
	LastRecoveredAt       int64                `json:"lastRecoveredAt,omitempty"`
	LastResult            DeviceRecoveryResult `json:"lastResult,omitempty"`
	CooldownUntil         int64                `json:"cooldownUntil,omitempty"`
	ActiveLock            DeviceRecoveryLock   `json:"activeLock"`
	LastReason            string               `json:"lastReason,omitempty"`
	InCooldown            bool                 `json:"inCooldown"`
}

// DeviceServicesStatus contains service availability info.
type DeviceServicesStatus struct {
	NmcliAvailable          bool `json:"nmcliAvailable"`
	NetworkManagerAvailable bool `json:"networkManagerAvailable"`
}

// DeviceStatus is the complete device status structure.
type DeviceStatus struct {
	Mode              DeviceMode           `json:"mode"`
	Hostname          string               `json:"hostname"`
	Network           DeviceNetworkSummary `json:"network"`
	Runtime           DeviceRuntimeStatus  `json:"runtime"`
	Recovery          DeviceRecoveryStatus `json:"recovery"`
	RestartPending    bool                 `json:"restartPending"`
	LastRestartReason string               `json:"lastRestartReason,omitempty"`
	Services          DeviceServicesStatus `json:"services"`
}

// WifiNetwork represents a scanned WiFi network.
type WifiNetwork struct {
	SSID     string `json:"ssid"`
	Signal   int    `json:"signal"`
	Security string `json:"security"`
	InUse    bool   `json:"inUse"`
}

// SavedNetworkProfile represents a saved WiFi connection profile.
type SavedNetworkProfile struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid,omitempty"`
	Type        string `json:"type"`
	Autoconnect bool   `json:"autoconnect"`
	Active      bool   `json:"active"`
}

// NetworkDiagnostics contains network diagnostic results.
type NetworkDiagnostics struct {
	Connected         bool     `json:"connected"`
	InternetReachable bool     `json:"internetReachable"`
	DNSResolved       bool     `json:"dnsResolved"`
	NmcliState        string   `json:"nmcliState"`
	ActiveConnection  string   `json:"activeConnection,omitempty"`
	IPv4Address       string   `json:"ipv4Address,omitempty"`
	Gateway           string   `json:"gateway,omitempty"`
	Details           []string `json:"details"`
}

// WiFiConnectRequest is the request body for WiFi connection.
type WiFiConnectRequest struct {
	SSID     string `json:"ssid"`
	Password string `json:"password,omitempty"`
	Hidden   bool   `json:"hidden,omitempty"`
}

// RestartRequest is the request body for system restart.
type RestartRequest struct {
	Reason string `json:"reason"`
	Scope  string `json:"scope,omitempty"` // "app" or "device"
}

// Configuration keys for device state
const (
	KeyDeviceMode               = "device.mode"
	KeyNetworkProvisioned       = "device.network.provisioned"
	KeyLastSSID                 = "device.network.lastSsid"
	KeyHotspotSSID              = "device.network.hotspotSsid"
	KeyInterfaceName            = "device.network.interfaceName"
	KeyAPEnabled                = "device.network.apEnabled"
	KeyLastError                = "device.network.lastError"
	KeyHotspotProfile           = "device.network.hotspotProfile"
	KeyHotspotPassword          = "device.network.hotspotPassword"
	KeyRestartPending           = "device.restart.pending"
	KeyRestartReason            = "device.restart.reason"
	KeyOnboardingCompleted      = "onboarding.completed"
	KeyRecoveryEnabled          = "device.recovery.enabled"
	KeyRecoveryFailureThreshold = "device.recovery.failureThreshold"
	KeyRecoveryCooldownMs       = "device.recovery.cooldownMs"
	KeyRecoveryFailureStreak    = "device.recovery.degradedNetworkStreak"
	KeyRecoveryLastAttemptAt    = "device.recovery.lastAttemptAt"
	KeyRecoveryLastRecoveredAt  = "device.recovery.lastRecoveredAt"
	KeyRecoveryLastResult       = "device.recovery.lastResult"
	KeyRecoveryCooldownUntil    = "device.recovery.cooldownUntil"
	KeyRecoveryActiveLock       = "device.recovery.activeLock"
	KeyRecoveryLastReason       = "device.recovery.lastReason"
	KeyAuthCode                 = "device.auth.code"
	KeyDeviceId                 = "device.auth.deviceId"
)

// Default values
const (
	DefaultHotspotPrefix            = "MoonHub"
	DefaultNetworkInterface         = "wlan0"
	DefaultHotspotProfile           = "MoonHub Hotspot"
	DefaultRecoveryFailureThreshold = 3
	DefaultRecoveryCooldownMs       = 120000 // 2 minutes
	DefaultNetworkProbeTimeoutMs    = 15000
	DefaultNetworkProbeMaxRetries   = 2
)

// NormalizeMode converts a string to DeviceMode with fallback to provisioning.
func NormalizeMode(value string) DeviceMode {
	switch value {
	case string(ModeProvisioning), string(ModeConnecting), string(ModeOnboarding),
		string(ModeReady), string(ModeMaintenance), string(ModeError):
		return DeviceMode(value)
	default:
		return ModeProvisioning
	}
}
