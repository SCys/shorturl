package main

import "time"

// Backend interface for datastore
type Backend interface {
	Set(key string, value string, expire time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Init()
}
