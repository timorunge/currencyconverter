package currencyconverter

import (
	"fmt"
	"sort"
)

// SupportedCurrencies are the currencies which are supported.
var (
	SupportedCurrencies = []string{
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
		"HRK",
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
		"RUB",
		"SEK",
		"SGD",
		"THB",
		"TRY",
		"USD",
		"ZAR",
	}
)

// Currencies is the struct for the currencies.
type Currencies struct {
	SupportedCurrencies []string
}

// IsSupportedCurrency is checking if the given currency is supported.
func IsSupportedCurrency(currency string) error {
	c := &Currencies{SupportedCurrencies: SupportedCurrencies}
	return c.isSupportedCurrency(currency)
}

// IsSupportedCurrency is checking if the given currency is supported.
func (c *Currencies) IsSupportedCurrency(currency string) error {
	return c.isSupportedCurrency(currency)
}

// SetSupportedCurrencies is setting the supported currencies.
func (c *Currencies) SetSupportedCurrencies(currencies []string) {
	c.SupportedCurrencies = currencies
}

// isSupportedCurrency is checking if the given currency is supported.
func (c *Currencies) isSupportedCurrency(currency string) error {
	sort.Strings(c.SupportedCurrencies)
	i := sort.SearchStrings(c.SupportedCurrencies, currency)
	if i < len(c.SupportedCurrencies) && c.SupportedCurrencies[i] == currency {
		return nil
	}
	return fmt.Errorf("\"%s\" is not a supported currency", currency)
}
