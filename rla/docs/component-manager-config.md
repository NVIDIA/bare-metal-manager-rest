# Component Manager Configuration

This document explains the configuration files for the Component Manager system.

## Overview

The Component Manager configuration controls:
1. Which implementation to use for each component type (compute, NVL switch, power shelf)
2. Which API providers to enable and their settings

## Configuration Files

| File | Purpose |
|------|---------|
| `componentmanager.prod.yaml` | Production configuration using real API implementations |
| `componentmanager.test.yaml` | Testing/development configuration using mock implementations |

## Configuration Structure

### Component Managers

```yaml
component_managers:
  compute: <implementation>
  nvlswitch: <implementation>
  powershelf: <implementation>
```

Maps each component type to its implementation. Available implementations:

| Component Type | Available Implementations | Description |
|----------------|---------------------------|-------------|
| `compute` | `carbide`, `mock` | Manages compute nodes |
| `nvlswitch` | `carbide`, `mock` | Manages NVLink switches |
| `powershelf` | `psm`, `mock` | Manages power shelves |

### Providers

```yaml
providers:
  carbide:
    timeout: "<duration>"
  psm:
    timeout: "<duration>"
```

Configures API client providers. **A provider is enabled if its section is present** in the configuration.

| Provider | Used By | Description |
|----------|---------|-------------|
| `carbide` | compute, nvlswitch | Carbide API for machine management |
| `psm` | powershelf | Power Shelf Manager API |

#### Provider Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `timeout` | duration string | `1m` (carbide), `30s` (psm) | gRPC call timeout |

Duration strings use Go format: `30s`, `1m`, `2m30s`, etc.

## Examples

### Production Configuration

```yaml
# Uses real implementations connecting to external services
component_managers:
  compute: carbide
  nvlswitch: carbide
  powershelf: psm

providers:
  carbide:
    timeout: "1m"
  psm:
    timeout: "30s"
```

### Test Configuration

```yaml
# Uses mock implementations - no external dependencies
component_managers:
  compute: mock
  nvlswitch: mock
  powershelf: mock

# No providers section needed for mock implementations
```

### Mixed Configuration (e.g., partial testing)

```yaml
# Real power shelf management, mock compute/nvlswitch
component_managers:
  compute: mock
  nvlswitch: mock
  powershelf: psm

providers:
  psm:
    timeout: "30s"
```

## Provider Auto-Detection

If the `providers` section is omitted entirely, providers are automatically enabled based on the component manager implementations:

- If any component uses `carbide` → Carbide provider is enabled with defaults
- If any component uses `psm` → PSM provider is enabled with defaults

This allows minimal configuration:

```yaml
component_managers:
  compute: carbide
  nvlswitch: carbide
  powershelf: psm
# Providers auto-enabled based on implementations above
```

## Usage

Set the configuration file path via:

1. **Command line flag**: `--component-config <path>`
2. **Environment variable**: `COMPONENT_MANAGER_CONFIG=<path>`
3. **Default**: `configs/componentmanager.prod.yaml`
