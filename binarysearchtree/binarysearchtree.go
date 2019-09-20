package binarysearchtree

import (
	"errors"
	"fmt"
)

// Node is a leaf and a tree.
type Node struct {
	Left  *Node
	Right *Node
	Key   int
	Value interface{}
}

// Insert adds the given key and value to the tree.
func (n *Node) Insert(key int, value interface{}) error {
	if n == nil {
		return errors.New("Cannot insert a value into a nil tree")
	}
	switch {
	case key < n.Key:
		if n.Left == nil {
			n.Left = &Node{Key: key, Value: value}
			return nil
		}
		return n.Left.Insert(key, value)
	case key > n.Key:
		if n.Right == nil {
			n.Right = &Node{Key: key, Value: value}
			return nil
		}
		return n.Right.Insert(key, value)
	}
	n.Value = value
	return nil
}

// Visitor is a function that is called during traversal.
type Visitor func(key int, value interface{})

// InOrderTraverse calls Visitor on the left node, the current node,
// and then the right node.
func (n *Node) InOrderTraverse(v Visitor) {
	if n == nil {
		return
	}

	if n.Left != nil {
		n.Left.InOrderTraverse(v)
	}
	v(n.Key, n.Value)
	if n.Right != nil {
		n.Right.InOrderTraverse(v)
	}
}

// PreOrderTraverse calls Visitor for the current node, the left node,
// and the right node.
func (n *Node) PreOrderTraverse(v Visitor) {
	if n == nil {
		return
	}

	v(n.Key, n.Value)
	if n.Left != nil {
		n.Left.PreOrderTraverse(v)
	}
	if n.Right != nil {
		n.Right.PreOrderTraverse(v)
	}
}

// PostOrderTraverse calls Visitor for the left node, the right node,
// and the current node.
func (n *Node) PostOrderTraverse(v Visitor) {
	if n == nil {
		return
	}

	if n.Left != nil {
		n.Left.PostOrderTraverse(v)
	}
	if n.Right != nil {
		n.Right.PostOrderTraverse(v)
	}
	v(n.Key, n.Value)
}

// Min returns the node associated with the min key in the tree.
func (n *Node) Min() *Node {
	if n == nil {
		return nil
	}

	curr := n
	for {
		if curr.Left == nil {
			return curr
		}

		curr = curr.Left
	}
}

// Max returns the node associated with the max key in the tree.
func (n *Node) Max() *Node {
	if n == nil {
		return nil
	}

	curr := n
	for {
		if curr.Right == nil {
			return curr
		}

		curr = curr.Right
	}
}

// Search returns true if the given key is found within the tree.
func (n *Node) Search(key int) bool {
	if n == nil {
		return false
	}

	if key < n.Key {
		if n.Left == nil {
			return false
		}
		return n.Left.Search(key)
	}

	if key > n.Key {
		if n.Right == nil {
			return false
		}
		return n.Right.Search(key)
	}

	return true
}

// Exact retrieves a node from the tree for the specified key, or nil.
func (n *Node) Exact(key int) *Node {
	if n == nil {
		return nil
	}

	if key < n.Key {
		if n.Left == nil {
			return nil
		}
		return n.Left.Exact(key)
	}

	if key > n.Key {
		if n.Right == nil {
			return nil
		}
		return n.Right.Exact(key)
	}

	return n
}

// Nearest retrieves the nearest node from the tree for the specified key, or nil.
func (n *Node) Nearest(key int) *Node {
	if n == nil {
		return nil
	}

	if key < n.Key {
		if n.Left != nil {
			return n.Left.Nearest(key)
		}
	}

	if key > n.Key {
		if n.Right != nil {
			return n.Right.Nearest(key)
		}
	}

	return n
}

// Remove removes the node associated with the given key and returns it.
func (n *Node) Remove(key int) *Node {
	return remove(n, key)
}

func remove(n *Node, key int) *Node {
	if n == nil {
		return nil
	}

	if key < n.Key {
		n.Left = remove(n.Left, key)
		return n
	}

	if key > n.Key {
		n.Right = remove(n.Right, key)
		return n
	}

	if n.Left == nil && n.Right == nil {
		n = nil
		return nil
	}

	if n.Left == nil {
		n = n.Right
		return n
	}

	if n.Right == nil {
		n = n.Left
		return n
	}

	smallestRight := n.Right
	for {
		if smallestRight != nil && smallestRight.Left != nil {
			smallestRight = smallestRight.Left
		} else {
			break
		}
	}

	n.Key, n.Value = smallestRight.Key, smallestRight.Value
	n.Right = remove(n.Right, n.Key)
	return n
}

// String prints a visual representation of the tree
func (n *Node) String() {
	fmt.Println("Stringify")
	stringify(n, 0)
}

// internal recursive function to print a tree
func stringify(n *Node, level int) {
	if n == nil {
		return
	}
	format := ""
	for i := 0; i < level; i++ {
		format += "\t"
	}
	level++
	if n.Left != nil {
		stringify(n.Left, level)
	}
	fmt.Printf(format+"%d\n", n.Key)
	if n.Right != nil {
		stringify(n.Right, level)
	}
}
