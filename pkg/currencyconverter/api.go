package currencyconverter

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// APIEndpoint is the endpoint for the API.
//   More details can be found at
//   https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html
// APIHistoricalData is enabling the historical data part.
// APITimeout is the timeout for the API call.
const (
	APIEndpoint       = "https://www.ecb.europa.eu"
	APIHistoricalData = false
	APITimeout        = 5 * time.Second
)

// ErrAPINoRate is the error when there is a successful response from the API
// but it's empty...
// ErrAPIStatus is the error when the API resonse is != 2xx
var (
	ErrAPINoRate = errors.New("The API is not providing a single exchange rate")
	ErrAPIStatus = errors.New("Can not get communicate with the exchange rate API")
)

// API is the struct for the API.
type API struct {
	Endpoint       string
	ExchangeRates  ExchangeRates
	HistoricalData bool
	Timeout        time.Duration
}

// NewAPI is returning a new API struct.
func NewAPI() *API {
	return &API{
		Endpoint:       APIEndpoint,
		HistoricalData: APIHistoricalData,
		Timeout:        APITimeout,
	}
}

// Get is returning the exchange rates from the API.
func (a *API) Get() (ExchangeRates, error) {
	url := fmt.Sprintf("%s/stats/eurofxref/eurofxref-daily.xml", a.Endpoint)
	if a.HistoricalData {
		url = fmt.Sprintf("%s/stats/eurofxref/eurofxref-hist.xml", a.Endpoint)
	}
	return a.getExchangeRates(url)
}

// SetEndpoint is setting the endpoint for the API.
func (a *API) SetEndpoint(endpoint string) *API {
	a.Endpoint = endpoint
	return a
}

// SetHistoricalData is the historical data flag.
func (a *API) SetHistoricalData(historicalData bool) *API {
	a.HistoricalData = historicalData
	return a
}

// SetTimeout is setting the HTTP timeout for the API call.
func (a *API) SetTimeout(timeout time.Duration) *API {
	a.Timeout = timeout
	return a
}

// getExchangeRates is getting the exchange rates from the API.
func (a *API) getExchangeRates(url string) (ExchangeRates, error) {
	httpClient := &http.Client{Timeout: a.Timeout}
	resp, err := httpClient.Get(url)
	if err != nil {
		return NullExchangeRates, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return NullExchangeRates, err
		}
		if err := xml.Unmarshal(bytes, &a.ExchangeRates); err != nil {
			return NullExchangeRates, err
		}
		if len(a.ExchangeRates.Dates) == 0 {
			return NullExchangeRates, ErrAPINoRate
		}
		return a.ExchangeRates, nil
	}
	return NullExchangeRates, ErrAPIStatus
}
