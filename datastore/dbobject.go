package datastore

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

type DataType int
type ConstraintType int

const (
	FIELD_VARCHAR DataType = iota
	FIELD_CHAR
	FIELD_TEXT
	FIELD_INT
	FIELD_INT_AUTO
	FIELD_BIGINT
	FIELD_BIGINT_AUTO
	FIELD_TINYINT
	FIELD_SMALLINT
	FIELD_FLOAT
	FIELD_DOUBLE
	FIELD_DATE
	FIELD_TIMESTAMP
	FIELD_TIMESTAMPTZ
	FIELD_BOOLEAN
)

var dataTypeStrings = []string{
	"FIELD_VARCHAR",
	"FIELD_CHAR",
	"FIELD_TEXT",
	"FIELD_INT",
	"FIELD_INT_AUTO",
	"FIELD_BIGINT",
	"FIELD_BIGINT_AUTO",
	"FIELD_TINYINT",
	"FIELD_SMALLINT",
	"FIELD_FLOAT",
	"FIELD_DOUBLE",
	"FIELD_DATE",
	"FIELD_TIMESTAMP",
	"FIELD_TIMESTAMPTZ",
	"FIELD_BOOLEAN"}

var dataTypeKind = []reflect.Kind{
	reflect.String,
	reflect.String,
	reflect.String,
	reflect.Int,
	reflect.Int,
	reflect.Int64,
	reflect.Int64,
	reflect.Int8,
	reflect.Int16,
	reflect.Float32,
	reflect.Float64,
	reflect.Struct,
	reflect.Struct,
	reflect.Struct,
	reflect.Bool,
}

const (
	CONSTRAINT_UNIQUE ConstraintType = iota
	CONSTRAINT_CHECK
	CONSTRAINT_FOREIGN
	CONSTRAINT_NOTNULL
)

type Column struct {
	Name        string       `json:"name" yaml:"name"`
	Type        DataType     `json:"type" yaml:"type"`
	Length      int          `json:"length" yaml:"length"`
	AllowNull   bool         `json:"allow_null" yaml:"allow_null"`
	Default     interface{}  `json:"default" yaml:"default"`
	Constraints []Constraint `json:"contraints" yaml:"constraints"`
}

type Index struct {
	Name   string   `json:"name" yaml:"name"`
	Fields []string `json:"fields" yaml:"fields"`
}

type Constraint struct {
	Type  ConstraintType `json:"type" yaml:"type"`
	Value string         `json:"value" yaml:"value"`
}

type Table struct {
	Name        string       `json:"name" yaml:"name"`
	Schema      string       `json:"schema" yaml:"schema"`
	Columns     []Column     `json:"columns" yaml:"columns"`
	Indexes     []Index      `json:"indexes" yaml:"indexes"`
	PrimaryKey  []string     `json:"primary_key" yaml:"primary_key"`
	Constraints []Constraint `json:"constraints" yaml:"constraints"`
	Model       interface{}
}

func (t *Table) GetColumnNames() []string {
	columnNames := make([]string, len(t.Columns))
	for i := 0; i < len(t.Columns); i++ {
		columnNames[i] = t.Columns[i].Name
	}

	return columnNames
}

func (t *Table) GetColumn(name string) (*Column, bool) {
	for _, col := range t.Columns {
		if col.Name == name {
			return &col, true
		}
	}

	return nil, false
}

func (t *Table) GetPrimaryKeyColumnNames() []string {
	return t.PrimaryKey
}

func (t *Table) IsPrimary(colName string) bool {
	for _, col := range t.PrimaryKey {
		if col == colName {
			return true
		}
	}
	return false
}

func (t *Table) CheckValidity() error {
	if !t.primaryKeyIsValid() {
		return fmt.Errorf("primary key definition is invalid")
	}

	return nil
}

func (t *Table) primaryKeyIsValid() bool {
	for _, key := range t.PrimaryKey {
		if _, ok := t.GetColumn(key); !ok {
			return false
		}
	}
	return true
}

func (t *Table) CreateInsertQuery(value interface{}) (query string, args []interface{}, err error) {
	columns, args, err := t.createColumnsAndValues(value, true)
	if err != nil {
		return "", nil, err
	}

	tableName := fmt.Sprintf("%s.%s", t.Schema, t.Name)
	colString := ""
	for _, col := range columns {
		if len(colString) > 0 {
			colString += ","
		}
		colString += col.Name
	}
	qry := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)",
		tableName, colString)

	return sqlx.In(qry, args)
}

