package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func defaultConfigPath() string {
	if path := os.Getenv("PLACTHA_CONFIG"); path != "" {
		return path
	}

	if runtime.GOOS == "windows" {
		return filepath.Join("configs", "default", "plachta.conf")
	}

	return "/etc/plachta/plachta.conf"
}

type Config struct {
	Path string
}

func New() Config {
	return Config{
		Path: defaultConfigPath(),
	}
}

func (c Config) Exists() bool {
	_, err := os.Stat(c.Path)
	return err == nil
}

func (c Config) Read() (string, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c Config) Get(key string) (string, error) {
	data, err := c.Read()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		if strings.TrimSpace(parts[0]) == key {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", os.ErrNotExist
}

func (c Config) Set(key, value string) error {
	dir := filepath.Dir(c.Path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := c.Read()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(data, "\n")
	found := false

	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		if strings.TrimSpace(parts[0]) == key {
			lines[i] = key + "=" + value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, key+"="+value)
	}

	return os.WriteFile(
		c.Path,
		[]byte(strings.Join(lines, "\n")),
		0644,
	)
}
