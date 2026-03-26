package main

import "sync"

type LRU[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Set(key K, value V)
}

type Node[K comparable, V any] struct {
	Key   K
	Value V
	Prev  *Node[K, V]
	Next  *Node[K, V]
}

type LRUCache[K comparable, V any] struct {
	Storage  map[K]*Node[K, V]
	Capacity int
	Head     *Node[K, V]
	Tail     *Node[K, V]
	Mu       sync.Mutex
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		Capacity: capacity,
		Storage:  make(map[K]*Node[K, V], capacity),
		Mu: sync.Mutex{},
	}
}

func (c *LRUCache[K, V]) Set(key K, value V) {
	if c.Capacity == 0 {
		return
	}
	c.Mu.Lock()
	defer c.Mu.Unlock()
	_, ok := c.Storage[key]
	if !ok {
		c.Storage[key] = &Node[K, V]{
			Key:   key,
			Value: value,
			Prev:  c.Tail,
			Next:  nil,
		}
	} else {
		c.Storage[key].Value = value
	}
	if c.Head == c.Tail && ok {
		return
	}
	if c.Tail == nil {
		c.Tail = c.Storage[key]
		c.Head = c.Storage[key]
	} else {
		if c.Storage[key].Prev == nil {
			c.Head = c.Head.Next
			c.Head.Prev = nil
		}
		c.Storage[key].Prev = c.Tail
		c.Storage[key].Next = nil
		c.Tail.Next = c.Storage[key]
		c.Tail = c.Storage[key]
	}
	if len(c.Storage) > c.Capacity {
		toDelete := c.Head
		c.Head = toDelete.Next
		c.Head.Prev = nil
		delete(c.Storage, toDelete.Key)
	}
}

func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if c.Capacity == 0 {
		return value, false
	}
	_, ok = c.Storage[key]
	if !ok {
		return
	}
	value = c.Storage[key].Value
	if c.Head == c.Tail {
		return
	}
	if c.Storage[key].Prev == nil {
		c.Head = c.Head.Next
		c.Head.Prev = nil
	}
	c.Storage[key].Prev = c.Tail
	c.Tail.Next = c.Storage[key]
	c.Storage[key].Next = nil
	c.Tail = c.Storage[key]
	return
}