func (t *Table) CreateUpsertQuery(value interface{}, excludeUpdate ...string) (query string, args []interface{}, err error) {
	columns, cargs, err := t.createColumnsAndValues(value, false)
	if err != nil {
		return "", nil, err
	}

	tableName := fmt.Sprintf("%s.%s", t.Schema, t.Name)
	colString := ""
	update := ""
	for _, col := range columns {
		if len(colString) > 0 {
			colString += ","
		}
		colString += col.Name

		if t.IsPrimary(col.Name) {
			continue
		}

		excluded := false
		for _, exc := range excludeUpdate {
			fmt.Printf("testing exc: %s - colName: %s\n", exc, col.Name)
			if col.Name == exc {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		if len(update) > 0 {
			update += ","
		}
		update += fmt.Sprintf("%s = excluded.%s", col.Name, col.Name)
	}
	pkeyName := t.Name + "_pkey"
	qry := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?) ON CONFLICT ON CONSTRAINT %s DO UPDATE SET %s",
		tableName, colString, pkeyName, update)

	return sqlx.In(qry, cargs)
}

func (t *Table) CreateUpdateQuery(valueMap map[string]interface{}) (query string, args []interface{}, err error) {
	columns, argParam, err := t.createColumnsAndValues(valueMap, true)
	if err != nil {
		return "", nil, err
	}

	tableName := fmt.Sprintf("%s.%s", t.Schema, t.Name)
	primaries := t.GetPrimaryKeyColumnNames()

	var primaryValues []interface{}
	criterias := ""
	for _, col := range primaries {
		val, ok := valueMap[col]
		if !ok {
			return "", nil, fmt.Errorf("update query must include primary key column(s)")
		}
		if len(criterias) > 0 {
			criterias += " AND "
		}
		criterias += col + " = ?"
		primaryValues = append(primaryValues, val)
	}

	updates := ""
	for idx, col := range columns {
		if t.IsPrimary(col.Name) {
			continue
		}

		if len(updates) > 0 {
			updates += ","
		}
		updates += col.Name + " = ?"
		args = append(args, argParam[idx])
	}

	args = append(args, primaryValues...)
	qry := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName, updates, criterias)

	return qry, args, nil
}

func (t *Table) createColumnsAndValues(value interface{}, skipAutoColumn bool) (columns []Column, values []interface{}, err error) {
	dataType := reflect.TypeOf(value)
	dataVal := reflect.ValueOf(value)
	modelType := reflect.TypeOf(t.Model)
	valIsMap := false
	var valKeys []reflect.Value

	if dataVal.Kind() == reflect.Ptr {
		dataVal = dataVal.Elem()
	}

	if dataVal.Kind() == reflect.Map {
		valIsMap = true
		valKeys = dataVal.MapKeys()
		for _, key := range dataVal.MapKeys() {
			if key.Type().Kind() != reflect.String {
				return nil, nil, fmt.Errorf("value as map should have string key")
			}
		}
	}

	valType := dataVal.Type()
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}

	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if dataType.Kind() == reflect.Ptr {
		dataType = dataType.Elem()
	}

	if !valIsMap && modelType.Name() != dataType.Name() {
		return nil, nil, fmt.Errorf("value must be of %s type or a map, got %s type", modelType.Name(), valType.Name())
	}

	for i := 0; i < modelType.NumField(); i++ {
		var val interface{}
		field := modelType.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		colName := field.Tag.Get("db")
		if colName == "" {
			continue
		}

		column, ok := t.GetColumn(colName)
		if !ok {
			continue
		}

		if skipAutoColumn && column.IsAutoIncrement() {
			continue
		}

		if valIsMap {
			key := reflect.ValueOf(column.Name)
			keyFound := false
			for _, mk := range valKeys {
				if mk.Interface() == column.Name {
					keyFound = true
					break
				}
			}

			if !keyFound {
				continue
			}

			//mapValType := mapVal.Elem().Type()
			//valKind := column.Type.ToKind()
			mapVal := dataVal.MapIndex(key)

			//if mapValType.Kind() != valKind && mapValType.Kind() != reflect.Struct {
			//	return nil, nil, fmt.Errorf("cannot assign %s into %s for %s", mapValType.Kind(), valKind.String(), column.Name)
			//}
			val = mapVal.Interface()
		} else {
			val = dataVal.Field(i).Interface()
		}

		if v, ok := val.(driver.Valuer); ok {
			buffVal, err := v.Value()
			if err != nil {
				fmt.Println("Failed to get value from struct, field type :", dataType.Field(i).Type)
				continue
			}

			//skip for nil value
			if buffVal == nil {
				continue
			}

			val = buffVal
		}

		columns = append(columns, *column)
		values = append(values, val)
	}

	return columns, values, nil
}

