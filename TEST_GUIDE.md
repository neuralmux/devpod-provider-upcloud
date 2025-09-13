# Testing Guide for UpCloud DevPod Provider

## Overview

This guide covers all testing approaches for the UpCloud DevPod provider, from quick local tests to full integration testing. The provider includes comprehensive testing infrastructure with mock mode, BDD tests, and CI/CD integration.

## Quick Start (No UpCloud Account Required)

### 1. Development Setup
```bash
# Run initial setup (installs tools and dependencies)
./scripts/dev-setup.sh

# Build the provider
make build

# Run pre-push validation (comprehensive checks)
make pre-push
```

### 2. Mock Mode Testing
The provider supports a special "mock mode" that simulates API responses without making real calls:

```bash
# Run the automated test script
./test-local.sh

# Or manually with mock credentials
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="2xCPU-4GB"
export UPCLOUD_STORAGE="50"
export UPCLOUD_IMAGE="Ubuntu Server 22.04 LTS (Jammy Jellyfish)"
export MACHINE_ID="test-machine"
export MACHINE_FOLDER="/tmp/test"

# This will succeed in mock mode
./bin/devpod-provider-upcloud init
```

## Testing Framework

### Unit Tests
```bash
# Run unit tests only (excludes integration tests)
make test

# With coverage
make coverage
```

### BDD Tests (Godog)
```bash
# Run BDD tests with mock credentials
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
make bdd

# Run with integration tag
go test -v -tags=integration ./...
```

### Pre-Push Validation
```bash
# Run comprehensive validation (same as CI/CD)
make pre-push

# This runs:
# - Go module validation
# - Code formatting checks
# - Linting with golangci-lint
# - Unit tests
# - Race condition tests
# - Build tests
# - YAML validation
# - Security checks (if tools available)
```

## Testing with Real UpCloud Account

### Prerequisites
1. UpCloud account with API access enabled
2. API credentials from UpCloud Control Panel
3. SSH key pair for server access

### Step 1: Set Real Credentials
```bash
export UPCLOUD_USERNAME="your-real-username"
export UPCLOUD_PASSWORD="your-real-password"
```

### Step 2: Test Authentication
```bash
./bin/devpod-provider-upcloud init
```

If successful, you'll see: "Successfully initialized UpCloud provider"

### Step 3: Full Provider Test Script
```bash
# Run automated test with real API
./test-provider.sh
```

### Step 4: Manual Server Operations (Costs Money!)

⚠️ **WARNING**: These commands will create real servers and incur charges!

```bash
# Setup for server creation
mkdir -p /tmp/devpod-test/.ssh
ssh-keygen -t ed25519 -f /tmp/devpod-test/.ssh/id_devpod -N ''

export MACHINE_FOLDER="/tmp/devpod-test"
export MACHINE_ID="devpod-test-$(date +%s)"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="1xCPU-1GB"  # Smallest plan to minimize cost
export UPCLOUD_STORAGE="10"       # Minimum storage
export AGENT_PATH="/home/devpod/.devpod/devpod"

# Create server (takes ~1-2 minutes)
./bin/devpod-provider-upcloud create

# Check status
./bin/devpod-provider-upcloud status

# Stop server
./bin/devpod-provider-upcloud stop

# Start server
./bin/devpod-provider-upcloud start

# Execute command on server
export COMMAND="echo 'Hello from DevPod'"
./bin/devpod-provider-upcloud command

# Delete server (important to avoid charges!)
./bin/devpod-provider-upcloud delete
```

## Testing with DevPod CLI

### Option 1: Local Provider
```bash
# Add local provider
devpod provider add . --name upcloud-local

# Configure
devpod provider set-options upcloud-local \
  --option UPCLOUD_USERNAME="$UPCLOUD_USERNAME" \
  --option UPCLOUD_PASSWORD="$UPCLOUD_PASSWORD" \
  --option UPCLOUD_ZONE="de-fra1" \
  --option UPCLOUD_PLAN="2xCPU-4GB"

# Use provider
devpod provider use upcloud-local

# Create workspace
devpod up github.com/loft-sh/devpod-quickstart
```

