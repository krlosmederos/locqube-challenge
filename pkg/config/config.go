package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	CriteriaWeights struct {
		PropertyType float64 `json:"property_type"`
		Bedrooms     float64 `json:"bedrooms"`
		Bathrooms    float64 `json:"bathrooms"`
		Size         float64 `json:"size"`
		Recency      float64 `json:"recency"`
		Status       float64 `json:"status"`
	} `json:"criteria_weights"`

	TimeScores struct {
		ThreeMonths float64 `json:"three_months"`
		SixMonths   float64 `json:"six_months"`
		NineMonths  float64 `json:"nine_months"`
	} `json:"time_scores"`

	StatusScores struct {
		Sold    float64 `json:"sold"`
		Pending float64 `json:"pending"`
		Active  float64 `json:"active"`
	} `json:"status_scores"`

	MinSalesCount int `json:"min_sales_count"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig loads the configuration from the JSON file
func LoadConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}

		configPath := filepath.Join("config", "application.json")

		file, err := os.ReadFile(configPath)
		if err != nil {
			panic("Failed to read config file: " + err.Error())
		}

		if err := json.Unmarshal(file, config); err != nil {
			panic("Failed to parse config file: " + err.Error())
		}
	})

	return config, nil
}

// GetConfig returns the current configuration or loads it if it's not loaded
func GetConfig() *Config {
	if config == nil {
		_, err := LoadConfig()
		if err != nil {
			panic("Failed to load config: " + err.Error())
		}
	}
	return config
}

// ResetForTesting resets the config singleton for testing purposes
func ResetForTesting() {
	config = nil
	once = sync.Once{}
}
