package filters

import (
	"testing"
	"time"

	"github.com/krlosmederos/locqube-challenge/pkg/config"
	"github.com/krlosmederos/locqube-challenge/pkg/models"
)

func createTestConfig() *config.Config {
	return &config.Config{
		MinSalesCount: 3,
	}
}

func createTestProperty(city string, size float64, beds int, baths float64, status string, listingDate int64, statusChange int64) models.Property {
	return models.Property{
		Address: models.Address{
			City: city,
		},
		Size: size,
		Beds: beds,
		Baths: models.Bathroom{
			Total: baths,
		},
		Status:                status,
		ListingDate:           listingDate,
		StatusChangeTimestamp: statusChange,
		ListPrice:             500000,
		SalePrice:             600000,
		Style:                 "Colonial",
	}
}

func TestPropertyFilter_IsSimilarProperty(t *testing.T) {
	subject := createTestProperty("Danbury", 2000, 4, 2.5, "Active", 0, 0)
	filter := NewPropertyFilter(subject, createTestConfig())

	tests := []struct {
		name     string
		property models.Property
		want     bool
	}{
		{
			name:     "exact match",
			property: createTestProperty("Danbury", 2000, 4, 2.5, "Active", 0, 0),
			want:     true,
		},
		{
			name:     "different city",
			property: createTestProperty("Norwalk", 2000, 4, 2.5, "Active", 0, 0),
			want:     false,
		},
		{
			name:     "size too different",
			property: createTestProperty("Danbury", 3000, 4, 2.5, "Active", 0, 0),
			want:     false,
		},
		{
			name:     "bedrooms too different",
			property: createTestProperty("Danbury", 2000, 6, 2.5, "Active", 0, 0),
			want:     false,
		},
		{
			name:     "bathrooms too different",
			property: createTestProperty("Danbury", 2000, 4, 4.0, "Active", 0, 0),
			want:     false,
		},
		{
			name: "no price information",
			property: models.Property{
				Address: models.Address{City: "Danbury"},
				Size:    2000,
				Beds:    4,
				Baths:   models.Bathroom{Total: 2.5},
				Status:  "Active",
			},
			want: false,
		},
		{
			name:     "within acceptable size range",
			property: createTestProperty("Danbury", 2200, 4, 2.5, "Active", 0, 0),
			want:     true,
		},
		{
			name:     "within acceptable bedroom range",
			property: createTestProperty("Danbury", 2000, 5, 2.5, "Active", 0, 0),
			want:     true,
		},
		{
			name:     "within acceptable bathroom range",
			property: createTestProperty("Danbury", 2000, 4, 3.0, "Active", 0, 0),
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter.isSimilarProperty(tt.property); got != tt.want {
				t.Errorf("PropertyFilter.isSimilarProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPropertyFilter_SortByStatusAndRecency(t *testing.T) {
	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)
	twoMonthsAgo := now - (60 * 24 * 60 * 60)
	fourMonthsAgo := now - (120 * 24 * 60 * 60)
	sevenMonthsAgo := now - (210 * 24 * 60 * 60)

	subject := createTestProperty("Danbury", 2000, 4, 2.5, "Active", now, 0)
	filter := NewPropertyFilter(subject, createTestConfig())

	properties := []models.Property{
		createTestProperty("Danbury", 2000, 4, 2.5, "Active", twoMonthsAgo, 0),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", sevenMonthsAgo, oneMonthAgo),
		createTestProperty("Danbury", 2000, 4, 2.5, "Under Contract", fourMonthsAgo, 0),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", fourMonthsAgo, twoMonthsAgo),
	}

	sorted := filter.sortByStatusAndRecency(properties)

	if len(sorted) < 2 {
		t.Fatal("Expected at least 2 properties after sorting")
	}

	if sorted[0].Status != "Closed" || sorted[0].StatusChangeTimestamp != oneMonthAgo {
		t.Errorf("First property should be the most recent closed property")
	}

	if sorted[1].Status != "Closed" || sorted[1].StatusChangeTimestamp != twoMonthsAgo {
		t.Errorf("Second property should be the second most recent closed property")
	}
}

func TestPropertyFilter_GetMostRecentSoldProperties(t *testing.T) {
	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)
	twoMonthsAgo := now - (60 * 24 * 60 * 60)
	fourMonthsAgo := now - (120 * 24 * 60 * 60)
	sevenMonthsAgo := now - (210 * 24 * 60 * 60)
	tenMonthsAgo := now - (300 * 24 * 60 * 60)

	subject := createTestProperty("Danbury", 2000, 4, 2.5, "Active", now, 0)
	cfg := createTestConfig()
	cfg.MinSalesCount = 3
	filter := NewPropertyFilter(subject, cfg)

	soldProperties := []models.Property{
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", tenMonthsAgo, tenMonthsAgo),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", sevenMonthsAgo, sevenMonthsAgo),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", fourMonthsAgo, fourMonthsAgo),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", twoMonthsAgo, twoMonthsAgo),
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", oneMonthAgo, oneMonthAgo),
	}

	recent := filter.getMostRecentSoldProperties(soldProperties)

	if len(recent) != 3 {
		t.Errorf("Expected 3 recent sold properties, got %d", len(recent))
	}
}

