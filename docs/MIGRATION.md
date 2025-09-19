# Migration Guide: Upgrading to New Server Plans

> 💰 **Save 36-89% on your development costs by migrating to UpCloud's new Developer Plans**

## Quick Migration

If you're using the default plan (`2xCPU-4GB`), migrate with:

```bash
# Update to new default (saves €10/month)
devpod provider set-options upcloud --option UPCLOUD_PLAN=DEV-2xCPU-4GB

# Recreate your workspace
devpod stop <workspace>
devpod delete <workspace>
devpod up <workspace> --provider upcloud
```

## What's Changed?

### New Plan Categories (2024)

UpCloud introduced two new plan categories optimized for development:

1. **Developer Plans** (September 2024)
   - 36-89% cheaper than General Purpose
   - Optimized for development workloads
   - Perfect for DevPod workspaces

2. **Cloud Native Plans** (December 2024)
   - Pay only when powered on
   - Ideal for ephemeral workspaces
   - Great with auto-shutdown

### Default Plan Change

| | Old Default | New Default | Savings |
|---|---|---|---|
| **Plan ID** | `2xCPU-4GB` | `DEV-2xCPU-4GB` | - |
| **Monthly Cost** | €28 | €18 | €10 (36%) |
| **Annual Cost** | €336 | €216 | €120 |

## Migration Paths

### From General Purpose to Developer Plans

Find your current plan and its recommended replacement:

| Current Plan (GP) | Specs | Cost | → | New Plan (DEV) | Cost | Savings |
|-------------------|-------|------|---|----------------|------|---------|
| `1xCPU-1GB` | 1 CPU, 1GB | €7/mo | → | `DEV-1xCPU-1GB` | €4.50/mo | 36% |
| `1xCPU-2GB` | 1 CPU, 2GB | €14/mo | → | `DEV-1xCPU-2GB` | €8/mo | 43% |
| `2xCPU-4GB` | 2 CPU, 4GB | €28/mo | → | `DEV-2xCPU-4GB` | €18/mo | 36% |
| `4xCPU-8GB` | 4 CPU, 8GB | €56/mo | → | `DEV-2xCPU-8GB` | €25/mo | 55% |
| `6xCPU-16GB` | 6 CPU, 16GB | €112/mo | → | `DEV-2xCPU-16GB` | €35/mo | 69% |
| `8xCPU-32GB` | 8 CPU, 32GB | €224/mo | → | `DEV-2xCPU-16GB`* | €35/mo | 84% |

*Note: DEV-2xCPU-16GB has less CPU but is sufficient for most development workloads

### To Cloud Native (Pay-per-use)

If you want to minimize costs with aggressive auto-shutdown:

| Use Case | Current | → | Cloud Native | When to Use |
|----------|---------|---|--------------|-------------|
| Occasional development | Any plan | → | `CN-2xCPU-4GB` | < 10 hrs/week usage |
| Microservices | `1xCPU-1GB` | → | `CN-1xCPU-1GB` | Container development |
| CI/CD agents | `2xCPU-4GB` | → | `CN-2xCPU-2GB` | Build automation |

## Step-by-Step Migration

### Step 1: Check Current Configuration

```bash
# View current provider options
devpod provider options upcloud

# Check specific workspace plan
devpod list
```

### Step 2: Choose New Plan

```bash
# List all available plans
devpod-provider-upcloud plans --recommended

# Compare specific categories
devpod-provider-upcloud plans --category developer --detailed
```

### Step 3: Update Provider Configuration

```bash
# Update default plan for all new workspaces
devpod provider set-options upcloud --option UPCLOUD_PLAN=DEV-2xCPU-8GB
```

### Step 4: Migrate Existing Workspaces

