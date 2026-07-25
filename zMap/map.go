// zMap 包提供并发安全的 Map 实现
//
// zMap.Map 是基于 sync.Map 包装的并发安全 Map，增加了元素计数功能
// 它提供了与 sync.Map 完全兼容的接口，同时维护准确的元素数量
//
// 主要特点：
// - 并发安全：使用 sync.Map 作为底层实现，保持无锁读取特性
// - 元素计数：通过 atomic.Value 维护准确的元素数量
// - 接口兼容：与 sync.Map 完全兼容的方法签名
// - 易于使用：提供直观的 API，类似于普通 Map
// - 性能优异：保持 sync.Map 的性能优势
//
// 与 sync.Map 的区别：
// - zMap.Map 维护准确的元素计数（Len() 方法）
// - sync.Map 不提供 Len() 方法，需要手动遍历统计
// - zMap.Map 的接口与 sync.Map 完全兼容
//
// 使用场景：
// 1. 高并发读写场景：如游戏服务器玩家管理、缓存、会话管理等
// 2. 需要准确元素计数的场景：如统计在线玩家数量、缓存命中率等
// 3. 需要与 sync.Map 兼容的场景：便于迁移现有的 sync.Map 代码
// 4. 读多写少的场景：sync.Map 在这种场景下性能最优
//
// 不适合的场景：
// 1. 写多读少的场景：sync.Map 在这种场景下性能不如普通 Map 加锁
// 2. 需要排序的场景：zMap.Map 不保证元素顺序
// 3. 需要复杂查询的场景：如范围查询、模糊匹配等
// 4. 需要持久化的场景：zMap.Map 是内存 Map，不支持持久化
//
// 性能注意事项：
// 1. 避免频繁调用 Len() 方法：虽然使用 atomic.Value，但仍有性能开销
// 2. 避免频繁调用 Clear() 方法：需要遍历所有元素，性能开销较大
// 3. 避免在 Range 中修改 Map：可能导致遍历不完整或并发错误
// 4. 读多写少场景下性能最优：写多读少场景建议使用普通 Map 加锁
// 5. 键值类型建议使用基本类型：如 string、int、int64 等，避免使用复杂类型
package zMap

import (
	"sync"
	"sync/atomic"
)

// Map 是基于 sync.Map 包装的并发安全 Map
// 它提供了与 sync.Map 完全兼容的接口，同时维护准确的元素数量
//
// 底层使用 sync.Map，保持无锁读取特性
// 通过 atomic.Value 维护准确的元素计数
//
// 主要方法：
// - Get(key interface{}) (interface{}, bool)：获取指定键的值
// - Store(key, value interface{})：存储键值对
// - Delete(key interface{})：删除指定键
// - Len() int64：获取元素数量
// - Range(f func(key, value interface{}) bool)：遍历所有键值对
// - Clear()：清空所有元素
// - LoadOrStore(key, value interface{}) (interface{}, bool)：加载或存储
// - LoadAndDelete(key interface{}) (interface{}, bool)：加载并删除
// - CompareAndDelete(key, oldValue interface{}) bool：比较并删除
// - CompareAndSwap(key, oldValue, newValue interface{}) bool：比较并替换
type Map struct {
	sMap  sync.Map
	count atomic.Int64
}

// NewMap 创建一个新的并发安全 Map
// 注意事项：
// - 返回的 Map 实例是并发安全的
// - 不需要手动初始化，可直接使用
// - 不要对 Map 实例进行拷贝，可能导致并发问题
func NewMap() *Map {
	m := &Map{}
	return m
}

// Load 获取指定键的值
// 注意事项：
// - 这是一个并发安全的操作
// - 读操作是无锁的，性能极高
// - 返回的值需要进行类型断言，如 value.(string)、value.(int) 等
// - 建议在调用前先检查 exists，避免类型断言panic
func (m *Map) Load(key interface{}) (interface{}, bool) {
	return m.sMap.Load(key)
}

// Store 存储键值对
// 注意事项：
// - 这是一个并发安全的操作
// - 会自动维护元素计数：如果键不存在，计数加1
// - 会覆盖已存在的键值对
// - 键和值可以是任意类型，但建议使用基本类型
// - 写操作比读操作慢，适合读多写少的场景
func (m *Map) Store(key, value interface{}) {
	// OPT-13: 用原子 Swap 判断键是否已存在来维护计数。此前 `Load(!ok) 再 count.Add(1) 再 Store`
	// 是检查-再动作，两个 goroutine 并发 Store 同一新键会都判 !ok、都 +1 → Len() 多计。
	// Swap 一步原子完成"存值 + 返回是否已存在"，计数准确无竞态。
	if _, loaded := m.sMap.Swap(key, value); !loaded {
		m.count.Add(1)
	}
}

