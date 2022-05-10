package zMap

import (
	"sync"
	"sync/atomic"
)

type zMap struct {
	Map   sync.Map
	Count int32
}

func NewMap() zMap {
	return zMap{
		Count: 0,
	}
}

func (m *zMap) Get(key interface{}) (interface{}, bool) {
	return m.Map.Load(key)
}

func (m *zMap) Store(key, value interface{}) {
	if _, ok := m.Map.Load(key); !ok {
		atomic.AddInt32(&m.Count, 1)
	}
	m.Map.Store(key, value)
}

func (m *zMap) Delete(key interface{}) {
	if _, ok := m.Map.Load(key); ok {
		m.Map.Delete(key)
		atomic.AddInt32(&m.Count, -1)
	}
}

func (m *zMap) Len() int32 {
	return m.Count
}

func (m *zMap) Range(f func(key, value interface{}) bool) {
	m.Map.Range(f)
}
