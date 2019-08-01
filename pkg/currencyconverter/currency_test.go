package currencyconverter

import (
	"testing"
)

func TestCurrency(t *testing.T) {
	if err := IsSupportedCurrency(testBaseCurrency); err != nil {
		t.Error("Expected result is", testBaseCurrency, "but got", err)
	}

	c := Currencies{}
	c.SetSupportedCurrencies(SupportedCurrencies)

	if err := c.IsSupportedCurrency(testBaseCurrency); err != nil {
		t.Error("Expected result is", testBaseCurrency, "but got", err)
	}
}
