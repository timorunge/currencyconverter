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
const (
	DateFormat      = "latest|YYYY-MM-DD"
	DateFormatRegex = "^latest|(1999|20[0-9]{2})-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])$"
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
func (c *Converter) Convert() (float64, error) {
	exchangeRate, err := c.calculateExchangeRate()
	if err != nil {
		return float64(0), err
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
func (c *Converter) calculateExchangeRate() (float64, error) {
	f64 := float64(0)
	if err := c.isValidAmount(); err != nil {
		return f64, err
	}
	if err := c.Currencies.isSupportedCurrency(c.BaseCurrency); err != nil {
		return f64, err
	}
	if err := c.Currencies.isSupportedCurrency(c.TargetCurrency); err != nil {
		return f64, err
	}

	baseCurrencyRate, err := c.ExchangeRates.GetRate(c.Date, c.BaseCurrency)
	if err != nil {
		return f64, err
	}
	targetCurrencyRate, err := c.ExchangeRates.GetRate(c.Date, c.TargetCurrency)
	if err != nil {
		return f64, err
	}

	return baseCurrencyRate / targetCurrencyRate, nil
}

// isValidAmount is checking if the given amount is valid.
func (c *Converter) isValidAmount() error {
	if v := reflect.ValueOf(c.Amount).Kind(); v == reflect.Float64 {
		if c.Amount == float64(0) {
			return errors.New("Amount should be bigger than zero")
		}
		return nil
	}
	return errors.New("Amount is no float64 value")
}

// isValidDate is checking if the date is in the right format.
func isValidDate(date string) error {
	c := &Converter{Date: date}
	return c.isValidDate()
}

// isValidDate is checking if the date is in the right format and not on a weekend.
func (c *Converter) isValidDate() error {
	if !regexp.MustCompile(DateFormatRegex).MatchString(c.Date) {
		return fmt.Errorf("\"%s\" is no valid date in format \"%s\" (using regex \"%s\")",
			c.Date,
			DateFormat,
			DateFormatRegex)
	}
	if c.Date != "latest" {
		parsedTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", c.Date))
		if err != nil {
			return err
		}
		if parsedTime.After(time.Now()) {
			return fmt.Errorf("\"%s\" is in the future", c.Date)
		}
		if parsedTime.Weekday() == time.Saturday || parsedTime.Weekday() == time.Sunday {
			return fmt.Errorf("\"%s\" is a %s, date needs to be between %s and %s",
				c.Date,
				parsedTime.Weekday().String(),
				time.Monday.String(),
				time.Friday.String())
		}
	}
	return nil
}
