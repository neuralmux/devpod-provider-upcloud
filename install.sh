#!/bin/bash
# One-line installer for UpCloud DevPod Provider
# Usage: curl -fsSL https://raw.githubusercontent.com/neuralmux/devpod-provider-upcloud/main/install.sh | bash

set -e

echo "🚀 Installing UpCloud DevPod Provider..."

# Check for DevPod
if ! command -v devpod >/dev/null 2>&1; then
    echo "📦 Installing DevPod first..."
    if [[ "$OSTYPE" == "darwin"* ]] && command -v brew >/dev/null 2>&1; then
        brew install loft-sh/tap/devpod
    else
        curl -L -o /tmp/devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
        chmod +x /tmp/devpod
        sudo mv /tmp/devpod /usr/local/bin/
    fi
fi

# Add UpCloud provider
echo "➕ Adding UpCloud provider..."
devpod provider add github.com/neuralmux/devpod-provider-upcloud --name upcloud

# Check for credentials
if [ -z "$UPCLOUD_USERNAME" ] || [ -z "$UPCLOUD_PASSWORD" ]; then
    if [ ! -f "$HOME/.config/upcloud/config.json" ]; then
        echo ""
        echo "⚠️  Please set your UpCloud API credentials:"
        echo "   export UPCLOUD_USERNAME='your-username'"
        echo "   export UPCLOUD_PASSWORD='your-password'"
        echo ""
        echo "   Get API credentials at: https://hub.upcloud.com/account/api"
    fi
fi

echo "✅ UpCloud provider installed successfully!"
echo ""
echo "📝 Quick Start:"
echo "   devpod up github.com/your/repo --provider upcloud"
echo "   devpod up . --provider upcloud"
echo ""