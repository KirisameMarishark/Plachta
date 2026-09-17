package cli

import (
	"fmt"
	"os/exec"
	"strings"
)

func ShowStatus() {
	fmt.Println("Plachta Status")
	fmt.Println("------------------------------")

	showServiceStatus("xray")
}

func showServiceStatus(service string) {
	cmd := exec.Command("systemctl", "is-active", service)
	output, err := cmd.Output()

	status := strings.TrimSpace(string(output))

	if err != nil || status == "" {
		status = "inactive"
	}

	fmt.Printf("%-14s: %s\n", service, status)
}
