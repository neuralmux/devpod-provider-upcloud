# Developing a DevPod provider

DevPod providers are small CLI programs defined through a provider.yaml that DevPod interacts with, in order to bring up the workspace.

Providers are standalone programs that DevPod will call, parsing a manifest called provider.yaml that will instruct DevPod on how to interact with the program. The most important sections of a provider.yaml are:

- `exec`: Defines what commands DevPod should execute to interact with the environment
- `options`: Defines what Options the user can configure in the provider
- `binaries`: Defines what additional helper binaries are required to run the provider
- `agent`: Defines configuration options for the provider such as drivers, auto-inactivity timeout and credentials injection

## Provider.yaml

With the provider.yaml, DevPod will know how to call the provider's binary, in order to perform all the required actions to bring up or down an environment for a workspace.

Following is a minimal example manifest for a provider:

```yaml
name:  name-of-provider
version: version-number
description: quick description # Optional
icon: https://url-to-icon.com  # Shown in the Desktop App
options:
  # Options for the provider, DevPod will pass these as
  # ENV Variables when calling the provider
  OPTION_NAME:
    description: "option description"
    default: "value"
    required: true # or false
  AGENT_PATH:
    description: The path where to inject the DevPod agent to.
    default: /opt/devpod/agent
agent:
  path: ${AGENT_PATH}
exec:
  command: # Required: a command to execute on the remote machine or container
  init:    # Optional: a command to init the provider, login to an account or similar
  create:  # Optional: a command to create the machine
  delete:  # Optional: a command to delete the machine
  start:   # Optional: a command to start the machine
  stop:    # Optional: a command to stop the machine
  status:  # Optional: a command to get the machine's status
binaries:  # Optional binaries DevPod should download for this provider
  MY_BINARY: # Will be available as MY_BINARY environment variable in the exec section
    ...
```

## How DevPod Interacts with a Provider

DevPod uses the exec section within the provider.yaml to know what commands it needs to execute to interact with the provider environments. These "commands" are regular POSIX shell scripts that call a helper program or directly interact with the underlying system of the user to manage the environment.

In the exec section of the provider.yaml, the following commands are allowed:

- command: The only command that is required for a provider to work, which defines how to run a command in the environment. DevPod will use this command to inject itself into the environment and route all communication through the commands standard output and input. An example for a local development provider would be: sh -c "${COMMAND}".
- init: Optional command to check if options are defined correctly and the provider is ready to create environments. For example for the Docker provider, this command checks if Docker is installed and reachable locally.
- create: Optional command how to create a machine. If this command is defined, the provider will automatically be treated as a machine provider and DevPod will also expect delete to be defined.
- delete: Optional command how to delete a machine. Counter command to create.
- start: Optional command how to start a stopped machine. Only usable for machine providers.
- stop: Optional command how to stop a machine. Only usable for machine providers.
- status: Optional command how to retrieve the status of a machine. Expects one of the following statuses on standard output:
-- Running: Machine is running and ready
-- Busy: Machine is doing something and DevPod should wait (e.g. terminating, starting, stopping etc.)
-- Stopped: Machine is currently stopped
-- NotFound: Machine is not found

The init is used to perform actions in order to verify and validate the options passed (in this case try a dummy command on the selected host)

The command is used to access the remote environment. ${COMMAND} will be supplied by DevPod and is a placeholder for the command DevPod will execute. In case of Machines providers this will usually have to SSH on the VMs created.

## Provider Options

Inside the provider.yaml, you can specify options that DevPod can pass to the provider when calling it. Each option will be passed as an environment variable to the commands or can be used directly inside the agent section of a provider.

```yaml
...
...
options:
  MY_OPTION_NAME:
    description: "this is my option"
    default: "default_value"
    required: false
    password: true
...
...
```

### How Options Work

Options are variables needed for the provider to function properly, for example:

- User Accounts
- VMs images
- VMs sizes
- Account region

These options are parsed and validated by DevPod when Adding the provider and passed to the provider as environment variables.

It's the provider's job to retrieve them from the environment and validate them. It's recommended to make use of the init command that will be called by DevPod when options change to validate environment variables on the provider side.

You can check our example in the Devpod's AWS Provider where we parse and validate the variables:

```yaml
...
    diskSizeGB, err := fromEnvOrError("AWS_DISK_SIZE")
    if err != nil {
        return nil, err
    }

    retOptions.DiskSizeGB, err = strconv.Atoi(diskSizeGB)
    if err != nil {
        return nil, err
    }
...
```

