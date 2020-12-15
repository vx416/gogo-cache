package gocache

import "container/list"

func newLRU(size int) *LRU {
	return &LRU{
		size:  size,
		list:  list.New(),
		items: make(map[interface{}]*list.Element),
	}
}

// LRU latest recently used cache
type LRU struct {
	size    int
	list    *list.List
	items   map[interface{}]*list.Element
	onEvict EvictCallback
}

type lruEntry struct {
	key   interface{}
	value interface{}
}

func (c *LRU) Add(key, value interface{}) (evicted bool) {
	// check if element found
	if element, ok := c.items[key]; ok {
		c.list.MoveToFront(element)
		element.Value.(*lruEntry).value = value
		return false
	}

	ent := &lruEntry{key: key, value: value}
	element := c.list.PushFront(ent)
	c.items[key] = element

	evicted = c.list.Len() > c.size
	if evicted {
		c.RemoveOldest()
	}
	return evicted
}

func (c *LRU) Contains(key interface{}) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
	elem, ok := c.items[key]
	if ok {
		c.list.MoveToFront(elem)
		ent := elem.Value.(*lruEntry)
		if ent == nil {
			return nil, false
		}

		return ent.value, true
	}

	return nil, false
}

func (c *LRU) Len() int {
	return c.list.Len()
}

func (c *LRU) Remove(key interface{}) (value interface{}, ok bool) {
	elem, ok := c.items[key]
	if ok {
		c.removeElement(elem)
		return elem.Value.(lruEntry).value, ok
	}

	return nil, false
}

func (c *LRU) removeElement(elem *list.Element) {
	c.list.Remove(elem)
	ent := elem.Value.(*lruEntry)
	delete(c.items, ent.key)
	if c.onEvict != nil {
		c.onEvict(ent.key, ent.value)
	}
}

func (c *LRU) Resize(size int) {
	c.size = size
}

func (c *LRU) RemoveOldest() (key, value interface{}, ok bool) {
	rearElem := c.list.Back()
	if rearElem != nil {
		c.removeElement(rearElem)
		lruEntry := rearElem.Value.(*lruEntry)
		return lruEntry.key, lruEntry.value, true
	}

	return nil, nil, false
}
