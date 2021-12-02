package formatter

type Option struct {
	ThousandSeparator string
	DecimalSeparator  string
	Symbol            string
	DecimalDigit      int
}

type FormatOption func(*Option)

func WithDotThousandSeparator() FormatOption {
	return func(o *Option) {
		o.ThousandSeparator = "."
		o.DecimalSeparator = ","
	}
}

func WithCommaThousandSeparator() FormatOption {
	return func(o *Option) {
		o.ThousandSeparator = ","
		o.DecimalSeparator = "."
	}
}

func WithSymbol(symbol string) FormatOption {
	return func(o *Option) {
		o.Symbol = symbol
	}
}

func WithDecimalDigit(digit int) FormatOption {
	return func(o *Option) {
		o.DecimalDigit = digit
	}
}
