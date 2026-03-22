package provisioning

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandResult represents the result of a shell command execution.
type CommandResult struct {
	OK     bool
	Stdout string
	Stderr string
}

// CommandRunner is a function type for running shell commands.
type CommandRunner func(ctx context.Context, args ...string) CommandResult

// NetworkOperator handles network operations via nmcli.
type NetworkOperator struct {
	interfaceName string
	run           CommandRunner
}

// NewNetworkOperator creates a new network operator.
func NewNetworkOperator(iface string, runner CommandRunner) *NetworkOperator {
	if runner == nil {
		runner = defaultCommandRunner
	}
	return &NetworkOperator{
		interfaceName: iface,
		run:           runner,
	}
}

// CheckServices verifies nmcli and NetworkManager availability.
func (n *NetworkOperator) CheckServices(ctx context.Context) (nmcliOK, nmServiceOK bool) {
	nmcliResult := n.run(ctx, "nmcli", "--version")
	nmcliOK = nmcliResult.OK

	nmServiceResult := n.run(ctx, "systemctl", "is-active", "NetworkManager")
	nmServiceOK = nmServiceResult.OK && strings.TrimSpace(nmServiceResult.Stdout) == "active"

	return nmcliOK, nmServiceOK
}

// DiagnoseNetwork performs comprehensive network diagnostics.
func (n *NetworkOperator) DiagnoseNetwork(ctx context.Context) NetworkDiagnostics {
	diagnostics := NetworkDiagnostics{
		Details: make([]string, 0),
	}

	// Check networking state
	networkingState := n.run(ctx, "nmcli", "-t", "-f", "STATE", "networking")
	diagnostics.NmcliState = strings.TrimSpace(networkingState.Stdout)
	if networkingState.Stderr != "" {
		diagnostics.Details = append(diagnostics.Details, networkingState.Stderr)
	}

	// Get device info
	deviceShow := n.run(ctx, "nmcli", "-t", "-f",
		"GENERAL.STATE,GENERAL.CONNECTION,IP4.ADDRESS,IP4.GATEWAY",
		"device", "show", n.interfaceName)

	deviceFields := parseDeviceShow(deviceShow.Stdout)
	if deviceShow.Stderr != "" {
		diagnostics.Details = append(diagnostics.Details, deviceShow.Stderr)
	}

	diagnostics.ActiveConnection = deviceFields["GENERAL.CONNECTION"]
	diagnostics.IPv4Address = deviceFields["IP4.ADDRESS[1]"]
	if diagnostics.IPv4Address == "" {
		diagnostics.IPv4Address = deviceFields["IP4.ADDRESS[0]"]
	}
	diagnostics.Gateway = deviceFields["IP4.GATEWAY"]

	// Check connectivity
	diagnostics.Connected = networkingState.OK &&
		strings.Contains(strings.ToLower(networkingState.Stdout), "connected")

	// Test internet connectivity
	pingResult := n.run(ctx, "ping", "-c", "1", "-W", "2", "1.1.1.1")
	diagnostics.InternetReachable = pingResult.OK
	if pingResult.Stderr != "" && !pingResult.OK {
		diagnostics.Details = append(diagnostics.Details, pingResult.Stderr)
	}

	// Test DNS resolution
	dnsResult := n.run(ctx, "getent", "hosts", "api.openai.com")
	diagnostics.DNSResolved = dnsResult.OK && strings.TrimSpace(dnsResult.Stdout) != ""

	return diagnostics
}

// ScanWifiNetworks scans for available WiFi networks.
func (n *NetworkOperator) ScanWifiNetworks(ctx context.Context) ([]WifiNetwork, error) {
	result := n.run(ctx, "nmcli", "-t", "-f", "SSID,SIGNAL,SECURITY,IN-USE",
		"device", "wifi", "list", "--rescan", "yes")

	if !result.OK {
		return nil, fmt.Errorf("%s", result.Stderr)
	}

	return parseWifiList(result.Stdout), nil
}