Options will also be passed to the agent and can be used in the agent.exec section. This is very useful if you require certain information on the agent side to perform an auto-inactivity timeout.

### Option Configuration

Each option has a set of attributes that can modify how DevPod interprets it when configuring or adding the provider:

- description: Description shown in devpod provider options and in the Desktop App
- default: Default value of the option provided as a string. Can also reference other variables, e.g. ${MY_OTHER_VAR}-suffix
- required: Boolean if this option needs to be non-empty before using the provider. DevPod will ask in the CLI and make sure that this option is filled in the Desktop application.
- password: Boolean to indicate this is a sensitive value. Prevents this value from showing up in the devpod provider options command and will be a password field in the Desktop application.
- suggestions: An array of suggestions for this option. Will be shown as auto complete options in the DevPod desktop application
- command: A command to retrieve the option value automatically. Can also reference other variables in the command, e.g. echo ${MY_OTHER_VAR}-suffix. For compatibility reasons, this command will be executed in an emulated shell on Windows.
- local: If true, the option will be filled individually for each machine / workspace
- global: If true, the option will be reused for each machine / workspace
- cache: If non-empty, DevPod will re-execute the command after the given timeout. E.g. if this is 5m, DevPod will re-execute the command after 5 minutes to re-fill this value. This is useful if you want to store a token or something that expires locally in a variable.
- hidden: If true, DevPod will not show this option in the Desktop application or through devpod provider options. Can be used to calculate variables internally or save tokens or other things internally.

### Default Values

As the name implies, this is a default value for the option. It is always advisable to place a sensible default for any option.

You can also reference other options inside the default value, e.g. ${MY_OTHER_VAR}-suffix. DevPod will automatically figure out what options need to be resolved before this option.

If not specified, it defaults to an empty string.

### Required Options

If an option is required, and no default is set, DevPod will prompt the user for a value when adding the provider.

In the DevPod Desktop App, the required options will be displayed and prompted right in the Provider's "Add" page.

If not specified, it defaults to false.

### Password Options

If specified and true, the option's value will be treated as a secret, so it won't be shown when listing options.

Example:

```sh
~$ adevpod provider options civo

          NAME            | REQUIRED |          DESCRIPTION           |               DEFAULT                |                VALUE
----------------------------+----------+--------------------------------+--------------------------------------+---------------------------------------
AGENT_PATH                | false    | The path where to inject the   | /var/lib/toolbox/devpod              | /var/lib/toolbox/devpod
                          |          | DevPod agent to.               |                                      |
CIVO_API_KEY              | true     | The civo api key to use        |                                      | ********
CIVO_DISK_IMAGE           | false    | The disk image to use.         | d927ad2f-5073-4ed6-b2eb-b8e61aef29a8 | d927ad2f-5073-4ed6-b2eb-b8e61aef29a8

...
```

If not specified, it defaults to false.

### Options Suggestions

Suggestions are a list of possible values for the option. Suggested use-cases could be for regions/locations, VM sizes, etc.

If not specified, it defaults to empty and ignored.

### Command Options

The command option lets you define a possible value for an option based on a shell command launched on your machine. Can also reference other variables in the command, e.g. echo ${MY_OTHER_VAR}-suffix. For compatibility reasons, this command will be executed in an emulated shell on Windows.

One example would be to forward ENV variables from your machine to the provider, for example:

```yaml
  AWS_ACCESS_KEY_ID:
    description: The aws access key id
    required: false
    command: printf "%s" "${AWS_ACCESS_KEY_ID:-}"
  AWS_SECRET_ACCESS_KEY:
    description: The aws secret access key
    required: false
    command: printf "%s" "${AWS_SECRET_ACCESS_KEY:-}"
```

Or running an helper command (defined in the binaries section), and forwarding the result as the option's value:

```yaml
  AWS_TOKEN:
    local: true
    hidden: true
    cache: 5m
    description: "The AWS auth token to use"
    command: |-
      ${AWS_PROVIDER} token
```

If not specified, it defaults to empty and ignored.

### Built-In Options

There are a couple of predefined options from DevPod, that can be used within the default field of another option or in an option command. Some built-in options are only available for local options as the MACHINE_ID might not be available already. Predefined options:

