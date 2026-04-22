package zKeyWordFilter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
)

type KeyWord struct {
	Children map[rune]*KeyWord
	End      bool
}

func (k *KeyWord) SplitString(str string) {
	rs := []rune(str)
	n := k
	for _, v := range rs {
		if n.Children == nil {
			n.Children = make(map[rune]*KeyWord)
		}

		if _, ok := n.Children[v]; ok {
			n = n.Children[v]
			continue
		} else {
			n.Children[v] = &KeyWord{
				Children: nil,
			}
			n = n.Children[v]
		}
	}
	n.End = true
}

func (k *KeyWord) Find(r rune) *KeyWord {
	if k.Children == nil {
		return nil
	}
	if _, ok := k.Children[r]; ok {
		return k.Children[r]
	}
	return nil
}

func (k *KeyWord) Print(x int) {
	strPre := ""
	for i := 0; i < x; i++ {
		strPre += "  "
	}
	str := fmt.Sprintf("%d", len(k.Children))
	if len(k.Children) > 0 {
		fmt.Println(strPre + "--" + str)
	}

	for k, v := range k.Children {
		fmt.Println(strPre, " |", string(k))
		v.Print(x + 1)
	}
}

type DeaFilter struct {
	root *KeyWord
	mu   sync.RWMutex
}

// AddWord 添加敏感词（线程安全）
func (df *DeaFilter) AddWord(str string) {
	df.mu.Lock()
	defer df.mu.Unlock()

	if df.root == nil {
		df.root = &KeyWord{
			Children: nil,
		}
	}

	df.root.SplitString(str)
}

// Filter 过滤敏感词（线程安全）
func (df *DeaFilter) Filter(content string) string {
	df.mu.RLock()
	defer df.mu.RUnlock()

	if df.root == nil {
		return content
	}

	chars := []rune(content)
	beginKey := -1
	k := df.root
	for i, r := range chars {
		k = k.Find(r)
		if k == nil {
			k = df.root
			beginKey = -1
			k = k.Find(r)
			if k != nil {
				beginKey = i
			} else {
				k = df.root
			}
			continue
		}

		if beginKey == -1 {
			beginKey = i
		}

		if k.End {
			for index := beginKey; index <= i; index++ {
				chars[index] = '*'
			}

			if k.Children == nil {
				beginKey = -1
				k = df.root
			}

			continue
		}
	}

	return string(chars)
}

func (df *DeaFilter) Print() {
	df.mu.RLock()
	defer df.mu.RUnlock()
	if df.root != nil {
		df.root.Print(0)
	}
}

func NewFilter() *DeaFilter {
	return &DeaFilter{
		root: &KeyWord{
			Children: nil,
		},
	}
}

var (
	DefaultFilter *DeaFilter
	defaultMu     sync.RWMutex
)

func InitDefaultFilter() {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	DefaultFilter = NewFilter()
}

func ensureDefaultFilter() {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if DefaultFilter == nil {
		DefaultFilter = NewFilter()
	}
}

func AddWord(str string) {
	ensureDefaultFilter()
	DefaultFilter.AddWord(str)
}

func ParseFromFile(filename string) error {
	ensureDefaultFilter()

	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	bufReader := bufio.NewReader(fp)

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
		} else {
			DefaultFilter.AddWord(string(line))
		}
	}
	return nil
}

func PrintDefault() {
	ensureDefaultFilter()
	DefaultFilter.Print()
}

func Filter(content string) string {
	ensureDefaultFilter()
	return DefaultFilter.Filter(content)
}
