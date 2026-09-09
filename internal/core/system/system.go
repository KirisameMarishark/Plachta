package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Info struct {
	OS           string
	Architecture string
	Hostname     string
	Version      string
	Kernel       string
	IPv4         string
	IPv6         string
	MemoryMB     int
	CPUCount     int
}

func Get() Info {
	info := Info{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Hostname:     hostname(),
	}

	if value, err := osReleaseValue("ID"); err == nil {
		info.OS = value
	}

	if value, err := osReleaseValue("VERSION_ID"); err == nil {
		info.Version = value
	}

	if value, err := commandOutput("uname", "-r"); err == nil {
		info.Kernel = value
	}

	if value, err := localIPv4(); err == nil {
		info.IPv4 = value
	}

	if value, err := localIPv6(); err == nil {
		info.IPv6 = value
	}

	if value, err := memoryMB(); err == nil {
		info.MemoryMB = value
	}

	if value, err := cpuCount(); err == nil {
		info.CPUCount = value
	}

	return info
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	return name
}

func osReleaseValue(key string) (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	prefix := key + "="

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !strings.HasPrefix(line, prefix) {
			continue
		}

		value := strings.TrimPrefix(line, prefix)
		value = strings.Trim(value, `"`)

		return value, nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("key %q not found in /etc/os-release", key)
}

func cpuCount() (int, error) {
	output, err := commandOutput("nproc")
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(output)
	if err != nil {
		return 0, fmt.Errorf("invalid CPU count %q: %w", output, err)
	}

	return count, nil
}

func memoryMB() (int, error) {
	output, err := commandOutput("free", "-m")
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)

		if len(fields) < 2 || fields[0] != "Mem:" {
			continue
		}

		memory, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, fmt.Errorf("invalid memory value %q: %w", fields[1], err)
		}

		return memory, nil
	}

	return 0, fmt.Errorf("Mem row not found in free output")
}

func localIPv4() (string, error) {
	output, err := commandOutput(
		"sh",
		"-c",
		"ip -4 route get 1.1.1.1 2>/dev/null | awk '{print $7; exit}'",
	)
	if err != nil {
		return "", err
	}

	if output == "" {
		return "", fmt.Errorf("IPv4 address not found")
	}

	return output, nil
}

func localIPv6() (string, error) {
	output, err := commandOutput(
		"sh",
		"-c",
		"ip -6 route get 2606:4700:4700::1111 2>/dev/null | awk '{for (i=1; i<=NF; i++) if ($i==\"src\") {print $(i+1); exit}}'",
	)
	if err != nil {
		return "", err
	}

	if output == "" {
		return "", fmt.Errorf("IPv6 address not found")
	}

	return output, nil
}

func commandOutput(name string, args ...string) (string, error) {
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"command %q failed: %w: %s",
			name,
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return strings.TrimSpace(string(output)), nil
}
