package datastore

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Store interface {
	Init(table *Table) error
	GetDBHandler() *sqlx.DB
	CreatePGListener(channel string, eventCallback pq.EventCallbackType) (*pq.Listener, error)
	SendNotify(channel string, payload string) error
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	RunSqlQuery(qry string) (*DataSet, error)
	RunSqlQueryArgs(qry string, args ...interface{}) (*DataSet, error)
	RunSqlNonQuery(queries ...string) (lastResultID int64, err error)
	RunSqlNonQueryArgs(queries ...QueryArg) (lastResultID int64, err error)
}

func NewStore(db *sqlx.DB) (*DataStore, error) {
	ds := &DataStore{
		db: db,
		helper: PostgresHelper{},
	}

	db.SetConnMaxLifetime(time.Minute * 10)

	ds.db = db

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return ds, nil
}
