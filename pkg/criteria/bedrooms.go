package criteria

import (
	"math"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type Bedrooms struct {
	Property models.Property
	Subject  models.Property
	Weight   float64
}

func NewBedrooms(property, subject models.Property, weight float64) *Bedrooms {
	return &Bedrooms{
		Property: property,
		Subject:  subject,
		Weight:   weight,
	}
}

func (b *Bedrooms) Evaluate() (float64, error) {
	score := 0.0
	diff := math.Abs(float64(b.Property.Beds - b.Subject.Beds))
	if diff == 0 {
		score = 1.0
	} else {
		score = 1.0 / (1.0 + diff)
	}
	return score * b.Weight, nil
}
