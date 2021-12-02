package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/buraksezer/olric"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/mytrix-technology/mylibgo/utils/helper"

	"github.com/rs/xid"

	"github.com/mytrix-technology/mylibgo/datastore"
	"github.com/mytrix-technology/mylibgo/data-structure/pub-sub"
)

type RepoStore interface {
	Select(dest interface{}, options ...QueryOption) error
	Get(dest interface{}, id interface{}) error
	GetFromCache(k interface{}) ([]byte, error)
	GetFromDB(dest interface{}, id interface{}) error
	Put(id interface{}, value interface{}, tx ...*sqlx.Tx) (interface{}, error)
	Update(id interface{}, value map[string]interface{}, tx ...*sqlx.Tx) error
	SetIntoCacheAndPublish(id interface{}, value interface{}) error
	SetIntoCache(id interface{}, value interface{}) error
	DeleteFromCache(id interface{}) error
	PublishUpdate(id interface{}, value interface{}) error
	PublishDelete(id interface{}, value interface{}) error
	DumpCacheWithHexKey() map[string][]byte
	GetColumnName(modelFieldName string) string
	DB() *sqlx.DB
}

const (
	groupCacheTypeString = "GroupCache"
	freeCacheTypeString  = "FreeCache"
	olricCacheTypeString = "OlricCache"
)

var (
	ErrCache = fmt.Errorf("cache error")
)

type RepoConfig struct {
	Name         string
	TableModel   *datastore.Table
	DebugLogger  helper.DebugFieldLogger
	KeyColumn    *datastore.Column
	InitOnCreate bool
	CacheOption  *CacheOption
	PubsubClient pubsub.Client
	OlricDB      *olric.Olric
}

type Repository struct {
	Name           string
	DataStore      datastore.Store
	TableModel     *datastore.Table
	keyColumn      *datastore.Column
	cache          RepoCache
	cacheOption    *CacheOption
	pubsubClient   pubsub.Client
	debugLogger    helper.DebugFieldLogger
	modelColumnMap map[string]string
	subscriptions  map[string]pubsub.Subscription
	olricDB        *olric.Olric
	initOnCreate   bool
}

type itemPayload struct {
	ID   interface{}
	Data interface{}
}

// type DBListener struct {
// 	listener       *pq.Listener
// 	Channel        string
// 	NotifyCallback NotifyEventCallback
// }

type StoreOptions func(*RepoConfig) error
type NotifyEventCallback func(notify *pq.Notification)
type ListenEventCallback func()

type NotificationPayload struct {
	Sender string
	Data   interface{}
}

type CacheOption struct {
	Instance   string
	Size       int64
	ItemMaxAge time.Duration
	Type       string
}

func WithGroupCache(size int, age time.Duration) StoreOptions {
	return func(r *RepoConfig) error {
		if size < 1 {
			size = 100_000
		}
		instance := r.Name + "-" + xid.New().String()

		r.CacheOption = &CacheOption{
			Size:       100_000,
			ItemMaxAge: age,
			Type:       groupCacheTypeString,
			Instance:   instance,
		}

		return nil
	}
}

func InitDBStore(init bool) StoreOptions {
	return func(r *RepoConfig) error {
		r.InitOnCreate = init
		return nil
	}
}

func Cache(size int) StoreOptions {
	return func(r *RepoConfig) error {
		if size < 1 {
			size = 500_000
		}
		instance := xid.New().String()

		// cache := freecache.NewCache(size)

		// r.cache = cache

		if r.CacheOption == nil {
			r.CacheOption = &CacheOption{
				Size:       int64(size),
				ItemMaxAge: 300 * time.Second,
				Type:       freeCacheTypeString,
				Instance:   instance,
			}
		} else {
			r.CacheOption.Size = int64(size)
			r.CacheOption.Type = freeCacheTypeString
			r.CacheOption.Instance = instance
		}
		return nil
	}
}

func CacheItemMaxAge(dur time.Duration) StoreOptions {
	return func(r *RepoConfig) error {
		if r.CacheOption == nil {
			r.CacheOption = &CacheOption{
				Size:       1000,
				ItemMaxAge: dur,
			}
		} else {
			r.CacheOption.ItemMaxAge = dur
		}

		return nil
	}
}

