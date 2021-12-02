package repository

import (
	"encoding/hex"

	"github.com/coocood/freecache"

	"github.com/mytrix-technology/mylibgo/utils/helper"
)

type freecacheWrapper struct {
	cache *freecache.Cache
	debugLogger helper.DebugFieldLogger
}

func newFreeCache(size int, debugLogger helper.DebugFieldLogger) (RepoCache, error) {
	if debugLogger == nil {
		debugLogger = helper.CreateNoopFieldLogger()
	}
	cache := freecache.NewCache(size)

	return &freecacheWrapper{
		cache:       cache,
		debugLogger: debugLogger,
	}, nil
}

func (f *freecacheWrapper) Get(key []byte) (value []byte, err error) {
	return f.cache.Get(key)
}

func (f *freecacheWrapper)Set(key []byte, value []byte, expireSeconds int) error {
	return f.cache.Set(key, value, expireSeconds)
}

func (f *freecacheWrapper)Del(key []byte) bool {
	return f.cache.Del(key)
}

func (f *freecacheWrapper) Clear() {
	f.cache.Clear()
}

func (f *freecacheWrapper) DumpCacheWithHexKey() map[string][]byte {
	iter := f.cache.NewIterator()
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