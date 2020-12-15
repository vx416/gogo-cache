package gocache

type ARC struct {
	size      int
	partition int

	t1 *LRU
	b1 *LRU

	t2 *LFU
	b2 *LFU
}

func (c *ARC) Get(key interface{}) (value interface{}, ok bool) {
	if c.t1.Contains(key) {
		value, ok = c.t1.Remove(key)
		c.t2.Add(key, value)
		return value, ok
	}

	if value, ok := c.t2.Get(key); ok {
		return value, ok
	}

	if c.b1.Contains(key) {
		c.b1.Remove(key)
	}
	if c.b2.Contains(key) {
		c.b2.Remove(key)
	}

	return nil, false
}

func (c *ARC) Add(key, value interface{}) (evicted bool) {

	if c.t1.Contains(key) {
		c.t1.Remove(key)
		c.t2.Add(key, value)
		return false
	}

	if c.t2.Contains(key) {
		c.t2.Add(key, value)
		return false
	}

	// 擴充 t1 size
	if c.b1.Contains(key) {
		c.extendSize(false)

		if c.isOverflow(3, false) {
			c.migrate(false)
			evicted = true
		}
		c.b1.Remove(key)
		c.t2.Add(key, value)
		return evicted
	}

	// 擴充 t2 size
	if c.b2.Contains(key) {
		c.extendSize(true)

		if c.isOverflow(3, false) {
			c.migrate(true)
			evicted = true
		}
		c.b2.Remove(key)
		c.t2.Add(key, value)
		return evicted
	}

	if c.isOverflow(3, false) {
		c.migrate(false)
		evicted = true
	}

	if c.isOverflow(1, true) {
		c.b1.RemoveOldest()
		evicted = true
	}

	if c.isOverflow(2, false) {
		c.b2.RemoveOldest()
		evicted = true
	}

	c.t1.Add(key, value)

	return evicted
}

func (c *ARC) isOverflow(check int8, ghost bool) bool {
	switch check {
	case 1:
		if ghost {
			return c.b1.Len() >= c.partition
		}
		return c.t1.Len() >= c.partition
	case 2:
		if ghost {
			return c.b2.Len() > c.size-c.partition
		}
		return c.t2.Len() > c.size-c.partition
	default:
		if ghost {
			return c.b1.Len()+c.b2.Len() >= c.size
		}
		return c.t1.Len()+c.t2.Len() >= c.size
	}
}

func (c *ARC) extendSize(t2size bool) {
	if t2size {
		c.partition--
	} else {
		c.partition++
	}
}

func (c *ARC) migrate(b2ContainsKey bool) {
	t1Len := c.t1.Len()
	if t1Len > 0 && (t1Len > c.partition || (t1Len == c.partition && b2ContainsKey)) {
		key, _, ok := c.t1.RemoveOldest()
		if ok {
			c.b1.Add(key, key)
		}
	} else {
		key, _, ok := c.t2.RemoveOldest()
		if ok {
			c.b2.Add(key, key)
		}
	}
}
