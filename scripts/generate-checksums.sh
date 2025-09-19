#!/bin/bash
# Generate checksums for release binaries
# This script helps update provider.yaml with actual checksums after building

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

cd "$ROOT_DIR"

echo -e "${BLUE}Checksum Generator for UpCloud Provider${NC}"
echo "========================================"
echo ""

# Check if dist directory exists (created by goreleaser)
if [ ! -d "dist" ]; then
    echo -e "${YELLOW}Warning: dist directory not found.${NC}"
    echo "Run 'goreleaser build --snapshot --clean' first to build all platforms"
    echo ""
    echo "Installing goreleaser if not present..."
    if ! command -v goreleaser >/dev/null 2>&1; then
        if command -v brew >/dev/null 2>&1; then
            brew install goreleaser
        else
            echo "Please install goreleaser: https://goreleaser.com/install/"
            exit 1
        fi
    fi

    echo "Building release binaries..."
    goreleaser build --snapshot --clean
fi

# Function to calculate SHA256
calc_sha256() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$1" | awk '{print $1}'
    else
        shasum -a 256 "$1" | awk '{print $1}'
    fi
}

echo "Calculating checksums..."
echo ""

# Define the binaries we need
declare -A binaries=(
    ["linux_amd64"]="dist/upcloud-provider_linux_amd64_v1/devpod-provider-upcloud"
    ["linux_arm64"]="dist/upcloud-provider_linux_arm64/devpod-provider-upcloud"
    ["darwin_amd64"]="dist/upcloud-provider_darwin_amd64_v1/devpod-provider-upcloud"
    ["darwin_arm64"]="dist/upcloud-provider_darwin_arm64/devpod-provider-upcloud"
    ["windows_amd64"]="dist/upcloud-provider_windows_amd64_v1/devpod-provider-upcloud.exe"
)

# Store checksums
declare -A checksums

for platform in "${!binaries[@]}"; do
    binary="${binaries[$platform]}"
    if [ -f "$binary" ]; then
        checksum=$(calc_sha256 "$binary")
        checksums[$platform]=$checksum
        echo -e "${GREEN}✓${NC} $platform: $checksum"
    else
        echo -e "${RED}✗${NC} $platform: Binary not found at $binary"
    fi
done

echo ""
echo "Updating provider.yaml..."

# Backup provider.yaml
cp provider.yaml provider.yaml.bak

# Get version from provider.yaml or use default
VERSION=$(grep "^version:" provider.yaml | awk '{print $2}')
if [ -z "$VERSION" ]; then
    VERSION="0.1.0"
fi

# Update provider.yaml with actual checksums
# This is a simplified version - in production, you'd use proper YAML parsing
update_provider_yaml() {
    local temp_file="provider.yaml.tmp"
    local in_binaries=0
    local platform=""

    while IFS= read -r line; do
        # Check if we're in a binaries section
        if [[ "$line" =~ ^binaries: ]] || [[ "$line" =~ "UC_PROVIDER:" ]]; then
            in_binaries=1
        elif [[ "$line" =~ ^[a-z] ]] && [ "$in_binaries" -eq 1 ]; then
            in_binaries=0
        fi

        # Detect platform
        if [[ "$line" =~ "os: linux" ]]; then
            platform="linux_"
        elif [[ "$line" =~ "os: darwin" ]]; then
            platform="darwin_"
        elif [[ "$line" =~ "os: windows" ]]; then
            platform="windows_"
        fi

        if [[ "$line" =~ "arch: amd64" ]]; then
            platform="${platform}amd64"
        elif [[ "$line" =~ "arch: arm64" ]]; then
            platform="${platform}arm64"
        fi

        # Replace checksum placeholders
        if [[ "$line" =~ "checksum:" ]] && [ ! -z "$platform" ] && [ ! -z "${checksums[$platform]}" ]; then
            echo "      checksum: ${checksums[$platform]}"
            platform=""
        # Replace version placeholders
        elif [[ "$line" =~ "##VERSION##" ]]; then
            echo "${line//##VERSION##/v$VERSION}"
        else
            echo "$line"
        fi
    done < provider.yaml > "$temp_file"

    mv "$temp_file" provider.yaml
}

# Update the file
update_provider_yaml

echo ""
echo -e "${GREEN}✓${NC} provider.yaml updated with checksums"
echo ""
echo "Next steps:"
echo "1. Review the changes: git diff provider.yaml"
echo "2. Commit the updated provider.yaml"
echo "3. Tag the release: git tag v$VERSION"
echo "4. Push tags: git push --tags"
echo "5. GoReleaser will create the GitHub release automatically"
echo ""
echo -e "${BLUE}Note:${NC} For production releases, ensure the binary URLs in provider.yaml"
echo "      match your GitHub releases structure."