package currencyconverter

import (
	"testing"
	"time"
)

func TestAPIEndpoint(t *testing.T) {
	testAPIEndpoint := "http://www.ecb.europa.eu"

	a := NewAPI()
	a.SetEndpoint(testAPIEndpoint)
	if a.Endpoint != testAPIEndpoint {
		t.Error("Expected result is", testAPIEndpoint, "but got", a.Endpoint)
	}
}

func TestAPIHistoricalData(t *testing.T) {
	testAPIHistoricalData := true

	a := NewAPI()
	a.SetHistoricalData(testAPIHistoricalData)
	if a.HistoricalData != testAPIHistoricalData {
		t.Error("Expected result is", testAPIHistoricalData, "but got", a.HistoricalData)
	}
}

func TestAPITimeout(t *testing.T) {
	testAPITimeout := 10 * time.Second

	a := NewAPI()
	a.SetTimeout(testAPITimeout)
	if a.Timeout != testAPITimeout {
		t.Error("Expected result is", testAPITimeout, "but got", a.Timeout)
	}
}
