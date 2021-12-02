package slicer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceOfInt64_UnmarshalJSON(t *testing.T) {
	type TestType struct {
		Data SliceOfInt64 `json:"data"`
	}

	testCase := []struct{
		input string
		expect TestType
		isErr bool
	}{
		{
			input: `{"data": 1}`,
			expect: TestType{SliceOfInt64From(1)},
		},
		{
			input: `{"data": "1"}`,
			expect: TestType{},
			isErr: true,
		},
		{
			input: `{"data": [1,2,3,4,5]}`,
			expect: TestType{SliceOfInt64From(1,2,3,4,5)},
		},
		{
			input: `{"data": 20.5}`,
			expect: TestType{},
			isErr: true,
		},
	}

	for _, test := range testCase {
		var res TestType
		if err := json.Unmarshal([]byte(test.input), &res); err != nil {
			if !test.isErr {
				t.Error(err)
			}
		}
		assert.Equalf(t, test.expect, res, "expecting: %v, got %v", test.expect, res)
	}
}

func TestSliceOfInt_UnmarshalJSON(t *testing.T) {
	type TestType struct {
		Data SliceOfInt `json:"data"`
	}

	testCase := []struct{
		input string
		expect TestType
		isErr bool
	}{
		{
			input: `{"data": 1}`,
			expect: TestType{[]int{1}},
		},
		{
			input: `{"data": "1"}`,
			expect: TestType{},
			isErr: true,
		},
		{
			input: `{"data": [1,2,3,4,5]}`,
			expect: TestType{[]int{1,2,3,4,5}},
		},
		{
			input: `{"data": 20.5}`,
			expect: TestType{},
			isErr: true,
		},
		{
			input: `{"data": null}`,
			expect: TestType{Data: nil},
		},
		{
			input: `{"data": []}`,
			expect: TestType{Data: []int{}},
		},
		{
			input: `{"salah": true}`,
			expect: TestType{Data: nil},
		},
	}

	var i = 0
	for _, test := range testCase {
		i++
		var res TestType
		if err := json.Unmarshal([]byte(test.input), &res); err != nil {
			if !test.isErr {
				t.Error(err)
			}
		}
		assert.Equalf(t, test.expect, res, "expecting: %v, got %v", test.expect, res)
	}
}