- DEVPOD: Absolute path to the current DevPod CLI binary. Can be used to call a helper function within DevPod or any other DevPod command. Also available on the agent side.
- DEVPOD_OS: Current Operating system. Can be either: linux, darwin or windows
- DEVPOD_ARCH: Current operating system architecture. Can be either: amd64 or arm64.
- MACHINE_ID: The machine id that should be used. (Only available for local options, commands and machine providers)
- MACHINE_FOLDER: The machine folder that can be used to cache information locally. (Only available for local options, commands and machine providers)
- MACHINE_CONTEXT: The DevPod context this machine was created in. (Only available for local options, commands and machine providers)
- MACHINE_PROVIDER: The provider name that was used to create this machine. (Only available for local options, commands and machine providers)
- WORKSPACE_ID: The workspace id that should be used. (Only available for local options, commands and non-machine providers)
- WORKSPACE_FOLDER: The workspace folder that can be used to cache information locally. (Only available for local options, commands and non-machine providers)
- WORKSPACE_CONTEXT: The DevPod context this workspace was created in. (Only available for local options, commands and non-machine providers)
- WORKSPACE_PROVIDER: The provider name that was used to create this workspace. (Only available for local options, commands and non-machine providers)
- PROVIDER_ID: The provider name. (Only available for local options, commands and non-machine providers)
- PROVIDER_CONTEXT: The provider context. (Only available for local options, commands and non-machine providers)
- PROVIDER_FOLDER: The provider folder where the provider config is saved in, can be used to save global information about the provider such as global session tokens etc. (Only available for local options, commands and non-machine providers)

### Option Groups

You can organize your options in groups, for example:

```yaml
optionGroups:
  - options:
      - AWS_ACCESS_KEY_ID
      - AWS_SECRET_ACCESS_KEY
      - AWS_AMI
      - AWS_DISK_SIZE
      - AWS_INSTANCE_TYPE
      - AWS_VPC_ID
    name: "AWS options"
    defaultVisible: true
  - options:
      - AGENT_PATH
      - INACTIVITY_TIMEOUT
      - INJECT_DOCKER_CREDENTIALS
      - INJECT_GIT_CREDENTIALS
    name: "Agent options"
    defaultVisible: false
```

Options are easily grouped by listing them, each group has a name and a defaultVisible property, which is false by default. If defaultVisible is false, then an user will need to manually expand the option group in the Desktop App.

## Provider Binaries

The binaries section can be used to specify helper binaries DevPod should download that help the provider to accomplish its tasks.

An example of this type of provider are:

- devpod-provider-aws
- devpod-provider-azure
- devpod-provider-civo
- devpod-provider-digitalocean
- devpod-provider-gcloud

Each binary that is required is declared through:

```yaml
binaries:
  NAME:
    - os: # Which OS is this specific binary
      arch: # Binary arch
      path: # Remote (URL) or local path to binary
      checksum:  # sha sum of the binary
      archivePath: # If its an archive, the relative path to the binary. Supported archives are .tgz, .tar, .tar.gz, .zip
```

When Adding a provider, DevPod will match the binary for your OS and Arch and download the specific one for it.

Example of the binary section in a provider.yaml:

```yaml
binaries:
  AWS_PROVIDER:
    - os: linux
      arch: amd64
      path: https://github.com/loft-sh/devpod-provider-aws/releases/download/v0.0.1-alpha.15/devpod-provider-aws-linux-amd64
      checksum: d1e774419d90c3ed399963d9322d57bfdcee189767eabb076a2c2e926bfd9b8b
    - os: linux
      arch: arm64
      path: https://github.com/loft-sh/devpod-provider-aws/releases/download/v0.0.1-alpha.15/devpod-provider-aws-linux-arm64
      checksum: fa15c13e3f0619170d002f9dae3ef41c9949a4595a71c5efe364d89ada604cec
    - os: darwin
      arch: amd64
      path: https://github.com/loft-sh/devpod-provider-aws/releases/download/v0.0.1-alpha.15/devpod-provider-aws-darwin-amd64
      checksum: fb89d41f6ce3e01e953f3ffd18f85bd5a42dd633abafd5d586dc9d9b1322166c
    - os: darwin
      arch: arm64
      path: https://github.com/loft-sh/devpod-provider-aws/releases/download/v0.0.1-alpha.15/devpod-provider-aws-darwin-arm64
      checksum: 82b6713069fa061ea59941600ed32a15f73806a9af3074d67a20ed367d18b2aa
    - os: windows
      arch: amd64
      path: https://github.com/loft-sh/devpod-provider-aws/releases/download/v0.0.1-alpha.15/devpod-provider-aws-windows-amd64.exe
      checksum: 49bd899d439f38d4e8647102db1c18b7a0d5242b3c09c89071b20a5444e20a81
```

