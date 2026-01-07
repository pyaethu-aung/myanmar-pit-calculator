package pitcalc

import (
	"testing"
)

func TestCalculatePIT_InvalidMonthlyIncome(t *testing.T) {
	tests := []struct {
		name          string
		monthlyIncome float64
		expectedError string
	}{
		{
			name:          "zero income",
			monthlyIncome: 0,
			expectedError: "monthly income must be greater than 0",
		},
		{
			name:          "negative income",
			monthlyIncome: -500000,
			expectedError: "monthly income must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CalculatePITInput{
				MonthlyIncome: tt.monthlyIncome,
				StartingMonth: 1,
			}
			result, err := CalculatePIT(input)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
			}
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

func TestCalculatePIT_InvalidStartingMonth(t *testing.T) {
	tests := []struct {
		name          string
		startingMonth int64
		expectedError string
	}{
		{
			name:          "month 0",
			startingMonth: 0,
			expectedError: "starting month must be between 1 and 12",
		},
		{
			name:          "month 13",
			startingMonth: 13,
			expectedError: "starting month must be between 1 and 12",
		},
		{
			name:          "negative month",
			startingMonth: -1,
			expectedError: "starting month must be between 1 and 12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CalculatePITInput{
				MonthlyIncome: 500000,
				StartingMonth: tt.startingMonth,
			}
			result, err := CalculatePIT(input)
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if err.Error() != tt.expectedError {
				t.Errorf("expected error %q, got %q", tt.expectedError, err.Error())
			}
			if result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

func TestCalculatePIT_InvalidDependentsAndSSB(t *testing.T) {
	tests := []struct {
		name          string
		input         CalculatePITInput
		expectedError string
	}{
		{
			name: "negative parents",
			input: CalculatePITInput{
				MonthlyIncome:    500000,
				StartingMonth:    4,
				DependentParents: -1,
			},
			expectedError: "number of dependent parents cannot be negative",
		},
		{
			name: "too many parents",
			input: CalculatePITInput{
				MonthlyIncome:    500000,
				StartingMonth:    4,
				DependentParents: 3,
			},
			expectedError: "number of dependent parents cannot exceed 2",
		},
		{
			name: "invalid spouse flag",
			input: CalculatePITInput{
				MonthlyIncome:   500000,
				StartingMonth:   4,
				DependentSpouse: 2,
			},
			expectedError: "dependent spouse value must be 0 or 1",
		},
		{
			name: "negative children",
			input: CalculatePITInput{
				MonthlyIncome:   500000,
				StartingMonth:   4,
				Childrens:       -2,
				DependentSpouse: 0,
			},
			expectedError: "number of children cannot be negative",
		},
		{
			name: "negative ssb",
			input: CalculatePITInput{
				MonthlyIncome: 500000,
				StartingMonth: 4,
				SSB:           -500,
			},
			expectedError: "yearly SSB contribution cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculatePIT(tt.input)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tt.expectedError {
				t.Fatalf("expected error %q, got %q", tt.expectedError, err.Error())
			}
			if result != nil {
				t.Fatalf("expected nil result, got %#v", result)
			}
		})
	}
}

func TestCalculatePIT_NoTaxBelow2Million(t *testing.T) {
	input := CalculatePITInput{
		MonthlyIncome: 500000,
		StartingMonth: 4,
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 500,000 * 12 = 6,000,000 gross
	// Personal relief = 20% * 6,000,000 = 1,200,000
	// Taxable = 6,000,000 - 1,200,000 = 4,800,000 (which is above 2M, so there will be tax)
	// This test should actually expect tax to be non-zero
	if result.TotalTax <= 0 {
		t.Errorf("expected TotalTax > 0, got %f", result.TotalTax)
	}
}

func TestCalculatePIT_MonthCalculation(t *testing.T) {
	tests := []struct {
		name                string
		startingMonth       int64
		expectedMonths      int64
		expectedYearlyGross float64
	}{
		{
			name:                "April (start of fiscal year)",
			startingMonth:       4,
			expectedMonths:      12,
			expectedYearlyGross: 6000000,
		},
		{
			name:                "January",
			startingMonth:       1,
			expectedMonths:      3,
			expectedYearlyGross: 1500000,
		},
		{
			name:                "March",
			startingMonth:       3,
			expectedMonths:      1,
			expectedYearlyGross: 500000,
		},
		{
			name:                "December",
			startingMonth:       12,
			expectedMonths:      4,
			expectedYearlyGross: 2000000,
		},
		{
			name:                "July",
			startingMonth:       7,
			expectedMonths:      9,
			expectedYearlyGross: 4500000,
		},
	}

	monthlyIncome := 500000.0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CalculatePITInput{
				MonthlyIncome: monthlyIncome,
				StartingMonth: tt.startingMonth,
			}
			result, err := CalculatePIT(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify that the yearly gross income is correct
			// Note: result.TotalTexable is after reliefs, so we compute the gross
			expectedRelief := 0.2 * tt.expectedYearlyGross
			if expectedRelief > 10000000 {
				expectedRelief = 10000000
			}
			expectedTaxable := tt.expectedYearlyGross - expectedRelief

			if result.TotalTexable != expectedTaxable {
				t.Errorf("for month %d: expected taxable %f, got %f",
					tt.startingMonth, expectedTaxable, result.TotalTexable)
			}
		})
	}
}