// ListSavedNetworks returns saved WiFi connection profiles.
func (n *NetworkOperator) ListSavedNetworks(ctx context.Context, hotspotProfile string) ([]SavedNetworkProfile, error) {
	result := n.run(ctx, "nmcli", "-t", "-f", "NAME,UUID,TYPE,AUTOCONNECT,ACTIVE", "connection", "show")

	if !result.OK {
		return nil, fmt.Errorf("%s", result.Stderr)
	}

	profiles := parseConnectionShow(result.Stdout)

	// Filter out hotspot and non-wifi profiles
	filtered := make([]SavedNetworkProfile, 0, len(profiles))
	for _, p := range profiles {
		if p.Type == "802-11-wireless" && p.Name != hotspotProfile && p.Name != "Hotspot" {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

// ConnectToWiFi connects to a WiFi network.
func (n *NetworkOperator) ConnectToWiFi(ctx context.Context, ssid, password string, hidden bool) error {
	args := []string{"nmcli", "device", "wifi", "connect", ssid}
	if password != "" {
		args = append(args, "password", password)
	}
	if hidden {
		args = append(args, "hidden", "yes")
	}

	result := n.run(ctx, args...)
	if !result.OK {
		return fmt.Errorf("%s", result.Stderr)
	}

	show := n.run(ctx, "nmcli", "-t", "-f", "GENERAL.CONNECTION", "device", "show", n.interfaceName)
	connName := strings.TrimSpace(show.Stdout)
	if connName != "" && connName != "--" {
		n.run(ctx, "nmcli", "connection", "modify", connName, "connection.autoconnect", "yes")
	}

	return nil
}

// ForgetNetwork deletes a saved network profile.
func (n *NetworkOperator) ForgetNetwork(ctx context.Context, name string) error {
	result := n.run(ctx, "nmcli", "connection", "delete", "id", name)
	if !result.OK {
		return fmt.Errorf("%s", result.Stderr)
	}
	return nil
}

// EnsureHotspotProfile creates or updates the hotspot connection profile.
func (n *NetworkOperator) EnsureHotspotProfile(ctx context.Context, profile, ssid, password string) error {
	// Check if profile exists
	existing := n.run(ctx, "nmcli", "-t", "-f", "NAME", "connection", "show", "id", profile)

	if !existing.OK {
		// Create new profile
		createArgs := []string{
			"nmcli", "connection", "add", "type", "wifi",
			"ifname", n.interfaceName,
			"con-name", profile,
			"autoconnect", "no",
			"ssid", ssid,
		}
		createResult := n.run(ctx, createArgs...)
		if !createResult.OK {
			return fmt.Errorf("failed to create hotspot profile: %s", createResult.Stderr)
		}
	}

	// Modify profile with hotspot settings
	modifyArgs := []string{
		"nmcli", "connection", "modify", profile,
		"802-11-wireless.mode", "ap",
		"802-11-wireless.band", "bg",
		"ipv4.method", "shared",
		"ipv6.method", "ignore",
		"wifi-sec.key-mgmt", "wpa-psk",
		"wifi-sec.psk", password,
		"connection.autoconnect", "no",
		"connection.interface-name", n.interfaceName,
		"802-11-wireless.ssid", ssid,
	}

	modifyResult := n.run(ctx, modifyArgs...)
	if !modifyResult.OK {
		return fmt.Errorf("failed to update hotspot profile: %s", modifyResult.Stderr)
	}

	return nil
}

// EnableHotspot activates the hotspot connection.
func (n *NetworkOperator) EnableHotspot(ctx context.Context, profile string) error {
	result := n.run(ctx, "nmcli", "connection", "up", "id", profile)
	if !result.OK {
		return fmt.Errorf("%s", result.Stderr)
	}
	return nil
}

// DisableHotspot deactivates the hotspot connection.
func (n *NetworkOperator) DisableHotspot(ctx context.Context, profile string) error {
	result := n.run(ctx, "nmcli", "connection", "down", "id", profile)
	if !result.OK {
		// Try fallback profile name
		fallback := n.run(ctx, "nmcli", "connection", "down", "id", "Hotspot")
		if !fallback.OK {
			// Last resort: disconnect interface
			disconnect := n.run(ctx, "nmcli", "device", "disconnect", n.interfaceName)
			if !disconnect.OK {
				return fmt.Errorf("%s", result.Stderr)
			}
		}
	}
	return nil
}

// ReconnectSavedNetwork reconnects to a saved WiFi network.
func (n *NetworkOperator) ReconnectSavedNetwork(ctx context.Context, name string) error {
	result := n.run(ctx, "nmcli", "connection", "up", "id", name)
	if !result.OK {
		return fmt.Errorf("%s", result.Stderr)
	}

	// Enable autoconnect
	n.run(ctx, "nmcli", "connection", "modify", name, "connection.autoconnect", "yes")

	return nil
}

// TestInternet tests internet connectivity.
func (n *NetworkOperator) TestInternet(ctx context.Context) bool {
	result := n.run(ctx, "ping", "-c", "1", "-W", "2", "1.1.1.1")
	return result.OK
}

// Run executes a command with the internal runner (exposed for external use).
func (n *NetworkOperator) Run(ctx context.Context, args ...string) CommandResult {
	return n.run(ctx, args...)
}

// Helper functions

func parseDeviceShow(stdout string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := line[:idx]
		value := line[idx+1:]
		result[key] = value
	}
	return result
}

func parseWifiList(stdout string) []WifiNetwork {
	var networks []WifiNetwork
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 4 {
			continue
		}

		n := len(parts)
		inUse := parts[n-1] == "*"
		security := parts[n-2]
		if security == "" {
			security = "open"
		}
		signalStr := parts[n-3]
		ssid := strings.Join(parts[:n-3], ":")
		if ssid == "" {
			continue
		}

		signal := 0
		fmt.Sscanf(signalStr, "%d", &signal)

		networks = append(networks, WifiNetwork{
			SSID:     ssid,
			Signal:   signal,
			Security: security,
			InUse:    inUse,
		})
	}
	return networks
}

func parseConnectionShow(stdout string) []SavedNetworkProfile {
	var profiles []SavedNetworkProfile
	for _, line := range strings.Split(stdout, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 5 {
			continue
		}

		name := parts[0]
		if name == "" {
			continue
		}

		profiles = append(profiles, SavedNetworkProfile{
			Name:        name,
			UUID:        parts[1],
			Type:        parts[2],
			Autoconnect: parts[3] == "yes",
			Active:      parts[4] == "yes",
		})
	}
	return profiles
}

func defaultCommandRunner(ctx context.Context, args ...string) CommandResult {
	if len(args) == 0 {
		return CommandResult{OK: false, Stderr: "no command specified"}
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Env = os.Environ()

	stdout, err := cmd.Output()
	if err != nil {
		var stderr string
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = string(exitErr.Stderr)
		} else {
			stderr = err.Error()
		}
		return CommandResult{OK: false, Stdout: string(stdout), Stderr: stderr}
	}

	return CommandResult{OK: true, Stdout: string(stdout)}
}

// SanitizeHotspotSuffix creates a safe suffix from hostname.
func SanitizeHotspotSuffix(hostname string) string {
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, hostname)
	cleaned = strings.ToUpper(cleaned)
	if len(cleaned) < 4 {
		return "NODE"
	}
	return cleaned[len(cleaned)-4:]
}

// DefaultHotspotSSID generates the default hotspot SSID.
func DefaultHotspotSSID(hostname string) string {
	return fmt.Sprintf("%s-%s", DefaultHotspotPrefix, SanitizeHotspotSuffix(hostname))
}

// DefaultHostname returns the system hostname.
func DefaultHostname() string {
	if host, err := os.Hostname(); err == nil && host != "" {
		return host
	}
	return "moonhub-device"
}