### Binary Checksum

Each binary is also verified over an expected checksum. This is important to ensure that whatever binary is declared in provider.yaml is indeed executed on the machine.

### Use Binaries in Commands

DevPod will make the binary path available through an environment variable within the exec section. For example:

```yaml
binaries:
  MY_BINARY:
    ....

exec:
  init: ${MY_BINARY} init
  ....
```

### Use Binaries in Options

You can also use binaries within the option command attribute. For example:

```yaml
binaries:
  MY_BINARY:
    ....

options:
  MY_OPTION:
    command: ${MY_BINARY} retrieve-option
```

### Use Binaries on the Agent Side

You can also define binaries DevPod should install on the agent side through agent.binaries. These binaries can then be used within the agent.exec section to automatically stop a virtual machine if inactive. For example:

```yaml
agent:
  path: ${AGENT_PATH}
  binaries:
    GCLOUD_PROVIDER:
      - os: linux
        arch: amd64
        path: https://github.com/loft-sh/devpod-provider-gcloud/releases/download/v0.0.1-alpha.10/devpod-provider-gcloud-linux-amd64
        checksum: 38f92457507563ee56ea40a2ec40196d12ac2bbd50a924d76f55827e96e5f831
      - os: linux
        arch: arm64
        path: https://github.com/loft-sh/devpod-provider-gcloud/releases/download/v0.0.1-alpha.10/devpod-provider-gcloud-linux-arm64
        checksum: 48e8dfa20962f1c3eb1e3da17d57842a0e26155df2b94377bcdf5b8070d7b17e
  exec:
    shutdown: |-
      ${GCLOUD_PROVIDER} stop --raw
```

## Provider Agent

When DevPod connects through a Provider to an environment, it will inject itself into the environment to handle the following tasks:

- deploying the container
- forward credentials
- ssh server
- auto-shutdown after a period of inactivity

This counterpart is called the DevPod agent, which is available in the same DevPod binary under devpod agent. Within the provider.yaml you can configure certain parts of how local DevPod should interact with its agent counterpart.

### Agent Section

The following options are available in the agent section

```yaml
agent: # You can also use options within this section (see injectGitCredentials as an example)
  path: $\{DEVPOD\}
  driver: docker # Optional, default: docker
  inactivityTimeout: 10m
  containerInactivityTimeout: 10m
  injectGitCredentials: ${INJECT_GIT_CREDENTIALS}
  injectDockerCredentials: ${INJECT_DOCKER_CREDENTIALS}
  binaries:
    MY_BINARY:
      - os: linux
        arch: amd64
        path: https://url-to-binary.com
        checksum: shasum-of-binary
  exec:
    shutdown: |-
      ${MY_BINARY} stop
```

Breaking down the options:

- path: where to place the agent on the remote machine. Use ${DEVPOD} here if you want to use the local machine instead.
- driver: which driver to use to run container, check the Drivers section for more information
- inactivityTimeout: after how much time to shut down the machine. Use for machine providers
- containerInactivityTimeout: after how much time to shut down the container. Use for non-machine providers
- injectGitCredentials: whether to inject git credentials into the machine.
- injectDockerCredentials: whether to inject docker credentials into the machine.
- exec.shutdown: command to execute when shutting down the machine after DevPod has determined the inactivityTimeout. Option values will be available here as well. For example, you can reuse an option that stores a cloud api key within this command to terminate the machine.
- binaries: this section can be used to declare additional binaries to download on the machine to use in exec.shutdown

### Auto-Inactivity Stop

One of the most important features of DevPod is to make sure that developer environments use as little resources as possible when they are not used.

#### Non-Machine Providers

For non-machine providers, DevPod can automatically kill the container its running in by terminating the process with pid 1. This is useful for providers such as docker, kubernetes or ssh, where you don't want the container to be running if its not needed. The timeout can be configured through agent.containerInactivityTimeout. DevPod will then start a process within the container to keep track of activity and then kill itself when the user hasn't connected for the given duration. This will not erase any state within the container and instead only stop it. Then when the user wants to start working with the workspace again, DevPod will start the container again.

#### Machine Providers

