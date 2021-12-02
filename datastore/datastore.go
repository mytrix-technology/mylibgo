package datastore

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"strings"
	"time"

	"github.com/mytrix-technology/mylibgo/utils/helper"
)

type rowData map[string]interface{}

type DataSet struct {
	Columns []string
	Rows    []rowData
}

type DBHelper interface {
	FieldTypeString(DataType, int) string
	GenerateCreateSchemaQuery(name string) []string
	GenerateCreateTableQuery(table *Table) []string
	GenerateCheckSchemaQuery(name string) string
	GenerateCheckTableQuery(schema, name string) string
	GenerateInsertQuery(table *Table, cols []string, values []string) string
	BuildInsertQuery(schema, table string, columns []string, values [][]interface{}, returnfields ...string) string
}

type DBConfig struct {
	Host     string
	Port     int
	Dbname   string
	User     string
	Password string
}

type DataStore struct {
	db     *sqlx.DB
	config DBConfig
	helper DBHelper
	debugLogger helper.DebugFieldLogger
}

type QueryArg struct {
	Query string
	Args []interface{}
}

type DataStoreOption func(ds *DataStore) error

func DBDebugLogger(fieldLogger helper.DebugFieldLogger) DataStoreOption {
	return func(ds *DataStore) error {
		ds.debugLogger = fieldLogger
		return nil
	}
}

func NewPostgresDatastore(config DBConfig, options ...DataStoreOption) (*DataStore, error) {
	ds := &DataStore{
		config:      config,
		helper:      PostgresHelper{},
		debugLogger: noopFieldLogger,
	}

	for _, op := range options {
		op(ds)
	}

	db, err := sqlx.Open("postgres", createPostgresConnString(config))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 10)

	ds.db = db

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return ds, nil
}

func createPostgresConnString(config DBConfig) string {
	return fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
}

func NewMysqlDatastore(config DBConfig) (*DataStore, error) {
	ds := &DataStore{
		config: config,
		helper: MySqlHelper{},
	}

	db, err := sqlx.Open("mysql", createMysqlConnString(config))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 10)

	ds.db = db

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return ds, nil
}

func createMysqlConnString(config DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local&interpolateParams=true",
		config.User, config.Password, config.Host, config.Port, config.Dbname)
}

func (ds *DataStore) Init(table *Table) error {
	ds.debugLogger("event", "db operation", "msg", fmt.Sprintf("checking existing %s schema", table.Schema))
	if !ds.checkSchema(table.Schema) {
		ds.debugLogger("event", "db operation", "msg", fmt.Sprintf("cannot find existing schema, will create new %s schema", table.Schema))
		qry := ds.helper.GenerateCreateSchemaQuery(table.Schema)
		if _, err := ds.RunSqlNonQuery(qry...); err != nil {
			return err
		}
	}

	if !ds.checkTable(table.Schema, table.Name) {
		qry := ds.helper.GenerateCreateTableQuery(table)
		if _, err := ds.RunSqlNonQuery(qry...); err != nil {
			fmt.Printf("[DEBUG] failed sql: %s\n", qry)
			return err
		}
	}

	return nil
}

func (ds *DataStore) GetDBHandler() *sqlx.DB {
	return ds.db
}

func (ds *DataStore) CreatePGListener(channel string, eventCallback pq.EventCallbackType) (*pq.Listener, error) {
	connstr := createPostgresConnString(ds.config)
	minReconn := 10 * time.Second
    maxReconn := time.Minute

	l := pq.NewListener(connstr, minReconn, maxReconn, eventCallback)
	if err := l.Listen(channel); err != nil {
		return nil, err
	}

	return l, nil
}

func (ds *DataStore) SendNotify(channel string, payload string) error {
	qry := fmt.Sprintf("NOTIFY %s, %s", pq.QuoteIdentifier(channel), pq.QuoteLiteral(payload))
	if _, err := ds.db.Exec(qry); err != nil {
		return err
	}

	return nil
}

