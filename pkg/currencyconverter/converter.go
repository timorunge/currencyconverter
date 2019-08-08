package currencyconverter

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"
)

// DateFormat is the format of the date.
// DateFormatRegex is the regex for the date format.
// DateLatest is giving just the latest day.
const (
	DateFormat      = "latest|YYYY-MM-DD"
	DateFormatRegex = "^latest|(1999|20[0-9]{2})-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])$"
	DateLatest      = "latest"
)

// ErrAmountIsZero is the error message if the amount is 0.
// ErrNoF64 is the error message if the amount is no float64 value.
// ErrStrDateInFuture is the error message string if the date is in the future.
// ErrStrDateIsWeekend is the error message string if the date is on a weekend.
// ErrStrNoValidDateFormat is the error message string if the format is wrong.
var (
	ErrAmountIsZero         = errors.New("Amount should be bigger than zero")
	ErrNoF64                = errors.New("Amount is no float64 value")
	ErrStrDateInFuture      = "\"%s\" is in the future"
	ErrStrDateIsWeekend     = "\"%s\" is a %s, date needs to be a weekday between %s and %s"
	ErrStrNoValidDateFormat = "\"%s\" is no valid date in format \"%s\" (using regex \"%s\")"
)

// Converter is the struct for the converter.
type Converter struct {
	Amount         float64
	BaseCurrency   string
	Currencies     Currencies
	Date           string
	ExchangeRates  ExchangeRates
	TargetCurrency string
}

// NewConverter is returning a new converter struct.
func NewConverter() *Converter {
	return &Converter{}
}

// Convert is returning the converted amount in the target currency.
func (c *Converter) Convert() (f64 float64, err error) {
	exchangeRate, err := c.calculateExchangeRate()
	if err != nil {
		return
	}
	return c.Amount * exchangeRate, nil
}

// ExchangeRate is returning the exchange rate between the currencies.
func (c *Converter) ExchangeRate() (float64, error) {
	return c.calculateExchangeRate()
}

// IsValidAmount is checking if the amount is valid.
func IsValidAmount(amount float64) error {
	c := &Converter{Amount: amount}
	return c.isValidAmount()
}

// IsValidDate is checking if the date is in the right format.
func IsValidDate(date string) error {
	c := &Converter{Date: date}
	return c.isValidDate()
}

// IsValidAmount is checking if the amount is valid.
func (c *Converter) IsValidAmount() error {
	return c.isValidAmount()
}

// IsValidDate is checking if the date is in the right format.
func (c *Converter) IsValidDate() error {
	return c.isValidDate()
}

// SetAmount is setting the amount that should be used for the calculation.
func (c *Converter) SetAmount(amount float64) *Converter {
	c.Amount = amount
	return c
}

// SetBaseCurrency is setting the base currency.
func (c *Converter) SetBaseCurrency(baseCurrency string) *Converter {
	c.BaseCurrency = baseCurrency
	return c
}

// SetCurrencies is setting the supported currencies.
func (c *Converter) SetCurrencies(currencies Currencies) *Converter {
	c.Currencies = currencies
	return c
}

// SetDate is setting the date.
func (c *Converter) SetDate(date string) *Converter {
	c.Date = date
	return c
}

// SetExchangeRates is setting the exchange rates.
func (c *Converter) SetExchangeRates(exchangeRates ExchangeRates) *Converter {
	c.ExchangeRates = exchangeRates
	return c
}

// SetTargetCurrency is setting the target currency.
func (c *Converter) SetTargetCurrency(targetCurrency string) *Converter {
	c.TargetCurrency = targetCurrency
	return c
}

// calculateExchangeRate is doing the !(ultra complex) mathematical converting
// operation of calculating the exchange rate.
func (c *Converter) calculateExchangeRate() (f64 float64, err error) {
	if err = c.isValidAmount(); err != nil {
		return
	}
	if err = c.Currencies.isSupportedCurrency(c.BaseCurrency); err != nil {
		return
	}
	if err = c.Currencies.isSupportedCurrency(c.TargetCurrency); err != nil {
		return
	}

	baseCurrencyRate, err := c.ExchangeRates.GetRate(c.Date, c.BaseCurrency)
	if err != nil {
		return
	}
	targetCurrencyRate, err := c.ExchangeRates.GetRate(c.Date, c.TargetCurrency)
	if err != nil {
		return
	}

	return baseCurrencyRate / targetCurrencyRate, nil
}

// isValidAmount is checking if the given amount is valid.
func (c *Converter) isValidAmount() error {
	if v := reflect.ValueOf(c.Amount).Kind(); v == reflect.Float64 {
		if c.Amount == float64(0) {
			return ErrAmountIsZero
		}
		return nil
	}
	return ErrNoF64
}

// isValidDate is checking if the date is in the right format and not on a weekend.
func (c *Converter) isValidDate() error {
	if !regexp.MustCompile(DateFormatRegex).MatchString(c.Date) {
		return fmt.Errorf(ErrStrNoValidDateFormat, c.Date, DateFormat, DateFormatRegex)
	}
	if c.Date != DateLatest {
		parsedTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", c.Date))
		if err != nil {
			return err
		}
		if parsedTime.After(time.Now()) {
			return fmt.Errorf(ErrStrDateInFuture, c.Date)
		}
		if parsedTime.Weekday() == time.Saturday || parsedTime.Weekday() == time.Sunday {
			return fmt.Errorf(ErrStrDateIsWeekend, c.Date, parsedTime.Weekday().String(), time.Monday.String(), time.Friday.String())
		}
	}
	return nil
}