func TestPropertyFilter_GetMaxAgeForSales(t *testing.T) {
	filter := NewPropertyFilter(models.Property{}, createTestConfig())

	tests := []struct {
		name    string
		sales3M int
		sales6M int
		wantAge float64
	}{
		{
			name:    "enough recent sales",
			sales3M: 3,
			sales6M: 2,
			wantAge: 3.0,
		},
		{
			name:    "enough sales within 6 months",
			sales3M: 2,
			sales6M: 2,
			wantAge: 6.0,
		},
		{
			name:    "not enough recent sales",
			sales3M: 1,
			sales6M: 1,
			wantAge: 9.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filter.getMaxAgeForSales(tt.sales3M, tt.sales6M); got != tt.wantAge {
				t.Errorf("PropertyFilter.getMaxAgeForSales() = %v, want %v", got, tt.wantAge)
			}
		})
	}
}

func TestPropertyFilter_Filter(t *testing.T) {
	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)
	twoMonthsAgo := now - (60 * 24 * 60 * 60)
	fourMonthsAgo := now - (120 * 24 * 60 * 60)

	subject := createTestProperty("Danbury", 2000, 4, 2.5, "Active", now, 0)
	filter := NewPropertyFilter(subject, createTestConfig())

	listings := []models.Property{
		// Similar properties
		createTestProperty("Danbury", 2000, 4, 2.5, "Closed", twoMonthsAgo, twoMonthsAgo),
		createTestProperty("Danbury", 2100, 4, 2.5, "Closed", oneMonthAgo, oneMonthAgo),
		createTestProperty("Danbury", 1900, 4, 2.0, "Active", fourMonthsAgo, 0),
		// non-similar properties
		createTestProperty("Norwalk", 2000, 4, 2.5, "Closed", oneMonthAgo, oneMonthAgo),
		createTestProperty("Danbury", 3000, 4, 2.5, "Closed", oneMonthAgo, oneMonthAgo),
		createTestProperty("Danbury", 2000, 6, 2.5, "Closed", oneMonthAgo, oneMonthAgo),
		createTestProperty("Danbury", 2000, 4, 4.0, "Closed", oneMonthAgo, oneMonthAgo),
	}

	filtered := filter.Filter(listings)

	if len(filtered) != 3 {
		t.Errorf("Expected 3 comparable properties, got %d", len(filtered))
	}

	for _, prop := range filtered {
		if !filter.isSimilarProperty(prop) {
			t.Errorf("Filtered property %+v is not similar to subject", prop)
		}
	}

	var lastWasClosed bool
	for i, prop := range filtered {
		if i == 0 {
			lastWasClosed = prop.Status == "Closed"
			continue
		}
		if lastWasClosed && prop.Status != "Closed" {
			continue
		}
		if !lastWasClosed && prop.Status == "Closed" {
			t.Errorf("Found closed property after non-closed property in sorted results")
		}
		lastWasClosed = prop.Status == "Closed"
	}
}
