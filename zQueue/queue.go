package zQueue

import "sync"

type Node struct {
	data interface{}
	next *Node
}

type Queue struct {
	locker sync.Mutex
	front  *Node
	rear   *Node
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Enqueue(i interface{}) {
	q.locker.Lock()
	defer q.locker.Unlock()
	data := &Node{data: i}
	if q.rear != nil {
		q.rear.next = data
	} else {
		q.front = data
	}
	q.rear = data
}

func (q *Queue) Dequeue() (interface{}, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.front == nil {
		return nil, false
	}

	data := q.front.data
	q.front = q.front.next

	if q.front == nil {
		q.rear = nil
	}

	return data, true
}

func (q *Queue) Peek() (interface{}, bool) {
	q.locker.Lock()
	defer q.locker.Unlock()
	if q.front == nil {
		return nil, false
	}
	return q.front.data, true
}

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

func (q *Queue) IsEmpty() bool {
	q.locker.Lock()
	defer q.locker.Unlock()
	return q.front == nil
}

func (q *Queue) Empty() {
	q.locker.Lock()
	defer q.locker.Unlock()
	q.front = nil
	q.rear = nil
}

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

type RingQueue struct {
	mu       sync.Mutex
	data     []interface{}
	capacity int
	head     int
	tail     int
	count    int
}

func NewRingQueue(capacity int) *RingQueue {
	if capacity <= 0 {
		capacity = 1024
	}
	return &RingQueue{
		data:     make([]interface{}, capacity),
		capacity: capacity,
	}
}

func (rq *RingQueue) Enqueue(item interface{}) bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if rq.count >= rq.capacity {
		rq.grow()
	}

	rq.data[rq.tail] = item
	rq.tail = (rq.tail + 1) % rq.capacity
	rq.count++
	return true
}

func (rq *RingQueue) Dequeue() (interface{}, bool) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if rq.count == 0 {
		return nil, false
	}

	item := rq.data[rq.head]
	rq.data[rq.head] = nil
	rq.head = (rq.head + 1) % rq.capacity
	rq.count--
	return item, true
}

func (rq *RingQueue) Peek() (interface{}, bool) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if rq.count == 0 {
		return nil, false
	}
	return rq.data[rq.head], true
}

func (rq *RingQueue) IsEmpty() bool {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	return rq.count == 0
}

func (rq *RingQueue) Len() int {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	return rq.count
}

func (rq *RingQueue) Cap() int {
	return rq.capacity
}

func (rq *RingQueue) Clear() {
	rq.mu.Lock()
	defer rq.mu.Unlock()
	rq.data = make([]interface{}, rq.capacity)
	rq.head = 0
	rq.tail = 0
	rq.count = 0
}

func (rq *RingQueue) grow() {
	newCapacity := rq.capacity * 2
	newData := make([]interface{}, newCapacity)

	for i := 0; i < rq.count; i++ {
		newData[i] = rq.data[(rq.head+i)%rq.capacity]
	}

	rq.data = newData
	rq.capacity = newCapacity
	rq.head = 0
	rq.tail = rq.count
}
