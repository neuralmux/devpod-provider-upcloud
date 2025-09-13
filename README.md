# UpCloud Provider for DevPod

[![DevPod](https://img.shields.io/badge/DevPod-Provider-blue)](https://devpod.sh)
[![UpCloud](https://img.shields.io/badge/UpCloud-Compatible-purple)](https://upcloud.com)
[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?logo=go)](https://golang.org)

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

Enter "github.com/alexandrevilain/devpod-provider-ovhcloud" or "alexandrevilain/devpod-provider-ovhcloud" as source and fill all the needed information:

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

### Building from Source

```bash
# Clone the repository
git clone https://github.com/neuralmux/devpod-provider-upcloud.git
cd devpod-provider-upcloud

# Install dependencies
make deps

# Build the provider
make build

# Run tests
make test-all
```

### Testing Locally

```bash
# Set credentials
export UPCLOUD_USERNAME="your-username"
export UPCLOUD_PASSWORD="your-password"

# Test provider initialization
./bin/devpod-provider-upcloud init

# Install locally
make install
```

## Architecture

This provider uses the official [UpCloud Go SDK v8](https://github.com/UpCloudLtd/upcloud-go-api) to interact with the UpCloud API. It implements:

- **Server Management**: Full lifecycle control (create, start, stop, delete)
- **SSH Integration**: Automatic SSH key injection and secure connectivity
- **Error Handling**: Comprehensive error messages and retry logic
- **State Management**: Proper state tracking and transitions

## Contributions

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/neuralmux/devpod-provider-upcloud/issues)
- **DevPod**: [DevPod Documentation](https://devpod.sh/docs)
- **UpCloud**: [UpCloud Documentation](https://upcloud.com/docs)

## Disclaimer

This software is provided "as is" without warranty of any kind, express or implied. Use it at your own risk.