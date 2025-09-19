package config

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed server-plans.yaml
var serverPlansYAML []byte

// ServerPlans represents the complete server plans configuration
type ServerPlans struct {
	Version        string                   `yaml:"version"`
	LastUpdated    string                   `yaml:"last_updated"`
	DefaultPlan    string                   `yaml:"default_plan"`
	Categories     map[string]*PlanCategory `yaml:"categories"`
	SelectionRules SelectionRules           `yaml:"selection_rules"`
	Metadata       PlanMetadata             `yaml:"metadata"`
}

// PlanCategory represents a category of server plans
type PlanCategory struct {
	Name                 string       `yaml:"name"`
	Description          string       `yaml:"description"`
	RecommendedForDevpod bool         `yaml:"recommended_for_devpod"`
	BilledWhenOnOnly     bool         `yaml:"billed_when_on_only"`
	Icon                 string       `yaml:"icon"`
	Plans                []ServerPlan `yaml:"plans"`
}

// ServerPlan represents an individual server plan
type ServerPlan struct {
	ID           string            `yaml:"id"`
	DisplayName  string            `yaml:"display_name"`
	Description  string            `yaml:"description"`
	CPU          int               `yaml:"cpu"`
	RAM          int               `yaml:"ram"`
	Storage      int               `yaml:"storage"`
	PriceMonthly float32           `yaml:"price_monthly"`
	PriceHourly  float32           `yaml:"price_hourly"`
	UseCases     []string          `yaml:"use_cases"`
	Default      bool              `yaml:"default"`
	Recommended  bool              `yaml:"recommended"`
	Restrictions *PlanRestrictions `yaml:"restrictions,omitempty"`
}

// PlanRestrictions represents any restrictions on a plan
type PlanRestrictions struct {
	MaxPerAccount int `yaml:"max_per_account"`
}

// SelectionRules defines rules for plan selection
type SelectionRules struct {
	Minimum         MinimumRequirements `yaml:"minimum"`
	Recommendations PlanRecommendations `yaml:"recommendations"`
}

// MinimumRequirements defines minimum specs for DevPod
type MinimumRequirements struct {
	CPU     int `yaml:"cpu"`
	RAM     int `yaml:"ram"`
	Storage int `yaml:"storage"`
}

// PlanRecommendations contains recommended plans for different scenarios
type PlanRecommendations struct {
	Default     string            `yaml:"default"`
	ByLanguage  map[string]string `yaml:"by_language"`
	ByFramework map[string]string `yaml:"by_framework"`
	ByWorkload  map[string]string `yaml:"by_workload"`
}

// PlanMetadata contains metadata about the plans
type PlanMetadata struct {
	Provider         string   `yaml:"provider"`
	APIVersion       string   `yaml:"api_version"`
	Currency         string   `yaml:"currency"`
	BillingUnit      string   `yaml:"billing_unit"`
	RegionsAvailable []string `yaml:"regions_available"`
}

// LoadServerPlans loads the server plans from the embedded YAML file
func LoadServerPlans() (*ServerPlans, error) {
	var plans ServerPlans
	if err := yaml.Unmarshal(serverPlansYAML, &plans); err != nil {
		return nil, fmt.Errorf("failed to parse server plans: %w", err)
	}
	return &plans, nil
}

