package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RootNode struct {
	ID     int
	Name   string
	Amount float64
	Childs []interface{}
}

type ChildNode struct {
	ID     int
	Name   string
	Amount float64
}

var serializeTestData = RootNode{
	ID:     245,
	Name:   "Dua Empat Lima",
	Amount: 12345.678,
	Childs: []interface{}{
		ChildNode{
			ID:     2451,
			Name:   "Child 1",
			Amount: 12222.45,
		},
		ChildNode{
			ID:     2452,
			Name:   "Child 2",
			Amount: 22222.45,
		},
	},
}

func TestSerializeToBase64(t *testing.T) {
	encoded, err := SerializeToBase64(serializeTestData, ChildNode{})
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("encoded result is: %s", string(encoded))

	var result RootNode
	if err := DeserializeFromBase64(string(encoded), &result); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, result, serializeTestData, "expecting: %+v\nGot: %+v\n", serializeTestData, result)
}
