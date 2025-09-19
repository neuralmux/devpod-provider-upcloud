package config

import (
	"strings"
	"testing"
)

func TestLoadServerPlans(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	// Check basic structure
	if plans.Version == "" {
		t.Error("Version should not be empty")
	}

	if plans.DefaultPlan == "" {
		t.Error("DefaultPlan should not be empty")
	}

	if len(plans.Categories) == 0 {
		t.Error("Categories should not be empty")
	}
}

func TestGetPlanByID(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	tests := []struct {
		name      string
		planID    string
		wantError bool
	}{
		{
			name:      "Valid Developer Plan",
			planID:    "DEV-2xCPU-4GB",
			wantError: false,
		},
		{
			name:      "Valid General Purpose Plan",
			planID:    "2xCPU-4GB",
			wantError: false,
		},
		{
			name:      "Invalid Plan",
			planID:    "INVALID-PLAN",
			wantError: true,
		},
		{
			name:      "Cloud Native Plan",
			planID:    "CN-2xCPU-4GB",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan, category, err := plans.GetPlanByID(tt.planID)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error for plan %s, but got none", tt.planID)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for plan %s: %v", tt.planID, err)
				}
				if plan == nil {
					t.Errorf("Plan should not be nil for %s", tt.planID)
				}
				if category == nil {
					t.Errorf("Category should not be nil for %s", tt.planID)
				}
			}
		})
	}
}

func TestGetRecommendedPlans(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	recommended := plans.GetRecommendedPlans()

	if len(recommended) == 0 {
		t.Error("Should have at least one recommended plan")
	}

	// Check that developer plans are included
	hasDeveloperPlan := false
	for _, plan := range recommended {
		if strings.HasPrefix(plan.ID, "DEV-") {
			hasDeveloperPlan = true
			break
		}
	}

	if !hasDeveloperPlan {
		t.Error("Recommended plans should include at least one developer plan")
	}
}

func TestGetDeveloperPlans(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	devPlans := plans.GetDeveloperPlans()

	if len(devPlans) == 0 {
		t.Error("Should have developer plans")
	}

	for _, plan := range devPlans {
		if !strings.HasPrefix(plan.ID, "DEV-") {
			t.Errorf("Developer plan %s should start with DEV-", plan.ID)
		}
	}
}

func TestGetDefaultPlan(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	defaultPlan, err := plans.GetDefaultPlan()
	if err != nil {
		t.Errorf("Failed to get default plan: %v", err)
	}

	if defaultPlan == nil {
		t.Error("Default plan should not be nil")
	}

	if defaultPlan.ID != plans.DefaultPlan {
		t.Errorf("Default plan ID mismatch: got %s, want %s", defaultPlan.ID, plans.DefaultPlan)
	}

	// Default should be DEV-2xCPU-4GB as per configuration
	if defaultPlan.ID != "DEV-2xCPU-4GB" {
		t.Errorf("Expected default plan to be DEV-2xCPU-4GB, got %s", defaultPlan.ID)
	}
}

func TestValidatePlan(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	tests := []struct {
		name      string
		planID    string
		wantError bool
	}{
		{"Valid Developer Plan", "DEV-2xCPU-4GB", false},
		{"Valid Cloud Native Plan", "CN-2xCPU-4GB", false},
		{"Valid General Purpose Plan", "4xCPU-8GB", false},
		{"Invalid Plan", "NONEXISTENT", true},
		{"Empty Plan", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := plans.ValidatePlan(tt.planID)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidatePlan(%s) error = %v, wantError %v", tt.planID, err, tt.wantError)
			}
		})
	}
}

func TestGetPlanSuggestions(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	suggestions := plans.GetPlanSuggestions()

	if len(suggestions) == 0 {
		t.Error("Should have plan suggestions")
	}

	// Check that developer plans come first
	foundDevFirst := false
	for i, suggestion := range suggestions {
		if strings.HasPrefix(suggestion, "DEV-") {
			if i < 5 { // Developer plans should be in the first few suggestions
				foundDevFirst = true
			}
			break
		}
	}

	if !foundDevFirst {
		t.Error("Developer plans should appear first in suggestions")
	}
}

