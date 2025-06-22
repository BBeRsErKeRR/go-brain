package logger

import (
	"errors"
	"sync"
	"time"
)

var ErrEmptyStack = errors.New("empty stack")

type Stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []time.Time
}

func NewStack() *Stack {
	return &Stack{sync.Mutex{}, make([]time.Time, 0)}
}

func (s *Stack) Len() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.s)
}

func (s *Stack) Push(v time.Time) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *Stack) Pop() (time.Time, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return time.Time{}, ErrEmptyStack
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}
