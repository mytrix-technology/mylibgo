// Deprecated
package datastore

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/coocood/freecache"
	"github.com/lib/pq"

	"database/sql/driver"

	"github.com/mytrix-technology/mylibgo/utils/helper"
)

//type Repository interface{
//	Get(id interface{}) (interface{}, error)
//	Set(id interface{}, value interface{}) (interface{}, error)
//	Update(id interface{}, value interface{}) (interface{}, error)
//	GetDataStore() DataStore
//	DB() *sqlx.DB
//	setCacheSize(size int64) error
//}

type Repository struct {
	DataStore    *DataStore
	TableModel   *Table
	keyColumn    *Column
	cache        *freecache.Cache
	cacheOption  *CacheOption
	debugLogger  helper.DebugFieldLogger
}

type StoreOptions func(*Repository) error
type NotifyEventCallback func(notify *pq.Notification)
type ListenEventCallback func()

type NotificationPayload struct {
	Sender string
	Data   interface{}
}

type CacheOption struct {
	Size       int
	ItemMaxAge time.Duration
}

func Cache(size int) StoreOptions {
	return func(r *Repository) error {
		if size < 1 {
			size = 500_000
		}

		// cache := freecache.NewCache(size)

		// r.cache = cache

		if r.cacheOption == nil {
			r.cacheOption = &CacheOption{
				Size:       size,
				ItemMaxAge: 300 * time.Second,
			}
		} else {
			r.cacheOption.Size = size
		}
		return nil
	}
}

func CacheItemMaxAge(dur time.Duration) StoreOptions {
	return func(r *Repository) error {
		if r.cacheOption == nil {
			r.cacheOption = &CacheOption{
				Size:       1000,
				ItemMaxAge: dur,
			}
		} else {
			r.cacheOption.ItemMaxAge = dur
		}

		return nil
	}
}

func DebugLogger(debugLogger helper.DebugFieldLogger) StoreOptions {
	return func(r *Repository) error {
		r.debugLogger = debugLogger
		return nil
	}
}

func KeyFieldName(name string) StoreOptions {
	return func(r *Repository) error {
		if len(r.TableModel.PrimaryKey) > 1 {
			return fmt.Errorf("Cannot set keyFieldName for table def with compound primary key")
		}

		col, found := r.TableModel.GetColumn(name)
		if !found {
			return fmt.Errorf("cannot find field '%s' in TableModel for KeyFieldName option", name)
		}

		r.keyColumn = col
		return nil
	}
}

func (r *Repository) init() error {
	return r.DataStore.Init(r.TableModel)
}

func NewRepository(ds *DataStore, tableModel *Table, options ...StoreOptions) (*Repository, error) {
	r := &Repository{
		ds,
		tableModel,
		nil,
		nil,
		nil,
		noopFieldLogger,
	}

	if err := r.TableModel.CheckValidity(); err != nil {
		return nil, err
	}

	for _, opt := range options {
		if opt != nil {
			opt(r)
		}
	}

	if r.keyColumn == nil {
		if r.TableModel.PrimaryKey == nil {
			return nil, fmt.Errorf("cannot use empty primary key for keyFieldName. you need to set the KeyFieldName option")
		}

		if len(r.TableModel.PrimaryKey) == 0 {
			return nil, fmt.Errorf("cannot use empty primary key for keyFieldName. you need to set the KeyFieldName option")
		}

		if len(r.TableModel.PrimaryKey) == 1 {
			// return nil, fmt.Errorf("cannot use compound primary key %v for keyFieldName. you need to set the KeyFieldName option", r.TableModel.PrimaryKey)
		}

		r.keyColumn, _ = r.TableModel.GetColumn(r.TableModel.PrimaryKey[0])
	}

	if r.cacheOption != nil {
		r.cache = freecache.NewCache(r.cacheOption.Size)
	}

	if r.cache != nil && r.keyColumn == nil {
		return nil, fmt.Errorf("using cache for table def with compound primary key %v is currently unsupported", r.TableModel.PrimaryKey)
	}

	if err := r.init(); err != nil {
		return nil, err
	}

	return r, nil
}

// Get item from cache, if not exists then get from db
func (r *Repository) Get(dest interface{}, id interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("dest need to be a pointer to struct, gets %s", destType.Kind().String())
	}

	var key []byte
	keybuf := new(bytes.Buffer)

	if r.cache != nil {
		if err := gob.NewEncoder(keybuf).Encode(id); err != nil {
			return r.GetFromDB(dest, id)
		}

		key = keybuf.Bytes()
		cache, err := r.cache.Get(key)
		if err != nil {
			return r.GetFromDB(dest, id)
		}

		if err := gob.NewDecoder(bytes.NewReader(cache)).Decode(dest); err != nil {
			return r.GetFromDB(dest, id)
		}

		return nil
	}

	return r.GetFromDB(dest, id)
}

