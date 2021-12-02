package commons

import (
	"math"

	"github.com/leekchan/accounting"
)

func ToIDRFormat(i float64) string {
	return ToIDRFormatWithDecimal(i, 2)
}

func ToIDRFormatWithDecimal(f float64, precission int) string {
	ac := accounting.Accounting{
		Symbol:    "Rp.",
		Precision: precission,
		Thousand:  ".",
		Decimal:   ",",
	}
	return ac.FormatMoney(f)
}

func Round(val float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return float64(int64(val*pow)) / pow
}