**Option A: Clean Migration (Recommended)**
```bash
# 1. Stop workspace
devpod stop my-workspace

# 2. Export any important data if needed
devpod ssh my-workspace -- tar czf /tmp/backup.tar.gz /important/data
devpod ssh my-workspace -- cat /tmp/backup.tar.gz > backup.tar.gz

# 3. Delete old workspace
devpod delete my-workspace

# 4. Create with new plan
devpod up my-workspace --provider upcloud

# 5. Restore data if needed
cat backup.tar.gz | devpod ssh my-workspace -- tar xzf - -C /
```

**Option B: Side-by-side Migration**
```bash
# 1. Create new workspace with different name
devpod up my-workspace-new --provider upcloud

# 2. Transfer data
devpod ssh my-workspace -- tar czf - /path/to/data | \
  devpod ssh my-workspace-new -- tar xzf - -C /

# 3. Test new workspace
devpod ssh my-workspace-new

# 4. Delete old workspace when ready
devpod delete my-workspace
```

## Cost Analysis

### Monthly Cost Comparison

Example for a small team (5 developers):

| Scenario | Old Plans | New Plans | Monthly Savings |
|----------|-----------|-----------|-----------------|
| **Basic Team** | 5 × €28 = €140 | 5 × €18 = €90 | €50 (36%) |
| **Mixed Team** | 3 × €28 + 2 × €56 = €196 | 3 × €18 + 2 × €25 = €104 | €92 (47%) |
| **Pro Team** | 5 × €56 = €280 | 5 × €25 = €125 | €155 (55%) |

### With Auto-Shutdown

Assuming 8 hours/day, 22 days/month usage:

| Plan | Full Price | Effective Cost (33% usage) | Annual Savings |
|------|------------|---------------------------|----------------|
| DEV-2xCPU-4GB | €18/mo | €18/mo (fixed) | - |
| CN-2xCPU-4GB | €16/mo | €5.28/mo | €152.64 |

## Backward Compatibility

### What Still Works

✅ **All old plan IDs remain valid**
- `1xCPU-1GB`, `2xCPU-4GB`, etc. continue to work
- No breaking changes to existing configurations
- Existing workspaces continue running

### What's Deprecated

⚠️ **Old defaults are not recommended**
- Still functional but more expensive
- Consider migration for cost savings

## Common Questions

### Do I have to migrate?

No, existing plans continue to work. Migration is recommended for cost savings.

### Will migration cause downtime?

Yes, workspaces must be recreated. Plan for 5-10 minutes per workspace.

### Can I migrate without data loss?

Yes, follow the backup/restore process in Step 4 above.

### Can I test new plans first?

Yes, create a test workspace:
```bash
devpod up test-workspace --provider upcloud --provider-option UPCLOUD_PLAN=DEV-2xCPU-4GB
```

### What if the new plan doesn't work for me?

You can always switch back:
```bash
devpod provider set-options upcloud --option UPCLOUD_PLAN=2xCPU-4GB
```

## Troubleshooting

### "Invalid plan" Error

```bash
# Check exact plan ID
devpod-provider-upcloud plans | grep -i "your-plan"

# Ensure using correct ID (case-sensitive)
devpod provider set-options upcloud --option UPCLOUD_PLAN=DEV-2xCPU-4GB
```

### Workspace Won't Start

Check plan restrictions:
- `DEV-1xCPU-1GB-10GB`: Max 2 per account
- `DEV-1xCPU-1GB`: Max 5 per account

### Performance Issues

If new plan has less resources:
```bash
# Upgrade to next tier
devpod provider set-options upcloud --option UPCLOUD_PLAN=DEV-2xCPU-8GB
```

## Rollback Procedure

If you need to revert to old plans:

```bash
# 1. Update provider
devpod provider set-options upcloud --option UPCLOUD_PLAN=2xCPU-4GB

# 2. Recreate workspace
devpod delete <workspace>
devpod up <workspace> --provider upcloud
```

## Support

For migration assistance:
- Check provider status: `devpod provider list`
- View logs: `devpod provider logs upcloud`
- Test configuration: `devpod-provider-upcloud init`

---

*Migration guide version 1.0 - December 2024*