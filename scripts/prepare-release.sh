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
# Using regular variables instead of associative arrays for Bash 3 compatibility
# Note: GoReleaser uses the build ID "provider" so paths start with "provider_"
BINARY_LINUX_AMD64="dist/provider_linux_amd64_v1/devpod-provider-upcloud"
BINARY_LINUX_ARM64="dist/provider_linux_arm64_v8.0/devpod-provider-upcloud"
BINARY_DARWIN_AMD64="dist/provider_darwin_amd64_v1/devpod-provider-upcloud"
BINARY_DARWIN_ARM64="dist/provider_darwin_arm64_v8.0/devpod-provider-upcloud"
BINARY_WINDOWS_AMD64="dist/provider_windows_amd64_v1/devpod-provider-upcloud.exe"

# Store checksums in regular variables
CHECKSUM_LINUX_AMD64=""
CHECKSUM_LINUX_ARM64=""
CHECKSUM_DARWIN_AMD64=""
CHECKSUM_DARWIN_ARM64=""
CHECKSUM_WINDOWS_AMD64=""

# Check and calculate checksums for each platform
if [ -f "$BINARY_LINUX_AMD64" ]; then
    CHECKSUM_LINUX_AMD64=$(calc_sha256 "$BINARY_LINUX_AMD64")
    echo -e "  ${GREEN}✓${NC} linux_amd64: ${CHECKSUM_LINUX_AMD64:0:16}..."
else
    echo -e "  ${RED}✗${NC} linux_amd64: Binary not found at $BINARY_LINUX_AMD64"
fi

if [ -f "$BINARY_LINUX_ARM64" ]; then
    CHECKSUM_LINUX_ARM64=$(calc_sha256 "$BINARY_LINUX_ARM64")
    echo -e "  ${GREEN}✓${NC} linux_arm64: ${CHECKSUM_LINUX_ARM64:0:16}..."
else
    echo -e "  ${RED}✗${NC} linux_arm64: Binary not found at $BINARY_LINUX_ARM64"
fi

if [ -f "$BINARY_DARWIN_AMD64" ]; then
    CHECKSUM_DARWIN_AMD64=$(calc_sha256 "$BINARY_DARWIN_AMD64")
    echo -e "  ${GREEN}✓${NC} darwin_amd64: ${CHECKSUM_DARWIN_AMD64:0:16}..."
else
    echo -e "  ${RED}✗${NC} darwin_amd64: Binary not found at $BINARY_DARWIN_AMD64"
fi

if [ -f "$BINARY_DARWIN_ARM64" ]; then
    CHECKSUM_DARWIN_ARM64=$(calc_sha256 "$BINARY_DARWIN_ARM64")
    echo -e "  ${GREEN}✓${NC} darwin_arm64: ${CHECKSUM_DARWIN_ARM64:0:16}..."
else
    echo -e "  ${RED}✗${NC} darwin_arm64: Binary not found at $BINARY_DARWIN_ARM64"
fi

if [ -f "$BINARY_WINDOWS_AMD64" ]; then
    CHECKSUM_WINDOWS_AMD64=$(calc_sha256 "$BINARY_WINDOWS_AMD64")
    echo -e "  ${GREEN}✓${NC} windows_amd64: ${CHECKSUM_WINDOWS_AMD64:0:16}..."
else
    echo -e "  ${RED}✗${NC} windows_amd64: Binary not found at $BINARY_WINDOWS_AMD64"
fi
echo ""

# Prepare provider.yaml with actual URLs and checksums
echo -e "${BLUE}Updating provider.yaml...${NC}"

# Create a temporary provider.yaml with actual values
cp provider.yaml provider.yaml.tmp

# Update version
sed -i.bak "s/version: .*/version: $VERSION_NO_V/" provider.yaml.tmp

# Update binary URLs and checksums
# Update version placeholder
sed -i.bak "s|##VERSION##|$VERSION|g" provider.yaml.tmp

# Linux AMD64
if [ ! -z "$CHECKSUM_LINUX_AMD64" ]; then
    sed -i.bak "s|##CHECKSUM_LINUX_AMD64##|$CHECKSUM_LINUX_AMD64|g" provider.yaml.tmp
fi

# Linux ARM64
if [ ! -z "$CHECKSUM_LINUX_ARM64" ]; then
    sed -i.bak "s|##CHECKSUM_LINUX_ARM64##|$CHECKSUM_LINUX_ARM64|g" provider.yaml.tmp
fi

# Darwin AMD64
if [ ! -z "$CHECKSUM_DARWIN_AMD64" ]; then
    sed -i.bak "s|##CHECKSUM_DARWIN_AMD64##|$CHECKSUM_DARWIN_AMD64|g" provider.yaml.tmp
fi

# Darwin ARM64
if [ ! -z "$CHECKSUM_DARWIN_ARM64" ]; then
    sed -i.bak "s|##CHECKSUM_DARWIN_ARM64##|$CHECKSUM_DARWIN_ARM64|g" provider.yaml.tmp
fi

# Windows AMD64
if [ ! -z "$CHECKSUM_WINDOWS_AMD64" ]; then
    sed -i.bak "s|##CHECKSUM_WINDOWS_AMD64##|$CHECKSUM_WINDOWS_AMD64|g" provider.yaml.tmp
fi

# Clean up backup files
rm -f provider.yaml.tmp.bak

echo -e "${GREEN}✓ provider.yaml updated${NC}"
echo ""

# Copy binaries to release directory with correct naming
echo -e "${BLUE}Preparing release artifacts...${NC}"
mkdir -p release

# Copy and rename binaries
[ -f "$BINARY_LINUX_AMD64" ] && cp "$BINARY_LINUX_AMD64" "release/devpod-provider-upcloud-linux-amd64" 2>/dev/null || true
[ -f "$BINARY_LINUX_ARM64" ] && cp "$BINARY_LINUX_ARM64" "release/devpod-provider-upcloud-linux-arm64" 2>/dev/null || true
[ -f "$BINARY_DARWIN_AMD64" ] && cp "$BINARY_DARWIN_AMD64" "release/devpod-provider-upcloud-darwin-amd64" 2>/dev/null || true
[ -f "$BINARY_DARWIN_ARM64" ] && cp "$BINARY_DARWIN_ARM64" "release/devpod-provider-upcloud-darwin-arm64" 2>/dev/null || true
[ -f "$BINARY_WINDOWS_AMD64" ] && cp "$BINARY_WINDOWS_AMD64" "release/devpod-provider-upcloud-windows-amd64.exe" 2>/dev/null || true

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
        # Calculate checksum and append with filename
        checksum=$(calc_sha256 "$file")
        echo "$checksum  $file" >> checksums.txt
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