func CacheSyncClient(pubsubClient pubsub.Client) StoreOptions {
	return func(r *RepoConfig) error {
		r.PubsubClient = pubsubClient
		return nil
	}
}

func DebugLogger(debugLogger helper.DebugFieldLogger) StoreOptions {
	return func(r *RepoConfig) error {
		r.DebugLogger = debugLogger
		return nil
	}
}

func KeyFieldName(name string) StoreOptions {
	return func(r *RepoConfig) error {
		if len(r.TableModel.PrimaryKey) > 1 {
			return fmt.Errorf("cannot set keyFieldName for table def with compound primary key")
		}

		col, found := r.TableModel.GetColumn(name)
		if !found {
			return fmt.Errorf("cannot find field '%s' in TableModel for KeyFieldName option", name)
		}

		r.KeyColumn = col
		return nil
	}
}

func (r *Repository) init() error {
	return r.DataStore.Init(r.TableModel)
}

func NewRepoStore(ds datastore.Store, tableModel *datastore.Table, options ...StoreOptions) (RepoStore, error) {
	return NewRepository(ds, tableModel, options...)
}

func NewRepository(ds datastore.Store, tableModel *datastore.Table, options ...StoreOptions) (*Repository, error) {
	conf := &RepoConfig{
		Name:         tableModel.Schema + "." + tableModel.Name,
		TableModel:   tableModel,
		DebugLogger:  helper.CreateNoopFieldLogger(),
		KeyColumn:    nil,
		InitOnCreate: false,
		CacheOption:  nil,
		PubsubClient: nil,
	}

	if err := tableModel.CheckValidity(); err != nil {
		return nil, err
	}

	for _, opt := range options {
		if opt != nil {
			if err := opt(conf); err != nil {
				return nil, err
			}
		}
	}

	r := &Repository{
		conf.Name,
		ds,
		tableModel,
		conf.KeyColumn,
		nil,
		conf.CacheOption,
		conf.PubsubClient,
		conf.DebugLogger,
		nil,
		nil,
		conf.OlricDB,
		conf.InitOnCreate,
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
		opt := r.cacheOption
		if r.pubsubClient == nil && opt.Type == freeCacheTypeString {
			return nil, fmt.Errorf("cannot use cache without cache sync options")
		}

		if r.keyColumn == nil {
			return nil, fmt.Errorf("using cache for table def with compound primary key %v is currently unsupported", r.TableModel.PrimaryKey)
		}

		switch opt.Type {
		case freeCacheTypeString:
			cache, err := newFreeCache(int(r.cacheOption.Size), r.debugLogger)
			if err != nil {
				return nil, err
			}
			r.cache = cache
		case olricCacheTypeString:
			cache, err := newOlricCache(r.Name, r.olricDB, r.debugLogger)
			if err != nil {
				return nil, err
			}
			r.cache = cache
		}

		if r.TableModel.Model == nil {
			return nil, fmt.Errorf("cannot have nil model in TableModel for cached repository")
		}

		gob.Register(r.TableModel.Model)
	}

	modelMap, err := createStructFieldTagMap(tableModel.Model, "db")
	if err != nil {
		return nil, err
	}

	r.modelColumnMap = modelMap

	if r.initOnCreate {
		if err := r.init(); err != nil {
			return nil, err
		}
	}

	if err := r.startListening(); err != nil {
		return nil, fmt.Errorf("failed to start cache sync listener: %s", err)
	}
	return r, nil
}

