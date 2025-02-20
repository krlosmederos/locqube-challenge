package criteria

import (
	"testing"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestStatusEvaluate(t *testing.T) {
	statusScores := StatusScores{
		Sold:    1.0,
		Pending: 0.6,
		Active:  0.4,
	}

	tests := []struct {
		name           string
		propertyStatus string
		weight         float64
		expectedScore  float64
	}{
		{
			name:           "closed property",
			propertyStatus: "Closed",
			weight:         0.05,
			expectedScore:  0.05,
		},
		{
			name:           "pending property",
			propertyStatus: "Under Contract",
			weight:         0.05,
			expectedScore:  0.03,
		},
		{
			name:           "active property",
			propertyStatus: "Active",
			weight:         0.05,
			expectedScore:  0.02,
		},
		{
			name:           "unknown status",
			propertyStatus: "Unknown",
			weight:         0.05,
			expectedScore:  0.02,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{Status: tt.propertyStatus}
			subject := models.Property{}

			status := NewStatus(property, subject, tt.weight, statusScores)
			score, err := status.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
