package csv

import "fmt"

const maxRecords int = 1000

var errMaxRecordsExceeded error = fmt.Errorf("total maximum number records exceeded. maximum records number: %v", maxRecords)

func strVal(val interface{}) string {
	switch val.(type) {
	case float64:
		return fmt.Sprintf("%.0f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
