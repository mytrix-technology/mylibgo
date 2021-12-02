package repository

type RepoCache interface{
	Get(key []byte) (value []byte, err error)
	Set(key []byte, value []byte, expireSeconds int) error
	Del(key []byte) bool
	Clear()
	DumpCacheWithHexKey() map[string][]byte
}
