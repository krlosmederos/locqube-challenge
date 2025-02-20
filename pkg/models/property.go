package models

import "time"

type Property struct {
	ID                    string   `json:"id"`
	Address               Address  `json:"address"`
	Baths                 Bathroom `json:"baths"`
	Beds                  int      `json:"beds"`
	ListPrice             float64  `json:"listPrice"`
	SalePrice             float64  `json:"salePrice,omitempty"`
	Size                  float64  `json:"size"`
	Status                string   `json:"status"`
	Style                 string   `json:"style"`
	YearBuilt             int      `json:"yearBuilt"`
	ListingDate           int64    `json:"listingDate"`
	StatusChangeTimestamp int64    `json:"statusChangeTimestamp"`
	PropertyType          string   `json:"propertyType"`
}

type Address struct {
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
	Street string `json:"street"`
}

type Bathroom struct {
	Total float64 `json:"total"`
	Full  int     `json:"full"`
	Half  int     `json:"half"`
}

// GetPrice returns the price of the property
// if the property is sold, it returns the sale price
// otherwise, it returns the list price
func (p *Property) GetPrice() float64 {
	if p.SalePrice > 0 {
		return p.SalePrice
	}
	return p.ListPrice
}

// GetAgeInMonths returns the age of the property in months
// if the property is sold, it returns the age of the property when it was sold
func (p *Property) GetAgeInMonths() float64 {
	date := p.ListingDate
	if p.Status == "Closed" {
		date = p.StatusChangeTimestamp
	}
	now := time.Now().Unix()
	ageInMonths := float64(now-date) / (30 * 24 * 60 * 60)
	return ageInMonths
}
