package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/myanmar-pit-calculator/pkg/pitcalc"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {

	fmt.Println("=====================================")
	fmt.Println("   🇲🇲 Myanmar PIT Calculator (CLI)")
	fmt.Println("=====================================")

	monthlyIncome := inputInt(
		"Enter monthly income (MMK): ",
		validateMonthlyIncome,
	)

	startingMonth := inputInt(
		"Enter starting month (1 = Jan, 2 = Feb, ..., 12 = Dec): ",
		validateStartingMonth,
	)

	dependentParents := inputInt(
		"Enter number of dependent parents (1,000,000 MMK for each): ",
		validateDependentParents,
	)

	dependentSpouse := inputInt(
		"Do you have a dependent spouse? (1 = Yes, 0 = No): ",
		validateDependentSpouse,
	)

	childrens := inputInt(
		"Enter number of children (500,000 MMK for each): ",
		validateChildrens,
	)

	ssb := inputInt(
		"Enter total SSB contribution yearly (MMK): ",
		validateSSB,
	)

	output, err := pitcalc.CalculatePIT(
		pitcalc.CalculatePITInput{
			MonthlyIncome:    float64(monthlyIncome),
			StartingMonth:    startingMonth,
			DependentParents: dependentParents,
			DependentSpouse:  dependentSpouse,
			Childrens:        childrens,
			SSB:              float64(ssb),
		})
	if err != nil {

		fmt.Printf("Error in calculating PIT: %v\n", err)
	}
	fmt.Println("=====================================")
	fmt.Printf(
		"Total Taxable Income: %s\n", currencyFormat(output.TotalTexable))
	fmt.Printf("Total Reliefs: %s\n", currencyFormat(output.TotalRelief))
	fmt.Printf("Total Personal Income Tax: %s\n", currencyFormat(output.TotalTax))
	sort.Slice(output.TaxBreakdown, func(i, j int) bool {

		return output.TaxBreakdown[i].Start < output.TaxBreakdown[j].Start
	})
	for _, v := range output.TaxBreakdown {

		if v.Limit == math.Inf(1) {

			fmt.Printf(
				"  Above from %s: %s\n",
				currencyFormat(v.Start),
				currencyFormat(v.Amount))
		} else {

			fmt.Printf(
				"  Up to %s: %s\n",
				currencyFormat(v.Limit),
				currencyFormat(v.Amount))
		}
	}
	fmt.Println("=====================================")
}

func inputInt(prompt string, validate func(int) *string) int64 {

	errMessage := "❌ Invalid input, try again."

	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print(prompt)
		text, _ := reader.ReadString('\n')
		value, err := strconv.Atoi(strings.TrimSpace(text))
		validationErrMessage := validate(value)
		if err == nil && validationErrMessage == nil {

			return int64(value)
		} else if validationErrMessage != nil {

			errMessage = *validationErrMessage
		}
		fmt.Println(errMessage)
	}
}

func currencyFormat(amount float64) string {

	return message.NewPrinter(language.English).Sprintf("%.2f MMK", amount)
}

func validateMonthlyIncome(value int) *string {
	if value <= 0 {
		errMessage := "❌ Monthly income must be greater than 0."
		return &errMessage
	}
	return nil
}

func validateStartingMonth(value int) *string {
	if value < 1 || value > 12 {
		errMessage := "❌ Starting month must be between 1 and 12."
		return &errMessage
	}
	return nil
}

func validateDependentParents(value int) *string {
	if value < 0 {
		errMessage := "❌ Number of dependent parents cannot be negative."
		return &errMessage
	}
	if value > 2 {
		errMessage := "❌ Number of dependent parents cannot exceed 2."
		return &errMessage
	}
	return nil
}

func validateDependentSpouse(value int) *string {
	if value != 0 && value != 1 {
		errMessage := "❌ Invalid input. Please enter 1 for Yes or 0 for No."
		return &errMessage
	}
	return nil
}

func validateChildrens(value int) *string {
	if value < 0 {
		errMessage := "❌ Number of children cannot be negative."
		return &errMessage
	}
	return nil
}

func validateSSB(value int) *string {
	if value < 0 {
		errMessage := "❌ Yearly SSB contribution cannot be negative."
		return &errMessage
	}
	return nil
}
