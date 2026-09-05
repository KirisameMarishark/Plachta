package reality

import (
	"fmt"
	"os/exec"
	"strings"
)

func GenerateUUID() (string, error) {
	output, err := exec.Command("xray", "uuid").Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	uuid := strings.TrimSpace(string(output))
	if uuid == "" {
		return "", fmt.Errorf("generated UUID is empty")
	}

	return uuid, nil
}

func GenerateShortID() (string, error) {
	output, err := exec.Command("openssl", "rand", "-hex", "8").Output()
	if err != nil {
		return "", fmt.Errorf("failed to generate Reality short ID: %w", err)
	}

	shortID := strings.TrimSpace(string(output))
	if shortID == "" {
		return "", fmt.Errorf("generated Reality short ID is empty")
	}

	return shortID, nil
}

type GeneratedKeyPair struct {
	PrivateKey string
	PublicKey  string
	Hash32     string
}

func GenerateKeyPair() (GeneratedKeyPair, error) {
	output, err := exec.Command("xray", "x25519").CombinedOutput()
	if err != nil {
		return GeneratedKeyPair{}, fmt.Errorf(
			"failed to generate Reality keypair: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return parseKeyPairOutput(string(output))
}

func parseKeyPairOutput(output string) (GeneratedKeyPair, error) {
	var result GeneratedKeyPair

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "PrivateKey:"):
			result.PrivateKey = strings.TrimSpace(
				strings.TrimPrefix(line, "PrivateKey:"),
			)

		case strings.HasPrefix(line, "Password (PublicKey):"):
			result.PublicKey = strings.TrimSpace(
				strings.TrimPrefix(line, "Password (PublicKey):"),
			)

		case strings.HasPrefix(line, "Hash32:"):
			result.Hash32 = strings.TrimSpace(
				strings.TrimPrefix(line, "Hash32:"),
			)
		}
	}

	if result.PrivateKey == "" {
		return GeneratedKeyPair{}, fmt.Errorf(
			"Reality private key not found in xray output",
		)
	}

	if result.PublicKey == "" {
		return GeneratedKeyPair{}, fmt.Errorf(
			"Reality public key not found in xray output",
		)
	}

	if result.Hash32 == "" {
		return GeneratedKeyPair{}, fmt.Errorf(
			"Reality hash32 not found in xray output",
		)
	}

	return result, nil
}
