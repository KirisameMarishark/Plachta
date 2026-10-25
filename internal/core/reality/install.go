package reality

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	xrayServicePath       = "/etc/systemd/system/xray.service"
	xrayServiceDropInPath = "/etc/systemd/system/xray.service.d/99-plachta.conf"
	xrayBinaryPath        = "/usr/local/bin/xray"
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
const xrayServiceDropIn = `[Service]
ExecStart=
ExecStart=/usr/local/bin/xray run -config /etc/plachta/reality/config.json
`

func Install() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("reality install requires root privileges")
	}

	const binaryPath = "/usr/local/bin/xray"
	const rollbackPath = "/usr/local/bin/xray.rollback"

	// Record the original service state before installation.
	serviceWasActive := exec.Command(
		"systemctl", "is-active", "--quiet", "xray",
	).Run() == nil

	// Back up the existing Reality configuration.
	configBackupPath, err := backupRealityConfig()
	if err != nil {
		return err
	}

	// Back up the existing Xray binary, if present.
	var binaryBackup []byte
	var binaryMode os.FileMode
	binaryExists := false

	info, err := os.Stat(binaryPath)
	if err == nil {
		binaryBackup, err = os.ReadFile(binaryPath)
		if err != nil {
			return fmt.Errorf("back up Xray binary: %w", err)
		}

		binaryMode = info.Mode().Perm()
		binaryExists = true
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect existing Xray binary: %w", err)
	}

	// Run the installation.
	if err := installReality(); err != nil {
		var rollbackErrors []error

		// Restore the original configuration, or remove the newly generated
		// configuration if this was the first installation.
		if configBackupPath != "" {
			if restoreErr := restoreRealityConfig(configBackupPath); restoreErr != nil {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("restore Reality config: %w", restoreErr),
				)
			}
		} else {
			if removeErr := os.Remove(defaultRealityConfigPath); removeErr != nil &&
				!os.IsNotExist(removeErr) {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("remove newly generated Reality config: %w", removeErr),
				)
			}
		}

		// Restore the original binary, or remove the newly installed one
		// if no binary existed before installation.
		if binaryExists {
			if removeErr := os.Remove(rollbackPath); removeErr != nil &&
				!os.IsNotExist(removeErr) {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("remove stale rollback file: %w", removeErr),
				)
			} else {
				writeErr := os.WriteFile(rollbackPath, binaryBackup, binaryMode)
				if writeErr == nil {
					writeErr = os.Chmod(rollbackPath, binaryMode)
				}
				if writeErr == nil {
					writeErr = os.Rename(rollbackPath, binaryPath)
				}

				if writeErr != nil {
					_ = os.Remove(rollbackPath)
					rollbackErrors = append(
						rollbackErrors,
						fmt.Errorf("restore Xray binary: %w", writeErr),
					)
				}
			}
		} else {
			if removeErr := os.Remove(binaryPath); removeErr != nil &&
				!os.IsNotExist(removeErr) {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("remove newly installed Xray binary: %w", removeErr),
				)
			}
		}

		// Restore the original service's running state.
		if serviceWasActive {
			if restartErr := runSystemctl("restart", "xray"); restartErr != nil {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("restart Xray after rollback: %w", restartErr),
				)
			}
		}

		if len(rollbackErrors) > 0 {
			return fmt.Errorf(
				"%w; rollback encountered errors: %v",
				err,
				rollbackErrors,
			)
		}

		return err
	}

	return nil
}

func prepareEmbeddedXray() (string, func(), error) {
	if err := requireCommand("unzip"); err != nil {
		return "", nil, err
	}

	tmpDir, err := os.MkdirTemp("", "plachta-xray-*")
	if err != nil {
		return "", nil, fmt.Errorf("create temporary Xray directory: %w", err)
	}

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	zipPath := filepath.Join(tmpDir, "Xray-linux-64.zip")
	if err := os.WriteFile(zipPath, xrayPackage, 0600); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("write embedded Xray archive: %w", err)
	}

	if err := runCommand("unzip", "-o", zipPath, "-d", tmpDir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("extract embedded Xray archive: %w", err)
	}

	binaryPath := filepath.Join(tmpDir, "xray")
	if _, err := os.Stat(binaryPath); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("embedded Xray binary not found: %w", err)
	}

	if err := os.Chmod(binaryPath, 0755); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("set Xray permissions: %w", err)
	}

	return binaryPath, cleanup, nil
}

