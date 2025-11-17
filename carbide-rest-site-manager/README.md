# Cloud Site Manager

## Cert-Manager Manifest Generation

The `cert-manager.yaml` files are not stored in this repository. Generate them using:

```bash
make generate-cert-manager
```

This downloads the cert-manager v1.9.1 manifests from the [official GitHub release](https://github.com/cert-manager/cert-manager/releases/tag/v1.9.1) and places them in:
- `charts/cert-manager/templates/cert-manager.yaml`
- `test/manifests/cert-manager/templates/cert-manager.yaml`

### Manual Generation

```bash
curl -fsSL https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml \
  -o charts/cert-manager/templates/cert-manager.yaml

curl -fsSL https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml \
  -o test/manifests/cert-manager/templates/cert-manager.yaml
```

### Updating Version

Edit `scripts/generate-cert-manager-manifests.sh` and update `CERT_MANAGER_VERSION`.

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Code Generation

```bash
make codegen-update
```

