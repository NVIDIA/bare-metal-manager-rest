<!--
SPDX-FileCopyrightText: Copyright (c) 2026 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->

# BMM CLI

Command-line client for the NVIDIA Bare Metal Manager REST API. Commands are dynamically generated from the embedded OpenAPI spec at startup — zero manual command code.

## Install

```bash
make install-bmmcli
```

Installs to `$(go env GOPATH)/bin/bmmcli`. Override with `make install-bmmcli INSTALL_DIR=/usr/local/bin`.

## Quick Start

Generate a config file:

```bash
bmmcli init                    # writes ~/.bmm/config.yaml
```

Edit `~/.bmm/config.yaml` with your server URL, org, and auth settings, then:

```bash
bmmcli login                   # exchange credentials for a token
bmmcli site list               # list all sites
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
bmmcli login

# OIDC with explicit flags
bmmcli --token-url https://auth.example.com/token login --username admin@example.com

# NGC API key
bmmcli login --api-key nvapi-xxxx

# Keycloak shorthand
bmmcli --keycloak-url http://localhost:8080 login --username admin@example.com
```

Tokens are saved to `~/.bmm/config.yaml` with auto-refresh for OIDC.

## Usage

```bash
bmmcli site list
bmmcli site get <siteId>
bmmcli site create --name "SJC4"
bmmcli site create --data-file site.json
cat site.json | bmmcli site create --data-file -
bmmcli site delete <siteId>
bmmcli instance list --status provisioned --page-size 20
bmmcli instance list --all                # fetch all pages
bmmcli allocation constraint create <allocationId> --constraint-type SITE
bmmcli site list --output table
bmmcli --debug site list
```

## Command Structure

Commands follow `bmmcli <resource> [sub-resource] <action> [args] [flags]`.

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
bmmcli allocation list
bmmcli allocation constraint list
bmmcli allocation constraint create <allocationId>
```

## Shell Completion

```bash
# Bash
eval "$(bmmcli completion bash)"

# Zsh
eval "$(bmmcli completion zsh)"

# Fish
bmmcli completion fish > ~/.config/fish/completions/bmmcli.fish
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
bmmcli --config ~/.bmm/config.staging.yaml site list
```
