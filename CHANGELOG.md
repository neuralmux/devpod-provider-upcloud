# Changelog

All notable changes to the UpCloud DevPod Provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-12-18

### 🎉 Major Update: Server Plan Templating System

This release introduces a flexible templating system for server plans, bringing massive cost savings with UpCloud's new Developer and Cloud Native plans.

### Added

#### Server Plan Templating System
- **Flexible Configuration**: New YAML-based plan configuration system (`configs/server-plans.yaml`)
- **Plan Discovery Command**: `devpod-provider-upcloud plans` command to explore available server plans
  - `--recommended` flag to show DevPod-optimized plans
  - `--category` filter for plan categories
  - `--format json/yaml` for automation
- **Smart Plan Loader**: Dynamic plan validation with fallback to legacy mappings
- **Embedded Configuration**: Plans embedded at build time for single-binary distribution

#### New Server Plans Support
- **Developer Plans** (September 2024): 36-89% cost reduction for development workloads
  - Plans from €3-35/month
  - Optimized for DevPod workspaces
  - Includes storage in price
- **Cloud Native Plans** (December 2024): Pay-only-when-powered-on billing
  - Perfect for ephemeral workspaces
  - Ideal with auto-shutdown feature
  - Storage configured separately

#### Documentation
- **`docs/SERVER-PLANS.md`**: Comprehensive user guide for server plans
- **`docs/TEMPLATES.md`**: Technical documentation for the templating system
- **`docs/MIGRATION.md`**: Step-by-step migration guide

#### Developer Experience
- **Quick Start Scripts**: `quickstart.sh`, `install-local.sh`, `install.sh`
- **Credential Auto-Detection**: Automatic detection from UpCloud CLI config
- **Plan Recommendations**: Built-in suggestions by language, framework, and workload

### Changed
- **Default Plan**: Changed from `2xCPU-4GB` (€28/mo) to `DEV-2xCPU-4GB` (€18/mo) - 36% cost reduction
- **Provider Options**: Better organization with option groups
- **Plan Validation**: Enhanced error messages with suggestions

### Migration Notes
Users should review the [Migration Guide](docs/MIGRATION.md) for upgrading existing workspaces.

## [0.1.0] - 2024-09-13

### Added
- Initial UpCloud DevPod provider implementation
- Full UpCloud API integration using official SDK v8
- Server lifecycle management (create, start, stop, delete)
- SSH key injection and secure connectivity
- Comprehensive error handling with user-friendly messages
- Multi-zone support across UpCloud's global infrastructure
- Configurable server plans and storage options
- Cloud-init support for automatic environment setup
- BDD testing framework with Godog
- Comprehensive CI/CD pipeline with GitHub Actions
- Security scanning with CodeQL, Gosec, and Trivy
- Automatic dependency updates with Dependabot
- Cross-platform binary releases (Linux, macOS, Windows)
- Docker container support
- Comprehensive documentation and testing guides

### Infrastructure
- GitHub Actions CI/CD pipeline
- GoReleaser for automated releases
- Security scanning and vulnerability assessment
- Dependabot for dependency management
- Code quality checks with golangci-lint
- Automated testing on multiple platforms

---

## Release Process

Releases are automatically created when a tag is pushed:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This triggers:
1. CI tests and security scans
2. Multi-platform binary builds
3. GitHub release creation
4. Docker image publication
5. Provider manifest updates