For machine providers, killing just the container within the remote machine is typically not enough as VMs still generate costs even if they are unused. Hence DevPod provides a way to configure automatically shutting down or deleting an unused machine on the cloud provider side if a developer is currently not working anymore. DevPod will then restart or recreate it again, when the development should continue.

DevPod tries to make this as easy as possible for you, as it will automatically keep track of when a user is connected to a workspace or not and only needs the command to run when the machine should be stopped from the provider. This command can be defined through agent.exec.shutdown. All configured options are available in this command and helper binaries needed can be defined through agent.exec.binaries

Official providers that use this method of automatically stopping an inactive machine are:

- devpod-provider-azure: Just uses shutdown -t now as agent.exec.shutdown to shutdown an unused machine.
- devpod-provider-aws: Uses the local aws cli tool to generate a temporary token, which is then saved in a DevPod option. This token is then used within agent.exec.shutdown to shutdown the machine on the agent side with an AWS api call.
- devpod-provider-gcloud: Uses the local gcloud cli tool to generate a temporary token, which is then saved in a DevPod option. This token is then used within agent.exec.shutdown to shutdown the machine on the agent side with an Google Cloud api call.
- devpod-provider-digitalocean: Deletes the whole machine on inactivity as stopped machines are still billed by DigitalOcean. The local digital ocean token is reused on the agent side to make an API call to delete the whole machine and preserve the state in an extra volume.

## Drivers

In DevPod you can specify a Driver in the Agent's configuration.

A Driver indicates how DevPod deploys the workspace container.

There are two types of drivers:

- Docker driver
- Kubernetes driver

### Docker Driver

The Docker driver is the default driver that DevPod uses to deploy the workspace container.

This container (specified through a devcontainer.json), is executed through Docker inside the provider environment, for example in a VM in case of Machine Providers.

Some optional configs are available:

- path: where to find the Docker CLI or a replacement, such as the Podman
- install: whether to install Docker or not in the target environment

Example config:

```yaml
agent:
  containerInactivityTimeout: 300
  docker:
    path: /usr/bin/docker
    install: false
```

### Kubernetes Driver

Instead of Docker, DevPod is also able to use Kubernetes as a Driver, which allows you to deploy the workspace to a Kubernetes cluster instead. For example, this makes it possible to create a provider that spins up a remote Kubernetes cluster (or just a namespace), connect to it, and create a workspace there. DevPod also has a default Kubernetes provider that uses the local Kubernetes config file to deploy the workspace.

DevPod also allows building an image through Kubernetes with buildkit. DevPod will automatically determine if building is necessary or if a prebuild can be used. However, the buildRepository option needs to be specified for this to work.

The allowed options for the Kubernetes driver are:

- path: where to find the kubectl binary or a drop-in replacement
- namespace: which namespace to use (if empty will use current namespace or default)
- context: which kube context to use (if empty will use current kube context)
- config: path to which kube config to use (if empty will use default kube config at ~/.kube/config)
- clusterRole: If defined, DevPod will create a role binding for the given cluster role for the workspace container. This is useful if you need Kubernetes access within the workspace container
- serviceAccount: If defined, DevPod will use the given service account for the dev container
- buildRepository: If defined, DevPod will build and push images to the given repository. If empty, DevPod will not build any images. Make sure you have push permissions for the given repository locally.
- helperImage: The image DevPod will use to find out the cluster architecture. Defaults to alpine.
- buildkitImage: The buildkit image to use for building dev containers.
- buildkitPrivileged: If the buildkit pod should run as a privileged pod
- persistentVolumeSize: The default size for the persistent volume to use.
- createNamespace: If true, DevPod will try to create the namespace

#### Example Kubernetes Provider

Example Kubernetes provider that uses local kubectl to run a workspace in the current kube context:

```yaml
name: simple-kubernetes
version: v0.0.1
agent:
  containerInactivityTimeout: 300 # Pod will automatically kill itself after timeout
  path: ${DEVPOD}
  driver: kubernetes
  kubernetes:
    # path: /usr/bin/kubectl
    # namespace: my-namespace-for-devpod
    # context: default
    # clusterRole: ""
    # serviceAccount: ""
    buildRepository: "ghcr.io/my-user/my-repo"
    # helperImage: "ubuntu:latest"
    # buildkitImage: "moby/buildkit"
    # buildkitPrivileged: false
    persistentVolumeSize: 20Gi
    createNamespace: true
exec:
  command: |-
    ${DEVPOD} helper sh -c "${COMMAND}"
```

Then add the provider via devpod provider add ./simple-kubernetes.yaml