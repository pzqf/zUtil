package zMap

// TypedMap 是基于泛型的类型安全并发 Map
// 内部使用 zMap.Map，接口与 Map 保持一致
type TypedMap[K comparable, V any] struct {
	inner *Map
}

func NewTypedMap[K comparable, V any]() *TypedMap[K, V] {
	return &TypedMap[K, V]{
		inner: NewMap(),
	}
}

func (m *TypedMap[K, V]) Load(key K) (V, bool) {
	value, exists := m.inner.Load(key)
	if !exists {
		var zero V
		return zero, false
	}
	return value.(V), true
}

func (m *TypedMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

func (m *TypedMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

func (m *TypedMap[K, V]) Len() int64 {
	return m.inner.Len()
}

func (m *TypedMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (m *TypedMap[K, V]) Clear() {
	m.inner.Clear()
}

func (m *TypedMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	existing, loaded := m.inner.LoadOrStore(key, value)
	if !loaded {
		return value, false
	}
	return existing.(V), true
}

func (m *TypedMap[K, V]) LoadAndDelete(key K) (V, bool) {
	value, loaded := m.inner.LoadAndDelete(key)
	if !loaded {
		var zero V
		return zero, false
	}
	return value.(V), true
}

func (m *TypedMap[K, V]) CompareAndDelete(key K, oldValue V) bool {
	return m.inner.CompareAndDelete(key, oldValue)
}

func (m *TypedMap[K, V]) CompareAndSwap(key K, oldValue V, newValue V) bool {
	return m.inner.CompareAndSwap(key, oldValue, newValue)
}
