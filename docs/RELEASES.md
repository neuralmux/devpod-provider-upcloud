# Release Process Documentation

## Overview

The UpCloud DevPod Provider follows semantic versioning and uses automated GitHub Actions workflows for releases. This document describes the complete release process.

## Release Strategy

### Versioning

We follow [Semantic Versioning](https://semver.org/) (SemVer):

- **MAJOR** (X.0.0): Breaking API changes
- **MINOR** (0.X.0): New features, backward compatible
- **PATCH** (0.0.X): Bug fixes, backward compatible

Current version: `v0.1.0`

### Release Cycle

- **Patch releases**: As needed for bug fixes
- **Minor releases**: When new features are ready
- **Major releases**: Only for breaking changes

## Creating a Release

### Prerequisites

1. Ensure all changes are merged to `main` branch
2. All CI checks are passing
3. Documentation is updated
4. Local environment is clean:
   ```bash
   git status  # Should show clean working directory
   git pull origin main  # Ensure up to date
   ```

### Step-by-Step Release Process

#### 1. Prepare for Release

```bash
# Ensure you're on main branch
git checkout main
git pull origin main

# Run full test suite locally
make test-all

# Run linting
make lint

# Build binaries to verify
make build
```

#### 2. Create and Push Tag

```bash
# Create annotated tag
git tag -a v0.1.1 -m "Release v0.1.1: Bug fixes and improvements"

# Push tag to trigger release workflow
git push origin v0.1.1
```

#### 3. Monitor Release Workflow

```bash
# Watch the release workflow
gh run watch

# Or view in browser
gh run list --workflow=release.yml
```

#### 4. Verify Release

Once the workflow completes:

1. Check GitHub Releases page
2. Verify all artifacts are uploaded:
   - Binary for each platform
   - provider.yaml with correct checksums
3. Test installation:
   ```bash
   devpod provider delete upcloud
   devpod provider add github.com/neuralmux/devpod-provider-upcloud
   ```

## Automated Release Process

### What Happens When You Push a Tag

1. **Trigger**: Tag matching `v*` pattern is pushed
2. **Test Suite**: Full test suite runs (unit + BDD)
3. **Build Binaries**: GoReleaser builds for all platforms
4. **Calculate Checksums**: SHA256 for each binary
5. **Update provider.yaml**: Insert version and checksums
6. **Create GitHub Release**: Upload all artifacts
7. **Notification**: Success/failure message

### Release Artifacts

Each release includes:

```
devpod-provider-upcloud-linux-amd64      # Linux AMD64 binary
devpod-provider-upcloud-linux-arm64      # Linux ARM64 binary
devpod-provider-upcloud-darwin-amd64     # macOS Intel binary
devpod-provider-upcloud-darwin-arm64     # macOS Apple Silicon binary
devpod-provider-upcloud-windows-amd64.exe # Windows binary
provider.yaml                              # Provider manifest with checksums
```

### provider.yaml Updates

The release workflow automatically:
1. Updates version number
2. Inserts download URLs with tag
3. Calculates and inserts SHA256 checksums

Example:
```yaml
version: 0.1.1
binaries:
  UC_PROVIDER:
    - os: linux
      arch: amd64
      path: https://github.com/neuralmux/devpod-provider-upcloud/releases/download/v0.1.1/devpod-provider-upcloud-linux-amd64
      checksum: <calculated-sha256>
```

## Re-releasing

### When a Release Fails

If the release workflow fails after tagging:

#### Option 1: Re-tag Current Commit (Recommended for failures)

```bash
# Delete local and remote tag
git tag -d v0.1.1
git push origin :refs/tags/v0.1.1

# Fix issues
git add .
git commit -m "Fix release issues"
git push origin main

# Re-tag and push
git tag -a v0.1.1 -m "Release v0.1.1"
git push origin v0.1.1
```

#### Option 2: Create New Patch Version

```bash
# If release was partially successful, increment version
git tag -a v0.1.2 -m "Release v0.1.2: Fix release issues"
git push origin v0.1.2
```

### Deleting a Release

If you need to remove a release completely:

```bash
# Delete GitHub release
gh release delete v0.1.1 --yes

# Delete tag
git push origin :refs/tags/v0.1.1
git tag -d v0.1.1
```

## Manual Release Process

If automation fails, you can release manually:

### 1. Build Binaries Locally

```bash
# Install GoReleaser
brew install goreleaser/tap/goreleaser

# Build binaries
goreleaser build --clean --snapshot
```

### 2. Calculate Checksums

```bash
cd dist
for file in */devpod-provider-upcloud*; do
  sha256sum "$file" >> checksums.txt
done
```

### 3. Update provider.yaml

```bash
# Replace placeholders in provider.yaml
VERSION="0.1.1"
sed -i "s/##VERSION##/v$VERSION/g" provider.yaml
# Add checksums manually from checksums.txt
```

### 4. Create GitHub Release

```bash
# Create release with GitHub CLI
gh release create v0.1.1 \
  --title "v0.1.1" \
  --notes "Release notes here" \
  dist/*/devpod-provider-upcloud* \
  provider.yaml
```

## Release Checklist

Before releasing, ensure:

- [ ] All tests pass locally (`make test-all`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated
- [ ] CHANGELOG is updated (if maintaining one)
- [ ] Version numbers are consistent
- [ ] Previous release worked correctly

## Post-Release

After a successful release:

### 1. Verify Installation

```bash
# Test with DevPod
devpod provider add github.com/neuralmux/devpod-provider-upcloud@v0.1.1
devpod provider use upcloud

# Create test workspace
devpod up --provider upcloud github.com/some/repo
```

### 2. Update Documentation

- Update README if needed
- Update installation instructions
- Add release notes to documentation

### 3. Announce Release

- Create announcement if significant release
- Update any external documentation
- Notify users through appropriate channels

## Troubleshooting Releases

### Common Issues

#### 1. Tag Already Exists

**Error**: `tag already exists`

**Solution**:
```bash
git tag -d v0.1.1
git push origin :refs/tags/v0.1.1
```

#### 2. GoReleaser Fails

**Error**: `only version: 2 configuration files are supported`

**Solution**: Ensure `.goreleaser.yml` has `version: 2`

#### 3. Checksum Mismatch

**Error**: Users report checksum verification failures

**Solution**:
1. Verify checksums were calculated correctly
2. Check if binaries were modified after checksum calculation
3. Re-release with correct checksums

#### 4. Missing Artifacts

**Error**: Some binaries missing from release

**Solution**:
1. Check GoReleaser logs for build failures
2. Verify all platforms in `.goreleaser.yml`
3. Re-run release workflow or upload manually

### Debugging Release Workflow

```bash
# View detailed logs
gh run view <run-id> --log

# Download workflow artifacts for inspection
gh run download <run-id>

# Re-run failed workflow
gh run rerun <run-id>
```

## Version Management

### Where Versions Are Defined

1. **Git tags**: Source of truth for releases
2. **provider.yaml**: Updated automatically during release
3. **Binary metadata**: Embedded during build via ldflags

### Bumping Versions

We don't maintain version in files. Version is determined by git tag:

```bash
# For next release, just create appropriate tag
git tag -a v0.2.0 -m "Minor release: New features"
git push origin v0.2.0
```

## Release Notes

### Format

Release notes should include:

```markdown
## What's Changed
- Feature: Description
- Fix: Description
- Docs: Description

## Breaking Changes
- None (or list them)

## Contributors
- @username

**Full Changelog**: https://github.com/neuralmux/devpod-provider-upcloud/compare/v0.1.0...v0.1.1
```

### Generating Release Notes

GitHub can auto-generate release notes:

```bash
gh release create v0.1.1 --generate-notes
```

## Security Considerations

### Signing Releases

Currently not implemented. Future enhancement:
- GPG sign tags
- Sign binaries
- Provide signature verification

### Checksum Verification

All binaries include SHA256 checksums in provider.yaml for verification.

## Rollback Procedure

If a release has critical issues:

1. **Mark as Pre-release** on GitHub
2. **Document known issues** in release notes
3. **Direct users to previous version**:
   ```bash
   devpod provider add github.com/neuralmux/devpod-provider-upcloud@v0.1.0
   ```
4. **Fix issues and release patch version**

## Future Improvements

Planned enhancements:

- [ ] Automated changelog generation
- [ ] Semantic release automation
- [ ] Release candidate (RC) builds
- [ ] Binary signing
- [ ] Homebrew formula updates
- [ ] Container image releases
- [ ] Release metrics dashboard