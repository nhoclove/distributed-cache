package engine

import (
	"errors"
	"sync"
)

var (
	ErrNotExist = errors.New("key does not exist")
)

type String struct {
	store map[string][]byte
	sync.RWMutex
}

func New() *String {
	return &String{
		store: make(map[string][]byte, 1024),
	}
}

func (d *String) Set(key string, val []byte) error {
	d.Lock()
	defer d.Unlock()
	d.store[key] = val
	return nil
}

func (d *String) Get(key string) ([]byte, error) {
	d.RLock()
	defer d.RUnlock()
	val, ok := d.store[key]
	if !ok {
		return nil, ErrNotExist
	}
	return val, nil
}

func (d *String) Del(key string) error {
	d.Lock()
	defer d.Unlock()
	delete(d.store, key)
	return nil
}
