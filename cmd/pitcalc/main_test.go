package main

import (
	"strings"
	"testing"
)

func TestCurrencyFormat(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{
			name:     "zero amount",
			amount:   0,
			expected: "0.00 MMK",
		},
		{
			name:     "small amount",
			amount:   1000,
			expected: "1,000.00 MMK",
		},
		{
			name:     "million",
			amount:   1000000,
			expected: "1,000,000.00 MMK",
		},
		{
			name:     "large amount",
			amount:   5000000.50,
			expected: "5,000,000.50 MMK",
		},
		{
			name:     "decimal amount",
			amount:   100.5,
			expected: "100.50 MMK",
		},
		{
			name:     "negative amount",
			amount:   -1000000,
			expected: "-1,000,000.00 MMK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := currencyFormat(tt.amount)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestInputValidation(t *testing.T) {
	tests := []struct {
		name          string
		value         int
		validator     func(int) *string
		shouldPass    bool
		expectedError string
	}{
		{
			name:       "valid positive income",
			value:      500000,
			validator:  validateMonthlyIncome,
			shouldPass: true,
		},
		{
			name:          "zero income",
			value:         0,
			validator:     validateMonthlyIncome,
			shouldPass:    false,
			expectedError: "❌ Monthly income must be greater than 0.",
		},
		{
			name:          "negative income",
			value:         -100000,
			validator:     validateMonthlyIncome,
			shouldPass:    false,
			expectedError: "❌ Monthly income must be greater than 0.",
		},
		{
			name:       "valid month april",
			value:      4,
			validator:  validateStartingMonth,
			shouldPass: true,
		},
		{
			name:          "month 0",
			value:         0,
			validator:     validateStartingMonth,
			shouldPass:    false,
			expectedError: "❌ Starting month must be between 1 and 12.",
		},
		{
			name:          "month 13",
			value:         13,
			validator:     validateStartingMonth,
			shouldPass:    false,
			expectedError: "❌ Starting month must be between 1 and 12.",
		},
		{
			name:       "valid one parent",
			value:      1,
			validator:  validateDependentParents,
			shouldPass: true,
		},
		{
			name:       "valid two parents",
			value:      2,
			validator:  validateDependentParents,
			shouldPass: true,
		},
		{
			name:          "three parents exceeds limit",
			value:         3,
			validator:     validateDependentParents,
			shouldPass:    false,
			expectedError: "❌ Number of dependent parents cannot exceed 2.",
		},
		{
			name:          "negative parents",
			value:         -1,
			validator:     validateDependentParents,
			shouldPass:    false,
			expectedError: "❌ Number of dependent parents cannot be negative.",
		},
		{
			name:       "spouse yes 1",
			value:      1,
			validator:  validateDependentSpouse,
			shouldPass: true,
		},
		{
			name:       "spouse no 0",
			value:      0,
			validator:  validateDependentSpouse,
			shouldPass: true,
		},
		{
			name:          "spouse invalid 2",
			value:         2,
			validator:     validateDependentSpouse,
			shouldPass:    false,
			expectedError: "❌ Invalid input. Please enter 1 for Yes or 0 for No.",
		},
		{
			name:       "three children valid",
			value:      3,
			validator:  validateChildrens,
			shouldPass: true,
		},
		{
			name:       "zero children valid",
			value:      0,
			validator:  validateChildrens,
			shouldPass: true,
		},
		{
			name:          "negative children invalid",
			value:         -1,
			validator:     validateChildrens,
			shouldPass:    false,
			expectedError: "❌ Number of children cannot be negative.",
		},
		{
			name:       "ssb 72000 valid",
			value:      72000,
			validator:  validateSSB,
			shouldPass: true,
		},
		{
			name:       "ssb 0 valid",
			value:      0,
			validator:  validateSSB,
			shouldPass: true,
		},
		{
			name:          "ssb negative invalid",
			value:         -1000,
			validator:     validateSSB,
			shouldPass:    false,
			expectedError: "❌ Yearly SSB contribution cannot be negative.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.validator(tt.value)
			if tt.shouldPass {
				if result != nil {
					t.Errorf("expected no error, got %q", *result)
				}
			} else {
				if result == nil {
					t.Errorf("expected error, got nil")
				} else if *result != tt.expectedError {
					t.Errorf("expected error %q, got %q", tt.expectedError, *result)
				}
			}
		})
	}
}

func TestCurrencyFormatContainsMMK(t *testing.T) {
	result := currencyFormat(12345.67)
	if !strings.Contains(result, "MMK") {
		t.Errorf("expected format to contain 'MMK', got %q", result)
	}
}

func TestCurrencyFormatTwoDecimalPlaces(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		expected string
	}{
		{
			name:     "whole number",
			amount:   100,
			expected: "100.00 MMK",
		},
		{
			name:     "one decimal place",
			amount:   100.5,
			expected: "100.50 MMK",
		},
		{
			name:     "two decimal places",
			amount:   100.55,
			expected: "100.55 MMK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := currencyFormat(tt.amount)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
