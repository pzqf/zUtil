package zMap

import (
	"sync"
	"sync/atomic"
)

type Map struct {
	sMap  sync.Map
	count atomic.Value
}

func NewMap() *Map {
	m := &Map{}
	m.count.Store(int64(0))
	return m
}

func (m *Map) Get(key interface{}) (interface{}, bool) {
	return m.sMap.Load(key)
}

func (m *Map) Store(key, value interface{}) {
	if _, ok := m.sMap.Load(key); !ok {
		current := m.count.Load().(int64)
		m.count.Store(current + 1)
	}
	m.sMap.Store(key, value)
}

func (m *Map) Delete(key interface{}) {
	if _, ok := m.sMap.Load(key); ok {
		m.sMap.Delete(key)
		current := m.count.Load().(int64)
		m.count.Store(current - 1)
	}
}

func (m *Map) Len() int64 {
	return m.count.Load().(int64)
}

func (m *Map) Range(f func(key, value interface{}) bool) {
	m.sMap.Range(f)
}

func (m *Map) Clear() {
	var keys []interface{}

	m.sMap.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})

	for _, v := range keys {
		m.sMap.Delete(v)
	}

	m.count.Store(int64(0))
}
