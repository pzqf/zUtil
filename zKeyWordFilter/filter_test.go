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

	AddWord("abc")

	c := Filter("aabcabc")
	fmt.Println(c)
}
