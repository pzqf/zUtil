package zStack

import (
	"errors"
	"sync"
)

type Stack struct {
	top  int
	data []interface{}
}

func New(arraySize int) Stack {
	if arraySize <= 0 {
		arraySize = 10
	}
	return Stack{
		top:  0,
		data: make([]interface{}, arraySize),
	}
}

func (s *Stack) Push(value interface{}) error {
	if s.top == len(s.data) {
		return errors.New("stack full")
	}
	s.data[s.top] = value
	s.top += 1
	return nil
}

func (s *Stack) Pop() (interface{}, error) {
	if s.top == 0 {
		return 0, errors.New("stack empty")
	}
	value := s.data[s.top-1]
	s.top -= 1
	return value, nil
}

func (s *Stack) Peek() (interface{}, error) {
	if s.top == 0 {
		return nil, errors.New("stack empty")
	}
	return s.data[s.top-1], nil
}

func (s *Stack) Get() []interface{} {
	return s.data[:s.top]
}

func (s *Stack) IsEmpty() bool {
	return s.top == 0
}

func (s *Stack) Empty() {
	s.top = 0
}

// ConcurrentStack 线程安全的栈
type ConcurrentStack struct {
	stack Stack
	mu    sync.RWMutex
}

// NewConcurrent 创建线程安全的栈
func NewConcurrent(arraySize int) *ConcurrentStack {
	return &ConcurrentStack{
		stack: New(arraySize),
	}
}

func (cs *ConcurrentStack) Push(value interface{}) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.stack.Push(value)
}

func (cs *ConcurrentStack) Pop() (interface{}, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.stack.Pop()
}

func (cs *ConcurrentStack) Peek() (interface{}, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.stack.Peek()
}

func (cs *ConcurrentStack) Get() []interface{} {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.stack.Get()
}

func (cs *ConcurrentStack) IsEmpty() bool {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.stack.IsEmpty()
}

func (cs *ConcurrentStack) Empty() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.stack.Empty()
}

func (cs *ConcurrentStack) Len() int {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.stack.top
}