func (r *Repository) Select(dest interface{}, options ...QueryOption) error {
	qc := queryConfig{}
	for _, op := range options {
		op(&qc)
	}

	destType := reflect.TypeOf(dest)
	modelTyp := reflect.TypeOf(r.TableModel.Model)

	if modelTyp.Kind() == reflect.Ptr {
		modelTyp = modelTyp.Elem()
	}

	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer to a slice of %s or *%s", modelTyp.Name(), modelTyp.Name())
	}

	if destType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice of %s or *%s", modelTyp.Name(), modelTyp.Name())
	}

	columns := strings.Join(r.TableModel.GetColumnNames(), ",")
	var argParam []interface{}

	where := ""
	for _, f := range qc.filters {
		if _, ok := r.TableModel.GetColumn(f.Field); !ok {
			return fmt.Errorf("field %s in filter does not exists", f.Field)
		}
		crit, arg, err := f.Encode()
		if err != nil {
			return err
		}

		if len(where) > 0 {
			where += " AND "
		}
		where += crit
		if arg != nil {
			argParam = append(argParam, arg)
		}
	}

	if where != "" {
		where = " WHERE " + where
	}

	sort := ""
	if len(qc.sorts) > 0 {
		sort = qc.sorts.Encode()
		if sort != "" {
			sort = " ORDER BY " + qc.sorts.Encode()
		}
	}

	limit := ""
	if qc.limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", qc.limit)
	}

	offset := ""
	if qc.offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", qc.offset)
	}

	table := r.TableModel.Schema + "." + r.TableModel.Name
	qry := fmt.Sprintf("SELECT %s FROM %s%s%s%s%s", columns, table, where, sort, limit, offset)
	q, args, err := sqlx.In(qry, argParam...)
	if err != nil {
		return err
	}

	db := r.DataStore.GetDBHandler()
	q = db.Rebind(q)
	r.debugLogger("event", "db operation", "msg", "sending query", "query", q, "args", args)

	if err := r.DataStore.Select(dest, q, args...); err != nil {
		return err
	}

	idFieldName := ""
	for i := 0; i < modelTyp.NumField(); i++ {
		tag := modelTyp.Field(i).Tag.Get("db")
		if tag == r.keyColumn.Name {
			idFieldName = modelTyp.Field(i).Name
		}
	}

	if r.cache != nil {
		destVal := reflect.ValueOf(dest).Elem()
		if destVal.Len() > 0 {
			for i := 0; i < destVal.Len(); i++ {
				val := destVal.Index(i)
				_ = r.SetIntoCache(val.FieldByName(idFieldName).Interface(), val.Interface())
			}
		}
	}

	return nil
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

//Get item from cache
func (r *Repository) GetFromCache(k interface{}) ([]byte, error) {
	var key []byte
	keybuf := new(bytes.Buffer)

	if r.cache != nil {
		if err := gob.NewEncoder(keybuf).Encode(k); err != nil {
			return nil, err
		}

		key = keybuf.Bytes()
		cache, err := r.cache.Get(key)
		if err != nil {
			return nil, err
		}

		return cache, nil
	}

	return nil, fmt.Errorf("cache option must not be empty")
}

// Get item directly from db
func (r *Repository) GetFromDB(dest interface{}, id interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return fmt.Errorf("dest need to be a pointer to struct, gets %s", destType.Kind().String())
	}

	columns := strings.Join(r.TableModel.GetColumnNames(), ",")
	var argParam []interface{}
	qry := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s = ?", columns, r.TableModel.Schema, r.TableModel.Name, r.keyColumn.Name)
	argParam = append(argParam, id)

	if err := r.DataStore.Get(dest, qry, argParam...); err != nil {
		return err
	}

	if r.cache != nil {
		_ = r.SetIntoCache(id, dest)
	}

	return nil
}

func (r *Repository) Put(id interface{}, value interface{}, tx ...*sqlx.Tx) (interface{}, error) {
	q, args, err := r.TableModel.CreateInsertQuery(value)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to build the query. %s", err)
	}

	primaryKeys := r.TableModel.GetPrimaryKeyColumnNames()
	keyColumns := r.keyColumn.Name
	var ids = make([]interface{}, len(primaryKeys))
	if id == nil {
		keyColumns = strings.Join(primaryKeys, ",")
	} else {
		ids = make([]interface{}, 1)
	}

	var ptrToIds = make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		ptrToIds[i] = &ids[i]
	}

	q += fmt.Sprintf(" RETURNING %s", keyColumns)

	db := r.DataStore.GetDBHandler()
	var xtx *sqlx.Tx
	if len(tx) > 0 {
		xtx = tx[0]
	} else {
		xtx, err = db.Beginx()
		if err != nil {
			return nil, err
		}
		defer xtx.Rollback()
	}

	if xtx == nil {
		return nil, fmt.Errorf("failed to start db transaction")
	}

	q = xtx.Rebind(q)

	err = xtx.QueryRowx(q, args...).Scan(ptrToIds...)
	if err != nil {
		return nil, err
	}

	if len(tx) == 0 {
		if err := xtx.Commit(); err != nil {
			return nil, err
		}
	}

	if id == nil {
		id = ids[0]
		if len(ids) > 1 {
			strId := ""
			for _, key := range ids {
				if len(strId) > 0 {
					strId += "-"
				}
				strId += fmt.Sprintf("%v", key)
			}
			id = strId
		}
	}
	if r.cache != nil {
		_ = r.SetIntoCacheAndPublish(id, value)
	}
	return id, nil
}

