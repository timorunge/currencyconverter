package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []string
		wantCode exitCode
		wantOut  string
		wantErr  string
	}{
		{"Help", []string{"--help"}, exitSuccess, "--amount", ""},
		{"Version", []string{"--version"}, exitSuccess, "v1.2.3", ""},
		{"SupportedCurrencies", []string{"--supported-currencies"}, exitSuccess, "USD", ""},
		{"UnknownFlag", []string{"--nonexistent"}, exitUsageError, "", ""},
		{"PositionalArgs", []string{"foo", "bar"}, exitUsageError, "", "unexpected positional arguments"},
		{"InvalidCurrency", []string{"--from", "XXX"}, exitUsageError, "", "not a supported currency"},
		{"InvalidAmount", []string{"--amount", "0"}, exitUsageError, "", "amount must be greater than zero"},
		{"NegativeAmount", []string{"--amount", "-5"}, exitUsageError, "", "amount must be greater than zero"},
		{"InvalidDateFormat", []string{"--date", "not-a-date"}, exitUsageError, "", "not a valid date"},
		{"WeekendDate", []string{"--date", "2019-07-27"}, exitUsageError, "", "Saturday"},
		{"MultipleErrors", []string{"--from", "XXX", "--amount", "0"}, exitUsageError, "", "not a supported currency"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			code := run(context.Background(), &stdout, &stderr, tt.args, "v1.2.3")
			if code != tt.wantCode {
				t.Fatalf("exit code = %d, want %d; stderr: %s", code, tt.wantCode, stderr.String())
			}
			if tt.wantOut != "" && !strings.Contains(stdout.String(), tt.wantOut) {
				t.Fatalf("stdout missing %q; got: %s", tt.wantOut, stdout.String())
			}
			if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
				t.Fatalf("stderr missing %q; got: %s", tt.wantErr, stderr.String())
			}
		})
	}
}

func TestValidateOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		amount  float64
		from    string
		to      string
		date    string
		wantErr string
	}{
		{
			name: "ValidDefaults", amount: 1, from: "EUR", to: "USD", date: "latest",
		},
		{
			name: "ValidHistorical", amount: 100, from: "GBP", to: "JPY", date: "2019-07-31",
		},
		{
			name: "InvalidFrom", amount: 1, from: "XXX", to: "USD", date: "latest",
			wantErr: "not a supported currency",
		},
		{
			name: "InvalidTo", amount: 1, from: "EUR", to: "ZZZ", date: "latest",
			wantErr: "not a supported currency",
		},
		{
			name: "ZeroAmount", amount: 0, from: "EUR", to: "USD", date: "latest",
			wantErr: "amount must be greater than zero",
		},
		{
			name: "InvalidDate", amount: 1, from: "EUR", to: "USD", date: "bad",
			wantErr: "not a valid date",
		},
		{
			name: "MultipleErrors", amount: 0, from: "XXX", to: "ZZZ", date: "bad",
			wantErr: "not a supported currency",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateOptions(tt.amount, tt.from, tt.to, tt.date)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %q, want to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestRunJSONOutput(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := run(context.Background(), &stdout, &stderr, []string{
		"--amount", "100", "--from", "EUR", "--to", "EUR", "--json", "--no-cache",
	}, "v1.2.3")
	if code != exitSuccess {
		t.Fatalf("exit code = %d, want %d; stderr: %s", code, exitSuccess, stderr.String())
	}

	var result conversionResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v; raw: %s", err, stdout.String())
	}
	if result.Amount != 100 {
		t.Fatalf("Amount = %v, want 100", result.Amount)
	}
	if result.From != "EUR" || result.To != "EUR" {
		t.Fatalf("From/To = %s/%s, want EUR/EUR", result.From, result.To)
	}
	if result.Result != 100 {
		t.Fatalf("Result = %v, want 100 (EUR→EUR)", result.Result)
	}
	if result.Source != "ECB" {
		t.Fatalf("Source = %q, want %q", result.Source, "ECB")
	}
	if result.Date == "" {
		t.Fatal("Date is empty")
	}
}

func TestRunReverse(t *testing.T) {
	t.Parallel()

	run := func(args ...string) string {
		var stdout, stderr bytes.Buffer
		code := run(context.Background(), &stdout, &stderr, args, "v1.2.3")
		if code != exitSuccess {
			t.Fatalf("exit code = %d; stderr: %s", code, stderr.String())
		}
		return stdout.String()
	}

	normal := run("--amount", "100", "--from", "EUR", "--to", "EUR", "--no-cache")
	reversed := run("--amount", "100", "--from", "EUR", "--to", "EUR", "--reverse", "--no-cache")

	// EUR→EUR reversed is still EUR→EUR, so output should match.
	if normal != reversed {
		t.Fatalf("reverse of EUR→EUR should be identical:\nnormal:   %sreversed: %s", normal, reversed)
	}
}

func TestRoundTo4(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"Exact", 1.1151, 1.1151},
		{"RoundDown", 75.58603274120514, 75.586},
		{"RoundUp", 1.23456, 1.2346},
		{"Integer", 100, 100},
		{"Zero", 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := roundTo4(tt.input)
			if got != tt.want {
				t.Fatalf("roundTo4(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{"Integer", 100, "100"},
		{"TwoDecimals", 1490.89, "1490.89"},
		{"FourDecimals", 1.1151, "1.1151"},
		{"TrailingZeros", 1.5000, "1.5"},
		{"LongDecimal", 0.896780557797507, "0.8968"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := formatAmount(tt.input)
			if got != tt.want {
				t.Fatalf("formatAmount(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
