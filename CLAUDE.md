# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the UpCloud provider for DevPod - a CLI program that enables DevPod to create, manage, and run workspaces on UpCloud infrastructure. The provider is fully implemented with comprehensive CI/CD, testing, and documentation.

## Architecture

This provider follows the standard DevPod provider architecture with a Cobra-based CLI structure:

### Directory Structure
```
.
├── cmd/                    # Cobra commands for provider operations
│   ├── root.go            # Root command and CLI setup
│   ├── init.go            # Provider initialization  
│   ├── create.go          # Server creation with cloud-init
│   ├── delete.go          # Server deletion
│   ├── start.go           # Server start
│   ├── stop.go            # Server stop
│   ├── status.go          # Server status check
│   └── command.go         # SSH command execution
├── pkg/
│   ├── options/           # Environment variable parsing
│   │   └── options.go     # DevPod option management
│   └── upcloud/           # UpCloud API client implementation
│       ├── client.go      # Main API client with all operations
│       ├── constants.go   # UpCloud-specific constants and mappings
│       ├── mapper.go      # Resource mapping utilities
│       └── errors.go      # User-friendly error handling
├── features/              # BDD test specifications
│   ├── provider.feature   # Gherkin scenarios
│   └── step_definitions/  # Godog step implementations
├── scripts/               # Development and setup scripts
│   ├── dev-setup.sh      # Development environment setup
│   └── pre-push.sh       # Pre-push validation script
├── .github/workflows/     # CI/CD pipeline
│   ├── ci.yml            # Main CI pipeline
│   ├── release.yml       # Automated releases
│   ├── security.yml      # Security scanning
│   ├── codeql.yml        # Code analysis
│   ├── docker.yml        # Docker builds
│   └── dependabot-auto-merge.yml
├── provider.yaml          # DevPod provider manifest
├── Dockerfile            # Container image
├── .goreleaser.yml       # Release configuration
├── .golangci.yml         # Linting configuration
└── main.go               # Entry point
```

## Development Commands

### Build and Test Commands

```bash
# Run development setup (first time)
./scripts/dev-setup.sh

# Build the provider binary
make build

# Run pre-push validation (comprehensive checks)
make pre-push

# Run all tests (unit + BDD)
make test-all

# Run only unit tests (excludes integration tests)
make test

# Run BDD tests with Godog (requires integration tag)
make bdd

# Format code
make fmt

# Run linter (go vet)
make vet

# Run golangci-lint (if installed)
make lint

# Install dependencies
make deps

# Generate test coverage
make coverage

# Clean build artifacts
make clean
```

### Testing Modes

#### Mock Mode (No API Calls)
```bash
# Test with mock credentials
./test-local.sh

# Or manually
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
./bin/devpod-provider-upcloud init
```

#### Real API Testing
```bash
# Set real credentials
export UPCLOUD_USERNAME="your-username"
export UPCLOUD_PASSWORD="your-password"

# Test full provider
./test-provider.sh
```

### Local Development with DevPod

```bash
# Build and install locally
make build
make install

# Add as local provider
devpod provider add . --name upcloud-dev
devpod provider use upcloud-dev

# Create a workspace
devpod up . --provider upcloud-dev
```

## Implementation Status

### ✅ Core Features (Completed)
- Cobra-based CLI structure matching official providers
- Full UpCloud API integration using official SDK v8
- Mock mode for testing without API credentials
- Server lifecycle management (create, start, stop, delete, status)
- SSH key injection and secure connectivity
- Cloud-init user data for DevPod agent setup
- Comprehensive error handling with user-friendly messages
- Multi-zone support across UpCloud's global infrastructure
- Configurable server plans and storage options

### ✅ Testing (Completed)
- Unit tests for packages
- BDD tests with Godog framework
- Integration test support
- Mock mode for API-free testing
- Pre-push validation script
- CI/CD test automation

### ✅ CI/CD Pipeline (Completed)
- GitHub Actions workflows
- Multi-platform builds (Linux, macOS, Windows)
- Automated releases with GoReleaser
- Security scanning (CodeQL, Gosec, Trivy)
- Docker image builds
- Dependency updates with Dependabot
- Code quality checks with golangci-lint

### ✅ Documentation (Completed)
- Comprehensive README
- CI/CD pipeline guide
- Testing documentation
- Development setup scripts
- Provider manifest with all options

## API Integration Details

The provider uses the official UpCloud Go SDK v8 with:

### Authentication
- Username/password authentication
- Mock mode support (username="test", password="test")

### Server Management
- Create servers with customizable plans, zones, and storage
- SSH key injection during provisioning
- Public IPv4 address assignment
- Cloud-init user data support
- Automatic cleanup on creation failure
- Proper state waiting with timeouts
- Server tagging with machine ID

### Error Handling
- User-friendly error messages
- Retry logic for transient failures
- Proper cleanup on errors
- Detailed logging for debugging

## Provider Options

Options are passed as environment variables to the provider commands:

### Required Options
- `UPCLOUD_USERNAME`: API username
- `UPCLOUD_PASSWORD`: API password

### Optional Options
- `UPCLOUD_ZONE`: Deployment zone (default: de-fra1)
- `UPCLOUD_PLAN`: Server size (default: 2xCPU-4GB)
- `UPCLOUD_STORAGE`: Disk size in GB (default: 50)
- `UPCLOUD_IMAGE`: OS image (default: Ubuntu Server 22.04 LTS)

### DevPod Options
- `MACHINE_ID`: DevPod machine identifier
- `MACHINE_FOLDER`: Local folder with SSH keys
- `AGENT_PATH`: DevPod agent installation path
- `AGENT_DATA_PATH`: Agent data directory
- `INACTIVITY_TIMEOUT`: Auto-stop timeout
- `INJECT_GIT_CREDENTIALS`: Git credential injection
- `INJECT_DOCKER_CREDENTIALS`: Docker credential injection

## Important Implementation Notes

### Mock Mode
When `UPCLOUD_USERNAME="test"` and `UPCLOUD_PASSWORD="test"`, the provider operates in mock mode:
- No real API calls are made
- Operations simulate success
- Used for testing and development

### SSH Key Handling
- SSH keys are read from `{MACHINE_FOLDER}/id_devpod[.pub]`
- Public key is injected during server creation
- Private key is used for SSH connectivity

### Cloud-init Script
The provider automatically generates cloud-init user data that:
- Installs DevPod agent
- Configures SSH access
- Sets up the development environment

### Server Identification
Servers are tagged with `devpod-machine={MACHINE_ID}` for identification across operations.

## Common Issues and Solutions

### Testing Failures
- BDD tests require integration build tag: `go test -tags=integration`
- Mock credentials can be used for most testing scenarios
- Real API testing requires valid UpCloud credentials

### Build Issues
- Requires Go 1.25 or higher
- Run `make deps` to ensure all dependencies are installed
- Use `make pre-push` to validate before committing

### CI/CD Pipeline
- All PRs must pass CI checks
- Security scans run automatically
- Releases are triggered by version tags (v*)

## Code Style Guidelines

- Follow Go idioms and best practices
- Use meaningful variable and function names
- Add error handling for all operations
- Include user-friendly error messages
- Write tests for new functionality
- Update documentation when adding features
- Run `make pre-push` before committing

## Future Enhancements (Optional)

These are potential improvements but not required:
- Firewall configuration options
- Backup management
- Network isolation options
- Custom metadata support
- IPv6 support
- Floating IP management
- Load balancer integration