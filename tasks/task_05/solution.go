package main

type Cache[K comparable, V any] struct {
	Capacity int
	Storage map[K]V
}

func NewCache[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		Capacity: capacity,
		Storage: make(map[K]V, capacity),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	if c.Capacity == 0 {
		return
	}
	c.Storage[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	v, ok := c.Storage[key]
	return v, ok
}
