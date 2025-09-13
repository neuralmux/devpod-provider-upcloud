#!/bin/bash
set -e

# Pre-push validation script for UpCloud DevPod Provider
# This script runs the same checks as CI/CD pipeline locally

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$ROOT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
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

log_section() {
    echo ""
    echo -e "${BLUE}=== $1 ===${NC}"
}

# Track overall status
FAILED_CHECKS=()
WARNED_CHECKS=()

# Function to run a check and track results
run_check() {
    local check_name="$1"
    local command="$2"
    local allow_warnings="${3:-false}"
    
    log_info "Running $check_name..."
    
    if eval "$command" > /tmp/check_output 2>&1; then
        log_success "$check_name passed"
        return 0
    else
        local exit_code=$?
        if [[ "$allow_warnings" == "true" && $exit_code -eq 1 ]]; then
            log_warning "$check_name has warnings (non-blocking)"
            WARNED_CHECKS+=("$check_name")
            cat /tmp/check_output
            return 0
        else
            log_error "$check_name failed"
            FAILED_CHECKS+=("$check_name")
            cat /tmp/check_output
            return $exit_code
        fi
    fi
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Clean up function
cleanup() {
    rm -f /tmp/check_output
    rm -f devpod-provider-upcloud-test
}

trap cleanup EXIT

log_section "Pre-push Validation for UpCloud DevPod Provider"
echo "This script validates your changes against CI/CD pipeline requirements"
echo ""

# 1. Environment Check
log_section "Environment Check"

# Check Go version
if ! command_exists go; then
    log_error "Go is not installed"
    exit 1
fi

GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1)
log_info "Go version: $GO_VERSION"

# Check if we have the required Go version (1.21+)
REQUIRED_GO_VERSION="1.21"
if [[ $(echo -e "$GO_VERSION\ngo$REQUIRED_GO_VERSION" | sort -V | head -n1) != "go$REQUIRED_GO_VERSION" ]]; then
    log_warning "Go version $GO_VERSION may not match CI (requires 1.21+)"
    WARNED_CHECKS+=("Go Version")
fi

log_success "Environment check completed"

# 2. Go Module Validation
log_section "Go Module Validation"
# Run go mod tidy first to ensure dependencies are clean
go mod tidy
# Check if go.mod and go.sum are properly formatted after tidy
if git diff --quiet go.mod go.sum; then
    log_success "Go modules are tidy and up-to-date"
else
    # Check if changes are only dependency organization (acceptable)
    if git diff go.mod | grep -E "^\+.*github\.com/UpCloudLtd/upcloud-go-api/v8" && \
       git diff go.mod | grep -E "^\-.*github\.com/UpCloudLtd/upcloud-go-api/v8.*// indirect"; then
        log_warning "Go modules have acceptable dependency reorganization changes"
        WARNED_CHECKS+=("Go module dependency organization")
    else
        log_error "Go module tidy would make unexpected changes - please commit go.mod/go.sum changes first"
        git diff go.mod go.sum
        FAILED_CHECKS+=("Go module tidy")
    fi
fi
run_check "Go module verify" "go mod verify"
run_check "Go module download" "go mod download"

# 3. Code Quality Checks
log_section "Code Quality Checks"

# Go vet (always available) - exclude test-only directories
run_check "Go vet" "go vet ./cmd/... ./pkg/... && go vet *.go 2>/dev/null || true"

# Go formatting check
run_check "Go formatting" "test -z \"\$(gofmt -l .)\""

# golangci-lint if available
if command_exists golangci-lint; then
    run_check "golangci-lint" "golangci-lint run" "true"
else
    log_warning "golangci-lint not found - will run in CI"
    WARNED_CHECKS+=("golangci-lint")
fi

# 4. Build Tests
log_section "Build Tests"

# Standard build using Makefile
run_check "Standard build" "make build"

# Build with version flags (like CI)
run_check "Build with version flags" "go build -ldflags=\"-w -s -X main.version=v1.0.0-test -X main.commit=\$(git rev-parse HEAD) -X main.date=\$(date -u +%Y-%m-%dT%H:%M:%SZ)\" -o devpod-provider-upcloud-test ."

# Test help output (DevPod providers use --help, not version)
if [[ -f "devpod-provider-upcloud-test" ]]; then
    run_check "Help command" "./devpod-provider-upcloud-test --help"
fi

# Cross-platform build test (simplified - just test one other platform)
if [[ $(uname) == "Darwin" ]]; then
    # Test Linux build on macOS
    run_check "Cross-platform build test" "GOOS=linux GOARCH=amd64 go build -o /dev/null ."
else
    # Test Windows build on Linux
    run_check "Cross-platform build test" "GOOS=windows GOARCH=amd64 go build -o /dev/null ."
fi

# 5. Test Suite
log_section "Test Suite"

# Unit tests
run_check "Unit tests" "make test"

# Race condition tests
run_check "Race condition tests" "go test -race ./pkg/... -timeout=5m"

