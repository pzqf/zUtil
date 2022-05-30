package zList

import (
	"container/list"
	"sync"
)

type List struct {
	list   *list.List
	locker sync.Mutex
}

func New() *List {
	return &List{
		list:   list.New(),
		locker: sync.Mutex{},
	}
}

func (l *List) PushFront(v any) *list.Element {
	l.locker.Lock()
	defer l.locker.Unlock()
	return l.list.PushFront(v)
}

func (l *List) PushBack(v any) *list.Element {
	l.locker.Lock()
	defer l.locker.Unlock()
	return l.list.PushBack(v)
}

func (l *List) Front() *list.Element {
	l.locker.Lock()
	defer l.locker.Unlock()
	return l.list.Front()
}

func (l *List) Back() *list.Element {
	l.locker.Lock()
	defer l.locker.Unlock()
	return l.list.Back()
}

func (l *List) Len() int {
	return l.list.Len()
}

func (l *List) Range(f func(e *list.Element, value any) bool) {
	e := l.list.Front()
	for {
		if e == nil {
			break
		}
		n := e.Next()
		if !f(e, e.Value) {
			break
		}
		e = n
	}
}

func (l *List) Remove(e *list.Element) any {
	l.locker.Lock()
	defer l.locker.Unlock()
	l.list.Remove(e)
	return e.Value
}
