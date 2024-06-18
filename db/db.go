package db

import (
	"github.com/dgraph-io/badger/v4"
)

type KV struct {
	DB *badger.DB
}

func New(persPath string) KV {
	db, err := badger.Open(badger.DefaultOptions(persPath))
	if err != nil {
		panic(err)
	}
	return KV{DB: db}
}

func (kv KV) Set(key, value string) {
	err := kv.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		panic(err)
	}
}

func (kv KV) Get(key string) string {
	var value []byte
	err := kv.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		panic(err)
	}
	return string(value)
}

func (kv KV) Delete(key string) {
	err := kv.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err != nil {
		panic(err)
	}
}

func (kv KV) AllKeys() []string {
	keys := make([]string, 0)
	err := kv.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			keys = append(keys, string(k))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return keys
}

func (kv KV) Close() {
	kv.DB.Close()
}
