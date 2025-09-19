#!/bin/bash
# Local development installer for UpCloud DevPod Provider
# This script builds and installs the provider locally for testing

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'
BOLD='\033[1m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}${BOLD}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}${BOLD}  UpCloud Provider - Local Install         ${NC}"
echo -e "${BLUE}${BOLD}═══════════════════════════════════════════${NC}"
echo ""

# Check for Go
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}✗${NC} Go is not installed. Please install Go 1.25+"
    exit 1
fi

# Build the provider
echo -e "${GREEN}➜${NC} Building provider..."
make build

if [ ! -f "bin/devpod-provider-upcloud" ]; then
    echo -e "${RED}✗${NC} Build failed!"
    exit 1
fi

echo -e "${GREEN}✓${NC} Build successful"
echo ""

# Install locally
echo -e "${GREEN}➜${NC} Installing provider locally..."

# Check if DevPod is installed
if ! command -v devpod >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠${NC} DevPod is not installed"
    echo ""
    echo "Install DevPod with one of these methods:"
    echo "  brew install loft-sh/tap/devpod"
    echo "  curl -L -o devpod https://github.com/loft-sh/devpod/releases/latest/download/devpod-\$(uname -s | tr '[:upper:]' '[:lower:]')-\$(uname -m)"
    echo ""
    exit 1
fi

# Remove existing local provider if exists
if devpod provider list | grep -q "upcloud-local"; then
    echo -e "${BLUE}ℹ${NC} Removing existing upcloud-local provider..."
    devpod provider delete upcloud-local --force >/dev/null 2>&1 || true
fi

# Add the local provider
echo -e "${GREEN}➜${NC} Adding local provider..."
devpod provider add . --name upcloud-local

echo -e "${GREEN}✓${NC} Provider 'upcloud-local' installed"
echo ""

# Show configuration status
echo -e "${BOLD}Configuration Status:${NC}"
echo ""

# Check for credentials
if [ ! -z "$UPCLOUD_USERNAME" ] && [ ! -z "$UPCLOUD_PASSWORD" ]; then
    echo -e "${GREEN}✓${NC} Credentials found in environment"
elif [ -f "$HOME/.config/upcloud/config.json" ]; then
    echo -e "${GREEN}✓${NC} Credentials found in UpCloud CLI config"
else
    echo -e "${YELLOW}⚠${NC} No credentials found"
    echo ""
    echo "Set credentials with:"
    echo "  export UPCLOUD_USERNAME='your-api-username'"
    echo "  export UPCLOUD_PASSWORD='your-api-password'"
fi

echo ""
echo -e "${BOLD}Testing Commands:${NC}"
echo ""
echo "  ${BLUE}Test with mock mode:${NC}"
echo "    export UPCLOUD_USERNAME=test"
echo "    export UPCLOUD_PASSWORD=test"
echo "    devpod up . --provider upcloud-local --debug"
echo ""
echo "  ${BLUE}Test with real API:${NC}"
echo "    devpod up . --provider upcloud-local --debug"
echo ""
echo "  ${BLUE}View provider options:${NC}"
echo "    devpod provider options upcloud-local"
echo ""
echo "  ${BLUE}Run unit tests:${NC}"
echo "    make test"
echo ""
echo "  ${BLUE}Run all tests:${NC}"
echo "    make test-all"
echo ""
echo "  ${BLUE}Check logs:${NC}"
echo "    devpod provider logs upcloud-local"
echo ""
echo -e "${GREEN}${BOLD}Local provider ready for testing! 🔧${NC}"