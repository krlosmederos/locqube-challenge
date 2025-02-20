package models

import (
	"testing"
	"time"
)

func TestProperty_GetPrice(t *testing.T) {
	tests := []struct {
		name     string
		property Property
		expected float64
	}{
		{
			name: "sold property returns sale price",
			property: Property{
				ListPrice: 500000,
				SalePrice: 520000,
			},
			expected: 520000,
		},
		{
			name: "active property returns list price",
			property: Property{
				ListPrice: 500000,
				SalePrice: 0,
			},
			expected: 500000,
		},
		{
			name: "pending property with no sale price returns list price",
			property: Property{
				ListPrice: 500000,
				Status:    "Under Contract",
			},
			expected: 500000,
		},
		{
			name: "zero price property",
			property: Property{
				ListPrice: 0,
				SalePrice: 0,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.property.GetPrice()
			if got != tt.expected {
				t.Errorf("GetPrice() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProperty_GetAgeInMonths(t *testing.T) {
	now := time.Now().Unix()
	oneMonthAgo := now - (30 * 24 * 60 * 60)
	threeMonthsAgo := now - (90 * 24 * 60 * 60)
	sixMonthsAgo := now - (180 * 24 * 60 * 60)

	tests := []struct {
		name          string
		property      Property
		expectedRange struct {
			min float64
			max float64
		}
	}{
		{
			name: "active property one month old",
			property: Property{
				Status:      "Active",
				ListingDate: oneMonthAgo,
			},
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 0.9,
				max: 1.1,
			},
		},
		{
			name: "closed property with status change three months ago",
			property: Property{
				Status:                "Closed",
				ListingDate:           sixMonthsAgo,
				StatusChangeTimestamp: threeMonthsAgo,
			},
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 2.9,
				max: 3.1,
			},
		},
		{
			name: "pending property uses listing date",
			property: Property{
				Status:      "Under Contract",
				ListingDate: threeMonthsAgo,
			},
			expectedRange: struct {
				min float64
				max float64
			}{
				min: 2.9,
				max: 3.1,
			},
		},
		{
			name: "zero timestamp property",
			property: Property{
				Status:      "Active",
				ListingDate: 0,
			},
			expectedRange: struct {
				min float64
				max float64
			}{
				min: float64(now)/(30*24*60*60) - 0.1,
				max: float64(now)/(30*24*60*60) + 0.1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.property.GetAgeInMonths()
			if got < tt.expectedRange.min || got > tt.expectedRange.max {
				t.Errorf("GetAgeInMonths() = %v, want between %v and %v",
					got, tt.expectedRange.min, tt.expectedRange.max)
			}
		})
	}
}

func TestProperty_Address(t *testing.T) {
	tests := []struct {
		name     string
		address  Address
		wantCity string
		wantZip  string
	}{
		{
			name: "complete address",
			address: Address{
				Street: "123 Main St",
				City:   "Danbury",
				State:  "CT",
				Zip:    "06810",
			},
			wantCity: "Danbury",
			wantZip:  "06810",
		},
		{
			name: "minimal address",
			address: Address{
				City: "Danbury",
			},
			wantCity: "Danbury",
			wantZip:  "",
		},
		{
			name:     "empty address",
			address:  Address{},
			wantCity: "",
			wantZip:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.address.City != tt.wantCity {
				t.Errorf("Address.City = %v, want %v", tt.address.City, tt.wantCity)
			}
			if tt.address.Zip != tt.wantZip {
				t.Errorf("Address.Zip = %v, want %v", tt.address.Zip, tt.wantZip)
			}
		})
	}
}

func TestProperty_Bathroom(t *testing.T) {
	tests := []struct {
		name      string
		bathroom  Bathroom
		wantTotal float64
		wantFull  int
		wantHalf  int
	}{
		{
			name: "standard bathroom count",
			bathroom: Bathroom{
				Total: 2.5,
				Full:  2,
				Half:  1,
			},
			wantTotal: 2.5,
			wantFull:  2,
			wantHalf:  1,
		},
		{
			name: "whole number bathrooms",
			bathroom: Bathroom{
				Total: 3.0,
				Full:  3,
				Half:  0,
			},
			wantTotal: 3.0,
			wantFull:  3,
			wantHalf:  0,
		},
		{
			name:      "zero bathrooms",
			bathroom:  Bathroom{},
			wantTotal: 0,
			wantFull:  0,
			wantHalf:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.bathroom.Total != tt.wantTotal {
				t.Errorf("Bathroom.Total = %v, want %v", tt.bathroom.Total, tt.wantTotal)
			}
			if tt.bathroom.Full != tt.wantFull {
				t.Errorf("Bathroom.Full = %v, want %v", tt.bathroom.Full, tt.wantFull)
			}
			if tt.bathroom.Half != tt.wantHalf {
				t.Errorf("Bathroom.Half = %v, want %v", tt.bathroom.Half, tt.wantHalf)
			}
		})
	}
}
