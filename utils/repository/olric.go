package repository

import (
	"fmt"
	"time"

	"github.com/mytrix-technology/mylibgo/security/cryptography"
	"github.com/mytrix-technology/mylibgo/utils/helper"
	"github.com/buraksezer/olric"
	"github.com/buraksezer/olric/query"
	"github.com/rs/xid"
)

type olricCache struct {
	cache *olric.DMap
	debugLogger helper.DebugFieldLogger
}

func WithOlricCache(olricDb *olric.Olric, maxAge time.Duration) StoreOptions {
	return func(r *RepoConfig) error {
		instance := r.Name + "-" + xid.New().String()
		if olricDb == nil {
			return fmt.Errorf("olricdb is nil")
		}
		r.CacheOption = &CacheOption{
			Size:       100_000,
			ItemMaxAge: maxAge,
			Type:       olricCacheTypeString,
			Instance:   instance,
		}

		r.OlricDB = olricDb

		return nil
	}
}

func newOlricCache(name string, olricDb *olric.Olric, debugLogger helper.DebugFieldLogger) (RepoCache, error) {
	if debugLogger == nil {
		debugLogger = helper.CreateNoopFieldLogger()
	}

	cache, err := olricDb.NewDMap("name")
	if err != nil {
		return nil, fmt.Errorf("failed to create olric cache DMap. %s", err)
	}

	return &olricCache{
		cache:       cache,
		debugLogger: debugLogger,
	}, nil
}

func (o *olricCache) Get(key []byte) (value []byte, err error) {
	hexKey := cryptography.EncodeHex(key)
	_ = o.debugLogger("event", "get from cache", "hex-key", hexKey)
	rawVal, err := o.cache.Get(hexKey)
	if err != nil {
		return nil, err
	}

	if val, ok := rawVal.([]byte); ok {
		return val, nil
	}

	_ = o.debugLogger("event", "cache error", "error", "value in the cache is not []byte")
	return nil, fmt.Errorf("invalid value in the cache")
}

func (o *olricCache) Set(key []byte, value []byte, expireSeconds int) error {
	hexKey := cryptography.EncodeHex(key)
	expiry := time.Duration(expireSeconds) * time.Second
	withEx := "no expiry"
	if expireSeconds > 0 {
		withEx = fmt.Sprintf("%s", expiry.String())
	}
	_ = o.debugLogger("event", "set value into cache", "hex-key", hexKey, "value", fmt.Sprintf("%d bytes", len(value)),
		"expiry", withEx)
	if expireSeconds > 0 {
		return o.cache.PutEx(hexKey, value, expiry)
	}
	return o.cache.Put(hexKey, value)
}

func (o *olricCache) Del(key []byte) bool {
	hexKey := cryptography.EncodeHex(key)
	_ = o.debugLogger("event", "remove value from cache", "hex-key", hexKey)
	err := o.cache.Delete(hexKey)
	if err != nil {
		_ = o.debugLogger("failed to delete from cache with hex-key: %s", hexKey)
	}
	return err == nil
}

func (o *olricCache) Clear() {
	_ = o.debugLogger("event", "destroy olric cache")
	err := o.cache.Destroy()
	if err != nil {
		_ = o.debugLogger("failed to destroy olric cache. %s", err)
	}
}

func (o *olricCache) DumpCacheWithHexKey() map[string][]byte {
	q := query.M{
		"$onKey": query.M{
			"$regexMatch": "",
	   	},
	}

	row, err := o.cache.Query(q)
	if err != nil {
		_ = o.debugLogger("event", "cache iteration error", "error", err)
		return nil
	}
	defer row.Close()

	resultMap := make(map[string][]byte)
	err = row.Range(func(key string, value interface{}) bool {
		val, ok := value.([]byte)
		if !ok {
			_ = o.debugLogger("event", "cache iteration error", "error", "value is not []byte")
			return false
		}

		resultMap[key] = val
		return true
	})

	if err != nil {
		return nil
	}

	return resultMap
}
