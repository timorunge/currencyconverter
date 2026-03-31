// Currency conversion and input validation.

package converter

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"
)

const (
	dateFormat = "latest|YYYY-MM-DD"

	// DateLatest is the sentinel value for requesting the most recent rates.
	DateLatest = "latest"
)

var (
	errInvalidAmount = errors.New("amount must be greater than zero")

	dateFormatRegex = regexp.MustCompile(`^(latest|(1999|20[0-9]{2})-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01]))$`)
)

// Convert returns the converted amount in the target currency.
func Convert(date DailyRates, amount float64, from, to string) (float64, error) {
	rate, err := exchangeRate(date, from, to)
	if err != nil {
		return 0, err
	}
	return amount * rate, nil
}

// exchangeRate computes the cross-rate between two currencies.
// Both rates are relative to EUR, so from→to = toRate / fromRate.
func exchangeRate(date DailyRates, from, to string) (float64, error) {
	fromRate, err := date.getRate(from)
	if err != nil {
		return 0, err
	}
	toRate, err := date.getRate(to)
	if err != nil {
		return 0, err
	}
	return toRate / fromRate, nil
}

// ValidateAmount checks that the amount is positive.
func ValidateAmount(amount float64) error {
	if amount <= 0 || math.IsInf(amount, 0) || math.IsNaN(amount) {
		return errInvalidAmount
	}
	return nil
}

// ValidateDate checks that the date string matches the expected format
// and represents a valid weekday not in the future.
func ValidateDate(date string) error {
	if !dateFormatRegex.MatchString(date) {
		return fmt.Errorf("%q is not a valid date in format %q", date, dateFormat)
	}
	if date != DateLatest {
		parsedTime, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return fmt.Errorf("invalid date %q: %w", date, err)
		}
		if parsedTime.After(time.Now()) {
			return fmt.Errorf("%q is in the future", date)
		}
		if parsedTime.Weekday() == time.Saturday || parsedTime.Weekday() == time.Sunday {
			return fmt.Errorf("%q is a %s, date needs to be a weekday between %s and %s",
				date, parsedTime.Weekday(), time.Monday, time.Friday)
		}
	}
	return nil
}
