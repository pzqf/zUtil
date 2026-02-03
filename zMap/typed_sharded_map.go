package zMap

type TypedShardedMap[K comparable, V any] struct {
	inner *ShardedMap
}

func NewTypedShardedMap[K comparable, V any](numShards int) *TypedShardedMap[K, V] {
	return &TypedShardedMap[K, V]{
		inner: NewShardedMap(numShards),
	}
}

func NewTypedShardedMap32[K comparable, V any]() *TypedShardedMap[K, V] {
	return &TypedShardedMap[K, V]{
		inner: NewShardedMap32(),
	}
}

func (m *TypedShardedMap[K, V]) Load(key K) (V, bool) {
	value, exists := m.inner.Load(key)
	if !exists {
		var zero V
		return zero, false
	}
	return value.(V), true
}

func (m *TypedShardedMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

func (m *TypedShardedMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

func (m *TypedShardedMap[K, V]) Len() int64 {
	return m.inner.Len()
}

func (m *TypedShardedMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (m *TypedShardedMap[K, V]) Clear() {
	m.inner.Clear()
}

func (m *TypedShardedMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	existing, loaded := m.inner.LoadOrStore(key, value)
	if !loaded {
		return value, false
	}
	return existing.(V), true
}

func (m *TypedShardedMap[K, V]) LoadAndDelete(key K) (V, bool) {
	value, loaded := m.inner.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	return value.(V), true
}

func (m *TypedShardedMap[K, V]) CompareAndDelete(key K, oldValue V) bool {
	return m.inner.CompareAndDelete(key, oldValue)
}

func (m *TypedShardedMap[K, V]) CompareAndSwap(key K, oldValue V, newValue V) bool {
	return m.inner.CompareAndSwap(key, oldValue, newValue)
}
