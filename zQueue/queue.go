package zQueue

import (
	"sync"
)

// Node 队列节点结构
type Node struct {
	data interface{} // 节点数据
	next *Node       // 下一个节点指针
}

// Queue 队列结构
// 实现了线程安全的队列，支持O(1)时间复杂度的入队和出队操作
type Queue struct {
	locker sync.Mutex // 互斥锁，保证线程安全
	front  *Node      // 队首指针
	rear   *Node      // 队尾指针
}

// NewQueue 创建一个新的队列
func NewQueue() *Queue {
	return &Queue{
		locker: sync.Mutex{},
		front:  nil,
		rear:   nil,
	}
}

// Enqueue 向队列尾部添加元素
// 时间复杂度: O(1)
func (q *Queue) Enqueue(i interface{}) {
	q.locker.Lock()
	defer q.locker.Unlock()
	data := &Node{data: i}
	if q.rear != nil {
		q.rear.next = data
	} else {
		// 队列为空时，front也指向新节点
		q.front = data
	}
	q.rear = data
}

// Dequeue 从队列头部移除并返回元素
// 返回值: 元素数据和是否成功的标志
// 时间复杂度: O(1)
func (q *Queue) Dequeue() (interface{}, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.front == nil {
		return nil, false
	}

	data := q.front.data
	q.front = q.front.next

	// 如果队列为空，rear也置为nil
	if q.front == nil {
		q.rear = nil
	}

	return data, true
}

// Peek 返回队列头部元素但不移除
// 返回值: 元素数据和是否成功的标志
// 时间复杂度: O(1)
func (q *Queue) Peek() (interface{}, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.front == nil {
		return nil, false
	}
	return q.front.data, true
}

// Get 返回队列中所有元素的切片
// 元素顺序为从队首到队尾
// 时间复杂度: O(n)
func (q *Queue) Get() []interface{} {
	q.locker.Lock()
	defer q.locker.Unlock()
	var items []interface{}
	current := q.front
	for current != nil {
		items = append(items, current.data)
		current = current.next
	}
	return items
}

// IsEmpty 判断队列是否为空
// 返回值: 如果队列为空返回true，否则返回false
// 时间复杂度: O(1)
func (q *Queue) IsEmpty() bool {
	q.locker.Lock()
	defer q.locker.Unlock()
	return q.front == nil
}

// Empty 清空队列
// 时间复杂度: O(1)
func (q *Queue) Empty() {
	q.locker.Lock()
	defer q.locker.Unlock()
	q.front = nil
	q.rear = nil
}

// Length 返回队列中元素的数量
// 时间复杂度: O(n)
func (q *Queue) Length() int {
	q.locker.Lock()
	defer q.locker.Unlock()
	n := 0
	it := q.front
	for it != nil {
		n++
		it = it.next
	}
	return n
}
