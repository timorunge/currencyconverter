# currencyconverter

[![CI](https://github.com/timorunge/currencyconverter/actions/workflows/ci.yml/badge.svg)](https://github.com/timorunge/currencyconverter/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/timorunge/currencyconverter)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/timorunge/currencyconverter)](https://goreportcard.com/report/github.com/timorunge/currencyconverter)
[![License](https://img.shields.io/github/license/timorunge/currencyconverter)](LICENSE)
[![Release](https://img.shields.io/github/v/release/timorunge/currencyconverter)](https://github.com/timorunge/currencyconverter/releases)

A simple CLI tool written in Go for calculating exchange rates using
foreign exchange reference rates published by the
[European Central Bank](https://www.ecb.europa.eu/stats/eurofxref/).

## Installation

### Go install

```bash
go install github.com/timorunge/currencyconverter/cmd/currencyconverter@latest
```

### Binary download

Download pre-built binaries for your platform from the
[releases page](https://github.com/timorunge/currencyconverter/releases).

### Build from source

```bash
git clone https://github.com/timorunge/currencyconverter.git
cd currencyconverter
make build
```

## Quick Start

Convert 100 EUR to USD:

```bash
currencyconverter --amount 100 --from EUR --to USD
```

```
100 EUR ≈ 111.51 USD (2019-07-31, ECB)
```

Convert for a specific date:

```bash
currencyconverter --amount 500 --from USD --to GBP --date 2019-07-31
```

Swap currencies with `--reverse`:

```bash
currencyconverter --from EUR --to USD --reverse
```

Get JSON output:

```bash
currencyconverter --amount 100 --from EUR --to USD --json
```

```json
{"amount":100,"from":"EUR","to":"USD","result":111.51,"date":"2019-07-31","source":"ECB"}
```

## Usage

```
currencyconverter - Exchange rate calculator using ECB reference rates

Usage:
  currencyconverter [OPTIONS]

Options:
  -f, --from string            Currency to use as base (default "EUR")
  -t, --to string              Currency to convert to (default "USD")
  -a, --amount float           Amount to calculate (default 1)
  -d, --date string            Date for the calculation (latest|YYYY-MM-DD) (default "latest")
  -r, --reverse                Swap from and to values
      --json                   Output result as JSON
      --no-cache               Disable caching
  -h, --help                   Show this help message
      --version                Show the version of currencyconverter
      --supported-currencies   Show a list with all supported currencies
```

## Supported Currencies

AUD, BGN, BRL, CAD, CHF, CNY, CZK, DKK, EUR, GBP, HKD, HUF, IDR,
ILS, INR, ISK, JPY, KRW, MXN, MYR, NOK, NZD, PHP, PLN, RON, SEK,
SGD, THB, TRY, USD, ZAR

Run `currencyconverter --supported-currencies` for the full list.

## Caching

Exchange rates are cached locally as JSON files in the OS temp
directory. Latest rates are cached until the top of the next hour,
keeping network calls low while staying current with ECB updates.
Historical dates are cached indefinitely since their rates never
change. Disable caching with `--no-cache`.

## Development

```bash
make help     # Show all available targets
make check    # Run all quality gates (fmt, tidy, vet, lint, test)
make lint     # Run golangci-lint
make test     # Run tests with race detector
make build    # Build static binary
```

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)
