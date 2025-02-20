package criteria

import (
	"testing"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestSizeEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		propertySize  float64
		subjectSize   float64
		weight        float64
		expectedScore float64
	}{
		{
			name:          "exact match",
			propertySize:  2000,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.35,
		},
		{
			name:          "within 5% difference",
			propertySize:  2090,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.35,
		},
		{
			name:          "within 10% difference",
			propertySize:  2180,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.28,
		},
		{
			name:          "within 20% difference",
			propertySize:  2350,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.175,
		},
		{
			name:          "within 30% difference",
			propertySize:  2500,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.07,
		},
		{
			name:          "more than 30% difference",
			propertySize:  3000,
			subjectSize:   2000,
			weight:        0.35,
			expectedScore: 0.035,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{Size: tt.propertySize}
			subject := models.Property{Size: tt.subjectSize}

			size := NewSize(property, subject, tt.weight)
			score, err := size.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