// GetPlanByID finds a plan by its ID across all categories
func (s *ServerPlans) GetPlanByID(id string) (*ServerPlan, *PlanCategory, error) {
	for _, category := range s.Categories {
		for i := range category.Plans {
			if category.Plans[i].ID == id {
				return &category.Plans[i], category, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("plan %s not found", id)
}

// GetRecommendedPlans returns all plans recommended for DevPod
func (s *ServerPlans) GetRecommendedPlans() []ServerPlan {
	var plans []ServerPlan
	for _, category := range s.Categories {
		if category.RecommendedForDevpod {
			plans = append(plans, category.Plans...)
		} else {
			// Also include individually recommended plans from non-recommended categories
			for _, plan := range category.Plans {
				if plan.Recommended {
					plans = append(plans, plan)
				}
			}
		}
	}
	return plans
}

// GetDeveloperPlans returns only the developer category plans
func (s *ServerPlans) GetDeveloperPlans() []ServerPlan {
	if category, exists := s.Categories["developer"]; exists {
		return category.Plans
	}
	return []ServerPlan{}
}

// GetDefaultPlan returns the default plan
func (s *ServerPlans) GetDefaultPlan() (*ServerPlan, error) {
	plan, _, err := s.GetPlanByID(s.DefaultPlan)
	return plan, err
}

// ValidatePlan checks if a plan ID is valid
func (s *ServerPlans) ValidatePlan(planID string) error {
	_, _, err := s.GetPlanByID(planID)
	return err
}

// GetPlanSuggestions returns a list of plan IDs for provider.yaml suggestions
func (s *ServerPlans) GetPlanSuggestions() []string {
	var suggestions []string

	// First add developer plans (recommended)
	if devCategory, exists := s.Categories["developer"]; exists {
		for _, plan := range devCategory.Plans {
			suggestions = append(suggestions, plan.ID)
		}
	}

	// Then add cloud native plans
	if cnCategory, exists := s.Categories["cloud_native"]; exists {
		for _, plan := range cnCategory.Plans {
			suggestions = append(suggestions, plan.ID)
		}
	}

	// Finally add some general purpose plans
	if gpCategory, exists := s.Categories["general_purpose"]; exists {
		// Only add smaller general purpose plans
		for _, plan := range gpCategory.Plans {
			if plan.CPU <= 4 && plan.RAM <= 8192 {
				suggestions = append(suggestions, plan.ID)
			}
		}
	}

	return suggestions
}

// FormatPlanList formats plans for display in CLI
func (s *ServerPlans) FormatPlanList(detailed bool) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("UpCloud Server Plans (v%s)\n", s.Version))
	output.WriteString(fmt.Sprintf("Last Updated: %s\n", s.LastUpdated))
	output.WriteString(fmt.Sprintf("Default Plan: %s\n\n", s.DefaultPlan))

	// Sort categories for consistent display
	var categoryKeys []string
	for key := range s.Categories {
		categoryKeys = append(categoryKeys, key)
	}
	sort.Strings(categoryKeys)

	// Display categories in order
	for _, key := range categoryKeys {
		category := s.Categories[key]

		output.WriteString(fmt.Sprintf("%s %s\n", category.Icon, category.Name))
		if category.RecommendedForDevpod {
			output.WriteString("  ⭐ Recommended for DevPod\n")
		}
		if category.BilledWhenOnOnly {
			output.WriteString("  💰 Billed only when powered on\n")
		}
		output.WriteString(fmt.Sprintf("  %s\n\n", category.Description))

		// Display plans
		for _, plan := range category.Plans {
			if plan.Recommended || plan.Default {
				output.WriteString("  ★ ")
			} else {
				output.WriteString("     ")
			}

			output.WriteString(fmt.Sprintf("%-20s", plan.ID))
			output.WriteString(fmt.Sprintf(" %d CPU, %d GB RAM", plan.CPU, plan.RAM/1024))

			if plan.Storage > 0 {
				output.WriteString(fmt.Sprintf(", %d GB Storage", plan.Storage))
			}

			output.WriteString(fmt.Sprintf(" - €%.2f/month", plan.PriceMonthly))

			if plan.Default {
				output.WriteString(" [DEFAULT]")
			}
			if plan.Recommended {
				output.WriteString(" [RECOMMENDED]")
			}

			output.WriteString("\n")

			if detailed {
				output.WriteString(fmt.Sprintf("       %s\n", plan.Description))
				if len(plan.UseCases) > 0 {
					output.WriteString(fmt.Sprintf("       Use cases: %s\n", strings.Join(plan.UseCases, ", ")))
				}
				if plan.Restrictions != nil && plan.Restrictions.MaxPerAccount > 0 {
					output.WriteString(fmt.Sprintf("       ⚠️  Max %d per account\n", plan.Restrictions.MaxPerAccount))
				}
				output.WriteString("\n")
			}
		}

		output.WriteString("\n")
	}

	return output.String()
}

// GetPlanRecommendation returns a recommended plan based on criteria
func (s *ServerPlans) GetPlanRecommendation(language, framework, workload string) string {
	// Check workload first (highest priority)
	if workload != "" {
		if plan, exists := s.SelectionRules.Recommendations.ByWorkload[strings.ToLower(workload)]; exists {
			return plan
		}
	}

	// Check framework
	if framework != "" {
		if plan, exists := s.SelectionRules.Recommendations.ByFramework[strings.ToLower(framework)]; exists {
			return plan
		}
	}

	// Check language
	if language != "" {
		if plan, exists := s.SelectionRules.Recommendations.ByLanguage[strings.ToLower(language)]; exists {
			return plan
		}
	}

	// Return default
	return s.SelectionRules.Recommendations.Default
}

// IsValidRegion checks if a region is valid
func (s *ServerPlans) IsValidRegion(region string) bool {
	for _, r := range s.Metadata.RegionsAvailable {
		if r == region {
			return true
		}
	}
	return false
}

// GetRegions returns all available regions
func (s *ServerPlans) GetRegions() []string {
	return s.Metadata.RegionsAvailable
}
