package zQueue

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	q := Queue{}

	for i := 0; i < 10; i++ {
		q.Enqueue(i)
	}
	//q.Enqueue(1)
	//q.Enqueue(2)
	//q.Enqueue(3)
	//q.Enqueue(4)
	//q.Enqueue(5)

	fmt.Println("length:", q.Length())
	list := q.Get()
	for k, v := range list {
		fmt.Println(k, v)
	}

	for {
		d, ok := q.Dequeue()
		if !ok {
			fmt.Println("nothing data")
			break
		}

		fmt.Println("dequeue:", d.(int), "length:", q.Length())

		//fmt.Println("length:", q.Length())
		//list = q.Get()
		//for k, v := range list {
		//	fmt.Println(k, v)
		//}
	}

}
