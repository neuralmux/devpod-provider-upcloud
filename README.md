# UpCloud Provider for DevPod

The UpCloud provider for [Loft Labs' DevPod](https://github.com/loft-sh/devpod).

## What are Providers?

Providers are simple CLI programs that let DevPod create, manage and run the workspaces requested by the user. In the simplest form, a provider defines a command to create, delete and connect to a virtual machine in a cloud.

DevPod relies on the provider model in order to allow flexibility and adaptability for any backend of choice. Providers in DevPod are defined through a provider.yaml that defines the necessary options, configuration, binaries and commands needed to handle workspace creation.

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

## Usage

### Using the desktop app

Click on the `Workspaces` tab, then on the `Create Workspace` button and fill all the needed information.

## Contributions

Contributions are welcome! If you find any issues or want to add new features, please open an issue or submit a pull request on the GitHub repository.

## Disclaimer

This software is provided "as is" without warranty of any kind, express or implied. Use it at your own risk.