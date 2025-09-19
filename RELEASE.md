# Release Process Documentation

This document describes the release process for the UpCloud DevPod Provider, including local testing, GitHub releases, and troubleshooting.

## Table of Contents

- [Overview](#overview)
- [Release Pipeline Architecture](#release-pipeline-architecture)
- [Prerequisites](#prerequisites)
- [Local Release Testing](#local-release-testing)
- [Creating a GitHub Release](#creating-a-github-release)
- [Release Workflow](#release-workflow)
- [Versioning Strategy](#versioning-strategy)
- [Troubleshooting](#troubleshooting)
- [Post-Release Checklist](#post-release-checklist)

## Overview

The UpCloud DevPod Provider uses a fully automated release pipeline that:

1. Builds cross-platform binaries (Linux, macOS, Windows)
2. Calculates SHA256 checksums for verification
3. Updates provider.yaml with correct URLs and checksums
4. Creates GitHub releases with all artifacts
5. Publishes Docker images
6. Makes the provider installable via DevPod

## Release Pipeline Architecture

```
┌──────────────────┐
│   Git Tag Push   │
│    (v0.2.0)      │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ GitHub Actions   │
│  release.yml     │
└────────┬─────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌──────────┐ ┌──────────┐
│  Tests   │ │  Build   │
│ (make    │ │ (server  │
│  test)   │ │  plans)  │
└────┬─────┘ └────┬─────┘
     │            │
     └──────┬─────┘
            │
            ▼
┌──────────────────┐
│   GoReleaser     │
│ - Build binaries │
│ - Create archives│
│ - Generate notes │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ Update provider  │
│ - Calculate SHA  │
│ - Update URLs    │
└────────┬─────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌──────────┐ ┌──────────┐
│  GitHub  │ │  Docker  │
│ Release  │ │  Images  │
│ + Assets │ │  (ghcr)  │
└──────────┘ └──────────┘
```

## Prerequisites

### For Local Testing

- Go 1.21 or higher
- GoReleaser (`brew install goreleaser` or download from [goreleaser.com](https://goreleaser.com))
- Make
- Git with configured signing (optional but recommended)

### For GitHub Releases

- Write access to the repository
- Git tags permission
- GitHub Actions enabled

## Local Release Testing

### 1. Test the Release Build Locally

Use the prepare-release.sh script to test the release process:

```bash
# Test with a version
./scripts/prepare-release.sh v0.2.0

# This will:
# - Run tests
# - Build binaries with GoReleaser
# - Calculate checksums
# - Update provider.yaml
# - Create release/ directory with artifacts
```

### 2. Verify the Build

After running the script, check:

```bash
# List release artifacts
ls -la release/

# Expected output:
# devpod-provider-upcloud-linux-amd64
# devpod-provider-upcloud-linux-arm64
# devpod-provider-upcloud-darwin-amd64
# devpod-provider-upcloud-darwin-arm64
# devpod-provider-upcloud-windows-amd64.exe
# provider.yaml
# checksums.txt

# Verify checksums
cd release
shasum -a 256 -c checksums.txt

# Test the binary
./devpod-provider-upcloud-darwin-arm64 --version
```

### 3. Test Provider Installation Locally

```bash
# Install from local directory
devpod provider add . --name upcloud-test

# Verify installation
devpod provider list | grep upcloud-test

# Test initialization
devpod provider options upcloud-test
```

## Creating a GitHub Release

### Automatic Release (Recommended)

1. **Update Version in Files** (if needed):
   ```bash
   # Update CHANGELOG.md with release notes
   vim CHANGELOG.md

   # Commit changes
   git add CHANGELOG.md
   git commit -m "chore: prepare release v0.2.0"
   git push origin main
   ```

2. **Create and Push Tag**:
   ```bash
   # Create annotated tag
   git tag -a v0.2.0 -m "Release v0.2.0 - Server Plan Templates"

   # Push tag to trigger release
   git push origin v0.2.0
   ```

3. **Monitor Release**:
   ```bash
   # Watch GitHub Actions
   gh run watch

   # Or check in browser
   open https://github.com/neuralmux/devpod-provider-upcloud/actions
   ```

### Manual Release (If Needed)

If automatic release fails, you can create one manually:

```bash
# 1. Build release locally
./scripts/prepare-release.sh v0.2.0

# 2. Create GitHub release
gh release create v0.2.0 \
  --title "v0.2.0 - Server Plan Templates" \
  --notes-file CHANGELOG.md \
  --target main

# 3. Upload artifacts
cd release
gh release upload v0.2.0 \
  devpod-provider-upcloud-* \
  provider.yaml \
  checksums.txt
```

## Release Workflow

The automated workflow (`release.yml`) performs these steps:

1. **Trigger**: On push of tags matching `v*`

2. **Pre-Release Checks**:
   - Set up Go environment
   - Copy server-plans.yaml for embedding
   - Run all tests
   - Verify build

3. **Build Process**:
   - GoReleaser builds for all platforms
   - Creates archives and checksums
   - Generates release notes

4. **Asset Preparation**:
   - Calculate SHA256 for each binary
   - Update provider.yaml with:
     - Correct version
     - Binary URLs
     - Checksums

5. **Publication**:
   - Create GitHub release
   - Upload binaries
   - Upload provider.yaml
   - Build and push Docker images

## Versioning Strategy

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR.MINOR.PATCH** (e.g., 0.2.0)
- **MAJOR**: Breaking changes
- **MINOR**: New features, backward compatible
- **PATCH**: Bug fixes

### Version Locations

Version must be updated in:
1. `provider.yaml` - `version:` field
2. `CHANGELOG.md` - New version section
3. Git tag - `v` prefix (e.g., v0.2.0)

## Troubleshooting

### Build Failures

#### "pattern server-plans.yaml: no matching files found"

```bash
# Solution: Copy config file before building
cp configs/server-plans.yaml pkg/config/
make build
```

#### "no such file or directory: devpod-provider-upcloud"

```bash
# Solution: Build was not successful
make clean
make build
```

### Release Failures

#### "release already exists"

```bash
# Delete the release and tag
gh release delete v0.2.0 --yes
git push --delete origin v0.2.0
git tag -d v0.2.0

# Recreate and push
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0
```

#### "404 when downloading provider"

This means no release exists yet. Create the first release:

```bash
# Create initial release
git tag -a v0.1.0 -m "Initial release"
git push origin v0.1.0
```

#### "checksum mismatch"

```bash
# Rebuild and recalculate
./scripts/prepare-release.sh v0.2.0

# Verify checksums match
cat release/checksums.txt
```

### Testing Issues

#### "UPCLOUD_USERNAME is required"

```bash
# Use test credentials for mock mode
export UPCLOUD_USERNAME=test
export UPCLOUD_PASSWORD=test
./bin/devpod-provider-upcloud init
```

## Post-Release Checklist

After releasing, verify:

- [ ] **GitHub Release Created**
  ```bash
  gh release view v0.2.0
  ```

- [ ] **All Binaries Uploaded**
  ```bash
  gh release view v0.2.0 --json assets -q '.assets[].name'
  ```

- [ ] **provider.yaml Has Correct Checksums**
  ```bash
  curl -L https://github.com/neuralmux/devpod-provider-upcloud/releases/download/v0.2.0/provider.yaml
  ```

- [ ] **Installation Works**
  ```bash
  devpod provider delete upcloud-test --force 2>/dev/null || true
  devpod provider add github.com/neuralmux/devpod-provider-upcloud@v0.2.0
  ```

- [ ] **Docker Images Published**
  ```bash
  docker pull ghcr.io/neuralmux/devpod-provider-upcloud:v0.2.0
  ```

- [ ] **Quickstart Script Works**
  ```bash
  curl -fsSL https://raw.githubusercontent.com/neuralmux/devpod-provider-upcloud/main/quickstart.sh | bash
  ```

## Release Announcement Template

After successful release, announce it:

```markdown
### 🎉 UpCloud DevPod Provider v0.2.0 Released!

**Highlights:**
- 💰 36-89% cost savings with new Developer Plans
- 🎯 Smart plan selection with `devpod-provider-upcloud plans`
- 🚀 Flexible server plan templates

**Quick Start:**
curl -fsSL https://raw.githubusercontent.com/neuralmux/devpod-provider-upcloud/main/quickstart.sh | bash

**Upgrade:**
devpod provider delete upcloud --force
devpod provider add github.com/neuralmux/devpod-provider-upcloud@v0.2.0

[Full Changelog](https://github.com/neuralmux/devpod-provider-upcloud/blob/main/CHANGELOG.md)
```

## Emergency Procedures

### Rollback a Release

If a release has critical issues:

```bash
# 1. Mark as pre-release
gh release edit v0.2.0 --prerelease

# 2. Or delete entirely
gh release delete v0.2.0 --yes

# 3. Fix issues and re-release
git tag -d v0.2.0
git tag -a v0.2.0 -m "Release v0.2.0 (fixed)"
git push origin v0.2.0 --force
```

### Hotfix Process

For urgent fixes:

```bash
# 1. Create hotfix branch
git checkout -b hotfix/v0.2.1 v0.2.0

# 2. Make fixes
# ... edit files ...

# 3. Update version
vim provider.yaml  # Update to 0.2.1
vim CHANGELOG.md   # Add hotfix notes

# 4. Commit and tag
git commit -am "hotfix: critical bug fix"
git tag -a v0.2.1 -m "Hotfix v0.2.1"

# 5. Push to trigger release
git push origin hotfix/v0.2.1
git push origin v0.2.1
```

## Development Releases

For testing pre-releases:

```bash
# Create pre-release tag
git tag -a v0.3.0-beta.1 -m "Beta release v0.3.0-beta.1"
git push origin v0.3.0-beta.1

# Mark as pre-release on GitHub
gh release edit v0.3.0-beta.1 --prerelease
```

---

*Last updated: December 2024*