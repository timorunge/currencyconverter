package currencyconverter

import (
	"testing"
	"time"
)

var (
	testAPIEndpoint       = "http://www.ecb.europa.eu"
	testAPIHistoricalData = true
	testAPITimeout        = 10 * time.Second
)

func TestAPI(t *testing.T) {
	a := NewAPI(API{})

	a.SetEndpoint(testAPIEndpoint)
	if a.Endpoint != testAPIEndpoint {
		t.Error("Expected result is", testAPIEndpoint, "but got", a.Endpoint)
	}

	a.SetHistoricalData(testAPIHistoricalData)
	if a.HistoricalData != testAPIHistoricalData {
		t.Error("Expected result is", testAPIHistoricalData, "but got", a.HistoricalData)
	}

	a.SetTimeout(testAPITimeout)
	if a.Timeout != testAPITimeout {
		t.Error("Expected result is", testAPITimeout, "but got", a.Timeout)
	}
}
