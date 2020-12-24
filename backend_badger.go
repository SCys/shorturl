package main

import (
	"time"

	"github.com/dgraph-io/badger/v2"
	"iscys.com/shorturl/core"
)

// BackendBadger badger backend
type BackendBadger struct {
	db *badger.DB
}

func (b *BackendBadger) Init() {
	var err error

	b.db, err = badger.Open(badger.DefaultOptions(*dsnRaw))
	if err != nil {
		core.F("open database failed:%s", err.Error())
		return
	}

	core.I("badger is connected")
}

func (b *BackendBadger) Get(key string) (string, error) {
	var val string

	err := b.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(core.Bytes(key))
		if err == badger.ErrKeyNotFound {
			return core.ErrObjectNotFound
		} else if err != nil {
			core.E("badger get failed", err)
			return err
		}

		return item.Value(func(raw []byte) error {
			val = core.String(raw)
			return nil
		})
	})

	return val, err
}

func (b *BackendBadger) Set(key string, value string, expire time.Duration) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(core.Bytes(key), core.Bytes(value))
	})
}

func (b *BackendBadger) Del(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(core.Bytes(key))
	})
}
