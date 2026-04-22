package zList

import (
	"container/list"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
	assert.Equal(t, 0, l.Len())
}

func TestPushFront(t *testing.T) {
	l := New()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)

	assert.Equal(t, 3, l.Len())

	front := l.Front()
	assert.NotNil(t, front)
	assert.Equal(t, 3, front.Value)
}

func TestPushBack(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	assert.Equal(t, 3, l.Len())

	back := l.Back()
	assert.NotNil(t, back)
	assert.Equal(t, 3, back.Value)
}

func TestFrontBack(t *testing.T) {
	l := New()
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())

	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	assert.Equal(t, 1, l.Front().Value)
	assert.Equal(t, 3, l.Back().Value)
}

func TestRemove(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	elem := l.Front()
	val := l.Remove(elem)
	assert.Equal(t, 1, val)
	assert.Equal(t, 2, l.Len())

	assert.Equal(t, 2, l.Front().Value)
}

func TestRemoveNil(t *testing.T) {
	l := New()
	val := l.Remove(nil)
	assert.Nil(t, val)
}

func TestRange(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	var values []any
	l.Range(func(e *list.Element, value any) bool {
		values = append(values, value)
		return true
	})

	assert.Equal(t, []any{1, 2, 3}, values)
}

func TestRangeStopEarly(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	var values []any
	l.Range(func(e *list.Element, value any) bool {
		values = append(values, value)
		return value.(int) < 2
	})

	assert.Equal(t, []any{1, 2}, values)
}

func TestRangeEmpty(t *testing.T) {
	l := New()
	called := false
	l.Range(func(e *list.Element, value any) bool {
		called = true
		return true
	})
	assert.False(t, called)
}

func TestRangeWithDelete(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	l.PushBack(5)

	l.RangeWithDelete(func(e *list.Element, value any) (bool, bool) {
		v := value.(int)
		if v > 3 {
			return true, true
		}
		return true, false
	})

	assert.Equal(t, 3, l.Len())

	var values []any
	l.Range(func(e *list.Element, value any) bool {
		values = append(values, value)
		return true
	})
	assert.Equal(t, []any{1, 2, 3}, values)
}

func TestRangeWithDeleteAll(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	l.RangeWithDelete(func(e *list.Element, value any) (bool, bool) {
		return true, true
	})

	assert.Equal(t, 0, l.Len())
}

func TestRangeWithDeleteStopEarly(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	l.RangeWithDelete(func(e *list.Element, value any) (bool, bool) {
		return false, true
	})

	assert.Equal(t, 2, l.Len())
}

func TestClear(t *testing.T) {
	l := New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)

	l.Clear()
	assert.Equal(t, 0, l.Len())
	assert.Nil(t, l.Front())
	assert.Nil(t, l.Back())
}

func TestConcurrentPush(t *testing.T) {
	l := New()
	var wg sync.WaitGroup
	count := 1000

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(n int) {
			defer wg.Done()
			l.PushBack(n)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, count, l.Len())
}

func TestConcurrentPushAndRemove(t *testing.T) {
	l := New()
	var wg sync.WaitGroup
	count := 100

	for i := 0; i < count; i++ {
		l.PushBack(i)
	}

	wg.Add(count)
	e := l.Front()
	for i := 0; i < count; i++ {
		elem := e
		e = e.Next()
		go func(el *list.Element) {
			defer wg.Done()
			l.Remove(el)
		}(elem)
	}
	wg.Wait()

	assert.Equal(t, 0, l.Len())
}

func TestConcurrentRangeWithDelete(t *testing.T) {
	l := New()
	for i := 0; i < 100; i++ {
		l.PushBack(i)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		l.RangeWithDelete(func(e *list.Element, value any) (bool, bool) {
			return true, value.(int)%2 == 0
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		l.PushBack(100)
	}()

	wg.Wait()
}
