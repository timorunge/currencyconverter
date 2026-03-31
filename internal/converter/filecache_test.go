package converter

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileCacheFilePath(t *testing.T) {
	t.Parallel()
	c := NewFileCache(FileCacheConfig{
		Directory: "/tmp",
		Filename:  "today",
	})
	want := filepath.Join("/tmp", "today.json")
	if got := c.filePath(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFileCacheRoundTrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	c := NewFileCache(FileCacheConfig{
		Directory: dir,
		Enabled:   true,
		Filename:  "test-roundtrip",
	})

	want := DailyRates{
		Date: "2019-07-31",
		Rates: []Rate{
			{Currency: "USD", Rate: 1.1151},
			{Currency: "GBP", Rate: 0.91623},
		},
	}

	if err := c.Write(want); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := c.Get()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Date != want.Date {
		t.Fatalf("Date: got %q, want %q", got.Date, want.Date)
	}
	if len(got.Rates) != len(want.Rates) {
		t.Fatalf("Rates: got %d, want %d", len(got.Rates), len(want.Rates))
	}
	for i, r := range got.Rates {
		if r.Currency != want.Rates[i].Currency || r.Rate != want.Rates[i].Rate {
			t.Fatalf("Rate[%d]: got %v, want %v", i, r, want.Rates[i])
		}
	}
}

func TestFileCacheExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	rates := DailyRates{
		Date:  "2019-07-31",
		Rates: []Rate{{Currency: "USD", Rate: 1.1151}},
	}

	t.Run("Fresh", func(t *testing.T) {
		t.Parallel()
		c := NewFileCache(FileCacheConfig{
			Directory: dir,
			Enabled:   true,
			Filename:  "expiry-fresh",
			MaxAge:    time.Hour,
		})
		if err := c.Write(rates); err != nil {
			t.Fatalf("Write: %v", err)
		}
		if _, err := c.Get(); err != nil {
			t.Fatalf("expected fresh cache hit, got: %v", err)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		t.Parallel()
		c := NewFileCache(FileCacheConfig{
			Directory: dir,
			Enabled:   true,
			Filename:  "expiry-old",
			MaxAge:    time.Hour,
		})
		if err := c.Write(rates); err != nil {
			t.Fatalf("Write: %v", err)
		}
		// Backdate the file to 2 hours ago.
		past := time.Now().Add(-2 * time.Hour)
		if err := os.Chtimes(c.filePath(), past, past); err != nil {
			t.Fatalf("Chtimes: %v", err)
		}
		_, err := c.Get()
		if !errors.Is(err, errFileCacheExpired) {
			t.Fatalf("expected errFileCacheExpired, got: %v", err)
		}
	})

	t.Run("NoMaxAge", func(t *testing.T) {
		t.Parallel()
		c := NewFileCache(FileCacheConfig{
			Directory: dir,
			Enabled:   true,
			Filename:  "expiry-none",
		})
		if err := c.Write(rates); err != nil {
			t.Fatalf("Write: %v", err)
		}
		// Backdate the file -- should still be served with zero MaxAge.
		past := time.Now().Add(-48 * time.Hour)
		if err := os.Chtimes(c.filePath(), past, past); err != nil {
			t.Fatalf("Chtimes: %v", err)
		}
		if _, err := c.Get(); err != nil {
			t.Fatalf("expected cache hit with no MaxAge, got: %v", err)
		}
	})
}

func TestFileCacheGetErrors(t *testing.T) {
	t.Parallel()

	t.Run("Disabled", func(t *testing.T) {
		t.Parallel()
		c := NewFileCache(FileCacheConfig{})
		_, err := c.Get()
		if !errors.Is(err, errFileCacheDisabled) {
			t.Fatalf("got %v, want %v", err, errFileCacheDisabled)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		c := NewFileCache(FileCacheConfig{
			Directory: t.TempDir(),
			Enabled:   true,
			Filename:  "nonexistent",
		})
		_, err := c.Get()
		if err == nil {
			t.Fatal("expected error for missing cache file")
		}
	})
}
