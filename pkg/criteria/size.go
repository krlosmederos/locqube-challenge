package criteria

import (
	"math"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type Size struct {
	Property models.Property
	Subject  models.Property
	Weight   float64
}

func NewSize(property, subject models.Property, weight float64) *Size {
	return &Size{
		Property: property,
		Subject:  subject,
		Weight:   weight,
	}
}

func (s *Size) Evaluate() (float64, error) {
	sizeDiff := math.Abs(s.Property.Size-s.Subject.Size) / s.Subject.Size
	score := 0.0

	if sizeDiff <= 0.05 {
		score = 1.0
	} else if sizeDiff <= 0.10 {
		score = 0.8
	} else if sizeDiff <= 0.20 {
		score = 0.5
	} else if sizeDiff <= 0.30 {
		score = 0.2
	} else {
		score = 0.1
	}

	return score * s.Weight, nil
}
