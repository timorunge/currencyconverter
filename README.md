# currencyconverter

[![Go Report](https://goreportcard.com/badge/github.com/timorunge/currencyconverter)](https://goreportcard.com/report/github.com/timorunge/currencyconverter)
[![Build Status](https://travis-ci.org/timorunge/currencyconverter.svg?branch=master)](https://travis-ci.org/timorunge/currencyconverter)

`currencyconverter` is a simple CLI tool written in Go for calculating the
exchange rate using foreign exchange reference rates published by the
[European Central Bank](https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html).

## Install

You can use an
[official release](https://github.com/timorunge/currencyconverter/releases) of
`currencyconverter`. The tarballs for each release contain the
`currencyconverter` CLI applicaton.

Copy the binary in your `$PATH` or call it directly via
`$YOURDIR/currencyconverter`.

To get the latest version of `currencyconverter` just run `go get`.

```sh
go get github.com/timorunge/currencyconverter/cmd/currencyconverter
```

If `$GOPATH/bin` is not in your `$PATH` call `currencyconverter` directly via
`$GOPATH/bin/currencyconverter`.

Last but not least you also have the possibility to to clone this repository
and use [Mage](https://magefile.org/) to `test`, `build` and `install`
`currencyconverter`.

## Usage

```sh
Usage:
  currencyconverter [OPTIONS]

currencyconverter is calculating the exchange rate using foreign exchange reference rates published by the European Central Bank

Convert options:
  -a, --amount=               Amount to calculate (default: 1)
  -f, --from=                 Currency to use as base (default: EUR)
  -t, --to=                   Currency to convert too (default: USD)
  -r, --reverse               Swap from and to values

Cache options:
      --cache-directory=      Directory to store cached responses (default: Operating system default temp directory)
      --cache-timeout=        Timeout in minutes to invalidate the cache (default: 60)
      --no-cache              Disable caching

Help options:
      --supported-currencies  Show a list with all supported currencies
  -v, --version               Show the version of currencyconverter
  -h, --help                  Show this help message
 ```

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

## Author Information

- Timo Runge
