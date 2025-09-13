#!/bin/bash

# Test script for UpCloud DevPod Provider
# This script helps test the provider without DevPod

set -e

echo "========================================="
echo "UpCloud DevPod Provider Test Script"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if provider binary exists
if [ ! -f "./bin/devpod-provider-upcloud" ]; then
    echo -e "${RED}Error: Provider binary not found!${NC}"
    echo "Please run 'make build' first"
    exit 1
fi

# Function to print test header
print_test() {
    echo -e "\n${YELLOW}Testing: $1${NC}"
    echo "----------------------------------------"
}

# Function to check result
check_result() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ $1 passed${NC}"
    else
        echo -e "${RED}✗ $1 failed${NC}"
    fi
}

# Set test environment variables
export UPCLOUD_USERNAME="${UPCLOUD_USERNAME:-test-user}"
export UPCLOUD_PASSWORD="${UPCLOUD_PASSWORD:-test-pass}"
export UPCLOUD_ZONE="${UPCLOUD_ZONE:-de-fra1}"
export UPCLOUD_PLAN="${UPCLOUD_PLAN:-2xCPU-4GB}"
export UPCLOUD_STORAGE="${UPCLOUD_STORAGE:-50}"
export UPCLOUD_IMAGE="${UPCLOUD_IMAGE:-Ubuntu Server 22.04 LTS (Jammy Jellyfish)}"

# For testing individual commands
export MACHINE_ID="${MACHINE_ID:-devpod-test-$(date +%s)}"
export MACHINE_FOLDER="${MACHINE_FOLDER:-/tmp/devpod-test}"

echo "Configuration:"
echo "  UPCLOUD_USERNAME: ${UPCLOUD_USERNAME}"
echo "  UPCLOUD_PASSWORD: [hidden]"
echo "  UPCLOUD_ZONE: ${UPCLOUD_ZONE}"
echo "  UPCLOUD_PLAN: ${UPCLOUD_PLAN}"
echo "  UPCLOUD_STORAGE: ${UPCLOUD_STORAGE}"
echo "  UPCLOUD_IMAGE: ${UPCLOUD_IMAGE}"
echo "  MACHINE_ID: ${MACHINE_ID}"
echo ""

# Test 1: Help command
print_test "Help Command"
./bin/devpod-provider-upcloud --help
check_result "Help command"

# Test 2: Init command (tests authentication)
print_test "Init Command (Authentication Test)"
if ./bin/devpod-provider-upcloud init 2>&1 | grep -q "Authentication failed\|authentication test"; then
    echo -e "${YELLOW}Note: Init failed due to invalid credentials (expected for test)${NC}"
else
    check_result "Init command"
fi

# Test 3: Status command (without real server)
print_test "Status Command"
./bin/devpod-provider-upcloud status 2>&1 | grep -q "NotFound" && echo "Status returned NotFound (expected)"
check_result "Status command"

echo ""
echo "========================================="
echo "Basic command tests completed!"
echo ""
echo "To test with real UpCloud credentials:"
echo "1. Set real credentials:"
echo "   export UPCLOUD_USERNAME='your-real-username'"
echo "   export UPCLOUD_PASSWORD='your-real-password'"
echo ""
echo "2. Run init to test authentication:"
echo "   ./bin/devpod-provider-upcloud init"
echo ""
echo "3. To test server creation (WILL INCUR COSTS):"
echo "   # Create SSH key for testing"
echo "   mkdir -p /tmp/devpod-test/.ssh"
echo "   ssh-keygen -t ed25519 -f /tmp/devpod-test/.ssh/id_ed25519 -N ''"
echo "   "
echo "   # Set machine folder"
echo "   export MACHINE_FOLDER=/tmp/devpod-test"
echo "   export MACHINE_ID=devpod-test-\$(date +%s)"
echo "   "
echo "   # Create server"
echo "   ./bin/devpod-provider-upcloud create"
echo "   "
echo "   # Check status"
echo "   ./bin/devpod-provider-upcloud status"
echo "   "
echo "   # Delete server"
echo "   ./bin/devpod-provider-upcloud delete"
echo "========================================="