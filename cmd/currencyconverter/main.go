// Package main provides the currencyconverter CLI entry point.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"

	"github.com/timorunge/currencyconverter/internal/converter"
)

type exitCode int

const (
	exitSuccess exitCode = iota
	exitError
	exitUsageError
)

var version = "dev"

type conversionResult struct {
	Amount float64 `json:"amount"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Result float64 `json:"result"`
	Date   string  `json:"date"`
	Source string  `json:"source"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Stdout, os.Stderr, os.Args[1:], version)
	stop()
	os.Exit(int(code))
}

func run(ctx context.Context, stdout, stderr io.Writer, args []string, version string) exitCode {
	fs := pflag.NewFlagSet("currencyconverter", pflag.ContinueOnError)
	fs.SortFlags = false
	fs.SetOutput(stderr)
	fs.Usage = func() { printUsage(stdout, stderr, fs) }

	// Conversion flags.
	from := fs.StringP("from", "f", "EUR", "Currency to use as base")
	to := fs.StringP("to", "t", "USD", "Currency to convert to")
	amount := fs.Float64P("amount", "a", 1, "Amount to calculate")
	date := fs.StringP("date", "d", "latest", "Date for the calculation (latest|YYYY-MM-DD)")
	reverse := fs.BoolP("reverse", "r", false, "Swap from and to values")

	// Output flags.
	jsonOut := fs.Bool("json", false, "Output result as JSON")
	noCache := fs.Bool("no-cache", false, "Disable caching")

	// Informational flags.
	help := fs.BoolP("help", "h", false, "Show this help message")
	showVersion := fs.Bool("version", false, "Show the version of currencyconverter")
	supportedCurrencies := fs.Bool("supported-currencies", false, "Show a list with all supported currencies")

	if err := fs.Parse(args); err != nil {
		return exitUsageError
	}

	if *help {
		printUsage(stdout, stderr, fs)
		return exitSuccess
	}
	if len(fs.Args()) > 0 {
		_, _ = fmt.Fprintf(stderr, "error: unexpected positional arguments: %s\n", strings.Join(fs.Args(), " "))
		return exitUsageError
	}
	if *showVersion {
		printVersion(version, stdout)
		return exitSuccess
	}
	if *supportedCurrencies {
		printSupportedCurrencies(stdout)
		return exitSuccess
	}

	baseCurrency := strings.ToUpper(*from)
	targetCurrency := strings.ToUpper(*to)
	if *reverse {
		baseCurrency, targetCurrency = targetCurrency, baseCurrency
	}

	if err := validateOptions(*amount, baseCurrency, targetCurrency, *date); err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return exitUsageError
	}

	historical := *date != converter.DateLatest

	api := converter.NewAPI(
		converter.WithHistoricalData(historical),
	)
	// Historical dates are cached forever -- their rates never change.
	// Latest rates rotate by hour and expire after 1 hour, ensuring
	// cross-day freshness without spamming the temp directory.
	cacheFile := "currencyconverter-" + *date
	var maxAge time.Duration
	if !historical {
		cacheFile = "currencyconverter-latest-" + time.Now().Format("T15")
		maxAge = time.Hour
	}
	fileCache := converter.NewFileCache(converter.FileCacheConfig{
		Enabled:  !*noCache,
		Filename: cacheFile,
		MaxAge:   maxAge,
	})

	rates, err := getRates(ctx, stderr, api, fileCache, *date)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return exitError
	}

	rawResult, err := converter.Convert(rates, *amount, baseCurrency, targetCurrency)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return exitError
	}
	result := roundTo4(rawResult)

	if *jsonOut {
		err = json.NewEncoder(stdout).Encode(conversionResult{
			Amount: *amount,
			From:   baseCurrency,
			To:     targetCurrency,
			Result: result,
			Date:   rates.Date,
			Source: "ECB",
		})
	} else {
		_, err = fmt.Fprintf(stdout, "%s %s ≈ %s %s (%s, ECB)\n",
			formatAmount(*amount),
			baseCurrency,
			formatAmount(result),
			targetCurrency,
			rates.Date)
	}
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return exitError
	}

	return exitSuccess
}

func getRates(ctx context.Context, stderr io.Writer, api *converter.API, fileCache *converter.FileCache, date string) (converter.DailyRates, error) {
	cached, cacheErr := fileCache.Get()
	if cacheErr == nil {
		return cached, nil
	}
	if !converter.IsCacheMiss(cacheErr) {
		_, _ = fmt.Fprintf(stderr, "warning: cache read failed: %v\n", cacheErr)
	}
	apiRates, err := api.Get(ctx)
	if err != nil {
		return converter.DailyRates{}, err
	}
	daily, err := apiRates.GetDate(date)
	if err != nil {
		return converter.DailyRates{}, err
	}
	if writeErr := fileCache.Write(daily); writeErr != nil {
		_, _ = fmt.Fprintf(stderr, "warning: cache write failed: %v\n", writeErr)
	}
	return daily, nil
}

func roundTo4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func formatAmount(v float64) string {
	s := strconv.FormatFloat(v, 'f', 4, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func printUsage(w, orig io.Writer, fs *pflag.FlagSet) {
	fs.SetOutput(w)
	_, _ = fmt.Fprintf(w, "currencyconverter - Exchange rate calculator using ECB reference rates\n\nUsage:\n  currencyconverter [OPTIONS]\n\nOptions:\n")
	fs.PrintDefaults()
	fs.SetOutput(orig)
}

func printSupportedCurrencies(w io.Writer) {
	_, _ = fmt.Fprintln(w, "Supported currencies (ECB foreign exchange reference rates):")
	_, _ = fmt.Fprintf(w, "  - %s\n", strings.Join(converter.SupportedCurrencies(), "\n  - "))
	_, _ = fmt.Fprintf(w, "\nSource: https://www.ecb.europa.eu/stats/eurofxref/\n")
}

func printVersion(version string, w io.Writer) {
	_, _ = fmt.Fprintf(w, "currencyconverter %s\n", version)
}

func validateOptions(amount float64, baseCurrency, targetCurrency, date string) error {
	return errors.Join(
		converter.ValidateCurrency(baseCurrency),
		converter.ValidateCurrency(targetCurrency),
		converter.ValidateAmount(amount),
		converter.ValidateDate(date),
	)
}
