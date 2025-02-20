package filters

import (
	"math"
	"sort"

	"github.com/krlosmederos/locqube-challenge/pkg/config"
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type PropertyFilter struct {
	Subject models.Property
	Config  *config.Config
}

func NewPropertyFilter(subject models.Property, config *config.Config) *PropertyFilter {
	return &PropertyFilter{
		Subject: subject,
		Config:  config,
	}
}

// Filter returns a list of comparable properties based on the subject property
func (f *PropertyFilter) Filter(listings []models.Property) []models.Property {
	var comparableProperties []models.Property

	for _, prop := range listings {
		if !f.isSimilarProperty(prop) {
			continue
		}

		comparableProperties = append(comparableProperties, prop)
	}

	return f.sortByStatusAndRecency(comparableProperties)
}

func (f *PropertyFilter) isSimilarProperty(prop models.Property) bool {
	if prop.ListPrice == 0 && prop.SalePrice == 0 {
		return false
	}

	if prop.Address.City != f.Subject.Address.City {
		return false
	}

	sizeDiff := math.Abs(prop.Size-f.Subject.Size) / f.Subject.Size
	if sizeDiff > 0.20 {
		return false
	}

	bedDiff := math.Abs(float64(prop.Beds - f.Subject.Beds))
	if bedDiff > 1 {
		return false
	}

	bathDiff := math.Abs(prop.Baths.Total - f.Subject.Baths.Total)
	return bathDiff <= 0.5
}

func (f *PropertyFilter) sortByStatusAndRecency(properties []models.Property) []models.Property {
	var soldProperties, nonSoldProperties []models.Property

	for _, prop := range properties {
		if prop.Status == "Closed" && prop.StatusChangeTimestamp > 0 {
			soldProperties = append(soldProperties, prop)
		} else {
			nonSoldProperties = append(nonSoldProperties, prop)
		}
	}

	result := f.getMostRecentSoldProperties(soldProperties)

	result = append(result, nonSoldProperties...)

	return result
}

func (f *PropertyFilter) getMostRecentSoldProperties(soldProperties []models.Property) []models.Property {
	sort.Slice(soldProperties, func(i, j int) bool {
		return soldProperties[i].GetAgeInMonths() < soldProperties[j].GetAgeInMonths()
	})

	var sales3M, sales6M, sales9M int
	for _, p := range soldProperties {
		age := p.GetAgeInMonths()
		switch {
		case age <= 3:
			sales3M++
		case age <= 6:
			sales6M++
		case age <= 9:
			sales9M++
		}
	}

	maxAge := f.getMaxAgeForSales(sales3M, sales6M)

	var result []models.Property
	for _, p := range soldProperties {
		if p.GetAgeInMonths() <= maxAge {
			result = append(result, p)
		}
	}

	return result
}

func (f *PropertyFilter) getMaxAgeForSales(sales3M, sales6M int) float64 {
	if sales3M >= f.Config.MinSalesCount {
		return 3.0
	}
	if sales3M+sales6M >= f.Config.MinSalesCount {
		return 6.0
	}
	return 9.0
}
