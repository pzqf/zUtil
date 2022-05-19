package zQueue

type Node struct {
	data interface{}
	next *Node
}

type Queue struct {
	rear *Node
}

func (list *Queue) Enqueue(i interface{}) {
	data := &Node{data: i}
	if list.rear != nil {
		data.next = list.rear
	}
	list.rear = data
}

func (list *Queue) Dequeue() (interface{}, bool) {
	if list.rear == nil {
		return 0, false
	}
	if list.rear.next == nil {
		i := list.rear.data
		list.rear = nil
		return i, true
	}
	current := list.rear
	for {
		if current.next.next == nil {
			i := current.next.data
			current.next = nil
			return i, true
		}
		current = current.next
	}
}

func (list *Queue) Peek() (interface{}, bool) {
	if list.rear == nil {
		return 0, false
	}
	return list.rear.data, true
}

func (list *Queue) Get() []interface{} {
	var items []interface{}
	current := list.rear
	for current != nil {
		items = append(items, current.data)
		current = current.next
	}
	return items
}

func (list *Queue) IsEmpty() bool {
	return list.rear == nil
}

func (list *Queue) Empty() {
	list.rear = nil
}

func (list *Queue) Length() int {
	n := 0
	it := list.rear
	for {
		if it == nil {
			break
		}
		n = n + 1
		it = it.next
	}

	return n
}
