package zStack

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	value := "test data"
	stack := New(100)
	err := stack.Push(value)
	if err != nil {
		fmt.Println("stack push err", err)
	}
	d, err := stack.Pop()
	if err != nil {
		fmt.Println("stack pop err", err)
	}
	fmt.Println("pop data:", d)
}
