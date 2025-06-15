package main

import (
	"sync"
)

type KV struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

func (kv *KV) Set(key []byte, value []byte) error {
	kv.mu.Lock()
	defer func() {
		kv.mu.Unlock()
	}()

	kv.data[string(key)] = []byte(value)
	return nil
}

func (kv *KV) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer func() {
		kv.mu.RUnlock()
	}()

	data, ok := kv.data[string(key)]
	return data, ok
}
