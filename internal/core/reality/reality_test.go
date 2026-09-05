package reality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRealityConfigRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	config := `{
		"inbounds": [
			{
				"port": 443,
				"settings": {
					"clients": [
						{
							"id": "test-uuid"
						}
					]
				},
				"streamSettings": {
					"security": "reality",
					"realitySettings": {
						"privateKey": "test-private-key",
						"shortIds": [
							"test-short-id"
						],
						"serverNames": [
							"www.cloudflare.com"
						]
					}
				}
			}
		]
	}`

	if err := os.WriteFile(path, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := NewWithPath(path)

	if got, err := cfg.UUID(); err != nil || got != "test-uuid" {
		t.Fatalf("UUID() = %q, %v", got, err)
	}

	if got, err := cfg.PrivateKey(); err != nil || got != "test-private-key" {
		t.Fatalf("PrivateKey() = %q, %v", got, err)
	}

	if got, err := cfg.ShortID(); err != nil || got != "test-short-id" {
		t.Fatalf("ShortID() = %q, %v", got, err)
	}

	if got, err := cfg.ServerName(); err != nil || got != "www.cloudflare.com" {
		t.Fatalf("ServerName() = %q, %v", got, err)
	}

	if got, err := cfg.Port(); err != nil || got != 443 {
		t.Fatalf("Port() = %d, %v", got, err)
	}
}

func TestParseGeneratedKeyPair(t *testing.T) {
	output := `
PrivateKey: private-test
Password (PublicKey): public-test
Hash32: hash-test
`

	result, err := parseKeyPairOutput(output)
	if err != nil {
		t.Fatal(err)
	}

	if result.PrivateKey != "private-test" {
		t.Fatalf("PrivateKey = %q", result.PrivateKey)
	}

	if result.PublicKey != "public-test" {
		t.Fatalf("PublicKey = %q", result.PublicKey)
	}

	if result.Hash32 != "hash-test" {
		t.Fatalf("Hash32 = %q", result.Hash32)
	}
}
