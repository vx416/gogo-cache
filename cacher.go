package gocache

type EvictCallback func(key interface{}, value interface{})
