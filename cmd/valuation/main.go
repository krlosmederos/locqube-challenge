package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/krlosmederos/locqube-challenge/pkg/algorithm"
	"github.com/krlosmederos/locqube-challenge/pkg/config"
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func main() {
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	subject := models.Property{
		PropertyType: "Single-family home",
		Beds:         4,
		Baths: models.Bathroom{
			Total: 3.5,
			Full:  3,
			Half:  1,
		},
		Size: 2750,
		Address: models.Address{
			City:  "Danbury",
			State: "CT",
		},
	}

	data, err := os.ReadFile("data/market_listings_response.json")
	if err != nil {
		log.Fatalf("Error reading market listings: %v", err)
	}

	var listings []models.Property
	if err := json.Unmarshal(data, &listings); err != nil {
		log.Fatalf("Error parsing market listings: %v", err)
	}

	valuation := algorithm.NewValuation(subject, listings)
	estimatedValue := valuation.Calculate()

	fmt.Printf("Estimated Property Value: $%.2f\n", estimatedValue)
}
