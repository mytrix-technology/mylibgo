package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolean_UnmarshalJSON(t *testing.T) {
	testDb := []struct {
		input string
		expected bool
		isErr bool
	}{
		{`true`, true, false},
		{`false`, false, false},
		{`"TRue"`, true, false},
		{`"False"`, false, false},
		{`1`, true, false},
		{`0`,false,false},
		{`"1"`, true, false},
		{`"0"`, false,false},
		{`2`, false, true},
		{`"satu"`, false, true},
	}

	for _, test := range testDb {
		t.Logf("testing with input: %v", test.input)
		expected := BooleanFrom(test.expected)
		if test.isErr {
			expected.Valid = false
		}
		var b Boolean
		err := json.Unmarshal([]byte(test.input), &b)
		if test.isErr && err == nil {
			t.Errorf("expecting an error from unmarshal of string %q. %s", test.input, err)
			continue
		}

		if !test.isErr && err != nil {
			t.Errorf(err.Error())
			continue
		}

		assert.Equal(t, expected, b, "expecting %v, got %v", expected, b)
	}
}
