# UpCloud Provider for DevPod

[![DevPod](https://img.shields.io/badge/DevPod-Provider-blue)](https://devpod.sh)
[![UpCloud](https://img.shields.io/badge/UpCloud-Compatible-purple)](https://upcloud.com)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org)
[![CI/CD](https://img.shields.io/badge/CI%2FCD-GitHub_Actions-2088FF?logo=github-actions)](https://github.com/neuralmux/devpod-provider-upcloud/actions)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

The official UpCloud provider for [DevPod](https://github.com/loft-sh/devpod) enables developers to create cloud development environments on [UpCloud](https://upcloud.com) infrastructure.

## Features

- 🚀 **Quick Setup** - Create development environments in under 2 minutes
- 🔐 **SSH Key Authentication** - Secure access with automatic SSH key injection
- 🌍 **Global Zones** - Deploy to 13+ data centers worldwide
- 💾 **Flexible Storage** - Configurable SSD storage with MaxIOPS tier
- 🔄 **Full Lifecycle Management** - Create, start, stop, and delete servers
- 🎯 **Auto-configuration** - Cloud-init support for automatic environment setup
- 💰 **Cost-effective** - Auto-stop functionality to save on idle resources

## What are Providers?

Providers are CLI programs that let DevPod create, manage and run workspaces on different backends. They define commands to create, delete and connect to development environments in the cloud.

DevPod relies on the provider model for flexibility and adaptability. Providers are defined through a provider.yaml that specifies options, configuration, and commands needed for workspace management.

## Usage

To use this provider in your DevPod setup, you will need to do the following steps:

1. See the [DevPod documentation](https://devpod.sh/docs/managing-providers/add-provider)
   for how to add a provider
1. Use the reference `neuralmux/devpod-provider-upcloud` to download the latest
   release from GitHub
1. Set up an [API user](https://upcloud.com/docs/guides/getting-started-upcloud-api/)
   in UpCloud Control Panel. This will be used to manage resources.

## Installation

### Using the CLI

```
devpod provider add github.com/neuralmux/devpod-provider-upcloud
devpod provider use github.com/neuralmux/devpod-provider-upcloud
```

### Pre-requisites

The UpCloud provider needs to be configured with a set of information:

- API user credentials

#### UpCloud API

Behind the scenes, the UpCloud DevPod provider communicates with UpCloud API to provision and manage resources.

### Using the desktop app

Open DevPod app on your computer, then go to "Providers" tab, then click on "Add Provider" button. Then in "Confgiure Provider" popup click on "+" button to add a custom provider.

Enter "github.com/neuralmux/devpod-provider-upcloud" or "neuralmux/devpod-provider-upcloud" as source and fill all the needed information:

![Screencast demo](./assets/desktop-demo.gif)

## Configuration

### Required Options

| Option | Description | Environment Variable |
|--------|-------------|---------------------|
| API Username | Your UpCloud API username | `UPCLOUD_USERNAME` |
| API Password | Your UpCloud API password | `UPCLOUD_PASSWORD` |

### Optional Options

| Option | Description | Default | Environment Variable |
|--------|-------------|---------|---------------------|
| Zone | Data center location | `de-fra1` | `UPCLOUD_ZONE` |
| Plan | Server size | `2xCPU-4GB` | `UPCLOUD_PLAN` |
| Storage | Disk size in GB | `50` | `UPCLOUD_STORAGE` |
| Image | Operating system | `Ubuntu 22.04` | `UPCLOUD_IMAGE` |

### Available Zones

- 🇩🇪 **Europe**: de-fra1, fi-hel1, fi-hel2, nl-ams1, uk-lon1, es-mad1, pl-waw1, se-sto1
- 🇺🇸 **Americas**: us-nyc1, us-chi1, us-sjo1
- 🌏 **Asia-Pacific**: sg-sin1, au-syd1

### Available Plans

- **Development**: 1xCPU-1GB, 1xCPU-2GB, 2xCPU-4GB
- **Professional**: 4xCPU-8GB, 6xCPU-16GB, 8xCPU-32GB

## Usage

### Using the desktop app

1. Open DevPod desktop application
2. Go to the **Providers** tab
3. Click **Add Provider**
4. Select **Custom Provider** and enter: `github.com/neuralmux/devpod-provider-upcloud`
5. Configure your UpCloud API credentials
6. Click on **Workspaces** tab
7. Click **Create Workspace** and select UpCloud as the provider

## Development

### Prerequisites

- Go 1.25 or higher
- Make
- Git
- (Optional) golangci-lint for code quality
- (Optional) goreleaser for release builds

### Quick Start

```bash
# Clone the repository
git clone https://github.com/neuralmux/devpod-provider-upcloud.git
cd devpod-provider-upcloud

# Run development setup
./scripts/dev-setup.sh

# Build the provider
make build

# Run all tests
make test-all
```

### Development Workflow

```bash
# Before pushing changes
make pre-push

# Format code
make fmt

# Run linting
make lint

# Generate test coverage
make coverage
```

### Testing

#### Unit Tests
```bash
make test
```

#### BDD Tests (Godog)
```bash
make bdd
```

#### Local Testing with Mock Mode
```bash
# Test without real API calls
./test-local.sh

# Or manually with mock credentials
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
./bin/devpod-provider-upcloud init
```

#### Testing with Real Credentials
```bash
# Set real credentials
export UPCLOUD_USERNAME="your-username"
export UPCLOUD_PASSWORD="your-password"

# Test provider
./test-provider.sh
```

## Architecture

This provider uses the official [UpCloud Go SDK v8](https://github.com/UpCloudLtd/upcloud-go-api) to interact with the UpCloud API.

### Key Components

- **CLI Framework**: Cobra-based command structure
- **API Client**: Full UpCloud API integration with mock mode support
- **Server Management**: Complete lifecycle control (create, start, stop, delete, status)
- **SSH Integration**: Automatic SSH key injection and secure connectivity
- **Error Handling**: User-friendly error messages with recovery suggestions
- **Cloud-init**: Automatic DevPod agent installation and configuration

### Project Structure

```
├── cmd/                    # CLI commands (Cobra-based)
│   ├── root.go            # Root command setup
│   ├── init.go            # Provider initialization
│   ├── create.go          # Server creation
│   ├── delete.go          # Server deletion
│   ├── start.go           # Server start
│   ├── stop.go            # Server stop
│   ├── status.go          # Server status
│   └── command.go         # Command execution
├── pkg/                   
│   ├── upcloud/           # UpCloud API client
│   │   ├── client.go      # Main client implementation
│   │   ├── constants.go   # UpCloud-specific constants
│   │   ├── mapper.go      # Resource mapping utilities
│   │   └── errors.go      # Error handling
│   └── options/           # Configuration management
│       └── options.go     # Environment variable parsing
├── features/              # BDD test specifications
├── scripts/               # Development and CI/CD scripts
├── .github/workflows/     # GitHub Actions CI/CD
└── provider.yaml          # DevPod provider manifest
```

## CI/CD Pipeline

The project includes a comprehensive CI/CD pipeline with:

- **Automated Testing**: Unit, BDD, and integration tests
- **Security Scanning**: CodeQL, Gosec, Trivy vulnerability scanning
- **Quality Checks**: golangci-lint with multiple linters
- **Multi-platform Builds**: Linux, macOS, Windows (amd64/arm64)
- **Automated Releases**: Tag-based releases with GoReleaser
- **Docker Images**: Multi-arch container builds
- **Dependency Management**: Automated updates with Dependabot

See [CI_CD_GUIDE.md](CI_CD_GUIDE.md) for detailed pipeline documentation.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Set up development environment (`./scripts/dev-setup.sh`)
4. Make your changes
5. Run pre-push validation (`make pre-push`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation as needed
- Ensure all CI checks pass
- Use conventional commit messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/neuralmux/devpod-provider-upcloud/issues)
- **DevPod**: [DevPod Documentation](https://devpod.sh/docs)
- **UpCloud**: [UpCloud Documentation](https://upcloud.com/docs)

## Disclaimer

This software is provided "as is" without warranty of any kind, express or implied. Use it at your own risk.