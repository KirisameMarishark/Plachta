package reality

import (
	"encoding/json"
	"testing"
)

func TestBuildRealityConfig(t *testing.T) {
	cfg := buildRealityConfig(
		"uuid-test",
		GeneratedKeyPair{
			PrivateKey: "private-test",
			PublicKey:  "public-test",
			Hash32:     "hash-test",
		},
		"shortid-test",
	)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal config: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("generated config is invalid JSON: %v", err)
	}

	inbounds := decoded["inbounds"].([]interface{})
	inbound := inbounds[0].(map[string]interface{})

	if inbound["port"] != float64(443) {
		t.Fatalf("unexpected port: %v", inbound["port"])
	}

	if inbound["protocol"] != "vless" {
		t.Fatalf("unexpected protocol: %v", inbound["protocol"])
	}

	settings := inbound["settings"].(map[string]interface{})
	clients := settings["clients"].([]interface{})
	client := clients[0].(map[string]interface{})

	if client["id"] != "uuid-test" {
		t.Fatalf("unexpected UUID: %v", client["id"])
	}

	if client["flow"] != "xtls-rprx-vision" {
		t.Fatalf("unexpected flow: %v", client["flow"])
	}

	stream := inbound["streamSettings"].(map[string]interface{})

	if stream["network"] != "tcp" {
		t.Fatalf("unexpected network: %v", stream["network"])
	}

	if stream["security"] != "reality" {
		t.Fatalf("unexpected security: %v", stream["security"])
	}

	reality := stream["realitySettings"].(map[string]interface{})

	if reality["show"] != false {
		t.Fatalf("unexpected show: %v", reality["show"])
	}

	if reality["dest"] != "www.cloudflare.com:443" {
		t.Fatalf("unexpected destination: %v", reality["dest"])
	}

	if reality["privateKey"] != "private-test" {
		t.Fatalf("unexpected private key: %v", reality["privateKey"])
	}

	shortIDs := reality["shortIds"].([]interface{})
	if shortIDs[0] != "shortid-test" {
		t.Fatalf("unexpected short ID: %v", shortIDs[0])
	}

	serverNames := reality["serverNames"].([]interface{})
	if serverNames[0] != "www.cloudflare.com" {
		t.Fatalf("unexpected server name: %v", serverNames[0])
	}
}

func TestBuildClientEnv(t *testing.T) {
	got := buildClientEnv(
		"uuid-test",
		"public-test",
		"short-test",
	)

	want := "UUID=uuid-test\nPUBLIC_KEY=public-test\nSHORT_ID=short-test\nSERVER_NAME=www.cloudflare.com\nPORT=443\n"

	if got != want {
		t.Fatalf("unexpected client.env:\n%s\nwant:\n%s", got, want)
	}
}