func installReality() error {
	fmt.Println("Installing Reality...")

	if err := requireCommand("systemctl"); err != nil {
		return err
	}

	xrayPath, cleanupXray, err := prepareEmbeddedXray()
	if err != nil {
		return err
	}
	defer cleanupXray()

	if err := GenerateConfigWithBinary(xrayPath); err != nil {
		return fmt.Errorf("failed to generate Reality config: %w", err)
	}

	if err := validateGeneratedConfigWithBinary(xrayPath); err != nil {
		return err
	}

	config := New()
	uri, err := config.uriWithBinary(xrayPath)
	if err != nil {
		return fmt.Errorf("failed to generate Reality URI: %w", err)
	}

	fmt.Println("Reality URI:")
	fmt.Println(uri)

	if err := installXrayBinaryFromPath(xrayPath); err != nil {
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
		if os.IsNotExist(err) {
			return "", nil
		}

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

func installXrayBinaryFromPath(sourcePath string) error {
	const targetPath = "/usr/local/bin/xray"
	const tempPath = "/usr/local/bin/xray.new"

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("create Xray directory: %w", err)
	}

	if err := requireCommand("systemctl"); err != nil {
		return err
	}

	if err := os.Remove(tempPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove temporary Xray binary: %w", err)
	}

	if err := copyFile(sourcePath, tempPath); err != nil {
		return fmt.Errorf("copy Xray binary: %w", err)
	}

	if err := os.Chmod(tempPath, 0755); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("set Xray binary permissions: %w", err)
	}

	serviceWasActive := exec.Command("systemctl", "is-active", "--quiet", "xray").Run() == nil

	if serviceWasActive {
		if err := runSystemctl("stop", "xray"); err != nil {
			_ = os.Remove(tempPath)
			return fmt.Errorf("stop Xray service: %w", err)
		}
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		_ = os.Remove(tempPath)

		if serviceWasActive {
			if restartErr := runSystemctl("start", "xray"); restartErr != nil {
				return fmt.Errorf(
					"install Xray binary: %w; additionally failed to restart Xray: %v",
					err,
					restartErr,
				)
			}
		}

		return fmt.Errorf("install Xray binary: %w", err)
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

	return nil
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
	return validateGeneratedConfigWithBinary("xray")
}

func validateGeneratedConfigWithBinary(xrayPath string) error {
	output, err := exec.Command(xrayPath, "run", "-test", "-config", defaultRealityConfigPath).CombinedOutput()
	if err != nil {
		return fmt.Errorf("xray config validation failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func installXrayService() error {
	serviceWasActive := exec.Command("systemctl", "is-active", "--quiet", "xray").Run() == nil
	type fileBackup struct {
		path   string
		data   []byte
		mode   os.FileMode
		exists bool
	}

	paths := []string{
		xrayServicePath,
		xrayServiceDropInPath,
	}

	backups := make([]fileBackup, 0, len(paths))

	for _, path := range paths {
		backup := fileBackup{path: path}

		info, err := os.Stat(path)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to inspect service file %s: %w", path, err)
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to back up service file %s: %w", path, err)
			}

			backup.data = data
			backup.mode = info.Mode().Perm()
			backup.exists = true
		}

		backups = append(backups, backup)
	}

	rollback := func(cause error) error {
		var rollbackErrors []error

		for i := len(backups) - 1; i >= 0; i-- {
			backup := backups[i]

			if backup.exists {
				if err := os.WriteFile(backup.path, backup.data, backup.mode); err != nil {
					rollbackErrors = append(
						rollbackErrors,
						fmt.Errorf("restore %s: %w", backup.path, err),
					)
					continue
				}

				if err := os.Chmod(backup.path, backup.mode); err != nil {
					rollbackErrors = append(
						rollbackErrors,
						fmt.Errorf("restore permissions for %s: %w", backup.path, err),
					)
				}
			} else {
				if err := os.Remove(backup.path); err != nil && !os.IsNotExist(err) {
					rollbackErrors = append(
						rollbackErrors,
						fmt.Errorf("remove newly created file %s: %w", backup.path, err),
					)
				}
			}
		}

		if err := runSystemctl("daemon-reload"); err != nil {
			rollbackErrors = append(rollbackErrors, err)
		}

		if serviceWasActive {
			if err := runSystemctl("restart", "xray"); err != nil {
				rollbackErrors = append(
					rollbackErrors,
					fmt.Errorf("restart Xray after rollback: %w", err),
				)
			}
		}
		if len(rollbackErrors) > 0 {
			return fmt.Errorf("%w; rollback also encountered errors: %v", cause, errors.Join(rollbackErrors...))
		}

		return cause
	}

	if err := os.WriteFile(xrayServicePath, []byte(xrayService), 0644); err != nil {
		return rollback(fmt.Errorf("failed to write xray systemd service: %w", err))
	}

	if err := os.MkdirAll(filepath.Dir(xrayServiceDropInPath), 0755); err != nil {
		return rollback(fmt.Errorf("failed to create Xray service drop-in directory: %w", err))
	}

	if err := os.WriteFile(xrayServiceDropInPath, []byte(xrayServiceDropIn), 0644); err != nil {
		return rollback(fmt.Errorf("failed to write Plachta Xray service drop-in: %w", err))
	}

	if err := runSystemctl("daemon-reload"); err != nil {
		return rollback(err)
	}

	if err := runSystemctl("enable", "xray"); err != nil {
		return rollback(err)
	}

	if err := runSystemctl("restart", "xray"); err != nil {
		return rollback(err)
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

func copyFile(sourcePath, targetPath string) error {
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, data, 0644)
}
