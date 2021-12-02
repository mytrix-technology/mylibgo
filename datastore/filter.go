package datastore

import (
	"fmt"
)

type ComparisonString string

const (
	// EQ equal
	EQ ComparisonString = "="

	// NEQ not equal
	NEQ ComparisonString = "<>"

	// GT greater than
	GT ComparisonString = ">"

	// GTE greater than or equal
	GTE ComparisonString = ">="

	// LT less than
	LT ComparisonString = "<"

	// LTE less than or equal
	LTE ComparisonString = "<="

	// IN in array
	IN ComparisonString = "IN"

	//NOT_IN
	NOT_IN ComparisonString = "NOT IN"

	//IS
	IS ComparisonString = "IS"
)

type Filter struct {
	Field      string
	Comparison ComparisonString
	Value      interface{}
}

func (f *Filter) Encode() (text string, arg interface{}, err error) {
	value := "?"

	if f.Comparison == IN || f.Comparison == NOT_IN {
		value = "(?)"
	}

	if f.Comparison == IS || f.Value == nil {
		value = "NULL"
	}

	criteria := fmt.Sprintf("%s %s %s", f.Field, f.Comparison, value)
	return criteria, f.Value, nil
}
