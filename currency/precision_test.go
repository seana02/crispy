package currency

import "testing"

func TestPrecisionFor(t *testing.T) {
	cases := []struct {
		currency  string
		precision int32
	}{
		{"USD", 2},
		{"EUR", 2},
		{"JPY", 0},
		{"GBP", 2},
		{"CNY", 2},
		{"AUD", 2},
		{"CAD", 2},
	}

	for _, c := range cases {
		t.Run(c.currency, func(t *testing.T) {
			got := PrecisionFor(c.currency)

			if got != c.precision {
				t.Errorf("Wanted %v, got %v", c.precision, got)
			}
		})
	}
}
