#!/bin/bash
# Development environment setup for UpCloud DevPod Provider

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

echo -e "${BLUE}🛠️  UpCloud DevPod Provider - Development Environment Setup${NC}"
echo ""

cd "$ROOT_DIR"

# 1. Check Go installation
log_info "Checking Go installation..."
if ! command_exists go; then
    log_error "Go is not installed. Please install Go 1.25+ from https://golang.org/"
    exit 1
fi

GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1)
GO_MAJOR=$(echo "$GO_VERSION" | cut -d. -f1 | sed 's/go//')
GO_MINOR=$(echo "$GO_VERSION" | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 25 ]); then
    log_error "Go 1.25+ is required. Current version: $GO_VERSION"
    exit 1
fi

log_success "Go $GO_VERSION is installed"

# 2. Install development tools
log_info "Installing/updating development tools..."

# golangci-lint
if ! command_exists golangci-lint; then
    log_info "Installing golangci-lint..."
    if command_exists brew; then
        brew install golangci-lint
    else
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.0
    fi
    log_success "golangci-lint installed"
else
    log_success "golangci-lint already installed"
fi

# govulncheck
if ! command_exists govulncheck; then
    log_info "Installing govulncheck..."
    go install golang.org/x/vuln/cmd/govulncheck@latest
    log_success "govulncheck installed"
else
    log_success "govulncheck already installed"
fi

# gosec
if ! command_exists gosec; then
    log_info "Installing gosec..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    log_success "gosec installed"
else
    log_success "gosec already installed"
fi

# goreleaser (for local testing)
if ! command_exists goreleaser; then
    log_info "Installing goreleaser..."
    if command_exists brew; then
        brew install goreleaser/tap/goreleaser
    else
        log_warning "goreleaser not installed - install with: https://goreleaser.com/install/"
    fi
else
    log_success "goreleaser already installed"
fi

# Optional DevPod CLI for testing
if ! command_exists devpod; then
    log_warning "DevPod CLI not installed (recommended for testing)"
    log_info "Install from: https://devpod.sh/docs/getting-started/install"
fi

# Optional tools with warnings if not available
if ! command_exists yamllint; then
    log_warning "yamllint not installed (optional for provider.yaml validation)"
    if command_exists brew; then
        log_info "Install with: brew install yamllint"
    elif command_exists pip3; then
        log_info "Install with: pip3 install yamllint"
    fi
fi

# 3. Download Go modules
log_info "Downloading Go modules..."
go mod download
go mod tidy
log_success "Go modules downloaded and tidied"

# 4. Build the project
log_info "Building UpCloud DevPod provider..."
make build
log_success "Provider binary built successfully"

# 5. Install Git hooks (if available)
if [ -f "./scripts/install-git-hooks.sh" ]; then
    log_info "Installing Git hooks..."
    ./scripts/install-git-hooks.sh
    log_success "Git hooks installed"
else
    log_info "Creating basic Git pre-commit hook..."
    mkdir -p .git/hooks
    cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
# Pre-commit hook for UpCloud DevPod Provider

echo "Running pre-commit checks..."

# Run linter
if ! make vet; then
    echo "❌ Linting failed"
    exit 1
fi

# Run tests
if ! make test; then
    echo "❌ Tests failed"
    exit 1
fi

echo "✅ Pre-commit checks passed"
EOF
    chmod +x .git/hooks/pre-commit
    log_success "Basic Git pre-commit hook created"
fi

# 6. Run quick validation
log_info "Running initial validation..."
if make test > /dev/null 2>&1; then
    log_success "Initial tests passed"
else
    log_warning "Some tests failed (normal for fresh setup without real credentials)"
fi

# 7. Test with mock credentials
log_info "Testing provider with mock credentials..."
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
if ./bin/devpod-provider-upcloud init > /dev/null 2>&1; then
    log_success "Provider test mode working"
else
    log_warning "Provider test mode failed (check binary)"
fi

# 8. Setup summary
echo ""
echo -e "${GREEN}🎉 UpCloud DevPod Provider development environment setup complete!${NC}"
echo ""
echo -e "${BLUE}📦 Project Structure:${NC}"
echo "  cmd/             - CLI commands (Cobra-based)"
echo "  pkg/upcloud/     - UpCloud API client"
echo "  pkg/options/     - Environment variable parsing"
echo "  features/        - BDD tests with Godog"
echo "  provider.yaml    - DevPod provider manifest"
echo ""
echo -e "${BLUE}🚀 Available commands:${NC}"
echo "  make build       - Build provider binary"
echo "  make test        - Run unit tests"
echo "  make bdd         - Run BDD tests with Godog"
echo "  make test-all    - Run all tests"
echo "  make lint        - Run linter (go vet)"
echo "  make fmt         - Format code"
echo "  make clean       - Clean build artifacts"
echo "  make coverage    - Generate test coverage"
echo "  make help        - Show all available commands"
echo ""
echo -e "${BLUE}🧪 Testing the provider:${NC}"
echo "  ./test-local.sh                    - Test without API calls"
echo "  ./bin/devpod-provider-upcloud init - Test with real credentials"
echo ""
echo -e "${BLUE}🔧 Development workflow:${NC}"
echo "  1. Set UpCloud credentials: export UPCLOUD_USERNAME=... UPCLOUD_PASSWORD=..."
echo "  2. Make your changes to Go files"
echo "  3. Run 'make build' to compile"
echo "  4. Run 'make test-all' for full validation"
echo "  5. Test with './test-local.sh' or real server creation"
echo "  6. Git hooks will run automatically on commit"
echo ""
echo -e "${BLUE}📊 CI/CD Pipeline:${NC}"
echo "  • GitHub Actions workflows in .github/workflows/"
echo "  • Push to trigger CI, tag vX.Y.Z to release"
echo "  • Multi-platform binaries built automatically"
echo ""
echo -e "${BLUE}🛠️  Installed tools:${NC}"
echo "  ✅ Go $GO_VERSION"
echo "  ✅ golangci-lint (code linting)"
echo "  ✅ govulncheck (vulnerability scanning)"
echo "  ✅ gosec (security analysis)"
if command_exists goreleaser; then
    echo "  ✅ goreleaser (release builds)"
else
    echo "  ⚠️  goreleaser (optional - release builds)"
fi
if command_exists devpod; then
    echo "  ✅ DevPod CLI (provider testing)"
else
    echo "  ⚠️  DevPod CLI (recommended - provider testing)"
fi
if command_exists yamllint; then
    echo "  ✅ yamllint (YAML validation)"
else
    echo "  ⚠️  yamllint (optional - YAML validation)"
fi
echo "  ✅ Git pre-commit hooks"
echo ""
echo -e "${YELLOW}💡 Next steps:${NC}"
echo "  1. Get UpCloud API credentials from https://hub.upcloud.com/"
echo "  2. Export UPCLOUD_USERNAME and UPCLOUD_PASSWORD"
echo "  3. Run './bin/devpod-provider-upcloud init' to test"
echo "  4. Try creating a test server with './test-provider.sh'"
echo "  5. Happy coding! 🚀"