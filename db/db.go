package db

import "sync"

type DB struct {
	mut  sync.RWMutex
	hMap map[string][]byte
}

func (d *DB) Has(key string) bool {
	d.mut.RLock()
	_, ok := d.hMap[key]
	d.mut.RUnlock()
	return ok
}

func (d *DB) Set(key string, value []byte) {
	d.mut.Lock()
	d.hMap[key] = value
	d.mut.Unlock()
}

func (d *DB) Get(key string) ([]byte, bool) {
	d.mut.RLock()
	v, ok := d.hMap[key]
	d.mut.RUnlock()

	return v, ok
}
