package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/timorunge/currencyconverter/pkg/currencyconverter"
)

var (
	buildDate = "Unknown"
	gitCommit = "Unknown"
	parser    = flags.NewParser(nil, flags.PrintErrors|flags.PassDoubleDash)
	version   = "Unknown"

	cliConvertOptions struct {
		Amount         float64 `short:"a" long:"amount" description:"Amount to calculate" default:"1" required:"true"`
		BaseCurrency   string  `short:"f" long:"from" description:"Currency to use as base" default:"EUR" required:"true"`
		TargetCurrency string  `short:"t" long:"to" description:"Currency to convert too" default:"USD" required:"true"`
		Date           string  `short:"d" long:"date" description:"Date to use for the calculation - Format: \"latest|YYYY-MM-DD\"" default:"latest" hidden:"true"`
		Reverse        bool    `short:"r" long:"reverse" description:"Swap from and to values"`
	}
	cliCacheOptions struct {
		CacheDirectory string `long:"cache-directory" description:"Directory to store cached responses (default: Operating system temp directory)"`
		CacheTimeout   int    `long:"cache-timeout" description:"Timeout in minutes to invalidate the cache" default:"60"`
		NoCache        bool   `long:"no-cache" description:"Disable caching"`
	}
	cliHelpOptions struct {
		SupportedCurrencies func() `long:"supported-currencies" description:"Show a list with all supported currencies"`
		Version             func() `short:"v" long:"version" description:"Show the version of currencyconverter"`
		Help                func() `short:"h" long:"help" description:"Show this help message"`
	}
)

func init() {
	log.SetFlags(0)

	parser.LongDescription = fmt.Sprintf("%s is calculating the exchange rate using foreign exchange reference rates published by the European Central Bank", parser.Name)
	parser.Usage = "[OPTIONS]"

	parser.AddGroup("Convert options", "", &cliConvertOptions)
	parser.AddGroup("Cache options", "", &cliCacheOptions)
	parser.AddGroup("Help options", "", &cliHelpOptions)

	cliHelpOptions.Help = func() { printHelp() }
	cliHelpOptions.SupportedCurrencies = func() { printSupportedCurrencies() }
	cliHelpOptions.Version = func() { printVersion() }
}

func main() {
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	cliConvertOptions.BaseCurrency = strings.ToUpper(cliConvertOptions.BaseCurrency)
	cliConvertOptions.TargetCurrency = strings.ToUpper(cliConvertOptions.TargetCurrency)
	if cliConvertOptions.Reverse {
		cliConvertOptions.BaseCurrency = cliConvertOptions.TargetCurrency
		cliConvertOptions.TargetCurrency = cliConvertOptions.BaseCurrency
	}

	err = validateOptions()
	if err != nil {
		log.Fatal(err)
	}

	fileCache := currencyconverter.NewFileCache()
	if cliCacheOptions.CacheDirectory != "" {
		fileCache.SetDirectory(cliCacheOptions.CacheDirectory)
	}
	fileCache.SetEnabled(!cliCacheOptions.NoCache)
	fileCache.SetFilename(fmt.Sprintf("%s-%s", currencyconverter.FileCacheFilename, "latest"))
	fileCache.SetTimeout(time.Duration(cliCacheOptions.CacheTimeout) * time.Minute)

	api := currencyconverter.NewAPI()

	rates, err := getRates(api, fileCache)
	if err != nil {
		log.Fatal(err)
	}

	currencies := currencyconverter.Currencies{}
	currencies.SetSupportedCurrencies(currencyconverter.SupportedCurrencies)

	converter := currencyconverter.NewConverter()
	converter.SetAmount(cliConvertOptions.Amount)
	converter.SetBaseCurrency(cliConvertOptions.BaseCurrency)
	converter.SetCurrencies(currencies)
	converter.SetDate(cliConvertOptions.Date)
	converter.SetExchangeRates(rates)
	converter.SetTargetCurrency(cliConvertOptions.TargetCurrency)

	printResult(converter)
}

// getRates is getting the exchange rates.
func getRates(api *currencyconverter.API, fileCache *currencyconverter.FileCache) (currencyconverter.ExchangeRates, error) {
	if cliConvertOptions.Date != "latest" {
		api.SetHistoricalData(true)
		fileCache.SetFilename(fmt.Sprintf("%s-%s", currencyconverter.FileCacheFilename, cliConvertOptions.Date))
		fileCache.SetTimeout(365 * 24 * time.Hour)
	}
	rates, err := fileCache.Get()
	if err != nil {
		rates, err := api.Get()
		if err != nil {
			return rates, err
		}
		dailyRates, err := rates.GetDate(cliConvertOptions.Date)
		if err != nil {
			return rates, err
		}
		finalRates := *currencyconverter.NewExchangeRates()
		finalRates.AddDate(dailyRates)
		if err := fileCache.Write(finalRates); err != nil {
			log.Print(err)
		}
		return finalRates, err
	}
	return rates, err
}

// printHelp is printing the help on stdout.
func printHelp() {
	parser.WriteHelp(os.Stdout)
	os.Exit(0)
}

// printResult is printing the final result of the conversion.
func printResult(converter *currencyconverter.Converter) {
	amountTargetCurrency, err := converter.Convert()
	if err != nil {
		log.Fatal(err)
	}

	exchangeRate, err := converter.ExchangeRate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%f %s are corresponding %f %s (Rate: %f)\n",
		converter.Amount,
		converter.BaseCurrency,
		amountTargetCurrency,
		converter.TargetCurrency,
		exchangeRate)
	os.Exit(0)
}

// printSupportedCurrencies is printing the supported currencies on stdout.
func printSupportedCurrencies() {
	fmt.Println("Supported currencies:")
	sort.Strings(currencyconverter.SupportedCurrencies)
	fmt.Println(fmt.Sprintf("  - %s",
		strings.Join(currencyconverter.SupportedCurrencies, "\n  - ")))
	os.Exit(0)
}

// printVersion is printing the version information on stdout.
func printVersion() {
	fmt.Printf("Version:    %s\nGit commit: %s\nBuild at:   %s",
		version,
		gitCommit,
		buildDate)
	os.Exit(0)
}

// validateOptions is validating the CLI option values.
func validateOptions() error {
	if err := currencyconverter.IsSupportedCurrency(cliConvertOptions.BaseCurrency); err != nil {
		return err
	}
	if err := currencyconverter.IsSupportedCurrency(cliConvertOptions.TargetCurrency); err != nil {
		return err
	}
	if err := currencyconverter.IsValidAmount(cliConvertOptions.Amount); err != nil {
		return err
	}
	if err := currencyconverter.IsValidDate(cliConvertOptions.Date); err != nil {
		return err
	}
	if cliCacheOptions.CacheDirectory != "" {
		if err := currencyconverter.IsValidDirectory(cliCacheOptions.CacheDirectory); err != nil {
			return err
		}
	}
	return nil
}
