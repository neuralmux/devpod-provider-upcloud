# CI/CD Pipeline Guide

This document describes the comprehensive CI/CD infrastructure for the UpCloud DevPod Provider.

## 🚀 Overview

The CI/CD pipeline provides:
- **Automated Testing** - Unit tests, BDD tests, integration tests
- **Security Scanning** - CodeQL, Gosec, Trivy, OWASP dependency check
- **Quality Assurance** - Linting, code formatting, complexity checks
- **Automated Releases** - Multi-platform binaries, Docker images, GitHub releases
- **Dependency Management** - Automated updates with Dependabot

## 📁 Pipeline Structure

```
.github/
├── workflows/                  # GitHub Actions workflows
│   ├── ci.yml                 # Main CI pipeline
│   ├── release.yml            # Release automation
│   ├── security.yml           # Security scans
│   ├── codeql.yml            # Code analysis
│   ├── docker.yml            # Docker builds
│   └── dependabot-auto-merge.yml
├── ISSUE_TEMPLATE/            # Issue templates
│   ├── bug_report.yml
│   └── feature_request.yml
├── PULL_REQUEST_TEMPLATE/     # PR templates
├── codeql/                    # CodeQL configuration
├── dependabot.yml            # Dependency updates
```

## 🔄 Workflows

### 1. CI Pipeline (`ci.yml`)

**Triggers:** Push to `main`/`develop`, Pull Requests

**Jobs:**
- **Lint**: Code quality with golangci-lint
- **Test**: Unit tests with coverage
- **BDD Test**: Godog behavior-driven tests
- **Build**: Multi-platform binary compilation
- **Security**: Gosec and Trivy scans
- **Integration**: End-to-end testing

**Features:**
- Go module caching for faster builds
- Coverage reporting to Codecov
- Artifact uploads
- Test mode for API-free testing

### 2. Release Pipeline (`release.yml`)

**Triggers:** Git tags (`v*`)

**Process:**
1. Run full test suite
2. Update provider.yaml version
3. Build multi-platform binaries with GoReleaser
4. Create GitHub release with changelog
5. Build and push Docker images
6. Update provider.yaml with release URLs/checksums

**Platforms:**
- Linux (amd64, arm64)
- macOS (amd64, arm64)  
- Windows (amd64)

**Artifacts:**
- Compressed binaries
- Docker images
- Homebrew formula
- Updated provider.yaml

### 3. Security Pipeline (`security.yml`)

**Triggers:** Weekly schedule, Push to `main`, PRs

**Scans:**
- **CodeQL**: Static analysis for security vulnerabilities
- **Gosec**: Go-specific security scanner
- **Trivy**: Vulnerability scanner for dependencies
- **Nancy**: OSS Index vulnerability scanner
- **OWASP Dependency Check**: Known vulnerable dependencies

**Integration:**
- Results uploaded to GitHub Security tab
- SARIF format for standardized reporting
- Automatic issue creation for critical findings

### 4. CodeQL Analysis (`codeql.yml`)

**Triggers:** Push, PR, Weekly schedule

**Features:**
- Security-extended queries
- Custom configuration for Go projects
- Path filtering for relevant code
- Integration with GitHub Security

### 5. Docker Pipeline (`docker.yml`)

**Triggers:** Push to branches, Tags, PRs

**Process:**
1. Multi-platform Docker builds (amd64, arm64)
2. Push to GitHub Container Registry
3. Security scanning with Trivy
4. Automated testing of built images

**Images:**
- `ghcr.io/neuralmux/devpod-provider-upcloud:latest`
- `ghcr.io/neuralmux/devpod-provider-upcloud:v1.0.0`

### 6. Dependabot Auto-merge (`dependabot-auto-merge.yml`)

**Features:**
- Auto-approval of patch/minor updates
- Auto-merge for patch updates with passing CI
- Manual review required for major updates
- Failure notifications

## 🔧 Configuration Files

### GoReleaser (`.goreleaser.yml`)
- Multi-platform build configuration
- Archive generation with proper naming
- Changelog generation from commits
- Docker image builds
- Homebrew tap integration

