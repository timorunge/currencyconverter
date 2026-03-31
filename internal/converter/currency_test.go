package converter

import "testing"

func TestValidateCurrency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{"EUR", "EUR", false},
		{"USD", "USD", false},
		{"GBP", "GBP", false},
		{"Unknown", "XYZ", true},
		{"HRK", "HRK", true},
		{"Empty", "", true},
		{"Lowercase", "eur", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateCurrency(tt.currency)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateCurrency(%q) error = %v, wantErr %v", tt.currency, err, tt.wantErr)
			}
		})
	}
}
