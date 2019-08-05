package currencyconverter

import (
	"testing"
)

func TestCurrencyIsSupported(t *testing.T) {
	testCurrency := "EUR"
	if err := IsSupportedCurrency(testCurrency); err != nil {
		t.Error("Expected result is", testCurrency, "but got", err)
	}
}

func TestCurrencySupportedCurrencies(t *testing.T) {
	testCurrency := "EUR"

	c := Currencies{}
	c.SetSupportedCurrencies(SupportedCurrencies)
	if err := c.IsSupportedCurrency(testCurrency); err != nil {
		t.Error("Expected result is", testCurrency, "but got", err)
	}
}
