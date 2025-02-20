package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func testSetup(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	configDir := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	return tempDir, func() {
		os.Chdir(originalWd)
		os.RemoveAll(tempDir)
	}
}

func TestLoadConfig(t *testing.T) {
	config = nil
	once = sync.Once{}

	tempDir, cleanup := testSetup(t)
	defer cleanup()

	testConfig := Config{
		CriteriaWeights: struct {
			PropertyType float64 `json:"property_type"`
			Bedrooms     float64 `json:"bedrooms"`
			Bathrooms    float64 `json:"bathrooms"`
			Size         float64 `json:"size"`
			Recency      float64 `json:"recency"`
			Status       float64 `json:"status"`
		}{
			PropertyType: 0.3,
			Bedrooms:     0.05,
			Bathrooms:    0.05,
			Size:         0.35,
			Recency:      0.2,
			Status:       0.05,
		},
		TimeScores: struct {
			ThreeMonths float64 `json:"three_months"`
			SixMonths   float64 `json:"six_months"`
			NineMonths  float64 `json:"nine_months"`
		}{
			ThreeMonths: 1.0,
			SixMonths:   0.7,
			NineMonths:  0.4,
		},
		StatusScores: struct {
			Sold    float64 `json:"sold"`
			Pending float64 `json:"pending"`
			Active  float64 `json:"active"`
		}{
			Sold:    1.0,
			Pending: 0.6,
			Active:  0.4,
		},
		MinSalesCount: 3,
	}

	configData, err := json.MarshalIndent(testConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	configPath := filepath.Join(tempDir, "config", "application.json")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.CriteriaWeights.PropertyType != 0.3 {
		t.Errorf("Expected PropertyType weight 0.3, got %v", cfg.CriteriaWeights.PropertyType)
	}
	if cfg.CriteriaWeights.Size != 0.35 {
		t.Errorf("Expected Size weight 0.35, got %v", cfg.CriteriaWeights.Size)
	}
	if cfg.TimeScores.ThreeMonths != 1.0 {
		t.Errorf("Expected ThreeMonths score 1.0, got %v", cfg.TimeScores.ThreeMonths)
	}
	if cfg.StatusScores.Sold != 1.0 {
		t.Errorf("Expected Sold score 1.0, got %v", cfg.StatusScores.Sold)
	}
	if cfg.MinSalesCount != 3 {
		t.Errorf("Expected MinSalesCount 3, got %v", cfg.MinSalesCount)
	}

	cfg2, err := LoadConfig()
	if err != nil {
		t.Fatalf("Second LoadConfig() failed: %v", err)
	}
	if cfg != cfg2 {
		t.Error("LoadConfig() did not return the same instance")
	}
}

func TestGetConfig(t *testing.T) {
	// Reset shared state
	config = nil
	once = sync.Once{}

	tempDir, cleanup := testSetup(t)
	defer cleanup()

	minimalConfig := map[string]interface{}{
		"criteria_weights": map[string]float64{
			"property_type": 0.3,
			"bedrooms":      0.05,
			"bathrooms":     0.05,
			"size":          0.35,
			"recency":       0.2,
			"status":        0.05,
		},
		"time_scores": map[string]float64{
			"three_months": 1.0,
			"six_months":   0.7,
			"nine_months":  0.4,
		},
		"status_scores": map[string]float64{
			"sold":    1.0,
			"pending": 0.6,
			"active":  0.4,
		},
		"min_sales_count": 3,
	}

	configData, err := json.MarshalIndent(minimalConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal minimal config: %v", err)
	}

	configPath := filepath.Join(tempDir, "config", "application.json")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		t.Fatalf("Failed to write minimal config: %v", err)
	}

	cfg := GetConfig()
	if cfg == nil {
		t.Fatal("GetConfig() returned nil")
	}

	if cfg.CriteriaWeights.PropertyType != 0.3 {
		t.Errorf("Expected PropertyType weight 0.3, got %v", cfg.CriteriaWeights.PropertyType)
	}

	cfg2 := GetConfig()
	if cfg != cfg2 {
		t.Error("GetConfig() did not return the same instance")
	}
}

func TestConfigWeightsSum(t *testing.T) {
	cfg := &Config{
		CriteriaWeights: struct {
			PropertyType float64 `json:"property_type"`
			Bedrooms     float64 `json:"bedrooms"`
			Bathrooms    float64 `json:"bathrooms"`
			Size         float64 `json:"size"`
			Recency      float64 `json:"recency"`
			Status       float64 `json:"status"`
		}{
			PropertyType: 0.3,
			Bedrooms:     0.05,
			Bathrooms:    0.05,
			Size:         0.35,
			Recency:      0.2,
			Status:       0.05,
		},
	}

	sum := cfg.CriteriaWeights.PropertyType +
		cfg.CriteriaWeights.Bedrooms +
		cfg.CriteriaWeights.Bathrooms +
		cfg.CriteriaWeights.Size +
		cfg.CriteriaWeights.Recency +
		cfg.CriteriaWeights.Status

	if sum != 1.0 {
		t.Errorf("Criteria weights sum to %v, expected 1.0", sum)
	}
}
