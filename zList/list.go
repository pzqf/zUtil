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

func (l *List) PushFront(v any) {
	l.locker.Lock()
	defer l.locker.Unlock()
	l.list.PushFront(v)
}

func (l *List) PushBack(v any) {
	l.locker.Lock()
	defer l.locker.Unlock()
	l.list.PushBack(v)
}

func (l *List) Front() any {
	l.locker.Lock()
	defer l.locker.Unlock()
	i := l.list.Front()
	l.list.Remove(i)
	return i.Value
}

func (l *List) Back() any {
	l.locker.Lock()
	defer l.locker.Unlock()
	i := l.list.Back()
	l.list.Remove(i)
	return i.Value
}

func (l *List) Len() any {
	return l.list.Len()
}

func (l *List) Range(f func(value any) bool) {
	e := l.list.Front()
	for {
		if e == nil {
			break
		}

		if !f(e.Value) {
			break
		}
		e = e.Next()
	}
}
