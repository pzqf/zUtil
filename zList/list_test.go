package zList

import (
	"fmt"
	"sync"
	"testing"
)

func Test(t *testing.T) {
	l := New()
	wg := sync.WaitGroup{}
	count := 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(n int) {
			l.PushBack(n)
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Println("=========", l.Len())

	l.Range(func(value any) bool {
		if value.(int) > 50 {
			return false
		}
		fmt.Println(value)
		return true
	})
}
