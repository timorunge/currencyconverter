// ECB API client for fetching exchange rates.

package converter

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiEndpoint = "https://www.ecb.europa.eu"
	apiTimeout  = 10 * time.Second

	maxResponseSize = 10 << 20 // 10 MB
)

var errAPINoRate = errors.New("api: no exchange rate provided")

// API fetches exchange rates from the European Central Bank.
type API struct {
	endpoint       string
	historicalData bool
	client         *http.Client
}

type apiOption func(*API)

// withEndpoint overrides the ECB API base URL (used by tests).
func withEndpoint(endpoint string) apiOption {
	return func(a *API) { a.endpoint = endpoint }
}

// WithHistoricalData toggles between daily and full-history endpoints.
func WithHistoricalData(historical bool) apiOption {
	return func(a *API) { a.historicalData = historical }
}

func withTimeout(timeout time.Duration) apiOption {
	return func(a *API) {
		a.client = newHTTPClient(timeout)
	}
}

// NewAPI returns a new API with default settings.
func NewAPI(opts ...apiOption) *API {
	a := &API{
		endpoint: apiEndpoint,
		client:   newHTTPClient(apiTimeout),
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Get fetches exchange rates from the ECB API.
func (a *API) Get(ctx context.Context) (ExchangeRates, error) {
	url := a.endpoint + "/stats/eurofxref/eurofxref-daily.xml"
	if a.historicalData {
		url = a.endpoint + "/stats/eurofxref/eurofxref-hist.xml"
	}
	return a.getExchangeRates(ctx, url)
}

func (a *API) getExchangeRates(ctx context.Context, url string) (ExchangeRates, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ExchangeRates{}, fmt.Errorf("api: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return ExchangeRates{}, fmt.Errorf("api: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode/100 != 2 {
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return ExchangeRates{}, fmt.Errorf("api: unexpected status %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var rates ExchangeRates
	limited := io.LimitReader(resp.Body, maxResponseSize)
	if err := xml.NewDecoder(limited).Decode(&rates); err != nil {
		return ExchangeRates{}, fmt.Errorf("api: %w", err)
	}
	if len(rates.Dates) == 0 {
		return ExchangeRates{}, errAPINoRate
	}
	return rates, nil
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}
