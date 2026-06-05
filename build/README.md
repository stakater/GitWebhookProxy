# GitWebhookProxy Build System

This document describes the build system for GitWebhookProxy, which now supports multi-architecture builds and alternative container engines.

## Features

- **Multi-architecture support**: Build for both `linux/amd64` and `linux/arm64` platforms
- **Container engine flexibility**: Works with both Docker and Podman
- **Modern Go tooling**: Uses Go modules and Go's native embed package for embedding static assets
- **Efficient builds**: Uses multi-stage builds to minimize image size

## Build Options

### Basic Build

To build the binary locally:

```bash
make build
```

### Container Build

#### Multi-architecture Container (Recommended)

Build a multi-architecture container image:

```bash
make container
```

This will:
1. Create a manifest for the container image
2. Build the image for each platform specified in `PLATFORMS` (default: linux/amd64 and linux/arm64)
3. Push the manifest to the registry if `REGISTRY` is defined

You can customize the build with these variables:

- `PLATFORMS`: Space-separated list of platforms to build for (default: `linux/amd64 linux/arm64`)
- `DOCKER_IMAGE`: Image name (default: `stakater/gitwebhookproxy`)
- `DOCKER_TAG`: Image tag (default: `dev`)
- `REGISTRY`: Registry to push to (default: `docker.io`)
- `CONTAINER_CLI`: Container CLI to use (automatically detected, defaults to `docker` if `podman` is not available)

Example:

```bash
# Build for specific platforms
PLATFORMS="linux/amd64 linux/arm64 linux/ppc64le" make container

# Use a specific container engine
CONTAINER_CLI=podman make container

# Push to a specific registry
REGISTRY=docker.io/myorg make container
```

#### Legacy Single-architecture Container

For backward compatibility, you can still build a single-architecture container:

```bash
make binary-image
```

### Deployment

To build, push, and deploy:

```bash
# Multi-architecture deployment
make deploy

# Legacy single-architecture deployment
make deploy-legacy
```

## Dockerfile

The new build system uses `build/package/Dockerfile.multi`, which:

1. Uses a multi-stage build process
2. Supports multiple architectures via the `TARGETARCH` build argument
3. Uses the latest Alpine base image for security
4. Uses Go's native embed package instead of third-party libraries