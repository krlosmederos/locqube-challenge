package criteria

import (
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

type TimeScores struct {
	ThreeMonths float64
	SixMonths   float64
	NineMonths  float64
}

type Recency struct {
	Property   models.Property
	Subject    models.Property
	Weight     float64
	TimeScores TimeScores
}

func NewRecency(property, subject models.Property, weight float64, timeScores TimeScores) *Recency {
	return &Recency{
		Property:   property,
		Subject:    subject,
		Weight:     weight,
		TimeScores: timeScores,
	}
}

func (r *Recency) Evaluate() (float64, error) {
	ageInMonths := r.Property.GetAgeInMonths()
	score := 0.0

	switch {
	case ageInMonths <= 3:
		score = r.TimeScores.ThreeMonths
	case ageInMonths <= 6:
		score = r.TimeScores.SixMonths
	case ageInMonths <= 9:
		score = r.TimeScores.NineMonths
	default:
		score = 0.1
	}

	return score * r.Weight, nil
}
