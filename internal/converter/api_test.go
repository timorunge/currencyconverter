package converter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testXMLResponse = `<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<gesmes:subject>Reference rates</gesmes:subject>
	<gesmes:Sender><gesmes:name>European Central Bank</gesmes:name></gesmes:Sender>
	<Cube>
		<Cube time='2019-07-31'>
			<Cube currency='USD' rate='1.1151'/>
			<Cube currency='GBP' rate='0.91623'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`

func TestNewAPI(t *testing.T) {
	t.Parallel()

	t.Run("Defaults", func(t *testing.T) {
		t.Parallel()
		a := NewAPI()
		if a.endpoint != apiEndpoint {
			t.Fatalf("got %q, want %q", a.endpoint, apiEndpoint)
		}
		if a.historicalData {
			t.Fatal("expected historicalData to be false by default")
		}
		if a.client == nil {
			t.Fatal("expected client to be initialized")
		}
	})

	t.Run("withTimeout", func(t *testing.T) {
		t.Parallel()
		a := NewAPI(withTimeout(30 * time.Second))
		if a.client.Timeout != 30*time.Second {
			t.Fatalf("got %v, want %v", a.client.Timeout, 30*time.Second)
		}
	})
}

func TestAPIGet(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(testXMLResponse))
		}))
		defer srv.Close()

		rates, err := NewAPI(withEndpoint(srv.URL)).Get(context.Background())
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if len(rates.Dates) != 1 {
			t.Fatalf("got %d dates, want 1", len(rates.Dates))
		}
		if rates.Dates[0].Date != "2019-07-31" {
			t.Fatalf("got date %q, want %q", rates.Dates[0].Date, "2019-07-31")
		}
	})

	t.Run("ServerError", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer srv.Close()

		_, err := NewAPI(withEndpoint(srv.URL)).Get(context.Background())
		if err == nil {
			t.Fatal("expected error for 500 response")
		}
	})

	t.Run("EmptyRates", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<Cube></Cube>
</gesmes:Envelope>`))
		}))
		defer srv.Close()

		_, err := NewAPI(withEndpoint(srv.URL)).Get(context.Background())
		if err == nil {
			t.Fatal("expected error for empty rates")
		}
	})

	t.Run("CanceledContext", func(t *testing.T) {
		t.Parallel()
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(testXMLResponse))
		}))
		defer srv.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := NewAPI(withEndpoint(srv.URL)).Get(ctx)
		if err == nil {
			t.Fatal("expected error for canceled context")
		}
	})

	t.Run("HistoricalEndpoint", func(t *testing.T) {
		t.Parallel()
		var gotPath string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			_, _ = w.Write([]byte(testXMLResponse))
		}))
		defer srv.Close()

		_, err := NewAPI(withEndpoint(srv.URL), WithHistoricalData(true)).Get(context.Background())
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if gotPath != "/stats/eurofxref/eurofxref-hist.xml" {
			t.Fatalf("got path %q, want %q", gotPath, "/stats/eurofxref/eurofxref-hist.xml")
		}
	})
}
