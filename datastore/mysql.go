package datastore

import (
	"fmt"
	"strings"
)

type MySqlHelper struct{}

func (h MySqlHelper) FieldTypeString(ft DataType, length int) string {
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
		str = "INT"
	case FIELD_BIGINT:
		str = "BIGINT"
	case FIELD_BIGINT_AUTO:
		str = "BIGINT"
	case FIELD_TINYINT:
		str = "TINYINT"
	case FIELD_SMALLINT:
		str = "SMALLINT"
	case FIELD_FLOAT:
		str = "FLOAT"
	case FIELD_DOUBLE:
		str = "DOUBLE"
	case FIELD_DATE:
		str = "DATE"
	case FIELD_TIMESTAMP:
		str = "DATETIME"
	case FIELD_TIMESTAMPTZ:
		str = "DATETIME"
	case FIELD_BOOLEAN:
		str = "BOOLEAN"
	}

	if length > 0 {
		str += fmt.Sprintf("(%v)", length)
	}

	return str
}

func (h MySqlHelper) GenerateCreateSchemaQuery(schemaName string) []string {
	qries := []string{}
	//qries = append(qries, fmt.Sprintf("CREATE SCHEMA %s", strings.ToLower(schemaName)))
	return qries
}

func (h MySqlHelper) GenerateCreateTableQuery(table *Table) []string {
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

		switch col.Type {
		case FIELD_INT_AUTO:
			allowNull += " AUTO_INCREMENT"
		case FIELD_BIGINT_AUTO:
			allowNull += " AUTO_INCREMENT"
		}

		columns += fmt.Sprintf("%s %s %s", col.Name, fieldType, allowNull)
	}

	primary := ""
	if len(table.PrimaryKey) > 0 {
		for _, f := range table.PrimaryKey {
			if len(primary) > 0 {
				primary += ","
			}

			primary += f
		}
	}

	indexes := []string{}
	for _, i := range table.Indexes {
		indexes = append(indexes, fmt.Sprintf("INDEX (%s)", strings.Join(i.Fields, ",")))
	}

	qry := fmt.Sprintf("CREATE TABLE %s.%s (%s", table.Schema, table.Name, columns)
	if primary != "" {
		qry += ",PRIMARY KEY(" + primary + ")"
	}

	if len(indexes) > 0 {
		qry += "," + strings.Join(indexes, ",")
	}
	qry += ")"
	qries = append(qries, qry)

	return qries
}

func (h MySqlHelper) GenerateCheckSchemaQuery(name string) string {
	return ""
}

func (h MySqlHelper) GenerateCheckTableQuery(schema, name string) string {
	return fmt.Sprintf("select * from information_schema.tables where table_schema = '%s' and table_name = '%s';", schema, name)
}

func (h MySqlHelper) GenerateInsertQuery(table *Table, cols []string, values []string) string {
	return fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)", table.Schema, table.Name, strings.Join(cols,","), strings.Join(values, ","))
}

func (h MySqlHelper) BuildInsertQuery(schema, table string, columns []string, values [][]interface{}, returnfields ...string) string {
	return ""
}