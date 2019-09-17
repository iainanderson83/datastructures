package queue

import (
	"container/list"

	"github.com/iain.anderson83/datastructures/stack"
)

// Queue is an interface for a classic queue datastructure
type Queue interface {
	Enqueue(interface{})
	Dequeue() interface{}
}

type arrQueue struct {
	s []interface{}
}

func (a *arrQueue) Enqueue(v interface{}) {
	if v == nil {
		return
	}

	a.s = append(a.s, v)
}

func (a *arrQueue) Dequeue() interface{} {
	if len(a.s) == 0 {
		return nil
	}

	v := a.s[0]
	a.s[0] = nil // GC
	a.s = a.s[1:]
	return v
}

type listQueue struct {
	l *list.List
}

func newListQueue() *listQueue {
	return &listQueue{list.New()}
}

func (l *listQueue) Enqueue(v interface{}) {
	l.l.PushBack(v)
}

func (l *listQueue) Dequeue() interface{} {
	if l.l.Len() == 0 {
		return nil
	}
	e := l.l.Front()
	return l.l.Remove(e)
}

type stackQueue struct {
	s1 stack.Stack
	s2 stack.Stack
}

func (s *stackQueue) Enqueue(v interface{}) {
	if v == nil {
		return
	}

	s.s1.Push(v)
}

func (s *stackQueue) Dequeue() interface{} {
	if s.s1.Peek() == nil && s.s2.Peek() == nil {
		return nil
	}

	if s.s2.Peek() == nil {
		for {
			v := s.s1.Pop()
			if v == nil {
				break
			}
			s.s2.Push(v)
		}
	}

	return s.s2.Pop()
}
