package types

import (
	"encoding/json"
	"testing"
)

func TestSequenceUnmarshalJSON(t *testing.T) {
	testDb := []struct {
		input    string
		expected Sequence
		isErr    bool
	}{
		{`1`, 1, false}, {`"1"`, 1, false}, {`"1st"`, 1, false}, {`0`, 0, true}, {`"2nd"`, 2, false}, {`"20rd"`, 0, true}, {`"zisth"`, 0, true},
	}

	for _, test := range testDb {
		var s Sequence
		if err := json.Unmarshal([]byte(test.input), &s); err != nil {
			if !test.isErr {
				t.Error(err)
			}
		}

		if test.expected != s {
			t.Errorf("expecting %d, got %d", test.expected, s)
		}
	}
}

func TestSequenceString(t *testing.T) {
	testDb := []struct {
		input    Sequence
		expected string
	}{
		{1, "1st"}, {2, "2nd"}, {3, "3rd"}, {4, "4th"}, {21, "21st"}, {53, "53rd"}, {20, "20th"},
	}

	for _, test := range testDb {
		if test.expected != test.input.String() {
			t.Errorf("expecting %s, got %s", test.expected, test.input.String())
		}
	}
}
