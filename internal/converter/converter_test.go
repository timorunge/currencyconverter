package converter

import (
	"encoding/json"
	"math"
	"testing"
)

var testDateJSON = []byte(`{"Date":"2019-07-31","Rates":[{"Currency":"USD","Rate":1.1151},{"Currency":"JPY","Rate":121.04},{"Currency":"BGN","Rate":1.9558},{"Currency":"CZK","Rate":25.658},{"Currency":"DKK","Rate":7.4674},{"Currency":"GBP","Rate":0.91623},{"Currency":"HUF","Rate":326.48},{"Currency":"PLN","Rate":4.2912},{"Currency":"RON","Rate":4.7338},{"Currency":"SEK","Rate":10.6645},{"Currency":"CHF","Rate":1.1041},{"Currency":"ISK","Rate":134.7},{"Currency":"NOK","Rate":9.7778},{"Currency":"TRY","Rate":6.161},{"Currency":"AUD","Rate":1.6175},{"Currency":"BRL","Rate":4.218},{"Currency":"CAD","Rate":1.4662},{"Currency":"CNY","Rate":7.6743},{"Currency":"HKD","Rate":8.7289},{"Currency":"IDR","Rate":15639.3},{"Currency":"ILS","Rate":3.8951},{"Currency":"INR","Rate":76.6965},{"Currency":"KRW","Rate":1318.25},{"Currency":"MXN","Rate":21.2005},{"Currency":"MYR","Rate":4.6015},{"Currency":"NZD","Rate":1.6883},{"Currency":"PHP","Rate":56.685},{"Currency":"SGD","Rate":1.5261},{"Currency":"THB","Rate":34.273},{"Currency":"ZAR","Rate":15.8634}]}`)

func testDate(t *testing.T) DailyRates {
	t.Helper()
	var d DailyRates
	if err := json.Unmarshal(testDateJSON, &d); err != nil {
		t.Fatal(err)
	}
	return d
}

func TestConvert(t *testing.T) {
	t.Parallel()
	d := testDate(t)

	t.Run("EURtoUSD", func(t *testing.T) {
		t.Parallel()
		result, err := Convert(d, 1337, "EUR", "USD")
		if err != nil {
			t.Fatal(err)
		}
		if result != 1490.8887 {
			t.Fatalf("got %f, want %f", result, 1490.8887)
		}
	})

	t.Run("UnknownFrom", func(t *testing.T) {
		t.Parallel()
		_, err := Convert(d, 100, "XYZ", "USD")
		if err == nil {
			t.Fatal("expected error for unknown base currency")
		}
	})

	t.Run("UnknownTo", func(t *testing.T) {
		t.Parallel()
		_, err := Convert(d, 100, "EUR", "XYZ")
		if err == nil {
			t.Fatal("expected error for unknown target currency")
		}
	})
}

func TestValidateAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"Positive", 100, false},
		{"Zero", 0, true},
		{"Negative", -1, true},
		{"NaN", math.NaN(), true},
		{"Inf", math.Inf(1), true},
		{"NegInf", math.Inf(-1), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateAmount(%v) error = %v, wantErr %v", tt.amount, err, tt.wantErr)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"Latest", "latest", false},
		{"Valid", "2019-07-31", false},
		{"Empty", "", true},
		{"InvalidFormat", "not-a-date", true},
		{"GarbagePrefix", "garbage2019-07-31", true},
		{"LatestSuffix", "latestgarbage", true},
		{"ImpossibleDate", "2019-02-31", true},
		{"Weekend", "2019-07-27", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateDate(tt.date)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateDate(%q) error = %v, wantErr %v", tt.date, err, tt.wantErr)
			}
		})
	}
}
