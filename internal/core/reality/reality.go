package reality

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultConfigPath = "/etc/plachta/reality/config.json"

type Config struct {
	Path string
}

type xrayConfig struct {
	Inbounds []struct {
		Port           int            `json:"port"`
		Settings       settings       `json:"settings"`
		StreamSettings streamSettings `json:"streamSettings"`
	} `json:"inbounds"`
}

type settings struct {
	Clients []struct {
		ID string `json:"id"`
	} `json:"clients"`
}

type streamSettings struct {
	RealitySettings realitySettings `json:"realitySettings"`
}

type realitySettings struct {
	PrivateKey  string   `json:"privateKey"`
	ShortIDs    []string `json:"shortIds"`
	ServerNames []string `json:"serverNames"`
}

func New() Config {
	return Config{
		Path: defaultConfigPath,
	}
}

func (c Config) Read() (xrayConfig, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		return xrayConfig{}, err
	}

	var cfg xrayConfig

	if err := json.Unmarshal(data, &cfg); err != nil {
		return xrayConfig{}, err
	}

	if len(cfg.Inbounds) == 0 {
		return xrayConfig{}, fmt.Errorf("no inbounds found")
	}

	return cfg, nil
}

func (c Config) UUID() (string, error) {
	cfg, err := c.Read()
	if err != nil {
		return "", err
	}

	if len(cfg.Inbounds[0].Settings.Clients) == 0 {
		return "", fmt.Errorf("no Reality client found")
	}

	return cfg.Inbounds[0].Settings.Clients[0].ID, nil
}

func (c Config) PrivateKey() (string, error) {
	cfg, err := c.Read()
	if err != nil {
		return "", err
	}

	return cfg.Inbounds[0].StreamSettings.RealitySettings.PrivateKey, nil
}

func (c Config) ShortID() (string, error) {
	cfg, err := c.Read()
	if err != nil {
		return "", err
	}

	ids := cfg.Inbounds[0].StreamSettings.RealitySettings.ShortIDs

	if len(ids) == 0 {
		return "", fmt.Errorf("no Reality short ID found")
	}

	return ids[0], nil
}

func (c Config) ServerName() (string, error) {
	cfg, err := c.Read()
	if err != nil {
		return "", err
	}

	names := cfg.Inbounds[0].StreamSettings.RealitySettings.ServerNames

	if len(names) == 0 {
		return "", fmt.Errorf("no Reality server name found")
	}

	return names[0], nil
}

func (c Config) Port() (int, error) {
	cfg, err := c.Read()
	if err != nil {
		return 0, err
	}

	return cfg.Inbounds[0].Port, nil
}
