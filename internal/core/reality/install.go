package reality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	xrayServicePath = "/etc/systemd/system/xray.service"
	xrayBinaryPath  = "/usr/local/bin/xray"
)

const xrayService = `[Unit]
Description=Xray Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/xray run -config /etc/plachta/reality/config.json
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`

func Install() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("reality install requires root privileges")
	}

	backupPath, err := backupRealityConfig()
	if err != nil {
		return err
	}

	if err := installReality(); err != nil {
		if restoreErr := restoreRealityConfig(backupPath); restoreErr != nil {
			return fmt.Errorf("%w; additionally failed to restore backup: %v", err, restoreErr)
		}
		return err
	}

	return nil
}

func installReality() error {
	fmt.Println("Reality installation started.")
	fmt.Println("Xray will restart during installation.")
	fmt.Println("If this SSH connection uses the current Reality tunnel, the connection may disconnect.")
	fmt.Println("Reconnect after a few seconds to verify the installation.")

	if err := requireCommand("systemctl"); err != nil {
		return err
	}

	if err := installXrayBinary(); err != nil {
		return err
	}

	if err := GenerateConfig(); err != nil {
		return fmt.Errorf("failed to generate Reality config: %w", err)
	}

	if err := validateGeneratedConfig(); err != nil {
		return err
	}

	if err := installXrayService(); err != nil {
		return err
	}

	return nil
}

func backupRealityConfig() (string, error) {
	data, err := os.ReadFile(defaultRealityConfigPath)
	if err != nil {
		return "", fmt.Errorf("failed to read existing Reality config: %w", err)
	}

	backupPath := defaultRealityConfigPath + ".bak"
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to backup Reality config: %w", err)
	}

	return backupPath, nil
}

func restoreRealityConfig(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read Reality backup: %w", err)
	}

	if err := os.WriteFile(defaultRealityConfigPath, data, 0600); err != nil {
		return fmt.Errorf("failed to restore Reality config: %w", err)
	}

	return nil
}

func installXrayBinary() error {
	tmpDir, err := os.MkdirTemp("", "plachta-xray-*")
	if err != nil {
		return fmt.Errorf(
			"failed to create temporary directory: %w",
			err,
		)
	}
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(
		tmpDir,
		"Xray-linux-64.zip",
	)

	if err := os.WriteFile(
		zipPath,
		xrayPackage,
		0600,
	); err != nil {
		return fmt.Errorf(
			"failed to write embedded Xray package: %w",
			err,
		)
	}

	if err := os.MkdirAll(
		"/usr/local/bin",
		0755,
	); err != nil {
		return fmt.Errorf(
			"failed to create /usr/local/bin: %w",
			err,
		)
	}

	if err := runCommand(
		"unzip",
		"-o",
		zipPath,
		"-d",
		tmpDir,
	); err != nil {
		return fmt.Errorf(
			"failed to extract Xray package: %w",
			err,
		)
	}

	extractedBinary := filepath.Join(
		tmpDir,
		"xray",
	)

	if _, err := os.Stat(extractedBinary); err != nil {
		return fmt.Errorf(
			"Xray binary not found after extraction: %w",
			err,
		)
	}

	data, err := os.ReadFile(extractedBinary)
	if err != nil {
		return fmt.Errorf(
			"failed to read extracted Xray binary: %w",
			err,
		)
	}

	if err := runSystemctl("stop", "xray"); err != nil {
		return err
	}

	tmpBinary := filepath.Join(
		"/usr/local/bin",
		"xray.new",
	)

	if err := os.WriteFile(
		tmpBinary,
		data,
		0755,
	); err != nil {
		return fmt.Errorf(
			"failed to stage Xray binary: %w",
			err,
		)
	}

	if err := os.Chmod(tmpBinary, 0755); err != nil {
		_ = os.Remove(tmpBinary)

		return fmt.Errorf(
			"failed to set Xray binary permissions: %w",
			err,
		)
	}

	if err := os.Rename(
		tmpBinary,
		xrayBinaryPath,
	); err != nil {
		_ = os.Remove(tmpBinary)

		return fmt.Errorf(
			"failed to replace Xray binary: %w",
			err,
		)
	}

	if err := scheduleXrayRecovery(); err != nil {
		return err
	}

	return nil
}

func scheduleXrayRecovery() error {
	return runCommand(
		"systemd-run",
		"--unit=plachta-xray-recovery",
		"--collect",
		"/bin/systemctl",
		"start",
		"xray",
	)
}

func runCommand(name string, args ...string) error {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"%s failed: %w: %s",
			name,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

func requireCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("required command %q not found", name)
	}

	return nil
}

func validateGeneratedConfig() error {
	output, err := exec.Command(
		"xray",
		"run",
		"-test",
		"-config",
		defaultRealityConfigPath,
	).CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"xray config validation failed: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

func installXrayService() error {
	if err := os.WriteFile(xrayServicePath, []byte(xrayService), 0644); err != nil {
		return fmt.Errorf("failed to write xray systemd service: %w", err)
	}

	if err := runSystemctl("daemon-reload"); err != nil {
		return err
	}

	if err := runSystemctl("enable", "xray"); err != nil {
		return err
	}

	if err := runSystemctl("restart", "xray"); err != nil {
		return err
	}

	return nil
}

func runSystemctl(args ...string) error {
	output, err := exec.Command("systemctl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"systemctl %s failed: %w: %s",
			strings.Join(args, " "),
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}
