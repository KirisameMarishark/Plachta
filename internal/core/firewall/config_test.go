package firewall

import (
	"strings"
	"testing"
)

func TestRuleset(t *testing.T) {
	required := []string{
		"#!/usr/sbin/nft -f",
		"flush ruleset",
		"table inet filter",
		"chain input",
		"policy drop",
		"iif lo accept",
		"ct state established,related accept",
		"tcp dport 22 accept",
		"tcp dport 443 accept",
		"udp dport 443 accept",
		"ip protocol icmp accept",
		"ip6 nexthdr ipv6-icmp accept",
		"chain forward",
		"chain output",
		"policy accept",
	}

	for _, item := range required {
		if !strings.Contains(ruleset, item) {
			t.Fatalf("ruleset is missing %q", item)
		}
	}
}

func TestRulesetPath(t *testing.T) {
	if rulesetPath != "/etc/nftables.conf" {
		t.Fatalf("unexpected ruleset path: %q", rulesetPath)
	}
}
