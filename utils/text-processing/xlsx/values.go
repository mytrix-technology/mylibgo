package xlsx

import (
	"fmt"
	"time"
)

func NewEglsReportFilename(date time.Time) string {
	val := fmt.Sprintf("[name_file].xlsx", date.Year(), date.Month(), date.Day())
	return val
}
