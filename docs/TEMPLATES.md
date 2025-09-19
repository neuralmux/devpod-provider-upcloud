# Server Plan Templating System - Technical Documentation

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [System Components](#system-components)
- [Configuration File Structure](#configuration-file-structure)
- [Implementation Details](#implementation-details)
- [Maintenance Procedures](#maintenance-procedures)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## Architecture Overview

The server plan templating system provides a flexible, maintainable way to manage UpCloud server plans without hardcoding them in the provider source code.

### Design Principles

1. **Separation of Concerns**: Plan definitions separate from business logic
2. **Easy Updates**: YAML configuration instead of code changes
3. **Backward Compatibility**: Fallback to legacy mappings
4. **Type Safety**: Go structs with validation
5. **Embedded Configuration**: Single binary distribution

### System Flow

```
┌─────────────────┐
│  User Request   │
│  UPCLOUD_PLAN=  │
│  DEV-2xCPU-4GB  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ provider.yaml   │
│   Validates     │
│  against list   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ MapPlanName()   │
│ pkg/upcloud/    │
│   mapper.go     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│LoadServerPlans()│
│  pkg/config/    │
│   plans.go      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Embedded YAML   │
│server-plans.yaml│
│   (at build)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Validation &   │
│    Return       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  UpCloud API    │
│  Create Server  │
└─────────────────┘
```

## System Components

### 1. Configuration File (`configs/server-plans.yaml`)

The master configuration file containing all server plan definitions.

**Location**: `configs/server-plans.yaml`
**Embedded at**: Build time via `go:embed`
**Format**: YAML

### 2. Configuration Loader (`pkg/config/plans.go`)

Go package that:
- Embeds the YAML configuration
- Provides type-safe structs
- Implements validation logic
- Offers query methods

**Key Functions**:
```go
LoadServerPlans() (*ServerPlans, error)
GetPlanByID(id string) (*ServerPlan, *PlanCategory, error)
ValidatePlan(planID string) error
GetRecommendedPlans() []ServerPlan
```

### 3. Plan Mapper (`pkg/upcloud/mapper.go`)

Integration point that:
- Loads configuration
- Validates user input
- Falls back to legacy mappings
- Returns validated plan IDs

### 4. CLI Command (`cmd/plans.go`)

User interface that:
- Lists available plans
- Filters by criteria
- Formats output (table/json/yaml)
- Shows recommendations

### 5. Legacy Constants (`pkg/upcloud/constants.go`)

Backward compatibility layer containing hardcoded plan mappings for fallback.

## Configuration File Structure

### Root Structure

```yaml
version: "2024.12"                  # Configuration version
last_updated: "2024-12-18"          # Last modification date
default_plan: "DEV-2xCPU-4GB"       # Default plan ID

categories:                         # Plan categories
  developer:                        # Category key
    name: "Developer Plans"         # Display name
    description: "..."              # Description
    recommended_for_devpod: true    # DevPod recommendation flag
    billed_when_on_only: false      # Billing model flag
    icon: "🚀"                      # UI icon
    plans: [...]                    # Array of plans

selection_rules:                    # Plan selection logic
  minimum:                          # Minimum requirements
    cpu: 1
    ram: 1024
    storage: 10
  recommendations:                  # Recommendation mappings
    default: "DEV-2xCPU-4GB"
    by_language: {...}
    by_framework: {...}
    by_workload: {...}

metadata:                           # Provider metadata
  provider: "UpCloud"
  api_version: "v8"
  currency: "EUR"
  billing_unit: "hourly"
  regions_available: [...]
```

### Plan Definition

```yaml
plans:
  - id: "DEV-2xCPU-4GB"             # Unique identifier
    display_name: "Standard Dev"     # Human-readable name
    description: "..."               # Detailed description
    cpu: 2                           # CPU cores
    ram: 4096                        # RAM in MB
    storage: 60                      # Storage in GB
    price_monthly: 18.00             # Monthly price
    price_hourly: 0.025              # Hourly price
    use_cases:                       # Use case list
      - "Web development"
      - "API development"
    default: true                    # Is default plan
    recommended: true                # Is recommended
    restrictions:                    # Optional restrictions
      max_per_account: 5
```

## Implementation Details

### Embedding Configuration

The YAML file is embedded at compile time:

```go
//go:embed server-plans.yaml
var serverPlansYAML []byte
```

**Important**: The embedded file must be in the same package directory due to Go's embed restrictions. We copy it during build:
```bash
cp configs/server-plans.yaml pkg/config/
```

### Loading and Caching

Configuration is loaded on demand and parsed each time:

```go
func LoadServerPlans() (*ServerPlans, error) {
    var plans ServerPlans
    if err := yaml.Unmarshal(serverPlansYAML, &plans); err != nil {
        return nil, fmt.Errorf("failed to parse server plans: %w", err)
    }
    return &plans, nil
}
```

**Note**: No caching is implemented as the performance impact is minimal (~0.1ms per load).

### Validation Flow

1. **User provides plan ID** → `UPCLOUD_PLAN=DEV-2xCPU-4GB`
2. **MapPlanName() called** → Loads configuration
3. **ValidatePlan() executed** → Checks if plan exists
4. **Fallback to legacy** → If config fails, check hardcoded map
5. **Return validated ID** → Or error if invalid

### Error Handling

The system implements graceful degradation:

```go
func MapPlanName(plan string) (string, error) {
    plans, err := config.LoadServerPlans()
    if err != nil {
        // Fallback to legacy mapping
        if mappedPlan, ok := PlanMap[plan]; ok {
            return mappedPlan, nil
        }
        return "", fmt.Errorf("failed to load plans config: %w", err)
    }
    // ... validation logic
}
```

## Maintenance Procedures

### Adding New Plans

1. **Edit Configuration File**:
```yaml
# In configs/server-plans.yaml
plans:
  - id: "NEW-PLAN-ID"
    display_name: "New Plan"
    cpu: 4
    ram: 8192
    storage: 100
    price_monthly: 50.00
    price_hourly: 0.07
    use_cases: ["New use case"]
```

2. **Update Legacy Map** (optional, for backward compatibility):
```go
// In pkg/upcloud/constants.go
var PlanMap = map[string]string{
    // ...
    "NEW-PLAN-ID": "NEW-PLAN-ID",
}
```

3. **Copy Configuration for Embedding**:
```bash
cp configs/server-plans.yaml pkg/config/
```

4. **Build and Test**:
```bash
make build
make test
./bin/devpod-provider-upcloud plans | grep NEW-PLAN
```

### Updating Existing Plans

1. **Modify in YAML**:
   - Update prices
   - Change descriptions
   - Modify use cases
   - Adjust recommendations

2. **Version the Configuration**:
```yaml
version: "2024.13"  # Increment version
last_updated: "2024-12-20"  # Update date
```

3. **Test Changes**:
```bash
go test ./pkg/config/...
```

### Deprecating Plans

1. **Remove from Suggestions** in `provider.yaml`
2. **Keep in Configuration** for backward compatibility
3. **Add deprecation notice**:
```yaml
plans:
  - id: "OLD-PLAN"
    display_name: "Old Plan (DEPRECATED)"
    description: "⚠️ Deprecated - use NEW-PLAN instead"
```

### Changing Default Plan

1. **Update Configuration**:
```yaml
default_plan: "NEW-DEFAULT-ID"
```

2. **Update provider.yaml**:
```yaml
UPCLOUD_PLAN:
  default: NEW-DEFAULT-ID
```

3. **Update Documentation**

### Adding New Categories

1. **Define Category**:
```yaml
categories:
  new_category:
    name: "New Category"
    description: "Category description"
    icon: "🆕"
    plans: [...]
```

2. **Update CLI Help**:
```go
// In cmd/plans.go
plansCmd.Flags().StringVarP(&cmd.Category, "category", "c", "",
    "Filter by category (developer, cloud_native, general_purpose, high_cpu, high_memory, new_category)")
```

## Testing

### Unit Tests

Located in `pkg/config/plans_test.go`:

```bash
go test -v ./pkg/config/...
```

**Test Coverage**:
- Configuration loading
- Plan validation
- Recommendation logic
- Region validation
- Price comparisons
- Output formatting

### Integration Tests

Test with actual provider:

```bash
# Test plan listing
./bin/devpod-provider-upcloud plans

# Test plan validation
export UPCLOUD_PLAN=DEV-2xCPU-4GB
./bin/devpod-provider-upcloud init

# Test invalid plan
export UPCLOUD_PLAN=INVALID
./bin/devpod-provider-upcloud init  # Should error
```

### Adding Tests

When adding new plans or features:

```go
func TestNewPlanCategory(t *testing.T) {
    plans, err := LoadServerPlans()
    require.NoError(t, err)

    category, exists := plans.Categories["new_category"]
    assert.True(t, exists)
    assert.NotEmpty(t, category.Plans)
}
```

## Troubleshooting

### Common Issues

#### 1. "pattern server-plans.yaml: no matching files found"

**Cause**: Build-time embed failure
**Solution**: Ensure `server-plans.yaml` is copied to `pkg/config/`

```bash
cp configs/server-plans.yaml pkg/config/
make build
```

#### 2. "invalid plan: XYZ"

**Cause**: Plan not in configuration
**Debug**:
```bash
./bin/devpod-provider-upcloud plans | grep XYZ
```

#### 3. Configuration not updating

**Cause**: Binary not rebuilt after YAML changes
**Solution**:
```bash
cp configs/server-plans.yaml pkg/config/
make clean
make build
```

#### 4. Fallback to legacy plans

**Cause**: Configuration load failure
**Debug**:
```go
// Add logging in MapPlanName()
if err != nil {
    log.Printf("Config load failed: %v", err)
    // Fallback logic
}
```

### Debug Mode

Enable detailed logging:

```bash
export DEVPOD_DEBUG=true
./bin/devpod-provider-upcloud plans
```

### Validation Script

Create a validation script:

```bash
#!/bin/bash
# validate-plans.sh

echo "Validating server plans configuration..."

# Build
make build || exit 1

# Test loading
./bin/devpod-provider-upcloud plans > /dev/null || exit 1

# Test each category
for category in developer cloud_native general_purpose; do
    echo "Testing category: $category"
    ./bin/devpod-provider-upcloud plans --category $category > /dev/null || exit 1
done

# Test JSON output
./bin/devpod-provider-upcloud plans --format json | jq . > /dev/null || exit 1

echo "✅ All validations passed"
```

## Best Practices

1. **Always version changes** - Update version and last_updated fields
2. **Maintain backward compatibility** - Keep old plan IDs working
3. **Test before release** - Run full test suite after changes
4. **Document changes** - Update CHANGELOG.md
5. **Copy before build** - Ensure embedded file is current
6. **Validate pricing** - Double-check prices against UpCloud
7. **Update documentation** - Keep user docs in sync

## Future Enhancements

### Planned Improvements

1. **Dynamic Loading**: Fetch plans from UpCloud API
2. **Caching**: Cache parsed configuration
3. **Hot Reload**: Reload configuration without rebuild
4. **Plan Aliases**: Support multiple names per plan
5. **Custom Plans**: User-defined plan configurations
6. **Cost Calculator**: Estimate costs based on usage

### API Integration

Future version could fetch plans dynamically:

```go
func FetchPlansFromAPI(client *upcloud.Client) (*ServerPlans, error) {
    // Fetch current plans from UpCloud API
    // Transform to our format
    // Cache locally
}
```

### Configuration Validation

Add JSON schema validation:

```yaml
# schema/server-plans.schema.json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["version", "categories"],
  ...
}
```

---

*Technical documentation version 1.0 - December 2024*