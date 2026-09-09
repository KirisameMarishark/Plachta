package firewall

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Install() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("firewall install requires root privileges")
	}

	if err := requireDebian(); err != nil {
		return err
	}

	if err := requireCommand("systemctl"); err != nil {
		return err
	}

	if err := ensurePackage("nftables"); err != nil {
		return err
	}

	if err := GenerateConfig(); err != nil {
		return err
	}

	if err := runCommand("systemctl", "enable", "nftables"); err != nil {
		return fmt.Errorf("failed to enable nftables: %w", err)
	}

	if err := runCommand("systemctl", "start", "nftables"); err != nil {
		return fmt.Errorf("failed to start nftables: %w", err)
	}

	if err := runCommand("nft", "-f", rulesetPath); err != nil {
		return fmt.Errorf("failed to load firewall rules: %w", err)
	}

	return nil
}

func requireDebian() error {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return fmt.Errorf("cannot determine operating system: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, "ID=") {
			continue
		}

		id := strings.Trim(strings.TrimPrefix(line, "ID="), `"`)

		if id != "debian" {
			return fmt.Errorf("only Debian is supported")
		}

		return nil
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read /etc/os-release: %w", err)
	}

	return fmt.Errorf("operating system ID not found in /etc/os-release")
}

func ensurePackage(name string) error {
	cmd := exec.Command("dpkg", "-s", name)
	if err := cmd.Run(); err == nil {
		return nil
	}

	if err := runCommand("apt-get", "update"); err != nil {
		return fmt.Errorf("failed to update package index: %w", err)
	}

	if err := runCommand("apt-get", "install", "-y", name); err != nil {
		return fmt.Errorf("failed to install package %q: %w", name, err)
	}

	return nil
}

func requireCommand(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("required command %q not found: %w", name, err)
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return fmt.Errorf(
				"command %q failed: %w: %s",
				name,
				err,
				string(output),
			)
		}

		return fmt.Errorf("command %q failed: %w", name, err)
	}

	return nil
}
