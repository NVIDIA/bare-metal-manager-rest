<!--
SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->

# Carbide CLI

Command-line client for the NVIDIA Bare Metal Manager REST API. Commands are dynamically generated from the embedded OpenAPI spec at startup, so every API endpoint is available with zero manual command code.

## Prerequisites

- Go 1.25.4 or later
- Access to a running Bare Metal Manager REST API instance (local via `make kind-reset` or remote)

## Installation

### From the repo (recommended)

```bash
make carbide-cli
```

This builds and installs `carbidecli` to `$(go env GOPATH)/bin/carbidecli`. Override the destination with:

```bash
make carbide-cli INSTALL_DIR=/usr/local/bin
```

### With go install

```bash
go install ./cli/cmd/carbidecli
```

### Manual go build

```bash
go build -o /usr/local/bin/carbidecli ./cli/cmd/carbidecli
```

### Verify

```bash
carbidecli --version
```

## Quick Start

```bash
carbidecli init                    # writes ~/.bmm/config.yaml
```

Edit `~/.bmm/config.yaml` with your server URL, org, and auth settings, then:

```bash
carbidecli login                   # exchange credentials for a token
carbidecli site list               # list all sites
```

## Configuration

Config file: `~/.bmm/config.yaml`

```yaml
api:
  base: http://localhost:8388
  org: test-org
  name: carbide                # API path segment (default)

auth:
  # Option 1: Direct bearer token
  # token: eyJhbGciOi...

  # Option 2: OIDC provider (e.g. Keycloak)
  oidc:
    token_url: http://localhost:8080/realms/carbide-dev/protocol/openid-connect/token
    client_id: carbide-api
    client_secret: carbide-local-secret

  # Option 3: NGC API key
  # api_key:
  #   authn_url: https://authn.nvidia.com/token
  #   key: nvapi-xxxx
```

Flags and environment variables override config values:

| Flag | Env Var | Description |
|------|---------|-------------|
| `--base-url` | `BMM_BASE_URL` | API base URL |
| `--org` | `BMM_ORG` | Organization name |
| `--token` | `BMM_TOKEN` | Bearer token |
| `--token-url` | `BMM_TOKEN_URL` | OIDC token endpoint URL |
| `--keycloak-url` | `BMM_KEYCLOAK_URL` | Keycloak base URL (constructs token-url) |
| `--keycloak-realm` | `BMM_KEYCLOAK_REALM` | Keycloak realm (default: `carbide-dev`) |
| `--client-id` | `BMM_CLIENT_ID` | OAuth client ID |
| `--output`, `-o` | | Output format: `json` (default), `yaml`, `table` |

## Authentication

```bash
# OIDC (credentials from config, prompts for password if not stored)
carbidecli login

# OIDC with explicit flags
carbidecli --token-url https://auth.example.com/token login --username admin@example.com

# NGC API key
carbidecli login --api-key nvapi-xxxx

# Keycloak shorthand
carbidecli --keycloak-url http://localhost:8080 login --username admin@example.com
```

Tokens are saved to `~/.bmm/config.yaml` with auto-refresh for OIDC.

## Usage

```bash
carbidecli site list
carbidecli site get <siteId>
carbidecli site create --name "SJC4"
carbidecli site create --data-file site.json
cat site.json | carbidecli site create --data-file -
carbidecli site delete <siteId>
carbidecli instance list --status provisioned --page-size 20
carbidecli instance list --all                # fetch all pages
carbidecli allocation constraint create <allocationId> --constraint-type SITE
carbidecli site list --output table
carbidecli --debug site list
```

## Command Structure

Commands follow `carbidecli <resource> [sub-resource] <action> [args] [flags]`.

| Spec Pattern | CLI Action |
|---|---|
| `get-all-*` | `list` |
| `get-*` | `get` |
| `create-*` | `create` |
| `update-*` | `update` |
| `delete-*` | `delete` |
| `batch-create-*` | `batch-create` |
| `get-*-status-history` | `status-history` |
| `get-*-stats` | `stats` |

Nested API paths appear as sub-resource groups:

```
carbidecli allocation list
carbidecli allocation constraint list
carbidecli allocation constraint create <allocationId>
```

## Shell Completion

```bash
# Bash
eval "$(carbidecli completion bash)"

# Zsh
eval "$(carbidecli completion zsh)"

# Fish
carbidecli completion fish > ~/.config/fish/completions/carbidecli.fish
```

## Interactive TUI Mode

Launch an interactive terminal UI with environment selector:

```bash
carbidecli tui
```

The TUI reads all config files from `~/.bmm/` and lets you pick the target environment before running commands. You can also launch it with the `i` alias:

```bash
carbidecli i
```

To start the TUI with a specific config pre-selected:

```bash
carbidecli --config ~/.bmm/config.staging.yaml tui
```

## Multi-Environment Configs

Place multiple configs in `~/.bmm/`:

```
~/.bmm/config.yaml           # default (local dev)
~/.bmm/config.staging.yaml   # staging
~/.bmm/config.prod.yaml      # production
```

Select with `--config`:

```bash
carbidecli --config ~/.bmm/config.staging.yaml site list
```

## Troubleshooting

If `carbidecli` is not found after install, make sure `$(go env GOPATH)/bin` is in your PATH:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Use `--debug` on any command to see the full HTTP request and response for diagnosing issues:

```bash
carbidecli --debug site list
```
