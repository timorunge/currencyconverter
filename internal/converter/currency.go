// Supported currency list and validation.

package converter

import (
	"fmt"
	"slices"
)

// supportedCurrencies lists all currencies supported by the ECB feed.
var supportedCurrencies = []string{
	"AUD",
	"BGN",
	"BRL",
	"CAD",
	"CHF",
	"CNY",
	"CZK",
	"DKK",
	"EUR",
	"GBP",
	"HKD",
	"HUF",
	"IDR",
	"ILS",
	"INR",
	"ISK",
	"JPY",
	"KRW",
	"MXN",
	"MYR",
	"NOK",
	"NZD",
	"PHP",
	"PLN",
	"RON",
	"SEK",
	"SGD",
	"THB",
	"TRY",
	"USD",
	"ZAR",
}

// SupportedCurrencies returns a copy of the supported currency list.
func SupportedCurrencies() []string {
	return slices.Clone(supportedCurrencies)
}

// ValidateCurrency checks if the given currency is in the supported list.
func ValidateCurrency(currency string) error {
	if slices.Contains(supportedCurrencies, currency) {
		return nil
	}
	return fmt.Errorf("%q is not a supported currency", currency)
}
