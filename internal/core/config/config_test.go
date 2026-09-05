package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigGet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.conf")

	data := `# Plachta Configuration

DNS=1.1.1.1
MTU=1500
LOG_LEVEL=info
`

	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Path: path}

	tests := []struct {
		key  string
		want string
	}{
		{"DNS", "1.1.1.1"},
		{"MTU", "1500"},
		{"LOG_LEVEL", "info"},
	}

	for _, tt := range tests {
		got, err := cfg.Get(tt.key)
		if err != nil {
			t.Fatalf("Get(%q) returned error: %v", tt.key, err)
		}

		if got != tt.want {
			t.Fatalf("Get(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestConfigGetMissingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.conf")

	if err := os.WriteFile(path, []byte("DNS=1.1.1.1\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Path: path}

	_, err := cfg.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestConfigSetExistingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.conf")

	if err := os.WriteFile(
		path,
		[]byte("DNS=1.1.1.1\nMTU=1500\n"),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Path: path}

	if err := cfg.Set("DNS", "8.8.8.8"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)

	if !strings.Contains(content, "DNS=8.8.8.8") {
		t.Fatalf("updated DNS not found: %q", content)
	}

	if strings.Contains(content, "DNS=1.1.1.1") {
		t.Fatalf("old DNS value still exists: %q", content)
	}
}

func TestConfigSetNewKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.conf")

	if err := os.WriteFile(
		path,
		[]byte("DNS=1.1.1.1\n"),
		0644,
	); err != nil {
		t.Fatal(err)
	}

	cfg := Config{Path: path}

	if err := cfg.Set("MTU", "1500"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "MTU=1500") {
		t.Fatalf("new key not found: %q", string(data))
	}
}

func TestConfigSetCreatesMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.conf")

	cfg := Config{Path: path}

	if err := cfg.Set("DNS", "1.1.1.1"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "\nDNS=1.1.1.1" {
		t.Fatalf("unexpected config content: %q", string(data))
	}
}

func TestConfigInitCreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".config", "plachta", "config.conf")

	cfg := Config{Path: path}

	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := `# Plachta Configuration

DNS=1.1.1.1
MTU=1500
LOG_LEVEL=info
`

	if string(data) != want {
		t.Fatalf("unexpected config content: %q", string(data))
	}
}
