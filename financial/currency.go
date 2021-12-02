package formatter

import "strconv"

//ToCurrencyString format float64 into string with thousand separator and currency symbol
func ToCurrencyString(amount float64, symbol string, options ...FormatOption) string {
	opt := Option{
		ThousandSeparator: ",",
		DecimalSeparator:  ".",
		DecimalDigit:      2,
	}

	for _, op := range options {
		op(&opt)
	}

	sa := strconv.FormatInt(int64(amount), 10)
	frac := amount - float64(int64(amount))
	if symbol != "" {
		symbol = symbol + " "
	}

	for i := len(sa) - 3; i > 0; i -= 3 {
		sa = sa[:i] + opt.ThousandSeparator + sa[i:]
	}

	if frac > 0 {
		sfrac := strconv.FormatFloat(frac, 'f', opt.DecimalDigit, 64)
		sa = sa + opt.DecimalSeparator + sfrac[2:]
	}

	return symbol + sa
}