// Delete 删除指定键
// 注意事项：
// - 这是一个并发安全的操作
// - 会自动维护元素计数：如果键存在，计数减1
// - 如果键不存在，不会产生任何效果
// - 不会返回被删除的值，如需返回请使用 LoadAndDelete
func (m *Map) Delete(key interface{}) {
	if _, ok := m.sMap.LoadAndDelete(key); ok {
		m.count.Add(-1)
	}
}

// Len 获取 Map 的元素数量
// 注意事项：
// - 这是一个并发安全的操作
// - 使用 atomic.Value 维护，性能极高
// - 但仍有轻微性能，避免频繁调用
// - 返回的是瞬时值，可能在调用后立即改变
// - 如果需要精确统计，建议在临界区中调用
func (m *Map) Len() int64 {
	return m.count.Load()
}

// Range 遍历 Map 中的所有键值对
// 注意事项：
// - 这是一个并发安全的操作
// - 遍历期间可以进行读操作，但不建议进行写操作
// - 遍历期间如果有新的键值对插入，可能不会被遍历到
// - 遍历期间如果有键值对被删除，可能仍会被遍历到
// - 回调函数中如果修改 Map，可能导致遍历不完整
// - 键值对的顺序是随机的，不保证任何顺序
// - 如果回调函数返回 false，遍历会立即停止
// - 性能与元素数量成正比，大 Map 遍历可能耗时较长
func (m *Map) Range(f func(key, value interface{}) bool) {
	m.sMap.Range(f)
}

// Clear 清空 Map 中的所有元素
// 注意事项：
// - 这是一个并发安全的操作
// - 会重置元素计数为0
// - 性能与元素数量成正比，大 Map 清空可能耗时较长
// - 清空后可以继续使用，不需要重新创建实例
// - 如果在清空期间有其他 goroutine 进行写操作，可能导致部分元素未被清空
func (m *Map) Clear() {
	var keys []interface{}

	m.sMap.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})

	for _, v := range keys {
		m.sMap.Delete(v)
	}

	m.count.Store(0)
}

// LoadOrStore 如果键不存在则存储值并返回该值，否则返回已存在的值
// 注意事项：
// - 这是一个并发安全的操作
// - 会自动维护元素计数：如果键不存在，计数加1
// - 是原子操作，不会出现竞态条件
// - 适合用于并发安全的初始化
// - 返回的值需要进行类型断言
func (m *Map) LoadOrStore(key, value interface{}) (interface{}, bool) {
	existing, loaded := m.sMap.LoadOrStore(key, value)
	if !loaded {
		m.count.Add(1)
	}
	return existing, loaded
}

// LoadAndDelete 如果键存在则删除并返回值，否则返回零值
// 注意事项：
// - 这是一个并发安全的操作
// - 会自动维护元素计数：如果键存在，计数减1
// - 是原子操作，不会出现竞态条件
// - 适合用于需要同时删除和获取值的场景
// - 返回的值需要进行类型断言
// - 比 Delete() 多返回被删除的值
func (m *Map) LoadAndDelete(key interface{}) (interface{}, bool) {
	value, loaded := m.sMap.LoadAndDelete(key)
	if loaded {
		m.count.Add(-1)
	}
	return value, loaded
}

// CompareAndDelete 如果键存在且值匹配则删除，返回是否删除成功
// 注意事项：
// - 这是一个并发安全的操作
// - 会自动维护元素计数：如果删除成功，计数减1
// - 是原子操作，不会出现竞态条件
// - 适合用于需要条件删除的场景
// - 只有当键存在且值等于 oldValue 时才会删除
// - 比 Delete() 多了值匹配检查
func (m *Map) CompareAndDelete(key, oldValue interface{}) bool {
	deleted := m.sMap.CompareAndDelete(key, oldValue)
	if deleted {
		m.count.Add(-1)
	}
	return deleted
}

// CompareAndSwap 如果键存在且值匹配则替换，返回是否替换成功
// 注意事项：
// - 这是一个并发安全的操作
// - 不会改变元素计数（因为只是替换现有值）
// - 是原子操作，不会出现竞态条件
// - 适合用于需要条件替换的场景
// - 只有当键存在且值等于 oldValue 时才会替换
// - 比 Store() 多了值匹配检查
// - 不会新增元素，只会替换现有元素
// - 如果需要新增元素，请使用 Store() 或 LoadOrStore()
func (m *Map) CompareAndSwap(key, oldValue, newValue interface{}) bool {
	return m.sMap.CompareAndSwap(key, oldValue, newValue)
}
