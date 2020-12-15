package gocache

import (
	"container/list"
)

func newLFU(size int) *LFU {
	return &LFU{
		size:  size,
		list:  list.New(),
		items: make(map[interface{}]*list.Element),
	}
}

type LFU struct {
	size    int
	list    *list.List
	items   map[interface{}]*list.Element
	onEvict EvictCallback
}

type lfuEntry struct {
	key   interface{}
	value interface{}
	freq  int
}

func (c *LFU) Resize(size int) {
	c.size = size
}

func (c *LFU) Add(key, value interface{}) (evicted bool) {
	// check if element found
	if element, ok := c.items[key]; ok {
		ent := element.Value.(*lfuEntry)
		ent.freq++
		c.increment(element)
		ent.value = value
		return false
	}

	var element *list.Element
	ent := &lfuEntry{key: key, value: value}
	rearElem := c.list.Back()
	if rearElem != nil {
		element = c.list.InsertBefore(ent, c.list.Back())
	} else {
		element = c.list.PushBack(ent)
	}

	c.items[key] = element

	evicted = c.list.Len() > c.size
	if evicted {
		c.RemoveOldest()
	}
	return evicted
}

func (c *LFU) Get(key interface{}) (value interface{}, ok bool) {
	elem, ok := c.items[key]
	if ok {
		ent := elem.Value.(*lfuEntry)
		ent.freq++
		c.increment(elem)
		if ent == nil {
			return nil, false
		}

		return ent.value, true
	}

	return nil, false
}

func (c *LFU) Len() int {
	return c.list.Len()
}

func (c *LFU) Remove(key interface{}) (value interface{}, ok bool) {
	elem, ok := c.items[key]
	if ok {
		c.removeElement(elem)
		return elem.Value.(lfuEntry).value, ok
	}

	return nil, false
}

func (c *LFU) RemoveOldest() (key, value interface{}, ok bool) {
	rearElem := c.list.Back()
	if rearElem != nil {
		c.removeElement(rearElem)
		ent := rearElem.Value.(*lfuEntry)
		return ent.key, ent.value, true
	}

	return nil, nil, false
}

func (c *LFU) Contains(key interface{}) (ok bool) {
	_, ok = c.items[key]
	return ok
}

func (c *LFU) increment(ele *list.Element) {
	ent := ele.Value.(*lfuEntry)

	prevElemt := ele.Prev()

	for prevElemt != nil && prevElemt.Value.(*lfuEntry).freq < ent.freq {
		prevElemt = prevElemt.Prev()
	}

	if prevElemt == nil {
		c.list.MoveToFront(ele)
		return
	}
	c.list.MoveAfter(ele, prevElemt)
}

func (c *LFU) removeElement(elem *list.Element) {
	c.list.Remove(elem)
	ent := elem.Value.(*lfuEntry)
	delete(c.items, ent.key)
	if c.onEvict != nil {
		c.onEvict(ent.key, ent.value)
	}
}
