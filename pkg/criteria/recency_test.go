package criteria

import (
	"testing"
	"time"

	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func TestRecencyEvaluate(t *testing.T) {
	now := time.Now().Unix()
	twoMonthsAgo := now - (60 * 24 * 60 * 60)
	fiveMonthsAgo := now - (150 * 24 * 60 * 60)
	eightMonthsAgo := now - (240 * 24 * 60 * 60)
	tenMonthsAgo := now - (300 * 24 * 60 * 60)

	timeScores := TimeScores{
		ThreeMonths: 1.0,
		SixMonths:   0.7,
		NineMonths:  0.4,
	}

	tests := []struct {
		name          string
		listingDate   int64
		status        string
		statusChange  int64
		weight        float64
		expectedScore float64
	}{
		{
			name:          "recent listing within 3 months",
			listingDate:   twoMonthsAgo,
			status:        "Active",
			statusChange:  0,
			weight:        0.20,
			expectedScore: 0.20,
		},
		{
			name:          "listing within 6 months",
			listingDate:   fiveMonthsAgo,
			status:        "Active",
			statusChange:  0,
			weight:        0.20,
			expectedScore: 0.14,
		},
		{
			name:          "listing within 9 months",
			listingDate:   eightMonthsAgo,
			status:        "Active",
			statusChange:  0,
			weight:        0.20,
			expectedScore: 0.08,
		},
		{
			name:          "old listing beyond 9 months",
			listingDate:   tenMonthsAgo,
			status:        "Active",
			statusChange:  0,
			weight:        0.20,
			expectedScore: 0.02,
		},
		{
			name:          "closed property with recent status change",
			listingDate:   tenMonthsAgo,
			status:        "Closed",
			statusChange:  twoMonthsAgo,
			weight:        0.20,
			expectedScore: 0.20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			property := models.Property{
				ListingDate:           tt.listingDate,
				Status:                tt.status,
				StatusChangeTimestamp: tt.statusChange,
			}
			subject := models.Property{}

			recency := NewRecency(property, subject, tt.weight, timeScores)
			score, err := recency.Evaluate()

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !almostEqual(score, tt.expectedScore, 0.0001) {
				t.Errorf("expected score %v, got %v", tt.expectedScore, score)
			}
		})
	}
}
