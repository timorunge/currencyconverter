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

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	baseCurrency := strings.ToUpper(cliConvertOptions.BaseCurrency)
	targetCurrency := strings.ToUpper(cliConvertOptions.TargetCurrency)
	if cliConvertOptions.Reverse {
		baseCurrency, targetCurrency = targetCurrency, baseCurrency
	}

	if err := currencyconverter.IsValidAmount(cliConvertOptions.Amount); err != nil {
		log.Fatal(err)
	}
	if err := currencyconverter.IsSupportedCurrency(baseCurrency); err != nil {
		log.Fatal(err)
	}
	if err := currencyconverter.IsSupportedCurrency(targetCurrency); err != nil {
		log.Fatal(err)
	}
	if err := currencyconverter.IsValidDate(cliConvertOptions.Date); err != nil {
		log.Fatal(err)
	}

	fileCache := currencyconverter.NewFileCache()
	if cliCacheOptions.CacheDirectory != "" {
		fileCache.SetDirectory(cliCacheOptions.CacheDirectory)
		if err := fileCache.IsValidDirectory(); err != nil {
			log.Fatal(err)
		}
	}
	fileCache.SetEnabled(!cliCacheOptions.NoCache)
	fileCache.SetTimeout(time.Duration(cliCacheOptions.CacheTimeout) * time.Minute)
	fileCache.SetFilename(fmt.Sprintf("%s-%s", currencyconverter.FileCacheFilename, "latest"))

	api := currencyconverter.NewAPI()

	exchangeRates, err := func() (currencyconverter.ExchangeRates, error) {
		if cliConvertOptions.Date != "latest" {
			api.SetHistoricalData(true)
			fileCache.SetFilename(fmt.Sprintf("%s-%s", currencyconverter.FileCacheFilename, cliConvertOptions.Date))
			fileCache.SetTimeout(365 * 24 * time.Hour)
		}
		exchangeRates, err := fileCache.Get()
		if err != nil {
			exchangeRates, err := api.Get()
			if err != nil {
				return exchangeRates, err
			}
			dailyRates, err := exchangeRates.GetDate(cliConvertOptions.Date)
			if err != nil {
				return exchangeRates, err
			}
			finalRates := *currencyconverter.NewExchangeRates()
			finalRates.AddDate(dailyRates)
			if err := fileCache.Write(finalRates); err != nil {
				log.Print(err)
			}
			return finalRates, err
		}
		return exchangeRates, err
	}()
	if err != nil {
		log.Fatal(err)
	}

	currencies := currencyconverter.Currencies{}
	currencies.SetSupportedCurrencies(currencyconverter.SupportedCurrencies)

	converter := currencyconverter.NewConverter()
	converter.SetAmount(cliConvertOptions.Amount)
	converter.SetBaseCurrency(baseCurrency)
	converter.SetCurrencies(currencies)
	converter.SetDate(cliConvertOptions.Date)
	converter.SetExchangeRates(exchangeRates)
	converter.SetTargetCurrency(targetCurrency)

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
}

// printHelp is printing the help on stdout.
func printHelp() {
	parser.WriteHelp(os.Stdout)
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
