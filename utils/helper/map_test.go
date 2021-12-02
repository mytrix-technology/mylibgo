package helper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeysToArrayOfInterface(t *testing.T) {
	testDb := []struct{
		input interface{}
		expected []interface{}
		expectedType string
	}{
		{
			input:    map[int] bool {1: true, 2: true, 3:true},
			expected: []interface{}{1,2,3},
			expectedType: "int",
		},
		{
			input: map[string] bool{"satu": true, "dua": true, "tiga": true},
			expected: []interface{}{"satu", "dua", "tiga"},
			expectedType: "string",
		},
	}

	for _, test := range testDb {
		keys := MapKeysToArrayOfInterface(test.input)
		res := assert.Equal(t, len(keys), len(test.expected), fmt.Sprintf("expecting %d lenght slice, got %v", len(test.expected), keys))
		if res {
			for _, key := range keys {
				switch test.expectedType {
				case "string":
					if _, ok := key.(string); !ok {
						t.Errorf("expected key type string, got %T", key)
					}
					break
				case "int":
					if _, ok := key.(int); !ok {
						t.Errorf("expected key type int, got %T", key)
					}
					break
				}
			}
		}
	}
}
