package pitcalc

import (
	"fmt"
	"math"
)

// TaxBracket represents a tax bracket with an upper limit and a tax rate.
type TaxBracket struct {
	Start float64
	Limit float64
	Rate  float64
}

// CalculatePITInput holds the input parameters for calculating personal income
// tax.
type CalculatePITInput struct {
	MonthlyIncome    float64
	StartingMonth    int64
	DependentParents int64
	DependentSpouse  int64
	Childrens        int64
	SSB              float64
}

// CalculatePITOutput holds the output results from calculating personal income
// tax.
type CalculatePITOutput struct {
	TaxBreakdown []struct {
		Start  float64
		Limit  float64
		Rate   float64
		Amount float64
	}
	TotalRelief  float64
	TotalTexable float64
	TotalTax     float64
}

// CalculatePIT computes personal income tax for Myanmar.
func CalculatePIT(input CalculatePITInput) (*CalculatePITOutput, error) {

	// Validate input
	if input.MonthlyIncome <= 0 {

		return nil, fmt.Errorf("monthly income must be greater than 0")
	}
	if input.StartingMonth < 1 || input.StartingMonth > 12 {
		return nil, fmt.Errorf("starting month must be between 1 and 12")
	}
	if input.DependentParents < 0 {
		return nil, fmt.Errorf("number of dependent parents cannot be negative")
	}
	if input.DependentParents > 2 {
		return nil, fmt.Errorf("number of dependent parents cannot exceed 2")
	}
	if input.DependentSpouse < 0 || input.DependentSpouse > 1 {
		return nil, fmt.Errorf("dependent spouse value must be 0 or 1")
	}
	if input.Childrens < 0 {
		return nil, fmt.Errorf("number of children cannot be negative")
	}
	if input.SSB < 0 {
		return nil, fmt.Errorf("yearly SSB contribution cannot be negative")
	}

	// Determine months in the budget year (April=4, ..., March=3)
	var months int64
	if input.StartingMonth >= 4 {

		// Apr(4)->12, Dec(12)->4
		months = 16 - input.StartingMonth
	} else {

		// Jan(1)->3, Mar(3)->1
		months = 4 - input.StartingMonth
	}

	yearlyGrossIncome := input.MonthlyIncome * float64(months)

	// Reliefs
	personalRelief := 0.2 * float64(yearlyGrossIncome)
	if personalRelief > 10000000 {

		personalRelief = 10000000
	}
	parentRelief := float64(input.DependentParents) * 1000000
	spouseRelief := float64(input.DependentSpouse) * 1000000
	childRelief := float64(input.Childrens) * 500000.0
	totalRelief := personalRelief + parentRelief + spouseRelief + childRelief + input.SSB

	taxableIncome := yearlyGrossIncome - totalRelief
	if taxableIncome < 0 {

		taxableIncome = 0
	}

	output := CalculatePITOutput{
		TotalRelief:  totalRelief,
		TotalTexable: taxableIncome,
	}

	// Calculate tax per bracket
	output.TaxBreakdown = make([]struct {
		Start  float64
		Limit  float64
		Rate   float64
		Amount float64
	}, 0)

	remaining := taxableIncome
	previousLimit := 0.0
	for _, bracket := range brackets {

		if remaining <= 0 {

			break
		}

		upper := bracket.Limit - previousLimit
		part := math.Min(remaining, upper)

		tax := part * bracket.Rate

		output.TaxBreakdown = append(output.TaxBreakdown, struct {
			Start  float64
			Limit  float64
			Rate   float64
			Amount float64
		}{
			Start:  bracket.Start,
			Limit:  bracket.Limit,
			Rate:   bracket.Rate,
			Amount: tax,
		})
		output.TotalTax += tax
		remaining -= part
		previousLimit = bracket.Limit
	}

	return &output, nil
}

var brackets = []TaxBracket{

	{1, 2000000, 0.00},
	{2000001, 10000000, 0.05},
	{10000001, 30000000, 0.10},
	{30000001, 50000000, 0.15},
	{50000001, 70000000, 0.20},
	{70000001, math.Inf(1), 0.25},
}
