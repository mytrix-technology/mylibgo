package slicer

import (
	"bytes"
	"encoding/json"
)

var nullBytes = []byte("null")

type SliceOfString []string

func SliceOfStringFrom(values ...string) SliceOfString {
	return values
}

func (str *SliceOfString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*str = nil
		return nil
	}

	var value string
	var slice []string

	valErr := json.Unmarshal(data, &value)
	slErr := json.Unmarshal(data, &slice)

	if valErr != nil && slErr != nil {
		return slErr
	}

	if valErr == nil && slErr != nil {
		*str = SliceOfStringFrom(value)
		return nil
	}

	*str = slice
	return nil
}

func (str *SliceOfString) GetUnique() SliceOfString {
	valMap := make(map[string]struct{})
	var newSlice SliceOfString
	for _, val := range *str {
		if _, ok := valMap[val]; !ok {
			valMap[val] = struct{}{}
			newSlice = append(newSlice, val)
		}
	}
	return newSlice
}

func (str *SliceOfString) ValueOrZero() []string {
	if *str == nil {
		return []string{}
	}
	return *str
}
