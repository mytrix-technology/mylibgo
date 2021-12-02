package repository

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/mailgun/groupcache/v2" //nolint:typecheck
	"github.com/rs/xid"

	"github.com/mytrix-technology/mylibgo/networking/peer"
)

type GCacheOption func(*GroupCacheConfig)
type GroupCacheConfig struct {
	Size int64
	MaxAge time.Duration
	GetterFunc groupcache.GetterFunc
	Scheme string
	Port int
	Name string
	Instance string
}

type GroupCache struct {
	cache *groupcache.Group
	pool *groupcache.HTTPPool
}

func WithCacheSize(size int64) GCacheOption {
	return func(gc *GroupCacheConfig) {
		gc.Size = size
	}
}

func WithMaxAge(age time.Duration) GCacheOption {
	return func (gc *GroupCacheConfig) {
		gc.MaxAge = age
	}
}

func newGroupCache(name string, getterFunc groupcache.GetterFunc, options ...GCacheOption) (*GroupCache, error) {
	cfg := GroupCacheConfig{
		Scheme:   "http",
		Name:     name,
		Size: 100_000,
	}

	for _, opt := range options {
		opt(&cfg)
	}

	return NewGroupCacheWithConfig(cfg)
}

func NewGroupCacheWithConfig(conf GroupCacheConfig) (*GroupCache, error) {
	if err := validateGCacheConfig(&conf); err != nil {
		return nil, err
	}

	local := url.URL{
		Scheme:     conf.Scheme,
		Host:       fmt.Sprintf("localhost:%d", conf.Port),
	}

	pool := groupcache.NewHTTPPoolOpts(local.String(), &groupcache.HTTPPoolOptions{})
	group := groupcache.NewGroup(conf.Name, conf.Size, conf.GetterFunc)
	gcache := &GroupCache{
		cache: group,
		pool:  pool,
	}

	p, err := peer.New(conf.Instance, conf.Name, peer.SetNotify(makePoolRegisterer(gcache, conf.Scheme, conf.Port)))
	if err != nil {
		return nil, fmt.Errorf("failed to create peer listener for group cache. %w", err)
	}
	if err := p.Listen(); err != nil {
		return nil, fmt.Errorf("failed to listen on localhost port %d. %w", conf.Port, err)
	}

	return gcache, nil
}

func (gc *GroupCache) Get(key []byte) (value []byte, err error) {
	ctx := context.Background()
	err = gc.cache.Get(ctx, string(key), groupcache.AllocatingByteSliceSink(&value))
	return
}

func (gc *GroupCache) Set(key []byte, value []byte, expireSeconds int) error {
	ctx := context.Background()
	_ = gc.cache.Remove(ctx, string(key))
	_, err := gc.Get(key)
	return err
}

func validateGCacheConfig(cfg *GroupCacheConfig) error {
	instance := cfg.Name + "-" + xid.New().String()
	cfg.Instance = instance

	if cfg.Port < 1 {
		cfg.Port = 9000
	}

	if cfg.Scheme != "http" && cfg.Scheme != "https" {
		cfg.Scheme = "http"
	}

	return nil
}

func makePoolRegisterer(gcache *GroupCache, scheme string, port int) func(peer.Discovered) {
	return func(d peer.Discovered) {
		p := url.URL{
			Scheme:     scheme,
			Host:       fmt.Sprintf("%s:%d", d.Address, port),
		}

		gcache.pool.Set(p.String())
	}
}

func MakeRepoGCGetter(repo *Repository, expireDuration time.Duration) groupcache.GetterFunc {
	destType := reflect.TypeOf(repo.TableModel.Model)
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

	return func(ctx context.Context, key string, sink groupcache.Sink) error {
		dest := reflect.New(destType).Interface()
		if err := repo.Get(&dest, key); err != nil {
			return err
		}
		valbuf := new(bytes.Buffer)
		if err := gob.NewEncoder(valbuf).Encode(dest); err != nil {
			return fmt.Errorf("failed to set into groupcache. %w", err)
		}
		valgob := valbuf.Bytes()
		return sink.SetBytes(valgob, time.Now().Add(expireDuration))
	}
}