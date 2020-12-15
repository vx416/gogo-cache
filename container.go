package gocache

import "container/list"

func newContainer() *kvContainer {
	return &kvContainer{
		list:  list.New(),
		items: make(map[interface{}]*list.Element),
	}
}

type kvContainer struct {
	list  *list.List
	items map[interface{}]*list.Element
}
