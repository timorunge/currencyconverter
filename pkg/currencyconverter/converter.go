package currencyconverter

import (
	"errors"
	"reflect"
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
func NewConverter(converter Converter) *Converter {
	return &Converter{
		Amount:         converter.Amount,
		BaseCurrency:   converter.BaseCurrency,
		Currencies:     converter.Currencies,
		Date:           converter.Date,
		ExchangeRates:  converter.ExchangeRates,
		TargetCurrency: converter.TargetCurrency,
	}
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

// IsValidAmount is checking if the amount is valid.
func (c *Converter) IsValidAmount() error {
	return c.isValidAmount()
}

// SetAmount is setting the amount that should be used for the calculation.
func (c *Converter) SetAmount(amount float64) {
	c.Amount = amount
}

// SetBaseCurrency is setting the base currency.
func (c *Converter) SetBaseCurrency(baseCurrency string) {
	c.BaseCurrency = baseCurrency
}

// SetCurrencies is setting the supported currencies.
func (c *Converter) SetCurrencies(currencies Currencies) {
	c.Currencies = currencies
}

// SetDate is setting the date.
func (c *Converter) SetDate(date string) {
	c.Date = date
}

// SetExchangeRates is setting the exchange rates.
func (c *Converter) SetExchangeRates(exchangeRates ExchangeRates) {
	c.ExchangeRates = exchangeRates
}

// SetTargetCurrency is setting the target currency.
func (c *Converter) SetTargetCurrency(targetCurrency string) {
	c.TargetCurrency = targetCurrency
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
