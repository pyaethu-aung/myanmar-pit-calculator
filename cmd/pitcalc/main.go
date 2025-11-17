package main

import (
	"fmt"

	"github.com/myanmar-pit-calculator/pkg/pitcalc"
)

func main() {

	fmt.Println("=====================================")
	fmt.Println("   🇲🇲 Myanmar PIT Calculator (CLI)")
	fmt.Println("=====================================")

	if err := pitcalc.CalculatePIT(); err != nil {

		fmt.Printf("Error in calculating PIT: %v\n", err)
	}
}
