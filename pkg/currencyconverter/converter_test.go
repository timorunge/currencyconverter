package currencyconverter

import (
	"encoding/json"
	"testing"
)

func TestConverter(t *testing.T) {
	testConverterAmount := float64(1337)
	testConverterBaseCurrency := "EUR"
	testConverterExchangeRate := float64(1.1151)
	testConverterExchangeRateDate := "2019-07-31"
	testConverterRates := []byte(`{"Dates":[{"Date":"2019-07-31","Rates":[{"Currency":"USD","Rate":1.1151},{"Currency":"JPY","Rate":121.04},{"Currency":"BGN","Rate":1.9558},{"Currency":"CZK","Rate":25.658},{"Currency":"DKK","Rate":7.4674},{"Currency":"GBP","Rate":0.91623},{"Currency":"HUF","Rate":326.48},{"Currency":"PLN","Rate":4.2912},{"Currency":"RON","Rate":4.7338},{"Currency":"SEK","Rate":10.6645},{"Currency":"CHF","Rate":1.1041},{"Currency":"ISK","Rate":134.7},{"Currency":"NOK","Rate":9.7778},{"Currency":"HRK","Rate":7.3823},{"Currency":"RUB","Rate":70.8041},{"Currency":"TRY","Rate":6.161},{"Currency":"AUD","Rate":1.6175},{"Currency":"BRL","Rate":4.218},{"Currency":"CAD","Rate":1.4662},{"Currency":"CNY","Rate":7.6743},{"Currency":"HKD","Rate":8.7289},{"Currency":"IDR","Rate":15639.3},{"Currency":"ILS","Rate":3.8951},{"Currency":"INR","Rate":76.6965},{"Currency":"KRW","Rate":1318.25},{"Currency":"MXN","Rate":21.2005},{"Currency":"MYR","Rate":4.6015},{"Currency":"NZD","Rate":1.6883},{"Currency":"PHP","Rate":56.685},{"Currency":"SGD","Rate":1.5261},{"Currency":"THB","Rate":34.273},{"Currency":"ZAR","Rate":15.8634}]}]}`)
	testConverterResult := float64(1490.8887)
	testConverterTargetCurrency := "USD"

	currencies := &Currencies{}
	currencies.SetSupportedCurrencies(SupportedCurrencies)

	if err := IsValidAmount(testConverterAmount); err != nil {
		t.Error("Expected result is nil but got", err)
	}

	rates := &ExchangeRates{}
	json.Unmarshal(testConverterRates, rates)

	converter := NewConverter()
	converter.SetExchangeRates(*rates)

	converter.SetAmount(testConverterAmount)
	if converter.Amount != testConverterAmount {
		t.Error("Expected result is", testConverterAmount, "but got", converter.Amount)
	}

	if err := converter.IsValidAmount(); err != nil {
		t.Error("Expected result is nil but got", err)
	}

	converter.SetBaseCurrency(testConverterBaseCurrency)
	if converter.BaseCurrency != testConverterBaseCurrency {
		t.Error("Expected result is", testConverterBaseCurrency, "but got", converter.BaseCurrency)
	}

	converter.SetCurrencies(*currencies)
	if len(converter.Currencies.SupportedCurrencies) != len(SupportedCurrencies) {
		t.Error("Expected result is", len(SupportedCurrencies), "but got", len(converter.Currencies.SupportedCurrencies))
	}

	converter.SetDate(testConverterExchangeRateDate)
	if converter.Date != testConverterExchangeRateDate {
		t.Error("Expected result is", testConverterExchangeRateDate, "but got", converter.Date)
	}

	if err := converter.IsValidDate(); err != nil {
		t.Error("Expected result is nil but got", err)
	}

	converter.SetTargetCurrency(testConverterTargetCurrency)
	if converter.TargetCurrency != testConverterTargetCurrency {
		t.Error("Expected result is", testConverterTargetCurrency, "but got", converter.TargetCurrency)
	}

	f64, err := converter.ExchangeRate()
	if err != nil {
		t.Error(err)
	}
	if f64 != testConverterExchangeRate {
		t.Error("Expected result is", testConverterExchangeRate, "but got", f64)
	}

	f64, err = converter.Convert()
	if err != nil {
		t.Error(err)
	}
	if f64 != testConverterResult {
		t.Error("Expected result is", testConverterResult, "but got", f64)
	}
}
