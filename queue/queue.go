package queue

import (
	"container/list"

	"github.com/iainanderson83/datastructures/stack"
)

// Queue is an interface for a classic queue datastructure
type Queue interface {
	Enqueue(interface{})
	Dequeue() interface{}
	Len() int
}

type ArrQueue struct {
	s []interface{}
}

func (a *ArrQueue) Enqueue(v interface{}) {
	if v == nil {
		return
	}

	a.s = append(a.s, v)
}

func (a *ArrQueue) Dequeue() interface{} {
	if len(a.s) == 0 {
		return nil
	}

	v := a.s[0]
	a.s[0] = nil // GC
	a.s = a.s[1:]
	return v
}

func (a *ArrQueue) Len() int { return len(a.s) }

type ListQueue struct {
	l *list.List
}

func NewListQueue() *ListQueue {
	return &ListQueue{list.New()}
}

func (l *ListQueue) Enqueue(v interface{}) {
	l.l.PushBack(v)
}

func (l *ListQueue) Dequeue() interface{} {
	if l.l.Len() == 0 {
		return nil
	}
	e := l.l.Front()
	return l.l.Remove(e)
}

func (l *ListQueue) Len() int { return l.Len() }

type StackQueue struct {
	s1    stack.Stack
	s2    stack.Stack
	count int
}

func (s *StackQueue) Enqueue(v interface{}) {
	if v == nil {
		return
	}

	s.count++
	s.s1.Push(v)
}

func (s *StackQueue) Dequeue() interface{} {
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

	s.count--
	return s.s2.Pop()
}

func (s *StackQueue) Len() int { return s.count }