func (r *Repository) Update(id interface{}, value map[string]interface{}, tx ...*sqlx.Tx) error {
	q, args, err := r.TableModel.CreateUpdateQuery(value)
	if err != nil {
		return fmt.Errorf("repository: failed to build the query. %s", err)
	}

	db := r.DataStore.GetDBHandler()
	xtx := getTransaction(db, tx...)
	if xtx == nil {
		return fmt.Errorf("failed to start db transaction")
	}

	if len(tx) == 0 {
		defer xtx.Rollback()
	}

	q = xtx.Rebind(q)

	if _, err = xtx.Exec(q, args...); err != nil {
		return fmt.Errorf("failed to update %s. %s", r.Name, err.Error())
	}

	if len(tx) == 0 {
		if err := xtx.Commit(); err != nil {
			return err
		}
	}

	if r.cache != nil {
		model := r.TableModel.CreateNewModel()
		if err := r.GetFromDB(&model, id); err != nil {
			return fmt.Errorf("%w. failed to get data to update cache", ErrCache)
		}
		_ = r.SetIntoCacheAndPublish(id, model)
	}
	return nil
}

func (r *Repository) SetIntoCacheAndPublish(id interface{}, value interface{}) error {
	if r.cache != nil {
		if err := r.SetIntoCache(id, value); err != nil {
			return err
		}

		if r.cacheOption.Type == freeCacheTypeString {
			return r.PublishUpdate(id, value)
		}
	}

	return nil
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
		_ = r.debugLogger("event", "cache set", "key", id, "value", value)
		if err := r.cache.Set(key, valgob, expiry); err != nil {
			return err
		}
		_ = r.debugLogger("event", "cache set successful", "key", id, "value-bytes", fmt.Sprintf("%d bytes", len(valgob)))
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

func (r *Repository) PublishUpdate(id interface{}, value interface{}) error {
	subject := fmt.Sprintf("%s.%s.update", r.TableModel.Schema, r.TableModel.Name)
	return r.publish(subject, id, value)
}

func (r *Repository) PublishDelete(id interface{}, value interface{}) error {
	subject := fmt.Sprintf("%s.%s.delete", r.TableModel.Schema, r.TableModel.Name)
	return r.publish(subject, id, value)
}

func (r *Repository) GetDataStore() datastore.Store {
	return r.DataStore
}

func (r *Repository) GetTableModel() datastore.Table {
	return *r.TableModel
}

func (r *Repository) DB() *sqlx.DB {
	return r.DataStore.GetDBHandler()
}

func (r *Repository) publish(subject string, id, value interface{}) error {
	//if r.model == nil {
	//	gob.Register(value)
	//	r.model = value
	//}

	payload := itemPayload{
		ID:   id,
		Data: value,
	}
	payloadBuf := new(bytes.Buffer)
	if err := gob.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return fmt.Errorf("failed to serialize message payload: %s", err)
	}
	return r.pubsubClient.Publish(subject, payloadBuf.Bytes())
}

func (r *Repository) startListening() error {
	if r.cacheOption == nil {
		return nil
	}
	r.subscriptions = make(map[string]pubsub.Subscription)

	updateSubject := fmt.Sprintf("%s.%s.update", r.TableModel.Schema, r.TableModel.Name)
	deleteSubject := fmt.Sprintf("%s.%s.delete", r.TableModel.Schema, r.TableModel.Name)

	r.subscriptions[updateSubject] = r.pubsubClient.Listen(updateSubject).Subscribe(r.listenForUpdates)
	r.subscriptions[deleteSubject] = r.pubsubClient.Listen(deleteSubject).Subscribe(r.listenForDeletes)

	return nil
}