// Get item directly from db
func (r *Repository) GetFromDB(dest interface{}, id interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("dest need to be a pointer to struct, gets %s", destType.Kind().String())
	}

	idVal := fmt.Sprintf("%v", id)
	if r.keyColumn.Type == FIELD_VARCHAR || r.keyColumn.Type == FIELD_TEXT || r.keyColumn.Type == FIELD_CHAR {
		idVal = fmt.Sprintf("'%v'", id)
	}

	qry := fmt.Sprintf("SELECT * FROM %s.%s WHERE %s = %s", r.TableModel.Schema, r.TableModel.Name, r.keyColumn.Name, idVal)

	if err := r.DataStore.Get(dest, qry); err != nil {
		return err
	}

	if r.cache != nil {
		_ = r.SetIntoCache(id, dest)
	}

	return nil
}

func (r *Repository) Put(id interface{}, value interface{}) (interface{}, error) {
	dataStruct := reflect.TypeOf(value)
	dataVal := reflect.ValueOf(value)

	if dataStruct.Kind() == reflect.Slice {
		return nil, fmt.Errorf("value is a slice. wants a struct or pointer to struct")
	}

	if dataStruct.Kind() == reflect.Ptr {
		dataStruct = dataStruct.Elem()
	}

	if dataVal.Kind() == reflect.Ptr {
		dataVal = dataVal.Elem()
	}

	colsValue := map[string]string{}
	for i := 0; i < dataStruct.NumField(); i++ {

		if dataVal.Field(i).Kind() == reflect.Ptr {
			continue
		}

		colName := dataStruct.Field(i).Tag.Get("db")
		if colName == "" {
			colName = dataStruct.Field(i).Name
		}

		val := dataVal.Field(i).Interface()

		if dataVal.Field(i).Kind() == reflect.Struct {
			v, ok := dataVal.Field(i).Interface().(driver.Valuer)
			if !ok {
				fmt.Println("Field :", dataStruct.Field(i).Type, " was not driver valuer")
				continue
			}

			buffVal, err := v.Value()
			if err != nil {
				fmt.Println("Failed to get value from struct, field type :", dataStruct.Field(i).Type)
				continue
			}

			//skip for nil value
			if buffVal == nil {
				continue
			}

			val = buffVal
		}

		var colVal string

		if reflect.TypeOf(val).Kind() == reflect.String {
			buffVal := strings.ReplaceAll(fmt.Sprintf("%v", val), "'", "''")
			colVal = fmt.Sprintf("'%s'", buffVal)
		} else if reflect.TypeOf(val).Kind() == reflect.Float64 || reflect.TypeOf(val).Kind() == reflect.Float32 {
			decimal := dataStruct.Field(i).Tag.Get("decimal")
			if decimal != "" {
				colVal = fmt.Sprintf("%." + decimal + "f", val)
			} else {
				colVal = fmt.Sprintf("%.2f", val)
			}

		} else if reflect.TypeOf(val).Kind() == reflect.Int || reflect.TypeOf(val).Kind() == reflect.Int8 ||
			reflect.TypeOf(val).Kind() == reflect.Int16 || reflect.TypeOf(val).Kind() == reflect.Int32 ||
			reflect.TypeOf(val).Kind() == reflect.Int64 {
			colVal = fmt.Sprintf("%d", val)
		} else {
			colVal = fmt.Sprintf("%v", val)
		}

		colsValue[colName] = colVal
	}

	cols := []string{}
	values := []string{}
	for _, col := range r.TableModel.Columns {
		if col.Type == FIELD_INT_AUTO || col.Type == FIELD_BIGINT_AUTO {
			continue
		}

		val, ok := colsValue[col.Name]
		if !ok {
			continue
		}

		values = append(values, val)
		cols = append(cols, col.Name)
	}

	qry := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s) RETURNING %s", r.TableModel.Schema, r.TableModel.Name,
		strings.Join(cols, ","), strings.Join(values, ","), r.keyColumn.Name)

	err := r.DataStore.db.QueryRowx(qry).Scan(&id)
	if err != nil {
		return nil, err
	}

	if r.cache != nil {
		_ = r.SetIntoCache(id, value)
	}
	return id, nil
}

func (r *Repository) SetIntoCache(id interface{}, value interface{}) error {
	if r.cache != nil {
		var key []byte
		var valgob []byte

		keybuf := new(bytes.Buffer)
		if err := gob.NewEncoder(keybuf).Encode(id); err != nil {
			return err
		}
		key = keybuf.Bytes()

		if _, ok := value.([]byte); ok {
			valgob = value.([]byte)
		} else {
			valbuf := new(bytes.Buffer)
			if err := gob.NewEncoder(valbuf).Encode(value); err != nil {
				return err
			}
			valgob = valbuf.Bytes()
		}

		expiry := int(r.cacheOption.ItemMaxAge.Seconds())
		if err := r.cache.Set(key, valgob, expiry); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteFromCache(id interface{}) error {
	if r.cache != nil {
		var key []byte
		keybuf := new(bytes.Buffer)

		if err := gob.NewEncoder(keybuf).Encode(id); err != nil {
			return err
		}

		key = keybuf.Bytes()

		r.cache.Del(key)
	}

	return nil
}

func (r *Repository) DumpCacheWithHexKey() map[string][]byte {
	if r.cache == nil {
		return nil
	}

	iter := r.cache.NewIterator()
	cache := make(map[string][]byte)
	for {
		e := iter.Next()
		if e == nil {
			break
		}
		var val = e.Value
		var key = e.Key

		id := hex.EncodeToString(key)
		cache[id] = val
	}

	return cache
}
