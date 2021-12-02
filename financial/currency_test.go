package formatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCurrency(t *testing.T) {
	input := []float64{
		15700000.655,
		200,
		15000,
		2450.634678,
	}

	type testOut struct {
		options []FormatOption
		symbol  string
		result  []string
	}

	output := []testOut{
		{
			options: []FormatOption{WithDecimalDigit(3)},
			symbol:  "USD",
			result: []string{
				"USD 15,700,000.655",
				"USD 200",
				"USD 15,000",
				"USD 2,450.635",
			},
		},
		{
			options: []FormatOption{WithDotThousandSeparator()},
			symbol:  "IDR",
			result: []string{
				"IDR 15.700.000,65",
				"IDR 200",
				"IDR 15.000",
				"IDR 2.450,63",
			},
		},
	}

	for _, out := range output {
		for idx, in := range input {
			result := ToCurrencyString(in, out.symbol, out.options...)
			assert.Equal(t, out.result[idx], result, "expected %q, got %q", out.result[idx], result)
		}
	}
}
