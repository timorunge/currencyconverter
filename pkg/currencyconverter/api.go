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
// APITimeout is the timeout for the API call.
const (
	APIEndpoint = "https://www.ecb.europa.eu"
	APITimeout  = 5 * time.Second
)

// API is the struct for the API.
type API struct {
	Endpoint       string
	ExchangeRates  ExchangeRates
	HistoricalData bool
	Timeout        time.Duration
}

// NewAPI returning a new API struct.
func NewAPI(api API) *API {
	return &API{
		Endpoint:       api.Endpoint,
		ExchangeRates:  api.ExchangeRates,
		HistoricalData: api.HistoricalData,
		Timeout:        api.Timeout,
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
func (a *API) SetEndpoint(endpoint string) {
	a.Endpoint = endpoint
}

// SetHistoricalData is the historical data flag.
func (a *API) SetHistoricalData(historicalData bool) {
	a.HistoricalData = historicalData
}

// SetTimeout is setting the HTTP timeout for the API call.
func (a *API) SetTimeout(timeout time.Duration) {
	a.Timeout = timeout
}

// getExchangeRates is getting the exchange rates from the API.
func (a *API) getExchangeRates(url string) (ExchangeRates, error) {
	httpClient := &http.Client{Timeout: a.Timeout}
	resp, err := httpClient.Get(url)
	if err != nil {
		return a.ExchangeRates, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return a.ExchangeRates, err
		}
		if err := xml.Unmarshal(bytes, &a.ExchangeRates); err != nil {
			return a.ExchangeRates, err
		}
		if len(a.ExchangeRates.Dates) == 0 {
			return a.ExchangeRates, errors.New("Not getting a single exchange rate from the API")
		}
		return a.ExchangeRates, nil
	}
	return a.ExchangeRates, errors.New("Not getting exchange rates from the API")
}
