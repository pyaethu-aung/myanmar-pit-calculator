# Data Model: CLI UI State

## Entities

### `TaxCalculatorForm`
The primary state container during data collection.
- `Language`: `string` ("EN" or "MY", default "EN")
- `MonthlySalary`: `float64` (Input, > 0)
- `YearlyBonus`: `float64` (Input, >= 0)
- `SSBContribution`: `float64` (Input, >= 0, max 30,000 MMK/month)
- `HasSpouse`: `bool` (Input)
- `NumberOfChildren`: `int` (Input, >= 0)
- `ParentsLivingWith`: `int` (Input, 0-2)
- `LifeInsurancePremium`: `float64` (Input, >= 0)

### `TaxResultView`
The summary view and export data structure.
- `TotalIncome`: `float64`
- `TotalReliefs`: `float64`
- `TaxableIncome`: `float64`
- `TaxPayable`: `float64`
- `MonthlyTax`: `float64`
- `Breakdown`: List of (Tax Layer, Rate, Amount)
- `SelectedExportFormat`: `string` ("TXT", "JSON", "CSV", "None")

## Validation Rules
- Monthly Salary must be numeric and positive.
- Number of Children must be 0 or more.
- Parents Living with must be 0, 1, or 2.
- SSB is capped at 300,000 MMK per year (30,000 per month).
