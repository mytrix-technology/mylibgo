package repository

import (
	"testing"

	"github.com/mytrix-technology/mylibgo/datastore"
)

type TestStruct struct {
	ID int `db:"id"`
	Name string `db:"name"`
}

var TestStructTable = datastore.Table{
	Name:        "test_struct_table",
	Schema:      "test_schema",
	Columns:     []datastore.Column{
		{
			Name:        "id",
			Type:        datastore.FIELD_INT,
		},
		{
			Name:        "name",
			Type:        datastore.FIELD_VARCHAR,
			Length:      64,
		},
	},
	PrimaryKey:  []string{"id"},
	Model:       TestStruct{},
}
func TestRepository_buildInsertQueryWithArgs(t *testing.T) {

}
