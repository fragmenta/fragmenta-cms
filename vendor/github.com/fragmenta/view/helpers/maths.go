package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

// PRICES

// FIXME - move to currency type with concrete implementations per currency, as it'd be neater than funcs with multiple options.  currency.GBP.PriceToCents something like that?

// PriceToCentsString returns a price in cents as a string for use in params
func PriceToCentsString(p string) string {
	if p == "" {
		return "0" // Return 0 for blank price
	}

	return fmt.Sprintf("%d", PriceToCents(p))
}

// PriceToCents converts a price string in human friendly notation (£45 or £34.40) to a price in pence as an int64
func PriceToCents(p string) int {
	price := strings.Replace(p, "£", "", -1)
	price = strings.Replace(price, ",", "", -1) // assumed to be in thousands
	price = strings.Replace(price, " ", "", -1)

	var pennies int
	var err error
	if strings.Contains(price, ".") {
		// Split the string on . and rejoin with padded pennies
		parts := strings.Split(price, ".")

		if len(parts[1]) == 0 {
			parts[1] = "00"
		} else if len(parts[1]) == 1 {
			parts[1] = parts[1] + "0"
		}

		price = parts[0] + parts[1]

		pennies, err = strconv.Atoi(price)
	} else {
		pennies, err = strconv.Atoi(price)
		pennies = pennies * 100
	}
	if err != nil {
		fmt.Printf("Error converting price %s", price)
		pennies = 0
	}

	return pennies
}

// CentsToPrice converts a price in pence to a human friendly price including currency unit
// At present it assumes the currency is pounds, it should instead take an optional param for currency
// or not include it at all
func CentsToPrice(p int64) string {
	price := fmt.Sprintf("£%.2f", float64(p)/100.0)
	return strings.TrimSuffix(price, ".00") // remove zero pence at end if we have it
}

// CentsToPriceShort converts a price in pence to a human friendly price abreviated (no pence)
func CentsToPriceShort(p int64) string {
	if p >= 100000000000 { // If greater than £1b use b suffix
		return fmt.Sprintf("£%.2fb", float64(p)/100000000000.0)
	} else if p >= 100000000 { // If greater than £1m use m suffix
		return fmt.Sprintf("£%.2fm", float64(p)/100000000.0)
	} else if p >= 100000 { // If greater than £1k use k suffix
		return fmt.Sprintf("£%.1fk", float64(p)/100000.0)
	}
	return CentsToPrice(p)
}

// CentsToBase converts cents to the base currency unit, preserving cent display, with no currency
func CentsToBase(p int64) string {
	return fmt.Sprintf("%.2f", float64(p)/100.0)
}

// Mod returns a modulo b
func Mod(a int, b int) int {
	return a % b
}

// Add returns a + b
func Add(a int, b int) int {
	return a + b
}

// Subtract returns a - b
func Subtract(a int, b int) int {
	return a - b
}

// Odd returns true if a is odd
func Odd(a int) bool {
	return a%2 == 0
}
