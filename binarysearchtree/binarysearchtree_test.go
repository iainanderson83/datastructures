package binarysearchtree

import (
	"fmt"
	"testing"
	"time"
)

var bst *Node

func fillTree(bst *Node) {
	bst.Insert(4, "4")
	bst.Insert(10, "10")
	bst.Insert(2, "2")
	bst.Insert(6, "6")
	bst.Insert(1, "1")
	bst.Insert(3, "3")
	bst.Insert(5, "5")
	bst.Insert(7, "7")
	bst.Insert(9, "9")
	bst.Insert(11, "11")
}

// isSameSlice returns true if the 2 slices are identical
func isSameSlice(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestInOrderTraverse(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	var result []string
	bst.InOrderTraverse(func(key int, value interface{}) {
		result = append(result, fmt.Sprintf("%s", value))
	})
	if !isSameSlice(result, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}) {
		t.Errorf("Traversal order incorrect, got %v", result)
	}
}

func TestPreOrderTraverse(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	var result []string
	bst.PreOrderTraverse(func(key int, value interface{}) {
		result = append(result, fmt.Sprintf("%s", value))
	})
	if !isSameSlice(result, []string{"8", "4", "2", "1", "3", "6", "5", "7", "10", "9", "11"}) {
		t.Errorf("Traversal order incorrect, got %v instead of %v", result, []string{"8", "4", "2", "1", "3", "6", "5", "7", "10", "9", "11"})
	}
}

func TestPostOrderTraverse(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	var result []string
	bst.PostOrderTraverse(func(key int, value interface{}) {
		result = append(result, fmt.Sprintf("%s", value))
	})
	if !isSameSlice(result, []string{"1", "3", "2", "5", "7", "6", "4", "9", "11", "10", "8"}) {
		t.Errorf("Traversal order incorrect, got %v instead of %v", result, []string{"1", "3", "2", "5", "7", "6", "4", "9", "11", "10", "8"})
	}
}

func TestMin(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	if fmt.Sprintf("%s", bst.Min().Value) != "1" {
		t.Errorf("min should be 1")
	}
}

func TestMax(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	if fmt.Sprintf("%s", bst.Max().Value) != "11" {
		t.Errorf("max should be 11")
	}
}

func TestSearch(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	if !bst.Search(1) || !bst.Search(8) || !bst.Search(11) {
		t.Errorf("search not working")
	}
}

func TestRemove(t *testing.T) {
	bst = &Node{Key: 8, Value: "8"}
	fillTree(bst)

	bst.Remove(1)

	if fmt.Sprintf("%s", bst.Min().Value) != "2" {
		t.Errorf("min should be 2, got %s", bst.Min().Value)
	}
}

func TestDates(t *testing.T) {
	ranges := []struct {
		BeginDate time.Time
		EndDate   time.Time
		Value     float64
	}{
		{
			time.Date(2009, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2009, time.June, 1, 0, 0, 0, 0, time.UTC),
			60.5,
		},
		{
			time.Date(2009, time.June, 1, 0, 0, 0, 1, time.UTC),
			time.Date(2009, time.December, 1, 0, 0, 0, 0, time.UTC),
			65.5,
		},
		{
			time.Date(2009, time.December, 1, 0, 0, 0, 1, time.UTC),
			time.Date(2010, time.June, 1, 0, 0, 0, 0, time.UTC),
			70.5,
		},
		{
			time.Date(2010, time.June, 1, 0, 0, 0, 1, time.UTC),
			time.Date(2010, time.December, 1, 0, 0, 0, 0, time.UTC),
			75.5,
		},
		{
			time.Date(2010, time.December, 1, 0, 0, 0, 1, time.UTC),
			time.Date(2011, time.June, 1, 0, 0, 0, 0, time.UTC),
			80.5,
		},
	}

	middle := len(ranges) / 2
	bst := &Node{Key: int(ranges[middle].BeginDate.Unix()), Value: ranges[middle].Value}

	for i := range ranges {
		if i == middle {
			continue
		}
		bst.Insert(int(ranges[i].BeginDate.Unix()), ranges[i].Value)
	}

	out := bst.Nearest(int(time.Date(2009, time.June, 2, 0, 0, 0, 1, time.UTC).Unix()))
	if out.Value != 65.5 {
		t.Fatalf("expected %f, got %f", 65.5, out.Value)
	}
}
