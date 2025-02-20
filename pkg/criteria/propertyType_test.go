package criteria

import (
	"testing"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestPropertyTypeEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		propertyStyle string
		subjectStyle  string
		weight        float64
		expectedScore float64
	}{
		{
			name:          "exact match",
			propertyStyle: "Colonial",
			subjectStyle:  "Colonial",
			weight:        0.30,
			expectedScore: 0.30,
		},
		{
			name:          "different style",
			propertyStyle: "Ranch",
			subjectStyle:  "Colonial",
			weight:        0.30,
			expectedScore: 0.06,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{Style: tt.propertyStyle}
			subject := models.Property{Style: tt.subjectStyle}

			propertyType := NewPropertyType(property, subject, tt.weight)
			score, err := propertyType.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
