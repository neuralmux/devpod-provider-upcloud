#!/bin/bash
# UpCloud DevPod Provider - Quick Start Script
# This script provides the easiest way to get started with the UpCloud DevPod provider

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Provider details
PROVIDER_REPO="github.com/neuralmux/devpod-provider-upcloud"
PROVIDER_NAME="upcloud"

# Functions
print_header() {
    echo -e "${BLUE}${BOLD}╔════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}${BOLD}║  UpCloud DevPod Provider - Quick Start     ║${NC}"
    echo -e "${BLUE}${BOLD}╚════════════════════════════════════════════╝${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}➜${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

check_command() {
    if command -v "$1" >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

install_devpod() {
    print_step "Installing DevPod..."

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if check_command brew; then
            brew install loft-sh/tap/devpod
        else
            print_info "Installing DevPod using curl..."
            curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-darwin-$(uname -m)"
            chmod +x devpod
            sudo mv devpod /usr/local/bin/
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        print_info "Installing DevPod using curl..."
        curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-linux-$(uname -m)"
        chmod +x devpod
        sudo mv devpod /usr/local/bin/
    else
        print_error "Unsupported OS. Please install DevPod manually from: https://devpod.sh/docs/getting-started/install"
        exit 1
    fi

    print_success "DevPod installed successfully!"
}

check_credentials() {
    local has_creds=true

    # Check environment variables
    if [ -z "$UPCLOUD_USERNAME" ] && [ -z "$UPCLOUD_PASSWORD" ]; then
        # Check UpCloud CLI config
        if [ ! -f "$HOME/.config/upcloud/config.json" ]; then
            has_creds=false
        fi
    fi

    if [ "$has_creds" = false ]; then
        print_warning "UpCloud credentials not found in environment or CLI config"
        echo ""
        echo "Please set your UpCloud API credentials:"
        echo "  1. Create API user at: https://hub.upcloud.com/account/api"
        echo "  2. Export credentials:"
        echo "     export UPCLOUD_USERNAME='your-username'"
        echo "     export UPCLOUD_PASSWORD='your-password'"
        echo ""
        read -p "Press Enter after setting credentials, or Ctrl+C to exit..."
    else
        print_success "UpCloud credentials detected!"
    fi
}

setup_provider() {
    print_step "Setting up UpCloud provider..."

    # Check if provider already exists
    if devpod provider list | grep -q "$PROVIDER_NAME"; then
        print_info "Provider '$PROVIDER_NAME' already exists"
        read -p "Do you want to update it? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            devpod provider delete "$PROVIDER_NAME" --force
        else
            return 0
        fi
    fi

    # Add the provider - check if release exists
    print_info "Adding UpCloud provider..."

    # Try to add from GitHub release first
    if ! devpod provider add "$PROVIDER_REPO" --name "$PROVIDER_NAME" 2>/dev/null; then
        print_warning "No official release found. Checking for local development setup..."

        # Check if we're in the provider directory or can find it
        if [ -f "provider.yaml" ] && [ -f "main.go" ]; then
            print_info "Found local provider directory. Building and installing locally..."

            # Build the provider
            if check_command go; then
                make build || go build -o bin/devpod-provider-upcloud

                # Install locally
                devpod provider add . --name "$PROVIDER_NAME"
                print_success "Installed from local build!"
            else
                print_error "Go is not installed. Cannot build provider locally."
                echo ""
                echo "Options:"
                echo "1. Install Go from https://go.dev and run this script again"
                echo "2. Wait for an official release at $PROVIDER_REPO/releases"
                echo "3. Build the provider on another machine and copy the binary"
                exit 1
            fi
        else
            print_error "Could not find provider source code or official release."
            echo ""
            echo "This appears to be a pre-release version. Options:"
            echo ""
            echo "1. Clone and build locally:"
            echo "   git clone https://github.com/neuralmux/devpod-provider-upcloud.git"
            echo "   cd devpod-provider-upcloud"
            echo "   ./quickstart.sh"
            echo ""
            echo "2. Wait for the first official release at:"
            echo "   https://github.com/neuralmux/devpod-provider-upcloud/releases"
            echo ""
            exit 1
        fi
    else
        print_success "Provider added from GitHub release!"
    fi

    # Set as default provider
    devpod provider use "$PROVIDER_NAME"

    print_success "UpCloud provider configured successfully!"
}

show_next_steps() {
    echo ""
    echo -e "${GREEN}${BOLD}🎉 Setup Complete!${NC}"
    echo ""
    echo -e "${BOLD}Quick Commands:${NC}"
    echo ""
    echo "  ${BLUE}Create a workspace:${NC}"
    echo "    devpod up github.com/your-org/your-repo --provider upcloud"
    echo ""
    echo "  ${BLUE}Create from local folder:${NC}"
    echo "    devpod up . --provider upcloud"
    echo ""
    echo "  ${BLUE}List workspaces:${NC}"
    echo "    devpod list"
    echo ""
    echo "  ${BLUE}Connect to workspace:${NC}"
    echo "    devpod ssh <workspace-name>"
    echo ""
    echo "  ${BLUE}Stop workspace (save costs):${NC}"
    echo "    devpod stop <workspace-name>"
    echo ""
    echo "  ${BLUE}Delete workspace:${NC}"
    echo "    devpod delete <workspace-name>"
    echo ""
    echo -e "${BOLD}Configuration:${NC}"
    echo ""
    echo "  ${BLUE}View provider options:${NC}"
    echo "    devpod provider options upcloud"
    echo ""
    echo "  ${BLUE}List available server plans:${NC}"
    echo "    devpod-provider-upcloud plans --recommended"
    echo ""
    echo "  ${BLUE}Change server size:${NC}"
    echo "    devpod provider set-options upcloud --option UPCLOUD_PLAN=DEV-2xCPU-8GB"
    echo ""
    echo "  ${BLUE}Change deployment zone:${NC}"
    echo "    devpod provider set-options upcloud --option UPCLOUD_ZONE=us-nyc1"
    echo ""
    echo -e "${BOLD}Available Zones:${NC} de-fra1, fi-hel1, nl-ams1, uk-lon1, us-nyc1, us-chi1, sg-sin1, au-syd1"
    echo -e "${BOLD}Recommended Plans:${NC}"
    echo "  • DEV-1xCPU-1GB-10GB (€3/mo) - Minimal development"
    echo "  • DEV-2xCPU-4GB (€18/mo) - Standard development (default)"
    echo "  • DEV-2xCPU-8GB (€25/mo) - Professional development"
    echo "  • CN-2xCPU-4GB (€16/mo*) - Pay-per-use (*only when running)"
    echo ""
    echo -e "${GREEN}${BOLD}Happy coding! 🚀${NC}"
}

main() {
    clear
    print_header

    # Step 1: Check for DevPod
    print_step "Checking for DevPod installation..."
    if ! check_command devpod; then
        print_warning "DevPod not found"
        read -p "Would you like to install DevPod? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_devpod
        else
            print_error "DevPod is required. Please install it from: https://devpod.sh"
            exit 1
        fi
    else
        print_success "DevPod is installed ($(devpod version))"
    fi
    echo ""

    # Step 2: Check credentials
    print_step "Checking UpCloud credentials..."
    check_credentials
    echo ""

    # Step 3: Setup provider
    setup_provider
    echo ""

    # Step 4: Show next steps
    show_next_steps
}

# Run main function
main