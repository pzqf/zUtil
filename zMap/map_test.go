package zMap

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	m := NewMap()
	m.Store(1, "aaa")
	m.Store(2, "bbb")
	m.Store("c", "ccc")

	fmt.Println("len:", m.Len())
	m.Range(func(key, value interface{}) bool {
		fmt.Println("===", key, value)
		return true
	})

	fmt.Println(m.Get(1))
	fmt.Println(m.Get("c"))

	m.Delete(1)
	fmt.Println("len:", m.Len())
	m.Range(func(key, value interface{}) bool {
		fmt.Println("===", key, value)
		return true
	})

	m.Clear()

	m.Range(func(key, value interface{}) bool {
		fmt.Println("=====", key, value)
		return true
	})
	fmt.Println("len:", m.Len())
}
