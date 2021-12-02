package datastore

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

func connRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	return conn, nil
}

func poolRedis() redis.Conn {
	var pool *redis.Pool
	connPool := pool.Get()
	return connPool
}

func poolIdleRedis() redis.Conn {
	pool := &redis.Pool{
		Dial:            connRedis,
		TestOnBorrow:    nil,
		MaxIdle:         10,
		MaxActive:       100,
		IdleTimeout:     240 * time.Second,
		Wait:            false,
		MaxConnLifetime: 30,
	}
	connPool := pool.Get()
	return connPool
}

func setRedis(Key string, Value string) error {
	conn, _ := connRedis()
	defer conn.Close()

	_, errs := conn.Do("SET", Key, Value)
	if errs != nil {
		log.Fatal(errs)
	}

	return nil
}

func getRedis(Key string, Value string) (interface{},error) {
	conn, _ := connRedis()
	defer conn.Close()

	nameVal, errs := conn.Do("GET", Key, Value)
	if errs != nil {
		return "error", errs
	}

	return nameVal, nil
}
