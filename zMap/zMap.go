package zMap

import (
	"sync"
	"sync/atomic"
)

type zMap struct {
	sMap  sync.Map
	count int32
}

func NewMap() zMap {
	return zMap{
		count: 0,
	}
}

func (m *zMap) Get(key interface{}) (interface{}, bool) {
	return m.sMap.Load(key)
}

func (m *zMap) Store(key, value interface{}) {
	if _, ok := m.sMap.Load(key); !ok {
		atomic.AddInt32(&m.count, 1)
	}
	m.sMap.Store(key, value)
}

func (m *zMap) Delete(key interface{}) {
	if _, ok := m.sMap.Load(key); ok {
		m.sMap.Delete(key)
		atomic.AddInt32(&m.count, -1)
	}
}

func (m *zMap) Len() int32 {
	return m.count
}

func (m *zMap) Range(f func(key, value interface{}) bool) {
	m.sMap.Range(f)
}
