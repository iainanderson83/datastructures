package linkedlist

type linkedListNode struct {
	Val  interface{}
	Next *linkedListNode
}

type linkedList struct {
	root *linkedListNode
	last *linkedListNode
	len  int
}

func newLinkedList() *linkedList {
	return &linkedList{}
}

func (l *linkedList) Append(v interface{}) {
	if v == nil {
		return
	}

	node := &linkedListNode{Val: v}

	if l.len == 0 {
		l.root = node
		l.last = node
	} else {
		last := l.last
		last.Next = node
		l.last = node
	}

	l.len++
}

func (l *linkedList) Delete(v interface{}) {
	if v == nil || l.len == 0 {
		return
	}

	var prev *linkedListNode
	node := l.root

	for node.Val != v {
		if node.Next == nil {
			return
		}

		prev = node
		node = node.Next
	}

	prev.Next = node.Next
	l.len--
}

func (l *linkedList) Insert(v interface{}, less func(i, j interface{}) bool) {
	if less == nil {
		l.Append(v)
		l.len++
		return
	}

	node := &linkedListNode{Val: v}
	if l.len == 0 {
		l.root = node
		l.last = node
	} else {
		var prev *linkedListNode
		curr := l.root

		for less(curr.Val, node.Val) {
			prev = curr
			curr = curr.Next
		}

		prev.Next = node
		node.Next = curr
	}

	l.len++
}

func (l *linkedList) Front() *linkedListNode {
	return l.root
}

func (l *linkedList) Back() *linkedListNode {
	return l.last
}
