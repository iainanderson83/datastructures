package linkedlist

import "testing"

func TestLinkedList(t *testing.T) {
	l := newLinkedList()

	l.Append(1)
	l.Append(2)
	l.Append(3)

	front := l.Front()
	if front.Val != 1 {
		t.Fatalf("expected %d, got %d", 1, front.Val)
	}

	for front.Next != nil {
		front = front.Next
	}

	if front.Val != 3 {
		t.Fatalf("expected %d, got %d", 3, front.Val)
	}
}
