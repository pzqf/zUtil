package zMap

import (
	"hash/fnv"
)

type ShardedMap struct {
	shards    []*Map
	numShards int
}

func NewShardedMap(numShards int) *ShardedMap {
	if numShards <= 0 {
		numShards = 16
	}
	shards := make([]*Map, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewMap()
	}
	return &ShardedMap{
		shards:    shards,
		numShards: numShards,
	}
}

func NewShardedMap32() *ShardedMap {
	return NewShardedMap(32)
}

func (m *ShardedMap) GetShard(key interface{}) *Map {
	hash := m.hash(key)
	shardIndex := int(hash % int64(m.numShards))
	if shardIndex < 0 {
		shardIndex = -shardIndex
	}
	return m.shards[shardIndex]
}

func (m *ShardedMap) hash(key interface{}) int64 {
	switch k := key.(type) {
	case int64:
		return k
	case int32:
		return int64(k)
	case int:
		return int64(k)
	case string:
		h := fnv.New64a()
		h.Write([]byte(k))
		return int64(h.Sum64())
	case []byte:
		h := fnv.New64a()
		h.Write(k)
		return int64(h.Sum64())
	default:
		return int64(0)
	}
}

func (m *ShardedMap) Load(key interface{}) (interface{}, bool) {
	return m.GetShard(key).Load(key)
}

func (m *ShardedMap) Store(key, value interface{}) {
	m.GetShard(key).Store(key, value)
}

func (m *ShardedMap) Delete(key interface{}) {
	m.GetShard(key).Delete(key)
}

func (m *ShardedMap) Len() int64 {
	var total int64
	for _, shard := range m.shards {
		total += shard.Len()
	}
	return total
}

func (m *ShardedMap) Clear() {
	for _, shard := range m.shards {
		shard.Clear()
	}
}

func (m *ShardedMap) Range(f func(key, value interface{}) bool) {
	for _, shard := range m.shards {
		stopped := false
		shard.Range(func(key, value interface{}) bool {
			if !f(key, value) {
				stopped = true
				return false
			}
			return true
		})
		if stopped {
			return
		}
	}
}

func (m *ShardedMap) LoadOrStore(key, value interface{}) (interface{}, bool) {
	return m.GetShard(key).LoadOrStore(key, value)
}

func (m *ShardedMap) LoadAndDelete(key interface{}) (interface{}, bool) {
	return m.GetShard(key).LoadAndDelete(key)
}

func (m *ShardedMap) CompareAndDelete(key, oldValue interface{}) bool {
	return m.GetShard(key).CompareAndDelete(key, oldValue)
}

func (m *ShardedMap) CompareAndSwap(key, oldValue, newValue interface{}) bool {
	return m.GetShard(key).CompareAndSwap(key, oldValue, newValue)
}
