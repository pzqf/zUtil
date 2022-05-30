package zList

import (
	"container/list"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	l := New()
	//wg := sync.WaitGroup{}
	count := 10
	//wg.Add(count)
	for i := 0; i < count; i++ {
		//go func(n int) {
		l.PushBack(i)
		//	wg.Done()
		//}(i)
	}
	//	wg.Wait()
	l.Front()

	fmt.Println("=========", l.Len())

	l.Range(func(e *list.Element, value any) bool {
		if value.(int) > 5 {
			fmt.Println("remove", value)
			l.Remove(e)
		}
		return true
	})

	l.Range(func(e *list.Element, value any) bool {
		fmt.Println(value)
		return true
	})
}
