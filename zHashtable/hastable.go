package zHashtable

import "hash/fnv"

type TableItem struct {
	key  string
	data interface{}
	next *TableItem
}

type HashTable struct {
	table [256]*TableItem
}

func (ht *HashTable) Add(key string, value interface{}) {
	position := generateHash(key)
	if ht.table[position] == nil {
		ht.table[position] = &TableItem{key: key, data: value}
		return
	}
	current := ht.table[position]
	for current.next != nil {
		current = current.next
	}
	current.next = &TableItem{key: key, data: value}
}

func (ht *HashTable) Get(key string) (interface{}, bool) {
	position := generateHash(key)
	current := ht.table[position]
	for current != nil {
		if current.key == key {
			return current.data, true
		}
		current = current.next
	}
	return 0, false
}

func (ht *HashTable) Set(key string, value interface{}) bool {
	position := generateHash(key)
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
	position := generateHash(key)
	if ht.table[position] == nil {
		return false
	}
	if ht.table[position].key == key {
		ht.table[position] = ht.table[position].next
		return true
	}
	current := ht.table[position]
	for current.next != nil {
		if current.next.key == key {
			current.next = current.next.next
			return true
		}
		current = current.next
	}
	return false
}

func generateHash(s string) uint8 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(s))
	return uint8(hash.Sum32() % 256)
}
