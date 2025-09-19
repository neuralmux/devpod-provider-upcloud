#!/bin/bash
# Prepare release artifacts for UpCloud DevPod Provider
# This script builds binaries and prepares provider.yaml for release

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Get version from argument or prompt
VERSION="${1:-}"
if [ -z "$VERSION" ]; then
    echo -e "${YELLOW}Usage: $0 <version>${NC}"
    echo -e "Example: $0 v0.1.0"
    exit 1
fi

# Remove 'v' prefix for some operations
VERSION_NO_V="${VERSION#v}"

echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}  Preparing Release: ${VERSION}${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}"
echo ""

# Check for required tools
echo -e "${BLUE}Checking requirements...${NC}"
for tool in go goreleaser git; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        echo -e "${RED}✗ $tool is required but not installed${NC}"
        exit 1
    fi
done
echo -e "${GREEN}✓ All requirements met${NC}"
echo ""

# Clean previous builds
echo -e "${BLUE}Cleaning previous builds...${NC}"
rm -rf dist/
echo -e "${GREEN}✓ Cleaned${NC}"
echo ""

# Run tests
echo -e "${BLUE}Running tests...${NC}"
export UPCLOUD_USERNAME=test
export UPCLOUD_PASSWORD=test
if ! make test; then
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}"
echo ""

# Build with GoReleaser
echo -e "${BLUE}Building binaries with GoReleaser...${NC}"
goreleaser build --snapshot --clean
echo -e "${GREEN}✓ Binaries built${NC}"
echo ""

# Calculate checksums
echo -e "${BLUE}Calculating checksums...${NC}"

# Function to calculate SHA256
calc_sha256() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$1" | awk '{print $1}'
    else
        shasum -a 256 "$1" | awk '{print $1}'
    fi
}

# Expected binary paths from goreleaser
declare -A BINARIES=(
    ["linux_amd64"]="dist/devpod-provider-upcloud_linux_amd64_v1/devpod-provider-upcloud"
    ["linux_arm64"]="dist/devpod-provider-upcloud_linux_arm64/devpod-provider-upcloud"
    ["darwin_amd64"]="dist/devpod-provider-upcloud_darwin_amd64_v1/devpod-provider-upcloud"
    ["darwin_arm64"]="dist/devpod-provider-upcloud_darwin_arm64/devpod-provider-upcloud"
    ["windows_amd64"]="dist/devpod-provider-upcloud_windows_amd64_v1/devpod-provider-upcloud.exe"
)

# Store checksums
declare -A CHECKSUMS

for platform in "${!BINARIES[@]}"; do
    binary="${BINARIES[$platform]}"
    if [ -f "$binary" ]; then
        checksum=$(calc_sha256 "$binary")
        CHECKSUMS[$platform]=$checksum
        echo -e "  ${GREEN}✓${NC} $platform: ${checksum:0:16}..."
    else
        echo -e "  ${RED}✗${NC} $platform: Binary not found at $binary"
    fi
done
echo ""

# Prepare provider.yaml with actual URLs and checksums
echo -e "${BLUE}Updating provider.yaml...${NC}"

# Create a temporary provider.yaml with actual values
cp provider.yaml provider.yaml.tmp

# Update version
sed -i.bak "s/version: .*/version: $VERSION_NO_V/" provider.yaml.tmp

# Update binary URLs and checksums
# Linux AMD64
if [ ! -z "${CHECKSUMS[linux_amd64]}" ]; then
    sed -i.bak "s|##VERSION##|$VERSION|g" provider.yaml.tmp
    sed -i.bak "s|##CHECKSUM_LINUX_AMD64##|${CHECKSUMS[linux_amd64]}|g" provider.yaml.tmp
fi

# Linux ARM64
if [ ! -z "${CHECKSUMS[linux_arm64]}" ]; then
    sed -i.bak "s|##CHECKSUM_LINUX_ARM64##|${CHECKSUMS[linux_arm64]}|g" provider.yaml.tmp
fi

# Darwin AMD64
if [ ! -z "${CHECKSUMS[darwin_amd64]}" ]; then
    sed -i.bak "s|##CHECKSUM_DARWIN_AMD64##|${CHECKSUMS[darwin_amd64]}|g" provider.yaml.tmp
fi

# Darwin ARM64
if [ ! -z "${CHECKSUMS[darwin_arm64]}" ]; then
    sed -i.bak "s|##CHECKSUM_DARWIN_ARM64##|${CHECKSUMS[darwin_arm64]}|g" provider.yaml.tmp
fi

# Windows AMD64
if [ ! -z "${CHECKSUMS[windows_amd64]}" ]; then
    sed -i.bak "s|##CHECKSUM_WINDOWS_AMD64##|${CHECKSUMS[windows_amd64]}|g" provider.yaml.tmp
fi

# Clean up backup files
rm -f provider.yaml.tmp.bak

echo -e "${GREEN}✓ provider.yaml updated${NC}"
echo ""

# Copy binaries to release directory with correct naming
echo -e "${BLUE}Preparing release artifacts...${NC}"
mkdir -p release

# Copy and rename binaries
cp "${BINARIES[linux_amd64]}" "release/devpod-provider-upcloud-linux-amd64" 2>/dev/null || true
cp "${BINARIES[linux_arm64]}" "release/devpod-provider-upcloud-linux-arm64" 2>/dev/null || true
cp "${BINARIES[darwin_amd64]}" "release/devpod-provider-upcloud-darwin-amd64" 2>/dev/null || true
cp "${BINARIES[darwin_arm64]}" "release/devpod-provider-upcloud-darwin-arm64" 2>/dev/null || true
cp "${BINARIES[windows_amd64]}" "release/devpod-provider-upcloud-windows-amd64.exe" 2>/dev/null || true

# Copy provider.yaml
cp provider.yaml.tmp release/provider.yaml

echo -e "${GREEN}✓ Release artifacts prepared in 'release/' directory${NC}"
echo ""

# Create checksums file
echo -e "${BLUE}Creating checksums.txt...${NC}"
cd release
> checksums.txt
for file in devpod-provider-upcloud-*; do
    if [ -f "$file" ]; then
        calc_sha256 "$file" >> checksums.txt
    fi
done
cd ..
echo -e "${GREEN}✓ checksums.txt created${NC}"
echo ""

# Summary
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}  Release $VERSION Prepared Successfully!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}Release artifacts in:${NC} $ROOT_DIR/release/"
echo -e "${BLUE}Updated provider.yaml:${NC} $ROOT_DIR/provider.yaml.tmp"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review the updated provider.yaml.tmp"
echo "2. Test locally: devpod provider add . --name upcloud-test"
echo "3. Create and push tag: git tag $VERSION && git push origin $VERSION"
echo "4. GitHub Actions will create the release automatically"
echo ""
echo -e "${BLUE}Manual release (if needed):${NC}"
echo "1. Create release on GitHub: https://github.com/neuralmux/devpod-provider-upcloud/releases/new"
echo "2. Tag: $VERSION"
echo "3. Upload all files from release/ directory"
echo "4. Upload provider.yaml.tmp as provider.yaml"