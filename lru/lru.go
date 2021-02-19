package lru

import "container/list"

// lru cache struct,It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64
	nbytes    int64
	list      *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

// cache element struct
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// new Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes = int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.list.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {

	}
}

// Get
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 淘汰策略淘汰最尾部的
func (c *Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Delete
func (c Cache) Delete(key string) (bool, error) {
	return false, nil
}
