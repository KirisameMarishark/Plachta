package firewall

import (
	"fmt"
	"os"
	"os/exec"
)

func Install() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("firewall install requires root privileges")
	}

	if err := requireCommand("systemctl"); err != nil {
		return err
	}

	if err := requireCommand("nft"); err != nil {
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
