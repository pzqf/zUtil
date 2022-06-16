package zRand

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	for i := 0; i < 110; i++ {
		x := RandN(100)
		if x == 0 {
			fmt.Println("=====", x)
		}
		if x == 100 {
			fmt.Println("------", x)
		}
		//time.Sleep(time.Nanosecond)
		fmt.Println(x)
	}

	fmt.Println(RandInterval(9, 1))
	fmt.Println(RandInterval(1, 9))
	fmt.Println(RandInterval(1, 1))
	fmt.Println(RandInterval(-5, -3))

	fmt.Println(RandString(10))
	fmt.Println(RandString(10))
}
