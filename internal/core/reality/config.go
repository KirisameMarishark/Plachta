package reality

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	defaultRealityDir        = "/etc/plachta/reality"
	defaultRealityConfigPath = "/etc/plachta/reality/config.json"
	defaultRealityClientEnv  = "/etc/plachta/reality/client.env"

	defaultRealityPort        = 443
	defaultRealityDestination = "www.cloudflare.com:443"
	defaultRealityServerName  = "www.cloudflare.com"
)

type generatedRealityConfig struct {
	Log       realityLog        `json:"log"`
	Inbounds  []realityInbound  `json:"inbounds"`
	Outbounds []realityOutbound `json:"outbounds"`
}

type realityLog struct {
	LogLevel string `json:"loglevel"`
}

type realityInbound struct {
	Listen         string                `json:"listen"`
	Port           int                   `json:"port"`
	Protocol       string                `json:"protocol"`
	Settings       realityClientSettings `json:"settings"`
	StreamSettings realityStreamSettings `json:"streamSettings"`
}

type realityClientSettings struct {
	Clients    []realityClient `json:"clients"`
	Decryption string          `json:"decryption"`
}

type realityClient struct {
	ID   string `json:"id"`
	Flow string `json:"flow"`
}

type realityStreamSettings struct {
	Network         string          `json:"network"`
	Security        string          `json:"security"`
	RealitySettings realitySettings `json:"realitySettings"`
}

type realityOutbound struct {
	Protocol string `json:"protocol"`
}

func buildRealityConfig(uuid string, keyPair GeneratedKeyPair, shortID string) generatedRealityConfig {
	return generatedRealityConfig{
		Log: realityLog{
			LogLevel: "warning",
		},
		Inbounds: []realityInbound{
			{
				Listen:   "0.0.0.0",
				Port:     defaultRealityPort,
				Protocol: "vless",
				Settings: realityClientSettings{
					Clients: []realityClient{
						{
							ID:   uuid,
							Flow: "xtls-rprx-vision",
						},
					},
					Decryption: "none",
				},
				StreamSettings: realityStreamSettings{
					Network:  "tcp",
					Security: "reality",
					RealitySettings: realitySettings{
						Show:        false,
						Dest:        defaultRealityDestination,
						XVer:        0,
						ServerNames: []string{defaultRealityServerName},
						PrivateKey:  keyPair.PrivateKey,
						ShortIDs:    []string{shortID},
					},
				},
			},
		},
		Outbounds: []realityOutbound{
			{
				Protocol: "freedom",
			},
		},
	}
}

func GenerateConfig() error {
	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	keyPair, err := GenerateKeyPair()
	if err != nil {
		return err
	}

	shortID, err := GenerateShortID()
	if err != nil {
		return err
	}

	cfg := buildRealityConfig(uuid, keyPair, shortID)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Reality config: %w", err)
	}

	data = append(data, '\n')

	if err := os.MkdirAll(defaultRealityDir, 0755); err != nil {
		return fmt.Errorf("failed to create Reality directory: %w", err)
	}

	if err := os.WriteFile(defaultRealityConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write Reality config: %w", err)
	}

	clientEnv := buildClientEnv(uuid, keyPair.PublicKey, shortID)

	if err := os.WriteFile(defaultRealityClientEnv, []byte(clientEnv), 0644); err != nil {
		return fmt.Errorf("failed to write Reality client environment: %w", err)
	}

	return nil
}

func buildClientEnv(uuid, publicKey, shortID string) string {
	return fmt.Sprintf(
		"UUID=%s\nPUBLIC_KEY=%s\nSHORT_ID=%s\nSERVER_NAME=%s\nPORT=%d\n",
		uuid,
		publicKey,
		shortID,
		defaultRealityServerName,
		defaultRealityPort,
	)
}
