package cache

import (
    "sync"
)

type Cache struct {
    mu    sync.RWMutex
    store map[string]interface{}
}

func NewCache() *Cache {
    return &Cache{
        store: make(map[string]interface{}),
    }
}

func (c *Cache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.store[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    value, ok := c.store[key]
    return value, ok
}