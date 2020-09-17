package engine

import (
	"errors"
	"sync"
)

var (
	ErrNotExist = errors.New("key does not exist")
)

type Default struct {
	store map[string][]byte
	sync.RWMutex
}

func New() *Default {
	return &Default{
		store: make(map[string][]byte, 1024),
	}
}

func (d *Default) Set(key string, val []byte) error {
	d.Lock()
	defer d.Unlock()
	d.store[key] = val
	return nil
}

func (d *Default) Get(key string) ([]byte, error) {
	d.RLock()
	defer d.RUnlock()
	val, ok := d.store[key]
	if !ok {
		return nil, ErrNotExist
	}
	return val, nil
}

func (d *Default) Del(key string) error {
	d.Lock()
	defer d.Unlock()
	delete(d.store, key)
	return nil
}
