package zHashtable

import (
	"hash/fnv"
	"sync"
)

type TableItem struct {
	key  string
	data interface{}
	next *TableItem
}

type HashTable struct {
	mu       sync.RWMutex
	table    []*TableItem
	capacity int
	size     int
}

func NewHashTable() *HashTable {
	return &HashTable{
		table:    make([]*TableItem, 256),
		capacity: 256,
		size:     0,
	}
}

func (ht *HashTable) Add(key string, value interface{}) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	if float64(ht.size+1)/float64(ht.capacity) > 0.75 {
		ht.resizeLocked()
	}

	position := ht.generateHash(key)
	if ht.table[position] == nil {
		ht.table[position] = &TableItem{key: key, data: value}
		ht.size++
		return
	}

	current := ht.table[position]
	for current != nil {
		if current.key == key {
			current.data = value
			return
		}
		if current.next == nil {
			break
		}
		current = current.next
	}

	current.next = &TableItem{key: key, data: value}
	ht.size++
}

func (ht *HashTable) Get(key string) (interface{}, bool) {
	ht.mu.RLock()
	defer ht.mu.RUnlock()

	position := ht.generateHash(key)
	current := ht.table[position]
	for current != nil {
		if current.key == key {
			return current.data, true
		}
		current = current.next
	}
	return nil, false
}

func (ht *HashTable) Set(key string, value interface{}) bool {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	position := ht.generateHash(key)
	current := ht.table[position]
	for current != nil {
		if current.key == key {
			current.data = value
			return true
		}
		current = current.next
	}
	return false
}

func (ht *HashTable) Remove(key string) bool {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	position := ht.generateHash(key)
	if ht.table[position] == nil {
		return false
	}

	if ht.table[position].key == key {
		ht.table[position] = ht.table[position].next
		ht.size--
		return true
	}

	current := ht.table[position]
	for current.next != nil {
		if current.next.key == key {
			current.next = current.next.next
			ht.size--
			return true
		}
		current = current.next
	}

	return false
}

func (ht *HashTable) Size() int {
	ht.mu.RLock()
	defer ht.mu.RUnlock()
	return ht.size
}

func (ht *HashTable) generateHash(s string) int {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return int(hash.Sum32() % uint32(ht.capacity))
}

func (ht *HashTable) resizeLocked() {
	newCapacity := ht.capacity * 2
	newTable := make([]*TableItem, newCapacity)

	oldTable := ht.table
	oldCapacity := ht.capacity

	ht.table = newTable
	ht.capacity = newCapacity
	ht.size = 0

	for i := 0; i < oldCapacity; i++ {
		current := oldTable[i]
		for current != nil {
			next := current.next
			position := ht.generateHash(current.key)
			current.next = newTable[position]
			newTable[position] = current
			ht.size++
			current = next
		}
	}
}
