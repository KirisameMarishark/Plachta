package firewall

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Verify() error {
	if err := requireCommand("nft"); err != nil {
		return err
	}

	fmt.Println("Firewall Verify")
	fmt.Println("------------------------------")

	fmt.Println("✔ nft command found.")

	enabled, err := systemctlState("is-enabled", "nftables")
	if err != nil {
		fmt.Println("⚠ nftables service not enabled.")
	} else {
		fmt.Printf("✔ nftables service enabled: %s\n", enabled)
	}

	active, err := systemctlState("is-active", "nftables")
	if err != nil {
		fmt.Println("⚠ nftables service not running.")
	} else {
		fmt.Printf("✔ nftables service running: %s\n", active)
	}

	fmt.Println()
	fmt.Println("Ruleset")
	fmt.Println("------------------------------")

	if err := listRuleset(); err != nil {
		return err
	}

	return nil
}

func systemctlState(action string, service string) (string, error) {
	cmd := exec.Command("systemctl", action, service)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return "", fmt.Errorf(
				"systemctl %s %s failed: %w: %s",
				action,
				service,
				err,
				strings.TrimSpace(string(output)),
			)
		}

		return "", fmt.Errorf(
			"systemctl %s %s failed: %w",
			action,
			service,
			err,
		)
	}

	return strings.TrimSpace(string(output)), nil
}

func listRuleset() error {
	cmd := exec.Command("nft", "list", "ruleset")

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) > 0 {
			return fmt.Errorf(
				"failed to list nftables ruleset: %w: %s",
				err,
				strings.TrimSpace(string(output)),
			)
		}

		return fmt.Errorf("failed to list nftables ruleset: %w", err)
	}

	fmt.Print(string(output))
	return nil
}

func requireRootForVerify() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("firewall verify requires root privileges")
	}

	return nil
}