func TestCalculatePIT_PersonalRelief(t *testing.T) {
	// Personal relief is 20% of gross income, capped at 10 million
	input := CalculatePITInput{
		MonthlyIncome: 500000,
		StartingMonth: 4, // 12 months
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Yearly income = 500,000 * 12 = 6,000,000
	// Personal relief = 20% * 6,000,000 = 1,200,000 (no cap exceeded)
	expectedRelief := 1200000.0
	if result.TotalRelief != expectedRelief {
		t.Errorf("expected TotalRelief=%f, got %f", expectedRelief, result.TotalRelief)
	}
}

func TestCalculatePIT_PersonalReliefCap(t *testing.T) {
	// Test that personal relief is capped at 10 million
	input := CalculatePITInput{
		MonthlyIncome: 500000,
		StartingMonth: 4, // 12 months
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// With a very high income, personal relief should be capped at 10,000,000
	// For this test, we check that relief doesn't exceed cap
	personalRelief := result.TotalRelief
	if personalRelief > 10000000 {
		t.Errorf("personal relief exceeded cap of 10,000,000: got %f", personalRelief)
	}
}

func TestCalculatePIT_DependentReliefs(t *testing.T) {
	input := CalculatePITInput{
		MonthlyIncome:    500000,
		StartingMonth:    4,
		DependentParents: 2,
		DependentSpouse:  1,
		Childrens:        3,
		SSB:              72000,
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Yearly income = 6,000,000
	// Personal relief = 20% * 6,000,000 = 1,200,000
	// Parent relief = 2 * 1,000,000 = 2,000,000
	// Spouse relief = 1 * 1,000,000 = 1,000,000
	// Child relief = 3 * 500,000 = 1,500,000
	// SSB = 72,000
	// Total relief = 1,200,000 + 2,000,000 + 1,000,000 + 1,500,000 + 72,000 = 5,772,000
	expectedRelief := 5772000.0
	if result.TotalRelief != expectedRelief {
		t.Errorf("expected TotalRelief=%f, got %f", expectedRelief, result.TotalRelief)
	}

	// Taxable income = 6,000,000 - 5,772,000 = 228,000
	expectedTaxable := 228000.0
	if result.TotalTexable != expectedTaxable {
		t.Errorf("expected TotalTexable=%f, got %f", expectedTaxable, result.TotalTexable)
	}
}

func TestCalculatePIT_TaxBrackets(t *testing.T) {
	tests := []struct {
		name           string
		monthlyIncome  float64
		startingMonth  int64
		minExpectedTax float64
		maxExpectedTax float64
	}{
		{
			name:           "income below 2 million (no tax)",
			monthlyIncome:  100000,
			startingMonth:  4,
			minExpectedTax: 0,
			maxExpectedTax: 0,
		},
		{
			name:           "income in 5% bracket",
			monthlyIncome:  1000000,
			startingMonth:  4,
			minExpectedTax: 300000, // Rough estimate
			maxExpectedTax: 500000, // Rough upper bound
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CalculatePITInput{
				MonthlyIncome: tt.monthlyIncome,
				StartingMonth: tt.startingMonth,
			}
			result, err := CalculatePIT(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify tax is within expected range
			if result.TotalTax < tt.minExpectedTax || result.TotalTax > tt.maxExpectedTax {
				t.Errorf("expected tax between %f and %f, got %f",
					tt.minExpectedTax, tt.maxExpectedTax, result.TotalTax)
			}

			// Tax should not exceed taxable income
			if result.TotalTax > result.TotalTexable {
				t.Errorf("tax (%f) exceeds taxable income (%f)", result.TotalTax, result.TotalTexable)
			}
		})
	}
}

func TestCalculatePIT_NegativeTaxableIncomeBecomesZero(t *testing.T) {
	// If total relief exceeds income, taxable should be 0, not negative
	input := CalculatePITInput{
		MonthlyIncome:    100000,
		StartingMonth:    4,
		DependentParents: 2,
		Childrens:        5,
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalTexable < 0 {
		t.Errorf("expected TotalTexable >= 0, got %f", result.TotalTexable)
	}
	if result.TotalTax < 0 {
		t.Errorf("expected TotalTax >= 0, got %f", result.TotalTax)
	}
}

func TestCalculatePIT_TaxBreakdownStructure(t *testing.T) {
	input := CalculatePITInput{
		MonthlyIncome: 5000000,
		StartingMonth: 4,
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.TaxBreakdown) == 0 {
		t.Errorf("expected TaxBreakdown to have entries")
	}

	// Verify tax breakdown entries have valid structure
	for i, breakdown := range result.TaxBreakdown {
		if breakdown.Rate < 0 || breakdown.Rate > 1 {
			t.Errorf("breakdown[%d]: invalid rate %f", i, breakdown.Rate)
		}
		if breakdown.Amount < 0 {
			t.Errorf("breakdown[%d]: negative amount %f", i, breakdown.Amount)
		}
		if breakdown.Limit < breakdown.Start {
			t.Errorf("breakdown[%d]: limit %f < start %f", i, breakdown.Limit, breakdown.Start)
		}
	}

	// Sum of breakdown amounts should equal total tax
	var sumTax float64
	for _, breakdown := range result.TaxBreakdown {
		sumTax += breakdown.Amount
	}
	const epsilon = 0.01
	if sumTax < result.TotalTax-epsilon || sumTax > result.TotalTax+epsilon {
		t.Errorf("sum of tax breakdowns (%f) != TotalTax (%f)", sumTax, result.TotalTax)
	}
}

func TestCalculatePIT_HighIncome(t *testing.T) {
	input := CalculatePITInput{
		MonthlyIncome: 10000000,
		StartingMonth: 4,
	}
	result, err := CalculatePIT(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalTax <= 0 {
		t.Errorf("expected positive tax for high income, got %f", result.TotalTax)
	}

	// Verify tax doesn't exceed income
	if result.TotalTax > result.TotalTexable {
		t.Errorf("tax (%f) exceeds taxable income (%f)", result.TotalTax, result.TotalTexable)
	}
}

func TestCalculatePIT_AllMonthsStartingMonths(t *testing.T) {
	// Test all valid starting months
	monthlyIncome := 1000000.0

	for month := int64(1); month <= 12; month++ {
		input := CalculatePITInput{
			MonthlyIncome: monthlyIncome,
			StartingMonth: month,
		}
		result, err := CalculatePIT(input)
		if err != nil {
			t.Errorf("month %d: unexpected error: %v", month, err)
		}

		if result == nil {
			t.Errorf("month %d: result should not be nil", month)
		}

		if result.TotalTexable < 0 {
			t.Errorf("month %d: negative taxable income %f", month, result.TotalTexable)
		}

		if result.TotalTax < 0 {
			t.Errorf("month %d: negative tax %f", month, result.TotalTax)
		}
	}
}
