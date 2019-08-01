package currencyconverter

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
)

// DateFormat is the format of the date.
// DateFormatRegex is the regex for the date format.
const (
	DateFormat      = "latest|YYYY-MM-DD"
	DateFormatRegex = "^latest|(1999|20[0-9]{2})-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])$"
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

// IsValidDate is checking if the date is in the right format.
func IsValidDate(date string) error {
	return isValidDate(date)
}

// IsValidDate is checking if the date is in the right format.
func (d *Date) IsValidDate() error {
	return d.isValidDate()
}

// GetRate is returning one exchange rate for a currency in a date context.
func (d *Date) GetRate(currency string) (float64, error) {
	if currency == baseRate.Currency {
		return baseRate.Rate, nil
	}
	return d.getRate(currency)
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

// GetExchangeRates is returning all exchange rates at a date.
func (r *ExchangeRates) GetExchangeRates(date string) (Date, error) {
	return r.getDate(date)
}

// isValidDate is checking if the date is in the right format.
func isValidDate(date string) error {
	d := &Date{Date: date}
	return d.isValidDate()
}

// isValidDate is checking if the date is in the right format.
func (d *Date) isValidDate() error {
	if !regexp.MustCompile(DateFormatRegex).MatchString(d.Date) {
		return fmt.Errorf("\"%s\" is no valid date in format \"%s\" (using regex \"%s\")",
			d.Date,
			DateFormat,
			DateFormatRegex)
	}
	return nil
}

// getDate is retuning all exchange rates at a date.
func (r *ExchangeRates) getDate(date string) (Date, error) {
	if err := isValidDate(date); err != nil {
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
