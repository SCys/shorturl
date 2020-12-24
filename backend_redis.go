package main

import (
	"time"

	"github.com/go-redis/redis"
	"iscys.com/shorturl/core"
)

// BackendRedis redis backend
type BackendRedis struct {
	db *redis.Client
}

func (b *BackendRedis) Init() {
	options, err := redis.ParseURL(*dsnRaw)
	if err != nil {
		core.F("invalid redis server options", err)
	}

	b.db = redis.NewClient(options)

	if sts, err := b.db.Ping().Result(); err != nil {
		core.F("ping redis server failed", err)
	} else {
		core.I("ping redis:%s", sts)
	}

	core.I("redis server is connected")
	if err != nil {
		core.F("open database failed:%s", err.Error())
	}
}

func (b *BackendRedis) Get(key string) (string, error) {
	raw, err := b.db.Get(key).Result()
	if redis.Nil == err {
		return "", core.ErrObjectNotFound
	} else if err != nil {
		core.E("get origin id failed", err)
		return "", err
	}

	return raw, nil
}

func (b *BackendRedis) Set(key string, value string, expire time.Duration) error {
	_, err := b.db.Set(key, value, expire).Result()
	if err != nil {
		core.E("redis set failed", err)
	}

	return err
}

func (b *BackendRedis) Del(key string) error {
	_, err := b.db.Del(key).Result()
	if err != nil {
		core.E("redis del failed", err)
	}

	return err
}
