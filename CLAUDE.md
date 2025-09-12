# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the UpCloud provider for DevPod - a CLI program that enables DevPod to create, manage, and run workspaces on UpCloud infrastructure.

## Architecture

This provider follows the standard DevPod provider architecture with a Cobra-based CLI structure:

### Directory Structure
```
.
├── cmd/                 # Cobra commands for provider operations
│   ├── root.go         # Root command and CLI setup
│   ├── init.go         # Provider initialization  
│   ├── create.go       # Server creation
│   ├── delete.go       # Server deletion
│   ├── start.go        # Server start
│   ├── stop.go         # Server stop
│   ├── status.go       # Server status check
│   └── command.go      # SSH command execution
├── pkg/
│   ├── options/        # Environment variable parsing
│   └── upcloud/        # UpCloud API client implementation
├── provider.yaml       # DevPod provider manifest
└── main.go            # Entry point

## Development Commands

### Build and Test Commands

```bash
# Build the provider binary
make build

# Run all tests (unit + BDD)
make test-all

# Run only unit tests
make test

# Run only BDD tests with Godog
make bdd

# Format code
make fmt

# Run linter
make vet

# Install dependencies
make deps

# Install provider locally for testing
make install

# Clean build artifacts
make clean

# Generate test coverage
make coverage
```

### Local Development with DevPod

```bash
# Build and test the provider locally
make build
devpod provider add . --name upcloud-dev
devpod provider use upcloud-dev

# Create a workspace
devpod up . --provider upcloud-dev
```

### Running BDD Tests

Set required environment variables before running tests:
```bash
export UPCLOUD_USERNAME="your-api-username"
export UPCLOUD_PASSWORD="your-api-password"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="1xCPU-1GB"
export AGENT_PATH="/opt/devpod/agent"

# Run BDD tests
make bdd
```

## Implementation Status

### Completed
- ✅ Cobra-based CLI structure matching official providers
- ✅ Environment variable parsing for options
- ✅ Provider manifest (provider.yaml) with UpCloud configuration
- ✅ Basic command implementations (init, create, delete, start, stop, status, command)
- ✅ SSH connectivity framework for command execution
- ✅ Godog BDD testing framework setup

### TODO - UpCloud API Integration
The `pkg/upcloud/client.go` currently contains placeholder implementations. To complete the provider:

1. **Add UpCloud Go SDK**: `go get github.com/UpCloudLtd/upcloud-go-api/v6`
2. **Implement actual API calls** in `pkg/upcloud/client.go`:
   - Server creation with SSH keys
   - Server lifecycle management  
   - IP address retrieval
   - Status checking

## Provider Options

Options are passed as environment variables to the provider commands:
- `UPCLOUD_USERNAME`: API username (required)
- `UPCLOUD_PASSWORD`: API password (required)
- `UPCLOUD_ZONE`: Deployment zone (default: de-fra1)
- `UPCLOUD_PLAN`: Server size (default: 2xCPU-4GB)
- `UPCLOUD_STORAGE`: Disk size in GB (default: 50)
- `UPCLOUD_IMAGE`: OS image (default: Ubuntu 22.04)
- `MACHINE_ID`: DevPod machine identifier
- `MACHINE_FOLDER`: Local folder with SSH keys

## Testing the Provider

```bash
# Set test credentials
export UPCLOUD_USERNAME="your-username"
export UPCLOUD_PASSWORD="your-password"

# Test init command
./bin/devpod-provider-upcloud init

# Run BDD tests
make bdd
```