package money

import (
	"errors"
	"testing"
)

func TestParseDecimal(t *testing.T) {
	tt := map[string]struct {
		decimal  string
		expected Decimal
		err      error
	}{
		"2 decimal digits": {
			decimal:  "1.52",
			expected: Decimal{subunits: 152, precision: 2},
			err:      nil,
		},
		"no decimal digits": {
			decimal:  "1",
			expected: Decimal{subunits: 1, precision: 0},
			err:      nil,
		},
		"suffix 0 as decimal digits": {
			decimal:  "1.50",
			expected: Decimal{subunits: 150, precision: 2},
			err:      nil,
		},
		"prefix 0 as decimal digits": {
			decimal:  "1.08",
			expected: Decimal{subunits: 108, precision: 2},
			err:      nil,
		},
		"not a number": {
			decimal: "NaN",
			err:     ErrInvalidDecimal,
		},
		"empty string": {
			decimal: "",
			err:     ErrInvalidDecimal,
		},
		"too large": {
			decimal: "1234567890123",
			err:     ErrTooLarge,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := ParseDecimal(tc.decimal)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error %v, got %v", tc.err, err)
			}

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
