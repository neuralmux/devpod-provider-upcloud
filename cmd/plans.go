package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/neuralmux/devpod-provider-upcloud/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// PlansCmd holds the plans command flags
type PlansCmd struct {
	Detailed    bool
	Recommended bool
	Category    string
	Format      string
}

// NewPlansCmd defines the plans command
func NewPlansCmd() *cobra.Command {
	cmd := &PlansCmd{}

	plansCmd := &cobra.Command{
		Use:   "plans",
		Short: "List available server plans",
		Long: `List all available UpCloud server plans for DevPod workspaces.

This command shows server plans organized by category, with pricing and recommendations.
Developer plans are recommended for most DevPod use cases as they offer the best value.`,
		Example: `  # List all plans
  devpod-provider-upcloud plans

  # Show detailed information
  devpod-provider-upcloud plans --detailed

  # Show only recommended plans
  devpod-provider-upcloud plans --recommended

  # Show plans from a specific category
  devpod-provider-upcloud plans --category developer

  # Output as JSON
  devpod-provider-upcloud plans --format json`,
		RunE: func(_ *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}

	plansCmd.Flags().BoolVarP(&cmd.Detailed, "detailed", "d", false, "Show detailed plan information")
	plansCmd.Flags().BoolVarP(&cmd.Recommended, "recommended", "r", false, "Show only recommended plans")
	plansCmd.Flags().StringVarP(&cmd.Category, "category", "c", "", "Filter by category (developer, cloud_native, general_purpose, high_cpu, high_memory)")
	plansCmd.Flags().StringVarP(&cmd.Format, "format", "f", "table", "Output format (table, json, yaml)")

	return plansCmd
}

// Run executes the plans command
func (cmd *PlansCmd) Run() error {
	// Load server plans
	plans, err := config.LoadServerPlans()
	if err != nil {
		return fmt.Errorf("failed to load server plans: %w", err)
	}

	// Handle different output formats
	switch cmd.Format {
	case "json":
		return cmd.outputJSON(plans)
	case "yaml":
		return cmd.outputYAML(plans)
	default:
		return cmd.outputTable(plans)
	}
}

// outputTable outputs plans in table format
func (cmd *PlansCmd) outputTable(plans *config.ServerPlans) error {
	if cmd.Recommended {
		return cmd.outputRecommendedPlans(plans)
	}

	if cmd.Category != "" {
		return cmd.outputCategoryPlans(plans, cmd.Category)
	}

	// Output all plans
	fmt.Print(plans.FormatPlanList(cmd.Detailed))
	return nil
}

// outputRecommendedPlans shows only recommended plans
func (cmd *PlansCmd) outputRecommendedPlans(plans *config.ServerPlans) error {
	fmt.Println("🌟 Recommended Plans for DevPod")
	fmt.Println("================================")
	fmt.Println()

	recommendedPlans := plans.GetRecommendedPlans()
	defaultPlan, _ := plans.GetDefaultPlan()

	for _, plan := range recommendedPlans {
		if defaultPlan != nil && plan.ID == defaultPlan.ID {
			fmt.Printf("★ ")
		} else {
			fmt.Printf("  ")
		}

		fmt.Printf("%-20s", plan.ID)
		fmt.Printf(" %d CPU, %d GB RAM", plan.CPU, plan.RAM/1024)

		if plan.Storage > 0 {
			fmt.Printf(", %d GB Storage", plan.Storage)
		}

		fmt.Printf(" - €%.2f/month", plan.PriceMonthly)

		if defaultPlan != nil && plan.ID == defaultPlan.ID {
			fmt.Printf(" [DEFAULT]")
		}

		fmt.Println()

		if cmd.Detailed {
			fmt.Printf("    %s\n", plan.Description)
			if len(plan.UseCases) > 0 {
				fmt.Printf("    Use cases: %s\n", strings.Join(plan.UseCases, ", "))
			}
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println("💡 Tip: Developer plans offer the best value for development workspaces")
	fmt.Println("   Use --category developer to see all developer plans")

	return nil
}

// outputCategoryPlans shows plans from a specific category
func (cmd *PlansCmd) outputCategoryPlans(plans *config.ServerPlans, categoryName string) error {
	category, exists := plans.Categories[categoryName]
	if !exists {
		return fmt.Errorf("category '%s' not found. Available categories: developer, cloud_native, general_purpose, high_cpu, high_memory", categoryName)
	}

	fmt.Printf("%s %s\n", category.Icon, category.Name)
	fmt.Printf("=====================================\n")
	if category.RecommendedForDevpod {
		fmt.Println("⭐ Recommended for DevPod")
	}
	if category.BilledWhenOnOnly {
		fmt.Println("💰 Billed only when powered on")
	}
	fmt.Printf("%s\n\n", category.Description)

	defaultPlan, _ := plans.GetDefaultPlan()

	for _, plan := range category.Plans {
		if plan.Recommended || (defaultPlan != nil && plan.ID == defaultPlan.ID) {
			fmt.Printf("★ ")
		} else {
			fmt.Printf("  ")
		}

		fmt.Printf("%-20s", plan.ID)
		fmt.Printf(" %d CPU, %d GB RAM", plan.CPU, plan.RAM/1024)

		if plan.Storage > 0 {
			fmt.Printf(", %d GB Storage", plan.Storage)
		}

		fmt.Printf(" - €%.2f/month", plan.PriceMonthly)

		if defaultPlan != nil && plan.ID == defaultPlan.ID {
			fmt.Printf(" [DEFAULT]")
		}
		if plan.Recommended {
			fmt.Printf(" [RECOMMENDED]")
		}

		fmt.Println()

		if cmd.Detailed {
			fmt.Printf("    %s\n", plan.Description)
			if len(plan.UseCases) > 0 {
				fmt.Printf("    Use cases: %s\n", strings.Join(plan.UseCases, ", "))
			}
			if plan.Restrictions != nil && plan.Restrictions.MaxPerAccount > 0 {
				fmt.Printf("    ⚠️  Max %d per account\n", plan.Restrictions.MaxPerAccount)
			}
			fmt.Println()
		}
	}

	return nil
}

// outputJSON outputs plans as JSON
func (cmd *PlansCmd) outputJSON(plans *config.ServerPlans) error {
	// For JSON output, we'll use a simpler structure
	type SimplePlan struct {
		ID           string   `json:"id"`
		DisplayName  string   `json:"display_name"`
		Category     string   `json:"category"`
		CPU          int      `json:"cpu"`
		RAM          int      `json:"ram_mb"`
		Storage      int      `json:"storage_gb"`
		PriceMonthly float32  `json:"price_monthly_eur"`
		PriceHourly  float32  `json:"price_hourly_eur"`
		Recommended  bool     `json:"recommended"`
		Default      bool     `json:"default"`
		UseCases     []string `json:"use_cases,omitempty"`
	}

	var output []SimplePlan
	defaultPlanID := plans.DefaultPlan

	for categoryName, category := range plans.Categories {
		// Filter by category if specified
		if cmd.Category != "" && categoryName != cmd.Category {
			continue
		}

		for _, plan := range category.Plans {
			// Filter by recommended if specified
			if cmd.Recommended && !plan.Recommended && plan.ID != defaultPlanID {
				continue
			}

			output = append(output, SimplePlan{
				ID:           plan.ID,
				DisplayName:  plan.DisplayName,
				Category:     categoryName,
				CPU:          plan.CPU,
				RAM:          plan.RAM,
				Storage:      plan.Storage,
				PriceMonthly: plan.PriceMonthly,
				PriceHourly:  plan.PriceHourly,
				Recommended:  plan.Recommended || category.RecommendedForDevpod,
				Default:      plan.ID == defaultPlanID,
				UseCases:     plan.UseCases,
			})
		}
	}

	// Use json encoder to output
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// outputYAML outputs plans as YAML
func (cmd *PlansCmd) outputYAML(plans *config.ServerPlans) error {
	// For YAML, we'll output the raw configuration
	// This allows users to save and modify it if needed
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()

	// Filter if needed
	if cmd.Category != "" {
		filtered := &config.ServerPlans{
			Version:     plans.Version,
			LastUpdated: plans.LastUpdated,
			DefaultPlan: plans.DefaultPlan,
			Categories:  make(map[string]*config.PlanCategory),
			Metadata:    plans.Metadata,
		}

		if category, exists := plans.Categories[cmd.Category]; exists {
			filtered.Categories[cmd.Category] = category
		}

		return encoder.Encode(filtered)
	}

	if cmd.Recommended {
		// Create a filtered structure with only recommended plans
		filtered := &config.ServerPlans{
			Version:     plans.Version,
			LastUpdated: plans.LastUpdated,
			DefaultPlan: plans.DefaultPlan,
			Categories:  make(map[string]*config.PlanCategory),
			Metadata:    plans.Metadata,
		}

		for categoryName, category := range plans.Categories {
			if category.RecommendedForDevpod {
				filtered.Categories[categoryName] = category
			}
		}

		return encoder.Encode(filtered)
	}

	return encoder.Encode(plans)
}