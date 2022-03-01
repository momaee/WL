package stack

import "sync"

// Default data type for stack
type Item interface{}

// Stack data structure
type Stack struct {
	items []Item
	lock  sync.RWMutex
}

func (s *Stack) Len() int {
	return len(s.items)
}

// Push the new data to the top of the stack
func (s *Stack) Push(t Item) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

// Pop the last item from top of the stack
func (s *Stack) Pop() Item {
	s.lock.Lock()
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	s.lock.Unlock()
	return top

}

// Create new instance of stack
func (s *Stack) NewStack() *Stack {
	s.items = []Item{}
	return s
}