func TestGetPlanRecommendation(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	tests := []struct {
		name      string
		language  string
		framework string
		workload  string
		want      string
	}{
		{
			name:     "Python language",
			language: "python",
			want:     "DEV-2xCPU-8GB",
		},
		{
			name:      "React framework",
			framework: "react",
			want:      "DEV-2xCPU-4GB",
		},
		{
			name:     "ML workload",
			workload: "ml_development",
			want:     "HCPU-8xCPU-16GB",
		},
		{
			name:     "Microservices workload",
			workload: "microservices",
			want:     "CN-2xCPU-4GB",
		},
		{
			name: "Default (no criteria)",
			want: "DEV-2xCPU-4GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := plans.GetPlanRecommendation(tt.language, tt.framework, tt.workload)
			if got != tt.want {
				t.Errorf("GetPlanRecommendation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidRegion(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	tests := []struct {
		name   string
		region string
		want   bool
	}{
		{"Valid Frankfurt", "de-fra1", true},
		{"Valid Singapore", "sg-sin1", true},
		{"Valid New York", "us-nyc1", true},
		{"Invalid region", "invalid-region", false},
		{"Empty region", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := plans.IsValidRegion(tt.region); got != tt.want {
				t.Errorf("IsValidRegion(%s) = %v, want %v", tt.region, got, tt.want)
			}
		})
	}
}

func TestGetRegions(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	regions := plans.GetRegions()

	if len(regions) == 0 {
		t.Error("Should have available regions")
	}

	// Check for some expected regions
	expectedRegions := []string{"de-fra1", "us-nyc1", "sg-sin1"}
	for _, expected := range expectedRegions {
		found := false
		for _, region := range regions {
			if region == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected region %s not found in regions list", expected)
		}
	}
}

func TestPlanPricing(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	// Test that developer plans have correct pricing
	plan, _, err := plans.GetPlanByID("DEV-2xCPU-4GB")
	if err != nil {
		t.Fatalf("Failed to get DEV-2xCPU-4GB plan: %v", err)
	}

	if plan.PriceMonthly != 18.00 {
		t.Errorf("DEV-2xCPU-4GB monthly price should be €18, got €%.2f", plan.PriceMonthly)
	}

	// Check that developer plans are cheaper than general purpose
	devPlan, _, _ := plans.GetPlanByID("DEV-2xCPU-4GB")
	gpPlan, _, _ := plans.GetPlanByID("2xCPU-4GB")

	if devPlan.PriceMonthly >= gpPlan.PriceMonthly {
		t.Errorf("Developer plan should be cheaper than general purpose: DEV €%.2f >= GP €%.2f",
			devPlan.PriceMonthly, gpPlan.PriceMonthly)
	}
}

func TestPlanRestrictions(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	// Test minimal dev plan has restrictions
	plan, _, err := plans.GetPlanByID("DEV-1xCPU-1GB-10GB")
	if err != nil {
		t.Fatalf("Failed to get DEV-1xCPU-1GB-10GB plan: %v", err)
	}

	if plan.Restrictions == nil {
		t.Error("DEV-1xCPU-1GB-10GB should have restrictions")
	} else if plan.Restrictions.MaxPerAccount != 2 {
		t.Errorf("DEV-1xCPU-1GB-10GB should have max 2 per account, got %d",
			plan.Restrictions.MaxPerAccount)
	}
}

func TestFormatPlanList(t *testing.T) {
	plans, err := LoadServerPlans()
	if err != nil {
		t.Fatalf("Failed to load server plans: %v", err)
	}

	// Test basic formatting
	output := plans.FormatPlanList(false)
	if output == "" {
		t.Error("FormatPlanList should return non-empty output")
	}

	// Test that it includes version
	if !strings.Contains(output, plans.Version) {
		t.Error("Output should include version")
	}

	// Test detailed formatting
	detailedOutput := plans.FormatPlanList(true)
	if len(detailedOutput) <= len(output) {
		t.Error("Detailed output should be longer than basic output")
	}

	// Check for use cases in detailed output
	if !strings.Contains(detailedOutput, "Use cases:") {
		t.Error("Detailed output should include use cases")
	}
}