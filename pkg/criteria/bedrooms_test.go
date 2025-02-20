package criteria

import (
	"testing"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestBedroomsEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		propertyBeds  int
		subjectBeds   int
		weight        float64
		expectedScore float64
	}{
		{
			name:          "exact match",
			propertyBeds:  4,
			subjectBeds:   4,
			weight:        0.05,
			expectedScore: 0.05,
		},
		{
			name:          "one bedroom difference",
			propertyBeds:  3,
			subjectBeds:   4,
			weight:        0.05,
			expectedScore: 0.025,
		},
		{
			name:          "two bedrooms difference",
			propertyBeds:  2,
			subjectBeds:   4,
			weight:        0.05,
			expectedScore: 0.0166667,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{Beds: tt.propertyBeds}
			subject := models.Property{Beds: tt.subjectBeds}

			bedrooms := NewBedrooms(property, subject, tt.weight)
			score, err := bedrooms.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
