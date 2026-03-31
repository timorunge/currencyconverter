package converter

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

var testRatesXML = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender><gesmes:name>European Central Bank</gesmes:name></gesmes:Sender>
	<Cube>
		<Cube time='2019-07-31'>
			<Cube currency='USD' rate='1.1151'/>
			<Cube currency='GBP' rate='0.91623'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`)

func TestDailyRatesGetRate(t *testing.T) {
	t.Parallel()

	var d DailyRates
	if err := json.Unmarshal(testDateJSON, &d); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		date     *DailyRates
		currency string
		want     float64
		wantErr  bool
	}{
		{"EUR", &d, "EUR", 1, false},
		{"USD", &d, "USD", 1.1151, false},
		{"Unknown", &d, "XYZ", 0, true},
		{"ZeroRate", &DailyRates{Rates: []Rate{{Currency: "BAD", Rate: 0}}}, "BAD", 0, true},
		{"NegativeRate", &DailyRates{Rates: []Rate{{Currency: "BAD", Rate: -1}}}, "BAD", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rate, err := tt.date.getRate(tt.currency)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getRate(%q) error = %v, wantErr %v", tt.currency, err, tt.wantErr)
			}
			if !tt.wantErr && rate != tt.want {
				t.Fatalf("getRate(%q) = %f, want %f", tt.currency, rate, tt.want)
			}
		})
	}
}

func TestExchangeRatesXMLParsing(t *testing.T) {
	t.Parallel()

	var r ExchangeRates
	if err := xml.Unmarshal(testRatesXML, &r); err != nil {
		t.Fatal(err)
	}

	if len(r.Dates) != 1 {
		t.Fatalf("got %d dates, want 1", len(r.Dates))
	}
	if r.Dates[0].Date != "2019-07-31" {
		t.Fatalf("got date %q, want %q", r.Dates[0].Date, "2019-07-31")
	}
	if len(r.Dates[0].Rates) != 2 {
		t.Fatalf("got %d rates, want 2", len(r.Dates[0].Rates))
	}
}

func TestExchangeRatesGetDate(t *testing.T) {
	t.Parallel()

	t.Run("ExactMatch", func(t *testing.T) {
		t.Parallel()
		r := &ExchangeRates{
			Dates: []DailyRates{{Date: "2019-07-31"}},
		}
		d, err := r.GetDate("2019-07-31")
		if err != nil {
			t.Fatal(err)
		}
		if d.Date != "2019-07-31" {
			t.Fatalf("got %q, want %q", d.Date, "2019-07-31")
		}
	})

	t.Run("LatestPicksMostRecent", func(t *testing.T) {
		t.Parallel()
		r := &ExchangeRates{
			Dates: []DailyRates{
				{Date: "2019-07-29"},
				{Date: "2019-07-31"},
				{Date: "2019-07-30"},
			},
		}
		d, err := r.GetDate("latest")
		if err != nil {
			t.Fatal(err)
		}
		if d.Date != "2019-07-31" {
			t.Fatalf("got %q, want %q", d.Date, "2019-07-31")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		r := &ExchangeRates{
			Dates: []DailyRates{{Date: "2019-07-31"}},
		}
		_, err := r.GetDate("2019-07-30")
		if err == nil {
			t.Fatal("expected error for missing date")
		}
	})

	t.Run("LatestEmpty", func(t *testing.T) {
		t.Parallel()
		r := &ExchangeRates{}
		_, err := r.GetDate("latest")
		if err == nil {
			t.Fatal("expected error for empty exchange rates")
		}
	})
}
