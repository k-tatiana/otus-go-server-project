package cache

import (
	"container/list"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
}

type entry struct {
	key   string
	value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			delete(c.cache, oldest.Value.(*entry).key)
			c.list.Remove(oldest)
		}
	}

	elem := c.list.PushFront(&entry{key, value})
	c.cache[key] = elem
}

func (c *LRUCache) Prepare() {
	// This method can be used to initialize or pre-load the cache if needed.
	// Currently, it does nothing but is kept for interface consistency.
}
