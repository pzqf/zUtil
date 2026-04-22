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
	l.locker.Lock()
	defer l.locker.Unlock()
	return l.list.Len()
}

// Range 遍历链表，注意：回调中不能调用 Remove（会导致死锁）
// 如需在遍历时删除元素，请使用 RangeWithDelete
func (l *List) Range(f func(e *list.Element, value any) bool) {
	l.locker.Lock()
	defer l.locker.Unlock()
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

// Remove 移除指定元素，返回其值
func (l *List) Remove(e *list.Element) any {
	l.locker.Lock()
	defer l.locker.Unlock()
	if e == nil {
		return nil
	}
	l.list.Remove(e)
	return e.Value
}

// RangeWithDelete 安全的遍历并支持删除
// 回调返回 true 表示继续遍历，false 表示停止
// 第二个返回值标记是否需要删除当前元素
func (l *List) RangeWithDelete(f func(e *list.Element, value any) (bool, bool)) {
	l.locker.Lock()

	var toDelete []*list.Element
	e := l.list.Front()
	for {
		if e == nil {
			break
		}
		n := e.Next()
		cont, del := f(e, e.Value)
		if del {
			toDelete = append(toDelete, e)
		}
		if !cont {
			break
		}
		e = n
	}

	for _, elem := range toDelete {
		l.list.Remove(elem)
	}

	l.locker.Unlock()
}

// Clear 清空链表
func (l *List) Clear() {
	l.locker.Lock()
	defer l.locker.Unlock()
	l.list.Init()
}
