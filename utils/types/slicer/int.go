package slicer

import (
	"bytes"
	"encoding/json"
)

type SliceOfInt64 []int64

func SliceOfInt64From(values ...int64) SliceOfInt64 {
	return values
}

func (sint64 *SliceOfInt64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*sint64 = nil
		return nil
	}

	var value int64
	var slice []int64

	valErr := json.Unmarshal(data, &value)
	slErr := json.Unmarshal(data, &slice)

	if valErr != nil && slErr != nil {
		return slErr
	}

	if valErr == nil && slErr != nil {
		*sint64 = SliceOfInt64From(value)
		return nil
	}

	*sint64 = slice
	return nil
}

func (sint64 *SliceOfInt64) GetUnique() SliceOfInt64 {
	valMap := make(map[int64]struct{})
	var newSlice SliceOfInt64
	for _, val := range *sint64 {
		if _, ok := valMap[val]; !ok {
			valMap[val] = struct{}{}
			newSlice = append(newSlice, val)
		}
	}
	return newSlice
}

func (sint64 *SliceOfInt64) ValueOrZero() []int64 {
	if *sint64 == nil {
		return []int64{}
	}
	return *sint64
}

type SliceOfInt []int

func SliceOfIntFrom(values ...int) SliceOfInt {
	return values
}

func (sint *SliceOfInt) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*sint = nil
		return nil
	}

	var value int
	var slice []int

	valErr := json.Unmarshal(data, &value)
	slErr := json.Unmarshal(data, &slice)

	if valErr != nil && slErr != nil {
		return slErr
	}

	if valErr == nil && slErr != nil {
		*sint = SliceOfIntFrom(value)
		return nil
	}

	*sint = SliceOfIntFrom(slice...)
	return nil
}

func (sint *SliceOfInt) GetUnique() SliceOfInt {
	valMap := make(map[int]struct{})
	var newSlice SliceOfInt
	for _, val := range *sint {
		if _, ok := valMap[val]; !ok {
			valMap[val] = struct{}{}
			newSlice = append(newSlice, val)
		}
	}
	return newSlice
}

func (sint *SliceOfInt) ValueOrZero() []int {
	if *sint == nil {
		return []int{}
	}
	return *sint
}