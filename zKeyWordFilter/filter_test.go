package zKeyWordFilter

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	InitDefaultFilter()
	/*
		err := ParseFromFile("keyword.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	AddWord("张三")
	AddWord("张三四")
	AddWord("张三四五")

	//PrintDefault()

	c := Filter("测张三四试")
	fmt.Println(c)
}
