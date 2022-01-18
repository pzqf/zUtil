package zTree

import (
	"fmt"
	"testing"
)

type TestData struct {
	ID   int    `json:"id"`
	Pid  int    `json:"pid""`
	Name string `json:"name""`
}

func (td *TestData) GetId() int {
	return td.ID
}

func (td *TestData) GetFatherId() int {
	return td.Pid
}

func (td *TestData) GetData() interface{} {
	return td
}

func (td *TestData) IsRoot() bool {
	return td.Pid == 0
}

func TestTree(t *testing.T) {
	var listData = []*TestData{
		{ID: 1, Pid: 0, Name: "1"},
		{ID: 2, Pid: 1, Name: "2"},
		{ID: 3, Pid: 1, Name: "3"},
		{ID: 4, Pid: 2, Name: "4"},
	}

	var list []INode
	for _, v := range listData {
		list = append(list, v)
	}
	treeList := GenerateTree(list)
	fmt.Println(treeList)
}
