package datastore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var testTable = Table{
	Name: "testable",
	Schema: "testschema",
	Columns: []Column{
		Column{
			Name: "id",
			Type: FIELD_INT_AUTO,
			AllowNull: false,
		},
		Column{
			Name: "test_column",
			Type: FIELD_VARCHAR,
			Length: 16,
			AllowNull: true,
		},
	},
	PrimaryKey: []string{"id"},
}

var jsonData = `
{
  "name": "testable",
  "schema": "testschema",
  "columns": [
    {
      "name": "id",
      "type": "FIELD_INT_AUTO",
      "allow_null": false
    },
    {
      "name": "test_column",
      "type": "FIELD_VARCHAR",
      "length": 16,
      "allow_null": true
    }
  ],
  "primary_key": ["id"]
}
`

var yamlData = `
name: testable
schema: testschema
columns:
    - name: id
      type: FIELD_INT_AUTO
      allow_null: false

    - name: test_column
      type: "FIELD_VARCHAR"
      length: 16
      allow_null: true
primary_key:
  - id
`

func TestTableParsing(t *testing.T) {

	t.Run("Testing json parsing", func(t *testing.T){
		table, err := ParseTableFromJson([]byte(jsonData))
		if err != nil {
			t.Errorf("Error while parsing the data: %v", err)
		}

		assert.Equal(t, testTable, *table, fmt.Sprintf("got: %+v\n\nwant: %+v\n", table, testTable))
	})

	t.Run("Testing YAML parsing", func(t *testing.T){
		table, err := ParseTableFromYaml([]byte(yamlData))
		if err != nil {
			t.Errorf("Error while parsing the data: %v", err)
		}

		assert.Equal(t, testTable, *table, fmt.Sprintf("got: %+v\n\nwant: %+v\n", table, testTable))
	})
}