### Option 2: From GitHub
```bash
# Add from GitHub
devpod provider add github.com/neuralmux/devpod-provider-upcloud

# Configure and use as above
```

## CI/CD Integration

The project includes comprehensive CI/CD testing:

### GitHub Actions Workflows
- **ci.yml**: Main CI pipeline (tests, linting, builds)
- **security.yml**: Security scanning (Gosec, Trivy)
- **codeql.yml**: Code analysis
- **docker.yml**: Container builds
- **release.yml**: Automated releases

### Running CI Locally
```bash
# Simulate CI checks
./scripts/pre-push.sh

# This validates:
# - Go module tidiness
# - Code formatting
# - Linting
# - Unit tests
# - Race conditions
# - Multi-platform builds
# - Provider validation
# - Documentation
```

## Troubleshooting

### Common Issues

1. **"missing credentials" error**
   - Ensure UPCLOUD_USERNAME and UPCLOUD_PASSWORD are set
   - Check for typos in environment variables

2. **"authentication failed" error**
   - Verify credentials are correct
   - Ensure API access is enabled in UpCloud Control Panel
   - For testing, use username="test" and password="test"

3. **"invalid zone" error**
   - Use one of the valid zones: de-fra1, fi-hel1, us-nyc1, sg-sin1, etc.

4. **"invalid plan" error**
   - Use exact plan names: 1xCPU-1GB, 2xCPU-4GB, 4xCPU-8GB, etc.

5. **SSH key issues**
   - Ensure MACHINE_FOLDER contains .ssh/id_devpod or .ssh/id_ed25519
   - Keys must have proper permissions (600)

6. **BDD test failures**
   - BDD tests need integration tag: `go test -tags=integration`
   - Use mock credentials for testing

### Debug Mode
```bash
# Enable debug logging
export DEVPOD_DEBUG=true
export DEVPOD_LOG_LEVEL=debug

# Run commands
./bin/devpod-provider-upcloud init
```

### Check Logs
```bash
# DevPod logs location
~/.devpod/logs/

# Provider specific logs
~/.devpod/contexts/default/providers/upcloud/logs/
```

## Cost Management

### Minimize Costs
- Use smallest plan: `1xCPU-1GB`
- Use minimum storage: `10` GB
- Delete servers immediately after testing
- Use auto-stop feature: `INACTIVITY_TIMEOUT=10m`
- Use mock mode for development

### Estimated Costs
- 1xCPU-1GB: ~$0.007/hour (~$5/month)
- Storage: ~$0.00003/GB/hour
- Test session (1 hour): ~$0.01

### Auto-cleanup Script
```bash
# List all DevPod servers
upcloud server list | grep devpod-

# Delete all test servers
upcloud server list | grep devpod-test | awk '{print $1}' | xargs -I {} upcloud server delete {} --delete-storages
```

## Test Coverage

### What's Tested

#### Unit Tests
- Option parsing from environment
- Error handling
- Mock mode behavior

#### BDD Tests (features/provider.feature)
- Provider initialization
- Server creation
- Server lifecycle (start/stop)
- Server deletion
- Command execution
- Credential validation

#### Integration Tests
- Real API authentication
- Server provisioning
- SSH connectivity
- Cloud-init execution

#### CI/CD Tests
- Multi-platform builds
- Docker container builds
- Security scanning
- Code quality checks

## Development Testing Workflow

1. **Make changes to code**
2. **Format code**: `make fmt`
3. **Run unit tests**: `make test`
4. **Run BDD tests**: `make bdd`
5. **Run pre-push validation**: `make pre-push`
6. **Test with mock mode**: `./test-local.sh`
7. **Test with real API** (optional): `./test-provider.sh`
8. **Commit and push**

## Need Help?

1. Check the [CI/CD Guide](CI_CD_GUIDE.md) for pipeline details
2. Review [CLAUDE.md](CLAUDE.md) for implementation details
3. Check provider logs in `~/.devpod/logs/`
4. Verify all environment variables are set correctly
5. Ensure your UpCloud account has API access enabled
6. Check UpCloud service status at https://status.upcloud.com/

Remember to **always delete test servers** to avoid unnecessary charges!