package algorithm

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/krlosmederos/locqube-challenge/pkg/config"
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func setupTestConfig(t *testing.T) func() {
	t.Helper()

	// Create a temporary directory for test config
	tempDir, err := os.MkdirTemp("", "valuation_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create config directory
	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create test configuration
	testConfig := map[string]interface{}{
		"criteria_weights": map[string]float64{
			"property_type": 0.2,
			"bedrooms":      0.05,
			"bathrooms":     0.05,
			"size":          0.1,
			"recency":       0.5,
			"status":        0.1,
		},
		"time_scores": map[string]float64{
			"three_months": 1.0,
			"six_months":   0.5,
			"nine_months":  0.25,
		},
		"status_scores": map[string]float64{
			"sold":    1.0,
			"pending": 0.6,
			"active":  0.4,
		},
		"min_sales_count": 3,
	}

	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	configPath := filepath.Join(configDir, "application.json")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// Reset config singleton
	config.ResetForTesting()

	// Return cleanup function
	return func() {
		os.Chdir(originalWd)
		os.RemoveAll(tempDir)
	}
}

func createTestProperty(city string, size float64, beds int, baths float64, style string, status string, listingDate int64, statusChange int64, listPrice, salePrice float64) models.Property {
	return models.Property{
		Address: models.Address{
			City: city,
		},
		Size: size,
		Beds: beds,
		Baths: models.Bathroom{
			Total: baths,
		},
		Style:                 style,
		Status:                status,
		ListingDate:           listingDate,
		StatusChangeTimestamp: statusChange,
		ListPrice:             listPrice,
		SalePrice:             salePrice,
	}
}

func TestValuation_Calculate(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)
	twoMonthsAgo := now - (60 * 24 * 60 * 60)
	threeMonthsAgo := now - (90 * 24 * 60 * 60)

	tests := []struct {
		name          string
		subject       models.Property
		listings      []models.Property
		expectedValue float64
		tolerance     float64
	}{
		{
			name: "exact match properties",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings: []models.Property{
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					threeMonthsAgo, threeMonthsAgo, 600000, 600000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					twoMonthsAgo, twoMonthsAgo, 610000, 610000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					oneMonthAgo, oneMonthAgo, 620000, 620000,
				),
			},
			expectedValue: 610000,
			tolerance:     10000,
		},
		{
			name: "similar properties with variations",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings: []models.Property{
				createTestProperty(
					"Danbury", 1900, 4, 2.5, "Colonial", "Closed",
					threeMonthsAgo, threeMonthsAgo, 580000, 580000,
				),
				createTestProperty(
					"Danbury", 2100, 4, 2.5, "Colonial", "Closed",
					twoMonthsAgo, twoMonthsAgo, 620000, 620000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 3.0, "Colonial", "Closed",
					oneMonthAgo, oneMonthAgo, 600000, 600000,
				),
			},
			expectedValue: 600000,
			tolerance:     10000,
		},
		{
			name: "mix of sold and active listings",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings: []models.Property{
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					threeMonthsAgo, threeMonthsAgo, 600000, 600000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Under Contract",
					twoMonthsAgo, twoMonthsAgo, 610000, 0,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Active",
					oneMonthAgo, 0, 620000, 0,
				),
			},
			expectedValue: 605000,
			tolerance:     10000,
		},
		{
			name: "no similar properties",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings: []models.Property{
				createTestProperty(
					"Norwalk", 2000, 4, 2.5, "Colonial", "Closed",
					threeMonthsAgo, threeMonthsAgo, 600000, 600000,
				),
				createTestProperty(
					"Danbury", 3500, 6, 4.5, "Victorian", "Closed",
					twoMonthsAgo, twoMonthsAgo, 900000, 900000,
				),
			},
			expectedValue: 0,
			tolerance:     0,
		},
		{
			name: "properties with different recency",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings: []models.Property{
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					oneMonthAgo, oneMonthAgo, 620000, 620000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					twoMonthsAgo, twoMonthsAgo, 600000, 600000,
				),
				createTestProperty(
					"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
					threeMonthsAgo, threeMonthsAgo, 580000, 580000,
				),
			},
			expectedValue: 605000,
			tolerance:     10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuation := NewValuation(tt.subject, tt.listings)
			got := valuation.Calculate()

			if tt.expectedValue == 0 {
				if got != 0 {
					t.Errorf("Expected zero valuation, got %v", got)
				}
				return
			}

			diff := got - tt.expectedValue
			if diff < 0 {
				diff = -diff
			}

			if diff > tt.tolerance {
				t.Errorf("Valuation = %v, want %v (±%v)", got, tt.expectedValue, tt.tolerance)
			}
		})
	}
}

// Some tests here are failing because I should tune the weights and scores
// based on the data and the results (maybe I have to tune the tolerance as well)
func TestValuation_CalculateWeight(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)

	subject := createTestProperty(
		"Danbury", 2000, 4, 2.5, "Colonial", "Active",
		now, 0, 600000, 0,
	)

	tests := []struct {
		name          string
		comparable    models.Property
		expectedRange struct {
			min float64
			max float64
		}
	}{
		{
			name: "exact match property",
			comparable: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Closed",
				oneMonthAgo, oneMonthAgo, 600000, 600000,
			),
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.95,
				max: 1.00,
			},
		},
		{
			name: "similar property with small differences",
			comparable: createTestProperty(
				"Danbury", 2100, 4, 3.0, "Colonial", "Closed",
				oneMonthAgo, oneMonthAgo, 620000, 620000,
			),
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.80,
				max: 0.95,
			},
		},
		{
			name: "active listing",
			comparable: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				oneMonthAgo, 0, 600000, 0,
			),
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.85,
				max: 0.95,
			},
		},
		{
			name: "different property type",
			comparable: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Ranch", "Closed",
				oneMonthAgo, oneMonthAgo, 600000, 600000,
			),
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.65,
				max: 0.75,
			},
		},
		{
			name: "significant size difference",
			comparable: createTestProperty(
				"Danbury", 2500, 4, 2.5, "Colonial", "Closed",
				oneMonthAgo, oneMonthAgo, 600000, 600000,
			),
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.70,
				max: 0.85,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuation := NewValuation(subject, []models.Property{tt.comparable})
			weight, err := valuation.calculateWeight(tt.comparable)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			t.Logf("Test case '%s' got weight: %v", tt.name, weight)

			if weight < tt.expectedRange.min || weight > tt.expectedRange.max {
				t.Errorf("Weight = %v, want between %v and %v",
					weight, tt.expectedRange.min, tt.expectedRange.max)
			}
		})
	}
}

func TestValuation_EdgeCases(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	now := time.Now().Unix()

	tests := []struct {
		name          string
		subject       models.Property
		listings      []models.Property
		expectedValue float64
	}{
		{
			name: "empty listings",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings:      []models.Property{},
			expectedValue: 0,
		},
		{
			name: "nil listings",
			subject: createTestProperty(
				"Danbury", 2000, 4, 2.5, "Colonial", "Active",
				now, 0, 600000, 0,
			),
			listings:      nil,
			expectedValue: 0,
		},
		{
			name:          "zero value subject",
			subject:       models.Property{},
			listings:      []models.Property{},
			expectedValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuation := NewValuation(tt.subject, tt.listings)
			got := valuation.Calculate()

			if got != tt.expectedValue {
				t.Errorf("Valuation = %v, want %v", got, tt.expectedValue)
			}
		})
	}
}
