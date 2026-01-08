package zHashtable

import "hash/fnv"

// TableItem 哈希表的链表节点结构
type TableItem struct {
	key  string        // 键
	data interface{}   // 值
	next *TableItem    // 下一个节点指针
}

// HashTable 哈希表结构
type HashTable struct {
	table    []*TableItem // 哈希表数组
	capacity int          // 哈希表容量
	size     int          // 哈希表中元素数量
}

// NewHashTable 创建一个新的哈希表
// 默认初始容量为256
func NewHashTable() *HashTable {
	return &HashTable{
		table:    make([]*TableItem, 256),
		capacity: 256,
		size:     0,
	}
}

// Add 向哈希表中添加键值对
// 如果键已存在，则更新对应的值
// 当哈希表负载因子超过0.75时，会自动扩容
func (ht *HashTable) Add(key string, value interface{}) {
	// 检查是否需要扩容
	if float64(ht.size+1)/float64(ht.capacity) > 0.75 {
		ht.resize()
	}

	position := ht.generateHash(key)
	if ht.table[position] == nil {
		ht.table[position] = &TableItem{key: key, data: value}
		ht.size++
		return
	}

	// 检查是否已存在相同的键
	current := ht.table[position]
	for current != nil {
		if current.key == key {
			current.data = value // 更新已有键的值
			return
		}
		if current.next == nil {
			break
		}
		current = current.next
	}

	// 添加到链表末尾
	current.next = &TableItem{key: key, data: value}
	ht.size++
}

// Get 根据键从哈希表中获取值
// 返回值和是否存在的标志
func (ht *HashTable) Get(key string) (interface{}, bool) {
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

// Set 更新哈希表中指定键的值
// 如果键存在则更新并返回true，否则返回false
func (ht *HashTable) Set(key string, value interface{}) bool {
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

// Remove 从哈希表中删除指定键的元素
// 如果键存在则删除并返回true，否则返回false
func (ht *HashTable) Remove(key string) bool {
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

// generateHash 生成键的哈希值
// 使用fnv-1a算法，返回哈希表中的位置
func (ht *HashTable) generateHash(s string) int {
	hash := fnv.New32a()
	// hash.Write在写入字节切片时不会返回错误
	hash.Write([]byte(s))
	return int(hash.Sum32() % uint32(ht.capacity))
}

// resize 扩容哈希表并重新哈希所有元素
func (ht *HashTable) resize() {
	newCapacity := ht.capacity * 2
	newTable := make([]*TableItem, newCapacity)

	// 保存旧表信息
	oldTable := ht.table
	oldCapacity := ht.capacity

	// 更新哈希表结构
	ht.table = newTable
	ht.capacity = newCapacity
	ht.size = 0

	// 重新哈希所有元素
	for i := 0; i < oldCapacity; i++ {
		current := oldTable[i]
		for current != nil {
			// 保存当前节点的下一个节点
			next := current.next
			// 重新插入当前节点
			position := ht.generateHash(current.key)
			current.next = newTable[position]
			newTable[position] = current
			// 更新size
			ht.size++
			// 处理下一个节点
			current = next
		}
	}
}
