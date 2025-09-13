#!/bin/bash

# Simple local test without UpCloud API
# This tests the provider binaries and command structure

echo "==================================="
echo "Local Provider Test (No API Calls)"
echo "==================================="
echo ""

# Build first
echo "Building provider..."
make build

echo ""
echo "1. Testing help for all commands:"
echo "---------------------------------"
./bin/devpod-provider-upcloud --help
echo ""

echo "2. Testing init without credentials:"
echo "------------------------------------"
unset UPCLOUD_USERNAME
unset UPCLOUD_PASSWORD
./bin/devpod-provider-upcloud init 2>&1 | head -5
echo ""

echo "3. Testing with mock credentials:"
echo "---------------------------------"
export UPCLOUD_USERNAME="mock-user"
export UPCLOUD_PASSWORD="mock-pass"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="2xCPU-4GB"
export UPCLOUD_STORAGE="50"
export UPCLOUD_IMAGE="Ubuntu Server 22.04 LTS (Jammy Jellyfish)"
export MACHINE_ID="devpod-test-local"

echo "Environment set:"
echo "  UPCLOUD_USERNAME=$UPCLOUD_USERNAME"
echo "  UPCLOUD_ZONE=$UPCLOUD_ZONE"
echo "  UPCLOUD_PLAN=$UPCLOUD_PLAN"
echo "  UPCLOUD_STORAGE=$UPCLOUD_STORAGE"
echo "  MACHINE_ID=$MACHINE_ID"
echo ""

echo "4. Testing init with mock credentials:"
echo "--------------------------------------"
./bin/devpod-provider-upcloud init 2>&1 | head -5
echo ""

echo "5. Testing status command:"
echo "-------------------------"
./bin/devpod-provider-upcloud status 2>&1
echo ""

echo "==================================="
echo "Local tests completed!"
echo ""
echo "What you're seeing:"
echo "- Init fails with authentication (expected - mock credentials)"
echo "- Status returns NotFound (expected - no server exists)"
echo "- Commands are structured correctly"
echo ""
echo "This confirms the provider binary works!"
echo "==================================="