package firewall

import (
	"fmt"
	"os"
)

const rulesetPath = "/etc/nftables.conf"

const ruleset = `#!/usr/sbin/nft -f

flush ruleset

table inet filter {

    chain input {

        type filter hook input priority 0;

        policy drop;

        iif lo accept

        ct state established,related accept

        tcp dport 22 accept

        tcp dport 443 accept

        udp dport 443 accept

        ip protocol icmp accept

        ip6 nexthdr ipv6-icmp accept
    }

    chain forward {

        type filter hook forward priority 0;

        policy drop;
    }

    chain output {

        type filter hook output priority 0;

        policy accept;
    }
}
`

func GenerateConfig() error {
	if err := os.WriteFile(rulesetPath, []byte(ruleset), 0644); err != nil {
		return fmt.Errorf("failed to write firewall ruleset: %w", err)
	}

	return nil
}
