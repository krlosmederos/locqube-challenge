package algorithm

import (
	"fmt"
	"sync"

	"github.com/krlosmederos/locqube-challenge/pkg/config"
	"github.com/krlosmederos/locqube-challenge/pkg/criteria"
	"github.com/krlosmederos/locqube-challenge/pkg/filters"
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type CriteriaEvaluator interface {
	Evaluate() (float64, error)
}

type Valuation struct {
	Subject  models.Property
	Listings []models.Property
	Config   *config.Config
	filter   *filters.PropertyFilter
}

func NewValuation(subject models.Property, listings []models.Property) *Valuation {
	cfg := config.GetConfig()
	return &Valuation{
		Subject:  subject,
		Listings: listings,
		Config:   cfg,
		filter:   filters.NewPropertyFilter(subject, cfg),
	}
}

// Calculate calculates the valuation of the subject property
func (v *Valuation) Calculate() float64 {
	type priceResult struct {
		price  float64
		weight float64
	}

	filteredListings := v.filter.Filter(v.Listings)

	results := make(chan priceResult, len(filteredListings))
	var wg sync.WaitGroup

	for _, prop := range filteredListings {
		wg.Add(1)
		go func(comp models.Property) {
			defer wg.Done()

			weight, err := v.calculateWeight(comp)
			if err != nil {
				fmt.Printf("Error calculating weight for property %s: %v\n", comp.ID, err)
				return
			}

			price := comp.GetPrice()
			results <- priceResult{price: price, weight: weight}
		}(prop)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalWeight, weightedSum float64
	for result := range results {
		weightedSum += result.price * result.weight
		totalWeight += result.weight
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

func (v *Valuation) calculateWeight(comp models.Property) (float64, error) {
	criteriaList := []CriteriaEvaluator{
		criteria.NewPropertyType(comp, v.Subject, v.Config.CriteriaWeights.PropertyType),
		criteria.NewBedrooms(comp, v.Subject, v.Config.CriteriaWeights.Bedrooms),
		criteria.NewBathrooms(comp, v.Subject, v.Config.CriteriaWeights.Bathrooms),
		criteria.NewSize(comp, v.Subject, v.Config.CriteriaWeights.Size),
		criteria.NewRecency(comp, v.Subject, v.Config.CriteriaWeights.Recency, criteria.TimeScores{
			ThreeMonths: v.Config.TimeScores.ThreeMonths,
			SixMonths:   v.Config.TimeScores.SixMonths,
			NineMonths:  v.Config.TimeScores.NineMonths,
		}),
		criteria.NewStatus(comp, v.Subject, v.Config.CriteriaWeights.Status, criteria.StatusScores{
			Sold:    v.Config.StatusScores.Sold,
			Pending: v.Config.StatusScores.Pending,
			Active:  v.Config.StatusScores.Active,
		}),
	}

	var totalScore float64
	for _, c := range criteriaList {
		score, err := c.Evaluate()
		if err != nil {
			return 0, fmt.Errorf("error evaluating criteria: %v", err)
		}
		totalScore += score
	}

	return totalScore, nil
}