func (r *Repository) listenForUpdates(ev pubsub.Event) {
	type Notification struct {
		Source  string
		Payload itemPayload
	}

	switch ev.Type {
	case pubsub.EVENT_ERROR:
		_ = r.debug("event", "listener error", "msg", ev.Message.Payload)
		_ = r.debug("event", "cache sync error", "msg", fmt.Sprintf("clearing cache for %s.%s due to sync error", r.TableModel.Schema, r.TableModel.Name))
		r.cache.Clear()
	case pubsub.EVENT_LISTENER_CONNECTED:
		_ = r.debug("event", "listener connected", "msg", fmt.Sprintf("%s.%s update listeners connected", r.TableModel.Schema, r.TableModel.Name))
		return
	case pubsub.EVENT_NOTIFY:
		if ev.Message.Source == r.pubsubClient.Instance() {
			return
		}

		var payload itemPayload
		if err := gob.NewDecoder(bytes.NewReader(ev.Message.Payload)).Decode(&payload); err != nil {
			_ = r.debug("event", "listener error", "msg", fmt.Sprintf("failed to deserialize message payload: %s", err))
		}

		id := payload.ID
		_ = r.debug("event", "listener event", "msg", fmt.Sprintf("receives notification: %v", payload))
		if err := r.SetIntoCache(id, payload.Data); err != nil {
			_ = r.debug("event", "update cache", "msg", fmt.Sprintf("failed to set into cache for %s.%s ID %q: %s", r.TableModel.Schema, r.TableModel.Name, payload.ID, err))
			return
		}
		_ = r.debug("event", "update cache", "msg", fmt.Sprintf("cache for %s.%s updated", r.TableModel.Schema, r.TableModel.Name))
	}
}

func (r *Repository) listenForDeletes(ev pubsub.Event) {
	type Notification struct {
		Source  string
		Payload itemPayload
	}

	switch ev.Type {
	case pubsub.EVENT_ERROR:
		_ = r.debug("event", "listener error", "msg", ev.Message.Payload)
		_ = r.debug("event", "cache sync error", "msg", fmt.Errorf("clearing cache for %s.%s due to sync error", r.TableModel.Schema, r.TableModel.Name))
		r.cache.Clear()
	case pubsub.EVENT_LISTENER_CONNECTED:
		_ = r.debug("event", "listener connected", "msg", fmt.Sprintf("%s.%s delete listeners connected", r.TableModel.Schema, r.TableModel.Name))
	case pubsub.EVENT_NOTIFY:
		if ev.Message.Source == r.pubsubClient.Instance() {
			return
		}

		var payload itemPayload
		if err := gob.NewDecoder(bytes.NewReader(ev.Message.Payload)).Decode(&payload); err != nil {
			_ = r.debug("event", "listener error", "msg", fmt.Sprintf("failed to deserialize message payload: %s", err))
		}

		id := payload.ID
		_ = r.debug("event", "listener event", "msg", fmt.Sprintf("receives notification: %v", payload))
		_ = r.debug("event", "delete cache item", "msg", fmt.Sprintf("delete cache key (%v)", id))
		_ = r.DeleteFromCache(id)
		_ = r.debug("event", "update cache", "msg", fmt.Sprintf("cache for %s.%s updated", r.TableModel.Schema, r.TableModel.Name))
	}
}

func (r *Repository) DumpCacheWithHexKey() map[string][]byte {
	if r.cache == nil {
		return nil
	}

	return r.cache.DumpCacheWithHexKey()
}

func (r *Repository) debug(keyvals ...interface{}) error {
	return r.debugLogger(keyvals...)
}

func (r *Repository) GetColumnName(modelFieldName string) string {
	name := ""
	if colname, ok := r.modelColumnMap[modelFieldName]; ok {
		name = colname
	}

	return name
}

func noopFieldLogger(keyvals ...interface{}) error {
	return nil
}

func getTransaction(db *sqlx.DB, tx ...*sqlx.Tx) *sqlx.Tx {
	var xtx *sqlx.Tx
	var err error
	if len(tx) > 0 {
		xtx = tx[0]
	} else {
		xtx, err = db.Beginx()
		if err != nil {
			return nil
		}
	}

	return xtx
}

func createStructFieldTagMap(model interface{}, tag string) (map[string]string, error) {
	structMap := make(map[string]string)
	modelTyp := reflect.TypeOf(model)
	if modelTyp.Kind() == reflect.Ptr {
		modelTyp = modelTyp.Elem()
	}

	if modelTyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct")
	}

	for i := 0; i < modelTyp.NumField(); i++ {
		name := modelTyp.Field(i).Name
		structMap[name] = modelTyp.Field(i).Tag.Get(tag)
	}

	return structMap, nil
}