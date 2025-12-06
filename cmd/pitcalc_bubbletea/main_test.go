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
		{
			name:     "monthly salary example",
			amount:   500000,
			expected: "500,000.00 MMK",
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

func TestParseNumericInput(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedValue float64
		expectedError bool
	}{
		{
			name:          "simple number",
			input:         "500000",
			expectedValue: 500000,
			expectedError: false,
		},
		{
			name:          "number with comma",
			input:         "1,000,000",
			expectedValue: 1000000,
			expectedError: false,
		},
		{
			name:          "decimal number",
			input:         "100.50",
			expectedValue: 100.50,
			expectedError: false,
		},
		{
			name:          "decimal with comma",
			input:         "1,234.56",
			expectedValue: 1234.56,
			expectedError: false,
		},
		{
			name:          "empty string",
			input:         "",
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "whitespace only",
			input:         "   ",
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "with leading whitespace",
			input:         "  500000",
			expectedValue: 500000,
			expectedError: false,
		},
		{
			name:          "with trailing whitespace",
			input:         "500000  ",
			expectedValue: 500000,
			expectedError: false,
		},
		{
			name:          "invalid characters",
			input:         "abc123",
			expectedValue: 0,
			expectedError: true,
		},
		{
			name:          "negative number",
			input:         "-1000",
			expectedValue: -1000,
			expectedError: false,
		},
		{
			name:          "zero",
			input:         "0",
			expectedValue: 0,
			expectedError: false,
		},
		{
			name:          "ssb contribution",
			input:         "72000",
			expectedValue: 72000,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNumericInput(tt.input)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result, got nil")
				} else if *result != tt.expectedValue {
					t.Errorf("expected %f, got %f", tt.expectedValue, *result)
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

func TestParseNumericInputReturnsPointer(t *testing.T) {
	result, err := parseNumericInput("500000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected non-nil pointer result")
	}
}

func TestParseNumericInputLargeNumbers(t *testing.T) {
	input := "999999999.99"
	result, err := parseNumericInput(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Errorf("expected non-nil pointer result")
	}
	expected := 999999999.99
	if *result != expected {
		t.Errorf("expected %f, got %f", expected, *result)
	}
}

func TestParseNumericInputWithCommas(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{
			name:     "triple digit groups",
			input:    "1,234,567.89",
			expected: 1234567.89,
		},
		{
			name:     "double digit groups",
			input:    "1,234,567",
			expected: 1234567,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseNumericInput(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Errorf("expected non-nil pointer result")
			} else if *result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, *result)
			}
		})
	}
}
