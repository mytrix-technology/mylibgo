package datastore

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type expectedstruct struct {
	query string
	args []interface{}
}
var curTime = time.Now()
var testDB  = []struct{
		input interface{}
		expected expectedstruct
	}{
		{
			input:    &TestModel{
				ID:          1,
				Name:        "satu",
				CreatedTime: curTime,
				UpdatedTime: curTime,
			},
			expected: expectedstruct{
				"INSERT INTO sch.test_table (id,name,created_time,updated_time) VALUES (?, ?, ?, ?) ON CONFLICT ON CONSTRAINT test_table_pkey DO UPDATE SET name = excluded.name,updated_time = excluded.updated_time",
				[]interface{}{1,"satu", curTime, curTime},
			},
		},
		{
			input:    map[string]interface{}{
				"id": 1,
				"name": "satu",
				"created_time": curTime,
				"updated_time": curTime,
			},
			expected: expectedstruct{
				"INSERT INTO sch.test_table (id,name,created_time,updated_time) VALUES (?, ?, ?, ?) ON CONFLICT ON CONSTRAINT test_table_pkey DO UPDATE SET name = excluded.name,updated_time = excluded.updated_time",
				[]interface{}{1,"satu", curTime, curTime},
			},
		},
	}

var testDBMapOnly  = []struct{
		input map[string]interface{}
		expected expectedstruct
	}{
		{
			input:    map[string]interface{}{
				"id": 1,
				"name": "satu",
				"created_time": curTime,
				"updated_time": curTime,
			},
			expected: expectedstruct{
				"UPDATE sch.test_table SET name = ?,created_time = ?,updated_time = ? WHERE id = ?",
				[]interface{}{"satu", curTime, curTime, 1},
			},
		},
	}

func TestTable_CreateUpsertQuery(t *testing.T) {
	tb := TestModelTable
	for _, test := range testDB {
		q, args, err := tb.CreateUpsertQuery(test.input, "created_time")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.expected.args, args)
		assert.Equal(t, test.expected.query, q)
	}
}

func TestTable_CreateUpdateQuery(t *testing.T) {
	tb := TestModelTable
	for _, test := range testDBMapOnly {
		q, args, err := tb.CreateUpdateQuery(test.input)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, test.expected.args, args)
		assert.Equal(t, test.expected.query, q)
	}
}

func TestTable_createColumnsAndValues(t *testing.T) {
	tb := TestModelTable
	for _, test := range testDB {
		q, args, err := tb.createColumnsAndValues(test.input, false)
		if err != nil {
			t.Error(err)
		}
		t.Logf("columns: %+v\nargs: %+v", q, args)
	}
}

type TestModel struct {
	ID int `db:"id"`
	Name string `db:"name"`
	CreatedTime time.Time `db:"created_time"`
	UpdatedTime time.Time `db:"updated_time"`
}

var TestModelTable = &Table{
	Name:        "test_table",
	Schema:      "sch",
	Columns:     []Column{
		{
			Name:        "id",
			Type:        FIELD_INT_AUTO,
			AllowNull:   false,
		},
		{
			Name:        "name",
			Type:        FIELD_VARCHAR,
			Length:      30,
			AllowNull:   false,
		},
		{
			Name:        "created_time",
			Type:        FIELD_TIMESTAMPTZ,
			AllowNull:   false,
		},
		{
			Name:        "updated_time",
			Type:        FIELD_TIMESTAMPTZ,
			AllowNull:   false,
		},
	},
	Indexes:     []Index{
		{
			Name:   "name",
			Fields: []string{"name"},
		},
	},
	PrimaryKey:  []string{"id"},
	Model:       TestModel{},
}
