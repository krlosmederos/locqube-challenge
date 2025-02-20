package criteria

import "github.com/krlosmederos/locqube-challenge/pkg/models"

type PropertyType struct {
	Property models.Property
	Subject  models.Property
	Weight   float64
}

func NewPropertyType(property, subject models.Property, weight float64) *PropertyType {
	return &PropertyType{
		Property: property,
		Subject:  subject,
		Weight:   weight,
	}
}

func (p *PropertyType) Evaluate() (float64, error) {
	score := 0.0
	if p.Property.Style == p.Subject.Style {
		score = 1.0
	} else {
		score = 0.2
	}
	return score * p.Weight, nil
}
