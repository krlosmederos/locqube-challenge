package criteria

import (
	"math"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type Bathrooms struct {
	Property models.Property
	Subject  models.Property
	Weight   float64
}

func NewBathrooms(property, subject models.Property, weight float64) *Bathrooms {
	return &Bathrooms{
		Property: property,
		Subject:  subject,
		Weight:   weight,
	}
}

func (b *Bathrooms) Evaluate() (float64, error) {
	score := 0.0
	diff := math.Abs(b.Property.Baths.Total - b.Subject.Baths.Total)
	if diff == 0 {
		score = 1.0
	} else {
		score = 1.0 / (1.0 + diff)
	}
	return score * b.Weight, nil
}
