package queue

import (
	"testing"

	"github.com/iainanderson83/datastructures/stack"
)

func TestQueues(t *testing.T) {
	in := []int{1, 3, 5, 4, 2}

	tests := map[string]Queue{
		"ArrQueue":   &ArrQueue{},
		"ListQueue":  NewListQueue(),
		"StackQueue": &StackQueue{s1: &stack.ArrStack{}, s2: &stack.ArrStack{}},
	}

	for name, q := range tests {
		t.Run(name, func(t *testing.T) {
			for _, i := range in {
				q.Enqueue(i)
			}

			for _, i := range in {
				if q.Dequeue() != i {
					t.Fatal("out of order")
				}
			}
		})
	}
}

func TestQueues2(t *testing.T) {
	in := []int{1, 3, 5, 4, 2}

	tests := map[string]Queue{
		"ArrQueue":   &ArrQueue{},
		"ListQueue":  NewListQueue(),
		"StackQueue": &StackQueue{s1: &stack.ArrStack{}, s2: &stack.ArrStack{}},
	}

	for name, q := range tests {
		t.Run(name, func(t *testing.T) {
			q.Enqueue(in[0])
			q.Enqueue(in[1])
			q.Enqueue(in[2])

			q.Dequeue()
			q.Dequeue()

			q.Enqueue(in[3])
			q.Enqueue(in[4])

			if v := q.Dequeue(); v != 5 {
				t.Fatalf("expected %d, got %d", 5, v)
			}

			if v := q.Dequeue(); v != 4 {
				t.Fatalf("expected %d, got %d", 4, v)
			}
		})
	}
}
