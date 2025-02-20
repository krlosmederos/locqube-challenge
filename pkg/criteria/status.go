package criteria

import "github.com/krlosmederos/locqube-challenge/pkg/models"

type StatusScores struct {
	Sold    float64
	Pending float64
	Active  float64
}

type Status struct {
	Property     models.Property
	Subject      models.Property
	Weight       float64
	StatusScores StatusScores
}

func NewStatus(property, subject models.Property, weight float64, statusScores StatusScores) *Status {
	return &Status{
		Property:     property,
		Subject:      subject,
		Weight:       weight,
		StatusScores: statusScores,
	}
}

func (s *Status) Evaluate() (float64, error) {
	var score float64
	switch s.Property.Status {
	case "Closed":
		score = s.StatusScores.Sold
	case "Under Contract":
		score = s.StatusScores.Pending
	case "Active":
		score = s.StatusScores.Active
	default:
		score = s.StatusScores.Active
	}

	return score * s.Weight, nil
}
