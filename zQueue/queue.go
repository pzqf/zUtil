package zQueue

import (
	"sync"
)

type Node struct {
	data interface{}
	next *Node
}

type Queue struct {
	locker sync.Mutex
	rear   *Node
}

func NewQueue() *Queue {
	return &Queue{
		locker: sync.Mutex{},
		rear:   nil,
	}
}

func (q *Queue) Enqueue(i interface{}) {
	q.locker.Lock()
	defer q.locker.Unlock()
	data := &Node{data: i}
	if q.rear != nil {
		data.next = q.rear
	}
	q.rear = data
}

func (q *Queue) Dequeue() (interface{}, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.rear == nil {
		return 0, false
	}
	if q.rear.next == nil {
		i := q.rear.data
		q.rear = nil
		return i, true
	}
	current := q.rear
	for {
		if current.next.next == nil {
			i := current.next.data
			current.next = nil
			return i, true
		}
		current = current.next
	}
}

func (q *Queue) Peek() (interface{}, bool) {
	if q.rear == nil {
		return 0, false
	}
	return q.rear.data, true
}

func (q *Queue) Get() []interface{} {
	var items []interface{}
	current := q.rear
	for current != nil {
		items = append(items, current.data)
		current = current.next
	}
	return items
}

func (q *Queue) IsEmpty() bool {
	return q.rear == nil
}

func (q *Queue) Empty() {
	q.rear = nil
}

func (q *Queue) Length() int {
	n := 0
	it := q.rear
	for {
		if it == nil {
			break
		}
		n = n + 1
		it = it.next
	}

	return n
}
