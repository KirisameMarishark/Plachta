package reality

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func (c Config) URI() (string, error) {
	uuid, err := c.UUID()
	if err != nil {
		return "", err
	}

	port, err := c.Port()
	if err != nil {
		return "", err
	}

	shortID, err := c.ShortID()
	if err != nil {
		return "", err
	}

	serverName, err := c.ServerName()
	if err != nil {
		return "", err
	}

	publicKey, err := c.PublicKey()
	if err != nil {
		return "", err
	}

	serverIP, err := publicIPv4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"vless://%s@%s:%d?type=tcp&security=reality&pbk=%s&sid=%s&sni=%s&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality",
		uuid,
		serverIP,
		port,
		publicKey,
		shortID,
		serverName,
	), nil
}

func (c Config) PublicKey() (string, error) {
	privateKey, err := c.PrivateKey()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("xray", "x25519", "-i", privateKey)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate Reality public key: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Password (PublicKey):") {
			return strings.TrimSpace(
				strings.TrimPrefix(line, "Password (PublicKey):"),
			), nil
		}

		if strings.HasPrefix(line, "Public key:") {
			return strings.TrimSpace(
				strings.TrimPrefix(line, "Public key:"),
			), nil
		}
	}

	return "", fmt.Errorf("Reality public key not found in xray output")
}

func publicIPv4() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("failed to get public IPv4: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"public IPv4 service returned HTTP %d",
			resp.StatusCode,
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read public IPv4: %w", err)
	}

	ip := strings.TrimSpace(string(data))

	if ip == "" {
		return "", fmt.Errorf("public IPv4 is empty")
	}

	return ip, nil
}
