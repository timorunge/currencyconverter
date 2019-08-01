package currencyconverter

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestExchangeRatesAPI(t *testing.T) {
	r := &ExchangeRates{}
	if err := xml.Unmarshal(testAPIExchangeRates, &r); err != nil {
		t.Error("Can not unmarshal xml")
	}

	exchangeRates, err := r.GetExchangeRates(testExchangeRateDate)
	if err != nil {
		t.Error("Expected result is a exchnage rate date but got", err)
	}
	if exchangeRates.Date != testExchangeRateDate {
		t.Error("Expected result is", testExchangeRateDate, "but got", exchangeRates.Date)
	}

	f64, err := exchangeRates.GetRate(testTargetCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value but got", err)
	}
	if f64 != testTargetCurrencyRate {
		t.Error("Expected result is", testTargetCurrencyRate, "but got", f64)
	}

	f64, err = r.GetRate(testExchangeRateDate, testTargetCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value but got", err)
	}
	if f64 != testTargetCurrencyRate {
		t.Error("Expected result is", testTargetCurrencyRate, "but got", f64)
	}
}

func TestExchangeRatesCache(t *testing.T) {
	r := &ExchangeRates{}
	if err := json.Unmarshal(testCacheExchangeRates, &r); err != nil {
		t.Error("Can not unmarshal json")
	}

	exchangeRates, err := r.GetExchangeRates(testExchangeRateDate)
	if err != nil {
		t.Error("Expected result is a exchange rate date but got", err)
	}
	if exchangeRates.Date != testExchangeRateDate {
		t.Error("Expected result is", testExchangeRateDate, "but got", exchangeRates.Date)
	}

	f64, err := exchangeRates.GetRate(testTargetCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value but got", err)
	}
	if f64 != testTargetCurrencyRate {
		t.Error("Expected result is", testTargetCurrencyRate, "but got", f64)
	}

	f64, err = r.GetRate(testExchangeRateDate, testTargetCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value but got", err)
	}
	if f64 != testTargetCurrencyRate {
		t.Error("Expected result is", testTargetCurrencyRate, "but got", f64)
	}
}
