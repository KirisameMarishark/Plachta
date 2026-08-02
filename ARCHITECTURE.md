# Plachta Architecture

## Overview

Plachta is a modular VPS deployment framework focused on:

- Security
- Maintainability
- Automation
- Low fingerprint
- Reproducible deployment

Every feature is implemented as an independent module.

---

# Core

Core libraries never depend on modules.

```
Core
├── CLI
├── Config
├── Logger
├── Error
├── Package
├── Service
├── Prerequisite
├── Metadata
└── System
```

---

# Modules

Each module must contain:

```
README.md
install.sh
verify.sh
config.sh
service.sh
```

Optional:

```
backup.sh
health.sh
uninstall.sh
```

Modules should never call each other directly.

Communication must go through Core libraries.

---

# Templates

Templates are read-only.

Examples:

- nftables.conf
- xray.json
- systemd.service

Modules generate runtime configuration from templates.

---

# Configuration

Configuration priority:

```
CLI
↓

User Config

↓

Default Config

↓

Built-in Defaults
```

---

# Plugin Principle

Every module should be:

- Installable
- Verifiable
- Replaceable
- Idempotent

Running install twice must not break the system.

---

# Error Handling

Only Core handles fatal errors.

Modules return status.

---

# Future

Plachta v1 focuses on:

- Debian 12
- Xray Reality
- TUIC
- nftables
- SSH Hardening
- Fail2ban
- DNS Optimization

Future versions may support:

- Hysteria2
- Sing-box
- Docker
- Web UI