func (t *Table) CreateNewModel() interface{} {
	modelType := reflect.TypeOf(t.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return reflect.New(modelType).Interface()
}

func (c *Column) IsAutoIncrement() bool {
	if c.Type == FIELD_BIGINT_AUTO || c.Type == FIELD_INT_AUTO {
		return true
	}
	return false
}

//type Table interface{
//	Name() string
//	Schema() string
//	Columns() []Column
//	ColumnNames() []string
//	Indexes() []Index
//	PrimaryKey() []string
//	AddColumns(columns ...Column) error
//	AddIndex(idx Index) error
//	AddPrimaryKey(columnNames ...string) error
//}

//type table struct {
//	name       string
//	schema     string
//	columns    []Column
//	indexes    []Index
//	primaryKey []string
//}

//func NewTable(schema, tablename string, Columns ...Column) Table {
//	return nil
//}
//
//func (t *table) Name() string {
//	return t.name
//}
//
//func (t *table) Schema() string {
//	return t.schema
//}
//
//func (t *table) Columns() []Column {
//	return t.columns
//}
//

//
//func (t *table) Indexes() []Index {
//	return t.indexes
//}
//
//func (t *table) PrimaryKey() []string {
//	return t.primaryKey
//}
//
//func (t *table) AddColumns(columns ...Column) error {
//	for _, c := range columns {
//		if t.columnExists(c.Name) {
//			return fmt.Errorf("column '%s' already exists", c.Name)
//		}
//	}
//	return nil
//}
//
//func (t *table) columnExists(columnName string) bool {
//	for _, c := range t.columns {
//		if c.Name == columnName {
//			return true
//		}
//	}
//	return false
//}

func ParseTableFromJson(data []byte) (*Table, error) {
	var table Table
	if err := json.Unmarshal(data, &table); err != nil {
		return nil, err
	}

	return &table, nil
}

func ParseTableFromJsonFile(fpath string) (*Table, error) {
	f, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var table Table
	if err := json.Unmarshal(f, &table); err != nil {
		return nil, err
	}

	return &table, nil
}

func ParseTableFromYaml(data []byte) (*Table, error) {
	var table Table
	if err := yaml.Unmarshal(data, &table); err != nil {
		return nil, err
	}

	return &table, nil
}

func ParseTableFromYamlFile(fpath string) (*Table, error) {
	f, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var table Table
	if err := yaml.Unmarshal(f, &table); err != nil {
		return nil, err
	}

	return &table, nil
}

func (d DataType) ToString() string {
	return dataTypeStrings[d]
}

func (d DataType) ToKind() reflect.Kind {
	return dataTypeKind[d]
}

func (d *DataType) UnmarshalJSON(data []byte) error {
	dtStr := strings.Replace(strings.ToUpper(string(data)), "\"", "", -1)
	for idx, val := range dataTypeStrings {
		if val == dtStr {
			*d = DataType(idx)
			return nil
		}
	}

	return fmt.Errorf("unknown datatype string: %s", dtStr)
}

func (d *DataType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	dtStr := strings.ToUpper(str)
	for idx, val := range dataTypeStrings {
		if val == dtStr {
			*d = DataType(idx)
			return nil
		}
	}

	return fmt.Errorf("unknown datatype string: %s", dtStr)
}

func (c *Column) UnmarshalJSON(data []byte) error {
	type col Column
	table := col{
		AllowNull: true,
	}

	_ = json.Unmarshal(data, &table)
	*c = Column(table)
	return nil
}

func (c *Column) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type col Column
	table := col{
		AllowNull: true,
	}

	if err := unmarshal(&table); err != nil {
		return err
	}
	*c = Column(table)
	return nil
}
