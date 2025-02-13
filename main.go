package main

import (
	"bytes"
	"errors"
)

const (
	degree      = 5
	maxChildren = 2 * degree
	maxItems    = maxChildren - 1
	minItems    = degree - 1
)

type item struct {
	key []byte
	val []byte
}

type node struct {
	items       [maxItems]*item
	children    [maxChildren]*node
	numItems    int
	numChildren int
}

func (n *node) isLeaf() bool {
	return n.numChildren == 0
}

type BTree struct {
	root *node
}

func NewBTree() *BTree {
	return &BTree{}
}

func (n *node) search(key []byte) (int, bool) {
	low, high := 0, n.numItems
	var mid int
	for low < high {
		mid = (low + high) / 2
		cmp := bytes.Compare(key, n.items[mid].key)
		switch {
		case cmp > 0:
			low = mid + 1
		case cmp < 0:
			high = mid
		case cmp == 0:
			return mid, true
		}
	}
	return low, false
}

func (t *BTree) Find(key []byte) ([]byte, error) {
	for next := t.root; next != nil; {
		pos, found := next.search(key)

		if found {
			return next.items[pos].val, nil
		}

		next = next.children[pos]
	}
	return nil, errors.New("key not found")
}

func (n *node) insertItemAt(pos int, i *item) {
	if pos < n.numItems {
		copy(n.items[pos+1:n.numItems+1], n.items[pos:n.numItems])
	}
	n.items[pos] = i
	n.numItems++
}

func (n *node) insertChildAt(pos int, c *node) {
	if pos < n.numChildren {
		copy(n.children[pos+1:n.numChildren+1], n.children[:][pos:n.numChildren])
	}
	n.children[pos] = c
	n.numChildren++
}

func (n *node) split() (*item, *node) {
	mid := minItems
	midItem := n.items[mid]

	newNode := &node{}
	copy(newNode.items[:], n.items[mid+1:])
	newNode.numItems = minItems

	if !n.isLeaf() {
		copy(newNode.children[:], n.children[mid+1:])
		newNode.numChildren = minItems + 1
	}

	for i, l := mid, n.numItems; i < l; i++ {
		n.items[i] = nil
		n.numItems--

		if !n.isLeaf() {
			n.children[i+1] = nil
			n.numChildren--
		}
	}

	return midItem, newNode
}

func (n *node) insert(item *item) bool {
	pos, found := n.search(item.key)
	if found {
		n.items[pos] = item
		return false
	}

	if n.isLeaf() {
		n.insertItemAt(pos, item)
		return true
	}

	if n.children[pos].numItems >= maxItems {
		midItem, newNode := n.children[pos].split()
		n.insertItemAt(pos, midItem)
		n.insertChildAt(pos+1, newNode)

		switch cmp := bytes.Compare(item.key, n.items[pos].key); {
		case cmp < 0:
			// The key we are looking for is still smaller than the key of the middle item that we took from the child,
			// so we can continue following the same direction.
		case cmp > 0:
			pos++
		default:
			n.items[pos] = item
			return true
		}

	}
	return n.children[pos].insert(item)
}

func (t *BTree) splitRoot() {
	newRoot := &node{}
	midItem, newNode := t.root.split()
	newRoot.insertItemAt(0, midItem)
	newRoot.insertChildAt(0, t.root)
	newRoot.insertChildAt(1, newNode)
	t.root = newRoot
}

func (t *BTree) Insert(key, val []byte) {
	i := &item{key, val}

	if t.root == nil {
		t.root = &node{}
	}

	if t.root.numItems >= maxItems {
		t.splitRoot()
	}

	t.root.insert(i)
}

func (n *node) removeItemAt(pos int) *item {
	removedItem := n.items[pos]
	n.items[pos] = nil

	if lastPos := n.numItems - 1; pos < lastPos {
		copy(n.items[pos:lastPos], n.items[pos+1:lastPos+1])
		n.items[lastPos] = nil
	}
	n.numItems--

	return removedItem
}

func (n *node) removeChildAt(pos int) *node {
	removedChild := n.children[pos]
	n.children[pos] = nil
	if lastPos := n.numChildren - 1; pos < lastPos {
		copy(n.children[pos:lastPos], n.children[pos+1:lastPos+1])
		n.children[lastPos] = nil
	}
	n.numChildren--
	return removedChild
}

func (n *node) fillChildAt(pos int) {
	switch {
	case pos > 0 && n.children[pos-1].numItems > minItems:
		left, right := n.children[pos-1], n.children[pos]
		copy(right.items[1:right.numItems+1], right.items[:right.numItems])
		right.items[0] = n.items[pos-1]
		right.numItems++
		if !right.isLeaf() {
			right.insertChildAt(0, left.removeChildAt(left.numChildren-1))
		}
		n.items[pos-1] = left.removeItemAt(left.numItems - 1)
	case pos < n.numChildren-1 && n.children[pos+1].numItems > minItems:
		left, right := n.children[pos], n.children[pos+1]
		left.items[left.numItems] = n.items[pos]
		left.numItems++
		if !left.isLeaf() {
			left.insertChildAt(left.numChildren, right.removeChildAt(0))
		}
		n.items[pos] = right.removeItemAt(0)
	default:
		if pos >= n.numItems {
			pos = n.numItems - 1
		}
		left, right := n.children[pos], n.children[pos+1]
		left.items[left.numItems] = n.removeItemAt(pos)
		left.numItems++
		copy(left.items[left.numItems:], right.items[:right.numItems])
		left.numItems += right.numItems
		if !left.isLeaf() {
			copy(left.children[left.numChildren:], right.children[:right.numChildren])
			left.numChildren += right.numChildren
		}
		n.removeChildAt(pos + 1)
		right = nil
	}
}

func (n *node) delete(key []byte, isSeekingSuccessor bool) *item {
	pos, found := n.search(key)
	var next *node
	if found {
		if n.isLeaf() {
			return n.removeItemAt(pos)
		}
		next, isSeekingSuccessor = n.children[pos+1], true
	} else {
		next = n.children[pos]
	}

	if n.isLeaf() && isSeekingSuccessor {
		return n.removeItemAt(0)
	}

	if next == nil {
		return nil
	}

	deletedItem := next.delete(key, isSeekingSuccessor)

	if found && isSeekingSuccessor {
		n.items[pos] = deletedItem
	}

	if next.numItems < minItems {
		if found && isSeekingSuccessor {
			n.fillChildAt(pos + 1)
		} else {
			n.fillChildAt(pos)
		}
	}
	return deletedItem
}

func (t *BTree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}
	deletedItem := t.root.delete(key, false)

	if t.root.numItems == 0 {
		if t.root.isLeaf() {
			t.root = nil
		} else {
			t.root = t.root.children[[0]]
		}
	}

	if deletedItem != nil {
		return true
	}
	return false
}
