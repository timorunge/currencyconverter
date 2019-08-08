package currencyconverter

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

func TestRatesDate(t *testing.T) {
	testRate := Rate{
		Currency: "EUR",
		Rate:     float64(1),
	}

	d := &Date{}
	d.AddRate(testRate)

	f64, err := d.GetRate(testRate.Currency)
	if err != nil {
		t.Error("Expected result is a float64 value", err)
	}

	if f64 != testRate.Rate {
		t.Error("Expected result is", testRate.Rate, "but got", f64)
	}

	r := d.ToNewExchangeRates()
	f64, err = r.Dates[0].GetRate(testRate.Currency)
	if err != nil {
		t.Error("Expected result is a float64 value", err)
	}

	if f64 != testRate.Rate {
		t.Error("Expected result is", testRate.Rate, "but got", f64)
	}
}

func TestRatesExchangeRatesJSON(t *testing.T) {
	testRateCurrency := "USD"
	testRateDate := "2019-07-31"
	testRateResult := float64(0.896780557797507)
	testRatesDate := []byte(`{"Date":"2019-07-31","Rates":[{"Currency":"USD","Rate":1.1151},{"Currency":"JPY","Rate":121.04},{"Currency":"BGN","Rate":1.9558},{"Currency":"CZK","Rate":25.658},{"Currency":"DKK","Rate":7.4674},{"Currency":"GBP","Rate":0.91623},{"Currency":"HUF","Rate":326.48},{"Currency":"PLN","Rate":4.2912},{"Currency":"RON","Rate":4.7338},{"Currency":"SEK","Rate":10.6645},{"Currency":"CHF","Rate":1.1041},{"Currency":"ISK","Rate":134.7},{"Currency":"NOK","Rate":9.7778},{"Currency":"HRK","Rate":7.3823},{"Currency":"RUB","Rate":70.8041},{"Currency":"TRY","Rate":6.161},{"Currency":"AUD","Rate":1.6175},{"Currency":"BRL","Rate":4.218},{"Currency":"CAD","Rate":1.4662},{"Currency":"CNY","Rate":7.6743},{"Currency":"HKD","Rate":8.7289},{"Currency":"IDR","Rate":15639.3},{"Currency":"ILS","Rate":3.8951},{"Currency":"INR","Rate":76.6965},{"Currency":"KRW","Rate":1318.25},{"Currency":"MXN","Rate":21.2005},{"Currency":"MYR","Rate":4.6015},{"Currency":"NZD","Rate":1.6883},{"Currency":"PHP","Rate":56.685},{"Currency":"SGD","Rate":1.5261},{"Currency":"THB","Rate":34.273},{"Currency":"ZAR","Rate":15.8634}]}`)

	d := &Date{}
	if err := json.Unmarshal(testRatesDate, &d); err != nil {
		t.Error("Can not unmarshal json")
	}
	json.Unmarshal(testRatesDate, d)

	f64, err := d.GetRate(testRateCurrency)
	if err != nil {
		t.Error("Expected result is a date object but got", err)
	}
	if f64 != testRateResult {
		t.Error("Expected result is", testRateDate, "but got", f64)
	}

	r := NewExchangeRates()
	r.AddDate(*d)

	dateObj, err := r.GetDate(testRateDate)
	if err != nil {
		t.Error("Expected result is a date object but got", err)
	}
	if dateObj.Date != testRateDate {
		t.Error("Expected result is", testRateDate, "but got", dateObj.Date)
	}

	f64, err = r.GetRate(testRateDate, testRateCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value", err)
	}

	if f64 != testRateResult {
		t.Error("Expected result is", testRateResult, "but got", f64)
	}
}

func TestRatesExchangeRatesXML(t *testing.T) {
	testRateCurrency := "USD"
	testRateDate := "2019-07-31"
	testRateResult := float64(0.896780557797507)
	testRatesExchangeRates := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender>
		<gesmes:name>European Central Bank</gesmes:name>
	</gesmes:Sender>
	<Cube>
		<Cube time='2019-07-31'>
			<Cube currency='USD' rate='1.1151'/>
			<Cube currency='JPY' rate='121.04'/>
			<Cube currency='BGN' rate='1.9558'/>
			<Cube currency='CZK' rate='25.658'/>
			<Cube currency='DKK' rate='7.4674'/>
			<Cube currency='GBP' rate='0.91623'/>
			<Cube currency='HUF' rate='326.48'/>
			<Cube currency='PLN' rate='4.2912'/>
			<Cube currency='RON' rate='4.7338'/>
			<Cube currency='SEK' rate='10.6645'/>
			<Cube currency='CHF' rate='1.1041'/>
			<Cube currency='ISK' rate='134.70'/>
			<Cube currency='NOK' rate='9.7778'/>
			<Cube currency='HRK' rate='7.3823'/>
			<Cube currency='RUB' rate='70.8041'/>
			<Cube currency='TRY' rate='6.1610'/>
			<Cube currency='AUD' rate='1.6175'/>
			<Cube currency='BRL' rate='4.2180'/>
			<Cube currency='CAD' rate='1.4662'/>
			<Cube currency='CNY' rate='7.6743'/>
			<Cube currency='HKD' rate='8.7289'/>
			<Cube currency='IDR' rate='15639.30'/>
			<Cube currency='ILS' rate='3.8951'/>
			<Cube currency='INR' rate='76.6965'/>
			<Cube currency='KRW' rate='1318.25'/>
			<Cube currency='MXN' rate='21.2005'/>
			<Cube currency='MYR' rate='4.6015'/>
			<Cube currency='NZD' rate='1.6883'/>
			<Cube currency='PHP' rate='56.685'/>
			<Cube currency='SGD' rate='1.5261'/>
			<Cube currency='THB' rate='34.273'/>
			<Cube currency='ZAR' rate='15.8634'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`)

	r := &ExchangeRates{}
	if err := xml.Unmarshal(testRatesExchangeRates, &r); err != nil {
		t.Error("Can not unmarshal xml")
	}

	dateObj, err := r.GetDate(testRateDate)
	if err != nil {
		t.Error("Expected result is a date object but got", err)
	}

	if dateObj.Date != testRateDate {
		t.Error("Expected result is", testRateDate, "but got", dateObj.Date)
	}

	f64, err := r.GetRate(testRateDate, testRateCurrency)
	if err != nil {
		t.Error("Expected result is a float64 value", err)
	}

	if f64 != testRateResult {
		t.Error("Expected result is", testRateResult, "but got", f64)
	}
}
