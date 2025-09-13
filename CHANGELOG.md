# Changelog

All notable changes to the UpCloud DevPod Provider will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

## [0.1.0] - 2024-XX-XX

### Added
- Initial release of UpCloud DevPod Provider
- Support for all UpCloud zones and plans
- Automatic SSH key management
- DevPod workspace provisioning
- Complete provider lifecycle implementation

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