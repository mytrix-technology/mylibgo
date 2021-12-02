package datastore

import (
	"fmt"
	"strings"
)

type PostgresHelper struct{}

func (h PostgresHelper) FieldTypeString(ft DataType, length int) string {
	str := ""
	switch ft {
	case FIELD_VARCHAR:
		str = "VARCHAR"
	case FIELD_CHAR:
		str = "CHAR"
	case FIELD_TEXT:
		str = "TEXT"
	case FIELD_INT:
		str = "INT"
	case FIELD_INT_AUTO:
		str = "SERIAL"
	case FIELD_BIGINT:
		str = "BIGINT"
	case FIELD_BIGINT_AUTO:
		str = "BIGSERIAL"
	case FIELD_TINYINT:
		str = "SMALLINT"
	case FIELD_SMALLINT:
		str = "SMALLINT"
	case FIELD_FLOAT:
		str = "FLOAT8"
	case FIELD_DOUBLE:
		str = "DOUBLE PRECISION"
	case FIELD_DATE:
		str = "DATE"
	case FIELD_TIMESTAMP:
		str = "TIMESTAMP"
	case FIELD_TIMESTAMPTZ:
		str = "TIMESTAMP WITH TIME ZONE"
	case FIELD_BOOLEAN:
		str = "BOOLEAN"
	}

	if length > 0 {
		str += fmt.Sprintf("(%v)", length)
	}

	return str
}

func (h PostgresHelper) GenerateCreateSchemaQuery(schemaName string) []string {
	qries := []string{}
	qries = append(qries, fmt.Sprintf("CREATE SCHEMA %s", strings.ToLower(schemaName)))
	return qries
}

func (h PostgresHelper) GenerateCreateTableQuery(table *Table) []string {
	qries := []string{}

	columns := ""
	for _, col := range table.Columns {
		if len(columns) > 0 {
			columns += ","
		}

		fieldType := h.FieldTypeString(col.Type, col.Length)
		allowNull := "NOT NULL"
		if col.AllowNull {
			allowNull = "DEFAULT NULL"
		}

		columns += fmt.Sprintf("%s %s %s", col.Name, fieldType, allowNull)
	}

	primary := ""
	if len(table.PrimaryKey) > 0 {
		primary = strings.Join(table.PrimaryKey, ",")
	}

	indexes := []string{}
	for _, i := range table.Indexes {
		indexes = append(indexes, fmt.Sprintf("CREATE INDEX %s_%s ON %s.%s USING btree (%s)", table.Name, i.Name, table.Schema, table.Name, strings.Join(i.Fields, ",")))
	}

	qry := fmt.Sprintf("CREATE TABLE %s.%s (%s", table.Schema, table.Name, columns)
	if primary != "" {
		qry += ",PRIMARY KEY(" + primary + ")"
	}

	if len(table.Constraints) > 0 {
		for _, con := range table.Constraints {
			switch con.Type {
			case CONSTRAINT_UNIQUE:
				qry += ", UNIQUE " + con.Value
			}
		}
	}

	qry += ")"
	qries = append(qries, qry)

	for _, i := range indexes {
		qries = append(qries, i)
	}
	return qries
}

func (h PostgresHelper) GenerateCheckSchemaQuery(name string) string {
	return fmt.Sprintf("select * from information_schema.schemata where schema_name = '%s';", name)
}

func (h PostgresHelper) GenerateCheckTableQuery(schema, name string) string {
	return fmt.Sprintf("select * from information_schema.tables where table_schema = '%s' and table_name = '%s';", schema, name)
}

func (h PostgresHelper) GenerateInsertQuery(table *Table, cols []string, values []string) string {
	return fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)", table.Schema, table.Name, strings.Join(cols,","), strings.Join(values, ","))
}

func (h PostgresHelper) BuildInsertQuery(schema, table string, columns []string, values [][]interface{}, returnfields ...string) string {
	return ""
}