### GolangCI-Lint (`.golangci.yml`)
- 20+ linters enabled
- Custom rules for DevPod providers
- Performance and security checks
- Code complexity analysis

### Dependabot (`.github/dependabot.yml`)
- Weekly Go module updates
- GitHub Actions updates
- Docker base image updates
- Automatic PR creation

### CodeQL Config (`.github/codeql/codeql-config.yml`)
- Security-focused queries
- Path filtering
- Custom rules for Go projects

## 📊 Quality Gates

### Required Checks (PRs)
- ✅ All tests pass (unit + BDD + integration)
- ✅ Code quality (linting, formatting)
- ✅ Security scans pass
- ✅ Build succeeds on all platforms
- ✅ No new vulnerabilities introduced

### Release Criteria
- ✅ All CI checks pass
- ✅ Security scans clear
- ✅ Manual testing completed
- ✅ Documentation updated
- ✅ Version tagged correctly

## 🚀 Usage

### Development Workflow

1. **Create Feature Branch**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make Changes**
   ```bash
   # Write code
   make test      # Run tests locally
   make lint      # Check code quality
   ```

3. **Push and Create PR**
   ```bash
   git push origin feature/new-feature
   # Create PR in GitHub
   ```

4. **CI Runs Automatically**
   - All checks must pass
   - Address any failures
   - Merge when approved

### Release Process

1. **Prepare Release**
   ```bash
   # Update CHANGELOG.md
   # Ensure all tests pass
   make test-all
   ```

2. **Create Tag**
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Automated Release**
   - CI creates GitHub release
   - Binaries built and uploaded
   - Docker images published
   - Provider manifest updated

### Hot Fixes

1. **Create Fix Branch from Main**
   ```bash
   git checkout main
   git checkout -b hotfix/critical-fix
   ```

2. **Apply Fix and Test**
   ```bash
   # Make minimal changes
   make test-all
   ```

3. **Release Patch Version**
   ```bash
   git tag v1.0.1
   git push origin v1.0.1
   ```

## 🛠️ Local Development

### Setup
```bash
# Install dependencies
make deps

# Run all checks locally
make lint
make test-all
make security  # If you have gosec installed
```

### Testing
```bash
# Unit tests
make test

# BDD tests
make bdd

# Integration tests
./test-local.sh
```

### Building
```bash
# Local build
make build

# Cross-platform (requires GoReleaser)
goreleaser build --snapshot --rm-dist
```

## 🔒 Security

### Secrets Management
- `GITHUB_TOKEN`: Automatic (GitHub provided)
- `UPCLOUD_USERNAME`: Optional for integration tests
- `UPCLOUD_PASSWORD`: Optional for integration tests

### Security Scanning Schedule
- **Daily**: Trivy vulnerability scans
- **Weekly**: Full security audit
- **On PR**: Security impact analysis
- **On Release**: Complete security review

### Vulnerability Response
1. **Critical**: Immediate hotfix release
2. **High**: Patch in next minor release
3. **Medium**: Address in next major release
4. **Low**: Technical debt tracking

## 📈 Monitoring

### Metrics Tracked
- Build success rate
- Test coverage percentage
- Security scan results
- Release frequency
- Dependency freshness

### Dashboards
- GitHub Actions status
- Security alerts (GitHub Security tab)
- Release metrics
- Issue/PR velocity

## 🆘 Troubleshooting

### Common Issues

1. **Build Failures**
   - Check Go version compatibility
   - Verify all dependencies are available
   - Review error logs in Actions tab

2. **Test Failures**
   - Ensure test credentials are set
   - Check for race conditions
   - Verify test environment setup

3. **Release Issues**
   - Confirm tag format (v1.0.0)
   - Check GoReleaser configuration
   - Verify all secrets are available

4. **Security Scan Failures**
   - Review vulnerability details
   - Update dependencies
   - Consider false positive exclusions

### Getting Help
- Check GitHub Actions logs
- Review configuration files
- Create issue with "ci/cd" label
- Contact maintainers for access issues

---

This CI/CD pipeline ensures high-quality, secure, and reliable releases of the UpCloud DevPod Provider! 🎉