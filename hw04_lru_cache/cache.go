package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

/*
type cacheItem struct {
	key   Key
	value interface{}
}
*/

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	if v, ok := lc.items[key]; ok {
		v.Value = value
		lc.queue.MoveToFront(v)
		return ok
	}
	lc.items[key] = lc.queue.PushFront(value)
	if lc.queue.Len() > lc.capacity {
		toDel := lc.queue.Back()
		lc.queue.Remove(toDel)
		for key, val := range lc.items {
			if val.Value == toDel.Value {
				delete(lc.items, key)
			}
		}
	}
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(lc.items[key])
		return lc.items[key].Value, ok
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
