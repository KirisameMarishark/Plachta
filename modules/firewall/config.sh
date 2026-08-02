#!/usr/bin/env bash

set -euo pipefail

RULESET="/etc/nftables.conf"

generate_firewall_config() {
cat > "${RULESET}" <<'EOF'
#!/usr/sbin/nft -f

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
EOF

log_success "Firewall configuration generated."
}