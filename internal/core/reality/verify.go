package reality

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type VerifyCheck struct {
	Name string
	OK   bool
}

type VerifyResult struct {
	Checks []VerifyCheck
	URI    string
}

func (r *VerifyResult) add(name string, ok bool) {
	r.Checks = append(r.Checks, VerifyCheck{
		Name: name,
		OK:   ok,
	})
}

func (r VerifyResult) Passed() int {
	n := 0

	for _, check := range r.Checks {
		if check.OK {
			n++
		}
	}

	return n
}

func (r VerifyResult) Failed() int {
	n := 0

	for _, check := range r.Checks {
		if !check.OK {
			n++
		}
	}

	return n
}

func (r VerifyResult) Valid() bool {
	return r.Failed() == 0
}

func (c Config) Verify() VerifyResult {
	result := VerifyResult{}

	// Configuration file.
	_, err := os.Stat(c.Path)
	result.add("Config exists", err == nil)

	if err != nil {
		return result
	}

	// Read configuration.
	cfg, err := c.Read()
	result.add("Config readable", err == nil)

	if err != nil {
		return result
	}

	inbound := cfg.Inbounds[0]

	// UUID.
	uuidOK := len(inbound.Settings.Clients) > 0 &&
		strings.TrimSpace(inbound.Settings.Clients[0].ID) != ""

	result.add("UUID found", uuidOK)

	// Private key.
	privateKey := inbound.StreamSettings.RealitySettings.PrivateKey
	privateKeyOK := strings.TrimSpace(privateKey) != ""

	result.add("PrivateKey found", privateKeyOK)

	// Short ID.
	shortID := inbound.StreamSettings.RealitySettings.ShortIDs
	shortIDOK := len(shortID) > 0 &&
		strings.TrimSpace(shortID[0]) != ""

	result.add("ShortID found", shortIDOK)

	// Server name.
	serverNames := inbound.StreamSettings.RealitySettings.ServerNames
	serverNameOK := len(serverNames) > 0 &&
		strings.TrimSpace(serverNames[0]) != ""

	result.add("ServerName found", serverNameOK)

	// Reality security.
	securityOK := inbound.StreamSettings.Security == "reality"
	result.add("Reality inbound", securityOK)

	// Port.
	portOK := inbound.Port > 0 && inbound.Port <= 65535
	result.add("Port valid", portOK)

	// Xray.
	xrayPath, err := exec.LookPath("xray")
	xrayOK := err == nil
	result.add("Xray installed", xrayOK)

	// Only perform xray private-key validation when xray exists.
	if xrayOK && privateKeyOK {
		cmd := exec.Command(xrayPath, "x25519", "-i", privateKey)

		output, err := cmd.CombinedOutput()

		keyOK := err == nil &&
			strings.Contains(string(output), "Password (PublicKey):")

		result.add("PrivateKey usable", keyOK)
	}

	// Linux runtime checks.
	//
	// Windows is the development environment.
	// Runtime checks are performed on Linux/VPS only.
	if runtime.GOOS != "windows" {
		validateLinuxVerify(&result, inbound.Port)
	}

	// Generate URI only when enough information is available.
	if xrayOK &&
		uuidOK &&
		privateKeyOK &&
		shortIDOK &&
		serverNameOK &&
		portOK {

		if uri, err := c.URI(); err == nil {
			result.URI = uri
		}
	}

	return result
}

func validateLinuxVerify(result *VerifyResult, port int) {
	// systemd.
	if _, err := exec.LookPath("systemctl"); err == nil {
		cmd := exec.Command("systemctl", "is-active", "--quiet", "xray")
		result.add("Xray service running", cmd.Run() == nil)
	}

	// TCP listener.
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("ss -lntp | grep -q ':%d '", port),
	)

	result.add(
		fmt.Sprintf("Port %d listening", port),
		cmd.Run() == nil,
	)
}
