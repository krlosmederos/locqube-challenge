package criteria

import (
	"testing"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestBathroomsEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		propertyBaths float64
		subjectBaths  float64
		weight        float64
		expectedScore float64
	}{
		{
			name:          "exact match",
			propertyBaths: 2.5,
			subjectBaths:  2.5,
			weight:        0.05,
			expectedScore: 0.05,
		},
		{
			name:          "half bath difference",
			propertyBaths: 2.0,
			subjectBaths:  2.5,
			weight:        0.05,
			expectedScore: 0.033333,
		},
		{
			name:          "one bath difference",
			propertyBaths: 1.5,
			subjectBaths:  2.5,
			weight:        0.05,
			expectedScore: 0.025,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{Baths: models.Bathroom{Total: tt.propertyBaths}}
			subject := models.Property{Baths: models.Bathroom{Total: tt.subjectBaths}}

			bathrooms := NewBathrooms(property, subject, tt.weight)
			score, err := bathrooms.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