//Get is wrapping the sqlx Get function
func (ds *DataStore) Get(dest interface{}, query string, args ...interface{}) error {
	query = ds.db.Rebind(query)
	return ds.db.Get(dest, query, args...)
}

//Select is wrapping the sqlx Select function
func (ds *DataStore) Select(dest interface{}, query string, args ...interface{}) error {
	return ds.db.Select(dest, query, args...)
}

func (ds *DataStore) checkSchema(schemaName string) bool {
	qry := ds.helper.GenerateCheckSchemaQuery(schemaName)
	ds.debugLogger("event", "db operation", "msg", fmt.Sprintf("sending query: (%s)", qry))

	set, err := ds.RunSqlQuery(qry)
	if err != nil {
		ds.debugLogger("event", "db error", "msg", err)
		return false
	}

	if len(set.Rows) > 0 {
		return true
	}
	return false
}

func (ds *DataStore) checkTable(schemaName string, tableName string) bool {
	qry := ds.helper.GenerateCheckTableQuery(schemaName, tableName)
	set, err := ds.RunSqlQuery(qry)
	if err != nil {
		return false
	}
	if len(set.Rows) > 0 {
		return true
	}
	return false
}

func (ds *DataStore) RunSqlQuery(qry string) (*DataSet, error) {
	//if err := ds.db.Ping(); err != nil {
	//	return nil, err
	//}

	rows, err := ds.db.Query(qry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := &DataSet{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	data.Columns = columns
	data.Rows = []rowData{}

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		_ = rows.Scan(valuePtrs...)
		dt := rowData{}

		for i, col := range columns {
			var v interface{}
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			dt[strings.ToUpper(col)] = v
		}
		data.Rows = append(data.Rows, dt)
	}

	return data, nil
}

func (ds *DataStore) RunSqlQueryArgs(qry string, args ...interface{}) (*DataSet, error) {
	//if err := ds.db.Ping(); err != nil {
	//	return nil, err
	//}

	rows, err := ds.db.Queryx(qry, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := &DataSet{}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	data.Columns = columns
	data.Rows = []rowData{}

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		_ = rows.Scan(valuePtrs...)

		dt := rowData{}

		for i, col := range columns {
			var v interface{}
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			dt[strings.ToUpper(col)] = v
		}
		data.Rows = append(data.Rows, dt)
	}

	return data, nil
}

// RunSqlNonQuery will return lastResultID for auto increment column in MySQL, otherwise it will return 0.
func (ds *DataStore) RunSqlNonQuery(queries ...string) (lastResultID int64, err error) {
	//if err = ds.db.Ping(); err != nil {
	//	return 0, err
	//}
	var res sql.Result
	var tx *sql.Tx
	tx, err = ds.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		// Rollback the transaction after the function returns.
		// If the transaction was already committed, this will do nothing.
		_ = tx.Rollback()
	}()

	for _, q := range queries {
		if q != "" {
			res, err = tx.Exec(q)
			if err != nil {
				return 0, err
			}
		}
	}

	lastResultID, err = res.LastInsertId()
	if err != nil {
		lastResultID = 0
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return
}

// RunSqlNonQueryArgs will return lastResultID for auto increment column in MySQL, otherwise it will return 0.
func (ds *DataStore) RunSqlNonQueryArgs(queries ...QueryArg) (lastResultID int64, err error) {
	//if err = ds.db.Ping(); err != nil {
	//	return 0, err
	//}
	var res sql.Result
	var tx *sql.Tx
	tx, err = ds.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		// Rollback the transaction after the function returns.
		// If the transaction was already committed, this will do nothing.
		_ = tx.Rollback()
	}()

	for _, q := range queries {
		res, err = tx.Exec(q.Query, q.Args...)
		if err != nil {
			return 0, err
		}
	}

	lastResultID, err = res.LastInsertId()
	if err != nil {
		lastResultID = 0
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return
}

func noopFieldLogger(keyvals ...interface{}) error {
	return nil
}
