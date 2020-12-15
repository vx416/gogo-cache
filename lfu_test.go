package gocache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLFU(t *testing.T) {
	cache := newLFU(3)
	key1 := "test"
	key2 := "test2"
	key3 := "test3"
	key4 := "test4"

	evicted := cache.Add(key1, 1)
	assert.False(t, evicted)

	evicted = cache.Add(key2, 1)
	assert.False(t, evicted)

	evicted = cache.Add(key3, 1)
	assert.False(t, evicted)

	val, ok := cache.Get(key1)
	assert.True(t, ok)
	assert.Equal(t, val, 1)

	val, ok = cache.Get(key2)
	assert.True(t, ok)
	assert.Equal(t, val, 1)

	lfuKey := cache.list.Back().Value.(*lfuEntry).key
	assert.Equal(t, lfuKey, key3)

	evicted = cache.Add(key4, 2)
	assert.True(t, evicted)

	ok = cache.Contains(key3)
	assert.False(t, ok)
}
