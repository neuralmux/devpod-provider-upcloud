# Contributing to UpCloud DevPod Provider

Thank you for your interest in contributing to the UpCloud DevPod Provider! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Code Style](#code-style)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Debugging](#debugging)
- [Project Structure](#project-structure)

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Git
- Make
- golangci-lint (for linting)
- GoReleaser (for building releases)
- DevPod CLI (for testing)
- GitHub CLI (`gh`) - optional but recommended

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone git@github.com:YOUR_USERNAME/devpod-provider-upcloud.git
   cd devpod-provider-upcloud
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/neuralmux/devpod-provider-upcloud.git
   ```

## Development Setup

### Initial Setup

Run the development setup script:

```bash
./scripts/dev-setup.sh
```

This will:
- Install required Go tools
- Set up pre-commit hooks
- Verify your environment

### Manual Setup

If you prefer manual setup:

```bash
# Install dependencies
go mod download

# Install golangci-lint
brew install golangci-lint

# Install GoReleaser
brew install goreleaser

# Install Godog for BDD tests
go install github.com/cucumber/godog/cmd/godog@latest

# Build the provider
make build
```

### Environment Variables

For testing, set these environment variables:

```bash
# Use test mode (no real API calls)
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"

# Or use real credentials for integration testing
export UPCLOUD_USERNAME="your-username"
export UPCLOUD_PASSWORD="your-password"
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or fixes

### 2. Make Your Changes

Follow the coding standards and ensure your changes are well-tested.

### 3. Run Tests and Linting

Before committing, always run:

```bash
# Run all checks
make pre-push

# Or individually:
make test      # Unit tests
make bdd       # BDD tests
make lint      # Linting
make fmt       # Format code
```

### 4. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add support for custom SSH ports

- Added SSHPort option to ServerConfig
- Updated connection logic to use custom port
- Added tests for port configuration"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Build/tool changes

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Code Style

### Go Code Standards

We follow standard Go conventions:

1. **Formatting**: Use `gofmt` (automatically done by `make fmt`)
2. **Naming**:
   - Exported names start with capital letters
   - Use camelCase for variables and functions
   - Use descriptive names
3. **Error Handling**:
   ```go
   if err != nil {
       return WrapError(err, "descriptive context")
   }
   ```
4. **Comments**:
   - Export functions must have comments
   - Complex logic should be documented

### Linting Rules

Our `.golangci.yml` configuration enforces:
- `govet` - Go vet examiner
- `ineffassign` - Detects ineffectual assignments
- `staticcheck` - Static analysis
- `unused` - Finds unused code

Run linting with:
```bash
make lint
```

### File Organization

- Keep files focused and single-purpose
- Group related functionality
- Use meaningful package names
- Follow the existing project structure

## Testing

### Test Requirements

All new features and bug fixes must include tests.

### Types of Tests

#### 1. Unit Tests

Location: `*_test.go` files alongside code

```go
func TestFunctionName(t *testing.T) {
    // Test implementation
}
```

Run with:
```bash
make test
```

#### 2. BDD Tests

Location: `features/` directory

Write Gherkin scenarios:
```gherkin
Feature: Server Management
  Scenario: Create a server
    Given I have valid credentials
    When I create a server
    Then the server should be running
```

Run with:
```bash
make bdd
```

#### 3. Integration Tests

Use build tag for integration tests:

```go
//go:build integration

func TestIntegration(t *testing.T) {
    // Integration test
}
```

### Test Coverage

View test coverage:
```bash
make coverage
```

We don't enforce a specific coverage percentage, but aim for good coverage of critical paths.

### Test Mode

The provider supports "test mode" for testing without real API calls:

```bash
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
./bin/devpod-provider-upcloud init
```

## Pull Request Process

### Before Submitting

1. **Update Documentation**: If you've added features or changed behavior
2. **Add Tests**: Include tests for new functionality
3. **Run Checks**: Ensure `make pre-push` passes
4. **Update CHANGELOG**: If maintaining one (currently we don't)

### PR Guidelines

1. **Title**: Clear and descriptive
2. **Description**:
   - What changes were made
   - Why were they needed
   - How to test
3. **Size**: Keep PRs focused and reasonable in size
4. **Reviews**: Be responsive to feedback

### PR Template

When you create a PR, use this template:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Refactoring

## Testing
- [ ] Unit tests pass
- [ ] BDD tests pass
- [ ] Tested manually

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

### CI Checks

All PRs must pass:
1. Linting (golangci-lint)
2. Unit tests
3. BDD tests
4. Build verification

## Debugging

### Common Issues

#### 1. Build Failures

```bash
# Clean and rebuild
make clean
make build
```

#### 2. Test Failures

```bash
# Run specific test
go test -v -run TestName ./pkg/...

# Debug BDD tests
go test -v -tags=integration ./features/... -godog.verbose
```

#### 3. Linting Issues

```bash
# Auto-fix some issues
make fmt

# View detailed linting output
golangci-lint run -v
```

### Debugging Tools

#### Local Provider Testing

```bash
# Test locally without installing
./test-local.sh

# Test with real credentials
./test-provider.sh
```

#### DevPod Integration

```bash
# Install local build
make install

# Test with DevPod
devpod provider add . --name upcloud-dev
devpod provider use upcloud-dev
```

#### Logging

Add debug logging:
```go
fmt.Fprintf(os.Stderr, "Debug: %v\n", variable)
```

## Project Structure

```
.
├── cmd/                 # CLI commands
│   ├── root.go         # Root command
│   ├── init.go         # Provider initialization
│   ├── create.go       # Create server
│   ├── delete.go       # Delete server
│   ├── start.go        # Start server
│   ├── stop.go         # Stop server
│   ├── status.go       # Server status
│   └── command.go      # Execute commands
├── pkg/
│   ├── options/        # Environment variable handling
│   ├── upcloud/        # UpCloud API client
│   └── config/         # Configuration management
├── features/           # BDD test scenarios
├── scripts/           # Development scripts
├── docs/              # Documentation
└── .github/           # GitHub Actions workflows
```

### Key Files

- `provider.yaml` - Provider manifest
- `Makefile` - Build and test commands
- `.golangci.yml` - Linting configuration
- `.goreleaser.yml` - Release configuration

## Make Commands

Available make targets:

```bash
make build        # Build binary
make test         # Run unit tests
make bdd          # Run BDD tests
make test-all     # Run all tests
make lint         # Run linter
make fmt          # Format code
make coverage     # Generate coverage report
make clean        # Clean build artifacts
make install      # Install locally
make pre-push     # Run all checks before pushing
```

## Getting Help

If you need help:

1. Check existing issues on GitHub
2. Review the documentation
3. Ask in discussions
4. Create an issue with:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

## Recognition

Contributors are recognized in:
- GitHub contributors page
- Release notes
- README acknowledgments (for significant contributions)

## Release Process

Only maintainers can create releases. See [RELEASES.md](docs/RELEASES.md) for the release process.

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

## Thank You!

Your contributions help make DevPod on UpCloud better for everyone. We appreciate your time and effort!