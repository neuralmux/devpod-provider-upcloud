# Testing Guide for UpCloud DevPod Provider

## Quick Test (No UpCloud Account Required)

### 1. Build and Basic Test
```bash
# Build the provider
make build

# Run local test (no API calls)
./test-local.sh
```

This will test that the provider binary works correctly without making any API calls.

### 2. Test Mode
Use special test credentials to bypass API calls:

```bash
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="2xCPU-4GB"
export UPCLOUD_STORAGE="50"
export UPCLOUD_IMAGE="Ubuntu Server 22.04 LTS (Jammy Jellyfish)"

# This will succeed in test mode
./bin/devpod-provider-upcloud init
```

## Testing with Real UpCloud Account

### Prerequisites
1. UpCloud account with API access enabled
2. API credentials from UpCloud Control Panel

### Step 1: Set Real Credentials
```bash
export UPCLOUD_USERNAME="your-real-username"
export UPCLOUD_PASSWORD="your-real-password"
```

### Step 2: Test Authentication
```bash
./bin/devpod-provider-upcloud init
```

If successful, you'll see: "Successfully initialized UpCloud provider"

### Step 3: Test Server Operations (Costs Money!)

⚠️ **WARNING**: These commands will create real servers and incur charges!

```bash
# Setup for server creation
mkdir -p /tmp/devpod-test/.ssh
ssh-keygen -t ed25519 -f /tmp/devpod-test/.ssh/id_ed25519 -N ''

export MACHINE_FOLDER="/tmp/devpod-test"
export MACHINE_ID="devpod-test-$(date +%s)"
export UPCLOUD_ZONE="de-fra1"
export UPCLOUD_PLAN="1xCPU-1GB"  # Smallest plan to minimize cost
export UPCLOUD_STORAGE="10"       # Minimum storage

# Create server (takes ~1-2 minutes)
./bin/devpod-provider-upcloud create

# Check status
./bin/devpod-provider-upcloud status

# Stop server
./bin/devpod-provider-upcloud stop

# Start server
./bin/devpod-provider-upcloud start

# Delete server (important to avoid charges!)
./bin/devpod-provider-upcloud delete
```

## Testing with DevPod

### Option 1: Local Provider
```bash
# Add local provider
devpod provider add . --name upcloud-local

# Configure
devpod provider set-options upcloud-local \
  --option UPCLOUD_USERNAME="$UPCLOUD_USERNAME" \
  --option UPCLOUD_PASSWORD="$UPCLOUD_PASSWORD"

# Use provider
devpod provider use upcloud-local

# Create workspace
devpod up github.com/loft-sh/devpod-quickstart
```

### Option 2: From GitHub
```bash
# Add from GitHub
devpod provider add github.com/neuralmux/devpod-provider-upcloud

# Configure and use as above
```

## Troubleshooting

### Common Issues

1. **"missing credentials" error**
   - Ensure UPCLOUD_USERNAME and UPCLOUD_PASSWORD are set
   - Check for typos in environment variables

2. **"authentication failed" error**
   - Verify credentials are correct
   - Ensure API access is enabled in UpCloud Control Panel
   - Check if account has sufficient permissions

3. **"invalid zone" error**
   - Use one of the valid zones listed in README.md
   - Example: de-fra1, us-nyc1, sg-sin1

4. **"invalid plan" error**
   - Use exact plan names: 1xCPU-1GB, 2xCPU-4GB, etc.

5. **SSH key issues**
   - Ensure MACHINE_FOLDER contains .ssh/id_ed25519 or .ssh/id_rsa
   - Keys must have proper permissions (600)

### Debug Mode
To see more detailed output:
```bash
# Enable debug logging
export DEVPOD_DEBUG=true
export DEVPOD_LOG_LEVEL=debug

# Run commands
./bin/devpod-provider-upcloud init
```

### Check Server in UpCloud Console
1. Log into https://hub.upcloud.com/
2. Go to Servers section
3. Look for servers with names starting with "devpod-"

## Cost Management

### Minimize Costs
- Use smallest plan: `1xCPU-1GB`
- Use minimum storage: `10` GB
- Delete servers immediately after testing
- Use auto-stop feature (INACTIVITY_TIMEOUT)

### Estimated Costs
- 1xCPU-1GB: ~$0.007/hour (~$5/month)
- Storage: ~$0.00003/GB/hour
- Test session (1 hour): ~$0.01

## Automated Testing

### Run BDD Tests
```bash
# With test credentials
export UPCLOUD_USERNAME="test"
export UPCLOUD_PASSWORD="test"
make bdd
```

### Run All Tests
```bash
make test-all
```

## Need Help?

1. Check the provider logs
2. Verify all environment variables are set correctly
3. Ensure your UpCloud account has API access enabled
4. Check UpCloud service status at https://status.upcloud.com/

Remember to **always delete test servers** to avoid unnecessary charges!