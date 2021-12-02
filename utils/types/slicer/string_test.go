package slicer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceOfString_UnmarshalJSON(t *testing.T) {
	type TestType struct {
		Data SliceOfString `json:"data"`
	}

	testCase := []struct{
		input string
		expect TestType
		isErr bool
	}{
		{
			input: `{"data": "satu"}`,
			expect: TestType{SliceOfStringFrom("satu")},
		},
		{
			input: `{"data": [1,2,3,4]}`,
			expect: TestType{},
			isErr: true,
		},
		{
			input: `{"data": ["satu","dua","tiga"]}`,
			expect: TestType{SliceOfStringFrom("satu","dua","tiga")},
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