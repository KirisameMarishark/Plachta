# Plachta Architecture

## Philosophy

Plachta is a modular VPS deployment framework.

Goals:

- Simplicity
- Maintainability
- Reproducibility
- Security
- Extensibility

Every feature is implemented as an independent module.

Core never depends on modules.

Modules never depend on each other.

---

# Repository Layout

```
cmd/
internal/
lib/
modules/
templates/
runtime/
docs/
```

---

# Core

Core provides common capabilities shared by all modules.

```
Core
├── CLI
├── Config
├── Logger
├── Error
├── Package
├── Service
├── Download
├── Metadata
├── System
└── Prerequisite
```

Core libraries must never contain protocol-specific logic.

---

# Module Layout

Every module must contain:

```
install.sh
verify.sh
config.sh
service.sh
README.md
```

Optional:

```
backup.sh
health.sh
show.sh
uninstall.sh
```

Modules never call other modules directly.

Modules communicate only through Core libraries.

---

# Internal Layout

Each protocol owns an internal implementation directory.

Example:

```
internal/
    reality/
        generate.sh
        read.sh
        write.sh
        uri.sh
        validate.sh

    tuic/

    hysteria2/
```

Responsibilities:

```
generate.sh
    Generate configuration

read.sh
    Load configuration

write.sh
    Modify configuration

uri.sh
    Generate client URI

validate.sh
    Validate configuration
```

No file should have multiple responsibilities.

---

# Single Source of Truth

Every module has exactly one configuration source.

For Reality:

```
/etc/plachta/reality/config.json
```

All commands must read from the same configuration.

Never duplicate runtime state.

Temporary files are allowed.

Persistent duplicated configuration is not.

---

# Templates

Templates are read-only.

Examples:

```
xray.json

nftables.conf

systemd.service
```

Modules generate runtime configuration from templates.

Templates must never be modified directly.

---

# Configuration Priority

Configuration precedence:

```
CLI Arguments

↓

User Configuration

↓

Module Defaults

↓

Built-in Defaults
```

---

# Plugin Principles

Every module must be:

- Installable
- Verifiable
- Replaceable
- Idempotent

Running install multiple times must never break an existing deployment.

---

# Error Handling

Core handles fatal errors.

Modules return status.

Never exit unexpectedly inside reusable libraries.

---

# Public Interface

Every module exposes the same interface.

Example:

```
plachta install reality

plachta verify reality

plachta show reality

plachta config reality
```

Future modules must follow the exact same command structure.

---

# Coding Principles

Each function should have one responsibility.

Avoid duplicated logic.

Avoid duplicated configuration parsing.

Shared functionality belongs in Core or internal/.

---

# Future Roadmap

Version 1

- Debian 12
- Xray Reality
- nftables
- SSH Hardening
- Fail2ban
- DNS Optimization

Version 2

- TUIC
- Hysteria2
- Sing-box

Version 3

- Subscription
- Auto IP Update
- Multiple Node Management

Version 4

- Docker
- Web UI
- REST API