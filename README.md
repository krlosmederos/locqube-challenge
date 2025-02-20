# Comparative Valuation Algorithm

This project implements a comparative valuation algorithm for real estate properties based on MLS listings data. The algorithm uses a weighted scaling model to determine property values based on various factors including property type, number of bedrooms/bathrooms, size, transaction recency, and status.

## Features

- Property valuation based on multiple weighted factors:
  - Property type (20%)
  - Number of bedrooms (5%)
  - Number of bathrooms (5%)
  - Property size (10%)
  - Transaction recency (50%)
  - Property status (10%)
- Recency-based weighting:
  - Last 3 months: 100% weight
  - 3-6 months: 50% weight
  - 6-9 months: 25% weight
- Status-based adjustments:
  - Closed sales: 100% weight
  - Under contract: 60% weight
  - Active listings: 40% weight
- Minimum comparable sales requirement (default: 3)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/krlosmederos/locqube-challenge.git
```

2. Navigate to the project directory:
```bash
cd locqube-challenge
```

3. Build the project:
```bash
go build -o bin/valuation cmd/valuation/main.go
```

## Configuration

The algorithm is configured through a JSON file located at `config/application.json`. Example configuration:

```json
{
  "criteria_weights": {
    "property_type": 0.2,
    "bedrooms": 0.05,
    "bathrooms": 0.05,
    "size": 0.1,
    "recency": 0.5,
    "status": 0.1
  },
  "time_scores": {
    "three_months": 1.0,
    "six_months": 0.5,
    "nine_months": 0.25
  },
  "status_scores": {
    "sold": 1.0,
    "pending": 0.6,
    "active": 0.4
  },
  "min_sales_count": 3
}
```

## Usage

Run the valuation program:

```bash
./bin/valuation
```

The program will:
1. Load the configuration from `config/application.json`
2. Read market listings from the data file
3. Calculate and display the estimated property value

## Project Structure

```
.
├── cmd/
│   └── valuation/
│       └── main.go           # Main application entry point
├── pkg/
│   ├── algorithm/
│   │   └── valuation.go      # Core valuation algorithm
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── criteria/             # Individual scoring criteria
│   │   ├── bathrooms.go
│   │   ├── bedrooms.go
│   │   ├── propertyType.go
│   │   ├── recency.go
│   │   ├── size.go
│   │   └── status.go
│   ├── filters/
│   │   └── comparable.go     # Property filtering logic
│   └── models/
│       └── property.go       # Data models
├── config/
│   └── application.json      # Application configuration
└── data/
    └── market_listings.json  # Sample market data
```

## Testing

The project includes comprehensive test coverage for all components:

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./pkg/models -v
go test ./pkg/algorithm -v
go test ./pkg/criteria -v
```

### Test Coverage

- Property Model Tests:
  - Price calculation for different property statuses
  - Age calculation in months
  - Address and bathroom information validation

- Valuation Algorithm Tests:
  - Exact property matches
  - Similar properties with variations
  - Mixed status listings
  - Different property types
  - Edge cases (no matches, empty listings)

- Criteria Tests:
  - Individual scoring for each criterion
  - Weight calculations
  - Time-based scoring
  - Status-based scoring

## Algorithm Details

The valuation algorithm works by:

1. Filtering comparable properties based on:
   - Same city
   - Similar size (within 20%)
   - Similar bedroom count (±1)
   - Similar bathroom count (±0.5)

2. Scoring each comparable property on:
   - Property type match
   - Bedroom count similarity
   - Bathroom count similarity
   - Size similarity
   - Transaction recency
   - Listing status

3. Calculating final weights by:
   - Combining individual criteria scores
   - Applying configured weights
   - Normalizing results

4. Computing the final valuation by:
   - Weighted average of comparable property prices
   - Minimum sales requirement validation
   - Recent sales prioritization