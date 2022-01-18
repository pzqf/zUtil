package zStack

import "errors"

//const arraySize = 10

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
