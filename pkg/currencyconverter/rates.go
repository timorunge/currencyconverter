package currencyconverter

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// baseRate is the base rate for all calculations.
var (
	baseRate = Rate{
		Currency: "EUR",
		Rate:     float64(1),
	}
)

// ExchangeRates is holding the exchange rates on a daily basis.
type ExchangeRates struct {
	Dates []Date `xml:"Cube>Cube"`
}

// Date is holding all exchange rates on a date.
type Date struct {
	Date  string `xml:"time,attr"`
	Rates []Rate `xml:"Cube"`
}

// Rate is holding the currency and the exchange rate.
type Rate struct {
	Currency string  `xml:"currency,attr"`
	Rate     float64 `xml:"rate,attr"`
}

// NewExchangeRates is returning a new exchange rate struct.
func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{}
}

// AddRate is adding a rate struct to a date.
func (d *Date) AddRate(rate Rate) *Date {
	d.Rates = append(d.Rates, rate)
	return d
}

// GetRate is returning one exchange rate for a currency in a date context.
func (d *Date) GetRate(currency string) (float64, error) {
	if currency == baseRate.Currency {
		return baseRate.Rate, nil
	}
	return d.getRate(currency)
}

// AddDate is adding a date struct to the exchange rate struct.
func (r *ExchangeRates) AddDate(date Date) *ExchangeRates {
	r.Dates = append(r.Dates, date)
	return r
}

// GetDate is returning all exchange rates at a date.
func (r *ExchangeRates) GetDate(date string) (Date, error) {
	return r.getDate(date)
}

// GetRate is returning one exchange rate at for a currency in a exchange rate context.
func (r *ExchangeRates) GetRate(date string, currency string) (float64, error) {
	if currency == baseRate.Currency {
		return baseRate.Rate, nil
	}
	dateExchangeRates, err := r.getDate(date)
	if err != nil {
		return float64(0), err
	}
	return dateExchangeRates.getRate(currency)
}

// getDate is retuning all exchange rates at a date.
func (r *ExchangeRates) getDate(date string) (Date, error) {
	if err := IsValidDate(date); err != nil {
		return Date{}, err
	}
	sort.Slice(r.Dates, func(i, j int) bool { return r.Dates[i].Date < r.Dates[j].Date })
	if date == "latest" {
		return r.Dates[len(r.Dates)-1], nil
	}
	i := sort.Search(len(r.Dates), func(i int) bool { return r.Dates[i].Date >= date })
	if i < len(r.Dates) && r.Dates[i].Date == date {
		return r.Dates[i], nil
	}
	return Date{}, fmt.Errorf("Can not find exchange rates for date \"%s\"", date)
}

// getRate is returning one exchange rate at a given date.
func (d *Date) getRate(currency string) (float64, error) {
	if currency == baseRate.Currency {
		return baseRate.Rate, nil
	}
	f64 := float64(0)
	sort.Slice(d.Rates, func(i, j int) bool { return d.Rates[i].Currency < d.Rates[j].Currency })
	i := sort.Search(len(d.Rates), func(i int) bool { return d.Rates[i].Currency >= currency })
	if i < len(d.Rates) && d.Rates[i].Currency == currency {
		if v := reflect.ValueOf(d.Rates[i].Rate).Kind(); v == reflect.Float64 {
			return baseRate.Rate / d.Rates[i].Rate, nil
		}
		return f64, errors.New("Rate is no float64 value")
	}
	return f64, fmt.Errorf("Can not find rate for \"%s\"", currency)
}
