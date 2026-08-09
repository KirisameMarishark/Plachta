package reality

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ValidationResult struct {
	Valid  bool
	Errors []string
}

func (c Config) Validate() ValidationResult {
	result := ValidationResult{
		Valid: true,
	}

	// 1. Check configuration file.
	if _, err := os.Stat(c.Path); err != nil {
		result.add("Reality config not found")
		return result
	}

	// 2. Check JSON validity.
	data, err := os.ReadFile(c.Path)
	if err != nil {
		result.add("cannot read Reality config")
		return result
	}

	var raw map[string]interface{}

	if err := json.Unmarshal(data, &raw); err != nil {
		result.add("invalid Reality JSON")
		return result
	}

	// 3. Read Reality values.
	cfg, err := c.Read()
	if err != nil {
		result.add(err.Error())
		return result
	}

	if len(cfg.Inbounds) == 0 {
		result.add("no inbounds found")
		return result
	}

	inbound := cfg.Inbounds[0]

	if len(inbound.Settings.Clients) == 0 {
		result.add("no Reality client found")
	} else if strings.TrimSpace(inbound.Settings.Clients[0].ID) == "" {
		result.add("Reality UUID is empty")
	}

	reality := inbound.StreamSettings.RealitySettings

	if strings.TrimSpace(reality.PrivateKey) == "" {
		result.add("Reality private key is empty")
	}

	if len(reality.ShortIDs) == 0 ||
		strings.TrimSpace(reality.ShortIDs[0]) == "" {
		result.add("Reality short ID is empty")
	}

	if len(reality.ServerNames) == 0 ||
		strings.TrimSpace(reality.ServerNames[0]) == "" {
		result.add("Reality server name is empty")
	}

	// 4. Check port.
	if inbound.Port <= 0 || inbound.Port > 65535 {
		result.add("invalid Reality port")
	}

	// 5. Linux runtime checks.
	//
	// These checks intentionally do not fail on Windows.
	// Windows is our development environment.
	if runtime.GOOS != "windows" {
		validateLinuxRuntime(&result)
	}

	return result
}

func validateLinuxRuntime(result *ValidationResult) {
	// Check xray.
	if _, err := exec.LookPath("xray"); err != nil {
		result.add("xray not found")
	}

	// Check systemd service.
	if _, err := exec.LookPath("systemctl"); err == nil {
		cmd := exec.Command("systemctl", "is-active", "--quiet", "xray")

		if err := cmd.Run(); err != nil {
			result.add("xray service is not active")
		}
	}

	// Check listening port.
	//
	// We use a TCP connection probe instead of relying on `ss`
	// output parsing. This keeps the Go implementation portable.
	cfg, err := New().Read()
	if err != nil || len(cfg.Inbounds) == 0 {
		return
	}

	port := cfg.Inbounds[0].Port

	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("127.0.0.1:%d", port),
		500_000_000,
	)

	if err != nil {
		result.add(
			fmt.Sprintf("Reality port %d is not accepting connections", port),
		)
		return
	}

	conn.Close()
}

func (r *ValidationResult) add(message string) {
	r.Valid = false
	r.Errors = append(r.Errors, message)
}