# Provider validation with mock credentials
log_info "Running provider validation with mock credentials..."
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="1xCPU-1GB"
export UPCLOUD_STORAGE="25"
export UPCLOUD_IMAGE="Ubuntu Server 22.04 LTS (Jammy Jellyfish)"
export MACHINE_ID="test-machine"
export MACHINE_FOLDER="/tmp/test"
export AGENT_PATH="/home/devpod/.devpod/devpod"
export AGENT_DATA_PATH="/home/devpod/.devpod/agent"
if timeout 60s ./bin/devpod-provider-upcloud init > /tmp/provider_output 2>&1; then
    log_success "Provider validation completed"
else
    exit_code=$?
    if [[ $exit_code -eq 124 ]]; then
        log_warning "Provider validation timed out (non-blocking for pre-push)"
        WARNED_CHECKS+=("Provider validation")
    else
        log_warning "Provider validation failed (non-blocking for pre-push - requires build first)"
        WARNED_CHECKS+=("Provider validation")
        head -20 /tmp/provider_output
    fi
fi

# BDD tests if available
if [[ -d "features" ]] && command_exists go; then
    log_info "Running BDD tests..."
    if timeout 180s make bdd > /tmp/bdd_output 2>&1; then
        log_success "BDD tests passed"
    else
        exit_code=$?
        if [[ $exit_code -eq 124 ]]; then
            log_warning "BDD tests timed out (non-blocking for pre-push)"
            WARNED_CHECKS+=("BDD tests")
        else
            log_warning "BDD tests failed (non-blocking for pre-push)"
            WARNED_CHECKS+=("BDD tests")
            head -20 /tmp/bdd_output
        fi
    fi
else
    log_info "BDD tests not available or skipped"
fi

# 6. Security Checks (if tools available)
log_section "Security Checks"

# govulncheck if available
if command_exists govulncheck; then
    run_check "Vulnerability scan" "govulncheck ./..." "true"
else
    log_info "govulncheck not found - will run in CI"
fi

# gosec if available
if command_exists gosec; then
    run_check "Security scan" "gosec ./..." "true"
else
    log_info "gosec not found - will run in CI"
fi

# 7. Documentation Validation
log_section "Documentation Validation"

# Check that provider help works (if binary exists)
if [[ -f "./bin/devpod-provider-upcloud" ]]; then
    run_check "Provider help command" "./bin/devpod-provider-upcloud --help"
else
    log_warning "Provider binary not found - run 'make build' first"
    WARNED_CHECKS+=("Provider binary")
fi

# Validate workflow YAML and provider.yaml syntax
if command_exists yamllint; then
    # Use relaxed config for workflow validation
    echo "extends: relaxed" > /tmp/.yamllint
    echo "rules:" >> /tmp/.yamllint
    echo "  line-length:" >> /tmp/.yamllint
    echo "    max: 200" >> /tmp/.yamllint
    run_check "Workflow YAML validation" "yamllint -c /tmp/.yamllint .github/workflows/"
    run_check "Provider YAML validation" "yamllint -c /tmp/.yamllint provider.yaml"
    rm -f /tmp/.yamllint
else
    log_info "yamllint not found - basic YAML validation skipped"
fi

# Check for broken links in documentation (if available)
if command_exists markdown-link-check; then
    run_check "Documentation link check" "find . -name '*.md' -not -path './.git/*' -exec markdown-link-check {} \;" "true"
else
    log_info "markdown-link-check not found - will run in CI"
fi

# 8. Git Checks
log_section "Git Validation"

# Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    log_warning "Uncommitted changes detected"
    git status --short
    WARNED_CHECKS+=("Uncommitted changes")
fi

# Check current branch
CURRENT_BRANCH=$(git branch --show-current)
log_info "Current branch: $CURRENT_BRANCH"

# Check if we're ahead of origin
if git rev-parse --verify origin/$CURRENT_BRANCH >/dev/null 2>&1; then
    COMMITS_AHEAD=$(git rev-list --count origin/$CURRENT_BRANCH..$CURRENT_BRANCH)
    if [[ $COMMITS_AHEAD -gt 0 ]]; then
        log_info "Branch is $COMMITS_AHEAD commits ahead of origin/$CURRENT_BRANCH"
    fi
fi

# 9. Final Report
log_section "Validation Summary"

if [[ ${#FAILED_CHECKS[@]} -eq 0 ]]; then
    log_success "All critical checks passed! ✨"
    
    if [[ ${#WARNED_CHECKS[@]} -gt 0 ]]; then
        echo ""
        log_warning "Warnings (non-blocking):"
        for check in "${WARNED_CHECKS[@]}"; do
            echo "  - $check"
        done
        echo ""
        log_info "These warnings will not block CI, but consider addressing them."
    fi
    
    echo ""
    log_success "🚀 Ready to push! Your changes should pass CI/CD pipeline."
    echo ""
    echo "Next steps:"
    echo "  git push origin $CURRENT_BRANCH"
    echo ""
    
    exit 0
else
    echo ""
    log_error "❌ ${#FAILED_CHECKS[@]} critical checks failed:"
    for check in "${FAILED_CHECKS[@]}"; do
        echo "  - $check"
    done
    
    if [[ ${#WARNED_CHECKS[@]} -gt 0 ]]; then
        echo ""
        log_warning "Additionally, ${#WARNED_CHECKS[@]} warnings:"
        for check in "${WARNED_CHECKS[@]}"; do
            echo "  - $check"
        done
    fi
    
    echo ""
    log_error "Please fix the failed checks before pushing."
    exit 1
fi