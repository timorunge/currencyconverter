// Exchange rate data types and lookups.

package converter

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var errNoRatesAvailable = errors.New("no exchange rates available")

// EUR is implicit in ECB data -- not included in the XML feed.
var baseRate = Rate{
	Currency: "EUR",
	Rate:     1,
}

// ExchangeRates holds exchange rates grouped by date.
type ExchangeRates struct {
	Dates []DailyRates `xml:"Cube>Cube"`
}

// DailyRates holds all exchange rates for a single date.
type DailyRates struct {
	Date  string `xml:"time,attr" json:"date"`
	Rates []Rate `xml:"Cube" json:"rates"`
}

// Rate holds a currency code and its exchange rate relative to EUR.
type Rate struct {
	Currency string  `xml:"currency,attr" json:"currency"`
	Rate     float64 `xml:"rate,attr" json:"rate"`
}

// GetDate returns the exchange rates for a specific date.
// Pass DateLatest to get the most recent available date.
func (r *ExchangeRates) GetDate(date string) (DailyRates, error) {
	if date == DateLatest {
		if len(r.Dates) == 0 {
			return DailyRates{}, errNoRatesAvailable
		}
		return slices.MaxFunc(r.Dates, func(a, b DailyRates) int {
			return strings.Compare(a.Date, b.Date)
		}), nil
	}
	for _, d := range r.Dates {
		if d.Date == date {
			return d, nil
		}
	}
	return DailyRates{}, fmt.Errorf("cannot find exchange rates for date %q", date)
}

// getRate returns the EUR-based rate for the given currency.
// EUR itself has a rate of 1. All other currencies return their
// rate relative to EUR as published by the ECB (e.g. USD=1.1151).
func (d *DailyRates) getRate(currency string) (float64, error) {
	if currency == baseRate.Currency {
		return baseRate.Rate, nil
	}
	for _, r := range d.Rates {
		if r.Currency == currency {
			if r.Rate <= 0 {
				return 0, fmt.Errorf("invalid rate %f for currency %q", r.Rate, currency)
			}
			return r.Rate, nil
		}
	}
	return 0, fmt.Errorf("cannot find rate for currency %q", currency)
}
