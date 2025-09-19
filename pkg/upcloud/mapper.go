package upcloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/neuralmux/devpod-provider-upcloud/pkg/config"
)

// MapPlanName maps the provider plan name to UpCloud plan name
func MapPlanName(plan string) (string, error) {
	// Load server plans configuration
	plans, err := config.LoadServerPlans()
	if err != nil {
		// Fallback to legacy plan mapping if config load fails
		if mappedPlan, ok := PlanMap[plan]; ok {
			return mappedPlan, nil
		}
		return "", fmt.Errorf("failed to load plans config: %w", err)
	}

	// Validate the plan exists in configuration
	if err := plans.ValidatePlan(plan); err != nil {
		// Check legacy map as fallback for backward compatibility
		if mappedPlan, ok := PlanMap[plan]; ok {
			return mappedPlan, nil
		}
		return "", fmt.Errorf("invalid plan: %s (use 'plans' command to list available plans)", plan)
	}

	// Plan is valid, return it
	return plan, nil
}

// MapImageToTemplate maps the OS image name to UpCloud template UUID
func MapImageToTemplate(imageName string) (string, error) {
	if templateUUID, ok := ImageMap[imageName]; ok {
		return templateUUID, nil
	}
	// If not found in map, check if it's already a UUID
	if strings.HasPrefix(imageName, "01000000-") {
		return imageName, nil
	}
	return "", fmt.Errorf("unknown image: %s", imageName)
}

// MapServerStateToStatus maps UpCloud server state to DevPod status
func MapServerStateToStatus(state string) string {
	switch state {
	case upcloud.ServerStateStarted:
		return StatusRunning
	case upcloud.ServerStateStopped:
		return StatusStopped
	case upcloud.ServerStateMaintenance, upcloud.ServerStateError:
		return StatusBusy
	default:
		return StatusBusy
	}
}

// ParseStorageSize parses the storage size string to integer GB
func ParseStorageSize(storageStr string) (int, error) {
	size, err := strconv.Atoi(storageStr)
	if err != nil {
		return 0, fmt.Errorf("invalid storage size: %s", storageStr)
	}
	if size < 10 || size > 2048 {
		return 0, fmt.Errorf("storage size must be between 10 and 2048 GB")
	}
	return size, nil
}

// GetStorageTier returns the appropriate storage tier
func GetStorageTier() string {
	// For DevPod workspaces, we want good performance
	return upcloud.StorageTierMaxIOPS
}

// ValidateZone checks if the zone is valid
func ValidateZone(zone string) error {
	// Try to load configuration for zone validation
	plans, err := config.LoadServerPlans()
	if err != nil {
		// Fallback to hardcoded zones if config fails
		validZones := []string{
			"de-fra1", "fi-hel1", "fi-hel2", "nl-ams1", "uk-lon1",
			"us-nyc1", "us-chi1", "us-sjo1", "sg-sin1", "au-syd1",
			"es-mad1", "pl-waw1", "se-sto1",
		}

		for _, validZone := range validZones {
			if zone == validZone {
				return nil
			}
		}
		return fmt.Errorf("invalid zone: %s", zone)
	}

	// Use configuration-based validation
	if !plans.IsValidRegion(zone) {
		return fmt.Errorf("invalid zone: %s (available: %s)", zone, strings.Join(plans.GetRegions(), ", "))
	}
	return nil
}

// GenerateHostname creates a hostname from machine ID
func GenerateHostname(machineID string) string {
	// Remove "devpod-" prefix if it exists
	hostname := strings.TrimPrefix(machineID, "devpod-")
	// Ensure hostname is valid (lowercase, alphanumeric, hyphens)
	hostname = strings.ToLower(hostname)
	// Truncate if too long (max 63 chars for hostname)
	if len(hostname) > 63 {
		hostname = hostname[:63]
	}
	return hostname
}

// FindServerByMachineID searches for a server by the DevPod machine ID
func FindServerByMachineID(servers []upcloud.Server, machineID string) *upcloud.Server {
	// Try exact match first
	for _, server := range servers {
		if server.Hostname == machineID || server.Title == machineID {
			return &server
		}
	}

	// Try without "devpod-" prefix
	hostname := GenerateHostname(machineID)
	for _, server := range servers {
		if server.Hostname == hostname {
			return &server
		}
	}

	return nil
}

// GetPublicIPv4 extracts the public IPv4 address from server details
func GetPublicIPv4(server *upcloud.ServerDetails) (string, error) {
	for _, iface := range server.Networking.Interfaces {
		if iface.Type == "public" {
			for _, ip := range iface.IPAddresses {
				if ip.Family == upcloud.IPAddressFamilyIPv4 {
					return ip.Address, nil
				}
			}
		}
	}
	return "", fmt.Errorf("no public IPv4 address found")
}
