package zMap

import (
	"sync"
	"sync/atomic"
)

type Map struct {
	sMap  sync.Map
	count int32
}

func NewMap() Map {
	return Map{
		count: 0,
	}
}

func (m *Map) Get(key interface{}) (interface{}, bool) {
	return m.sMap.Load(key)
}

func (m *Map) Store(key, value interface{}) {
	if _, ok := m.sMap.Load(key); !ok {
		atomic.AddInt32(&m.count, 1)
	}
	m.sMap.Store(key, value)
}

func (m *Map) Delete(key interface{}) {
	if _, ok := m.sMap.Load(key); ok {
		m.sMap.Delete(key)
		atomic.AddInt32(&m.count, -1)
	}
}

func (m *Map) Len() int32 {
	return m.count
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
}