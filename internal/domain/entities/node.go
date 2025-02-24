package entities

import "bytes"

const (
	degree      = 2
	maxChildren = 2 * degree
	maxItems    = maxChildren - 1
	minItems    = degree - 1
)

type Node struct {
	Items       [maxItems]*Item
	Children    [maxChildren]*Node
	NumItems    int
	NumChildren int
}

func (n *Node) InsertItemAt(pos int, i *Item) {
	if pos < n.NumItems {
		copy(n.Items[pos+1:n.NumItems+1], n.Items[pos:n.NumItems])
	}
	n.Items[pos] = i
	n.NumItems++
}

func (n *Node) search(key []byte) (int, bool) {
	low, high := 0, n.NumItems
	var mid int
	for low < high {
		mid = (low + high) / 2
		cmp := bytes.Compare(key, n.Items[mid].Key)
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

func (n *Node) insertChildAt(pos int, c *Node) {
	if pos < n.NumChildren {
		copy(n.Children[pos+1:n.NumChildren+1], n.Children[:][pos:n.NumChildren])
	}
	n.Children[pos] = c
	n.NumChildren++
}

func (n *Node) split() (*Item, *Node) {
	mid := minItems
	midItem := n.Items[mid]

	newNode := &Node{}
	copy(newNode.Items[:], n.Items[mid+1:])
	newNode.NumItems = minItems

	if !n.IsLeaf() {
		copy(newNode.Children[:], n.Children[mid+1:])
		newNode.NumChildren = minItems + 1
	}

	for i, l := mid, n.NumItems; i < l; i++ {
		n.Items[i] = nil
		n.NumItems--

		if !n.IsLeaf() {
			n.Children[i+1] = nil
			n.NumChildren--
		}
	}

	return midItem, newNode
}

func (n *Node) insert(item *Item) bool {
	pos, found := n.search(item.Key)
	if found {
		n.Items[pos] = item
		return false
	}

	if n.IsLeaf() {
		n.InsertItemAt(pos, item)
		return true
	}

	if n.Children[pos].NumItems >= maxItems {
		midItem, newNode := n.Children[pos].split()
		n.InsertItemAt(pos, midItem)
		n.insertChildAt(pos+1, newNode)

		switch cmp := bytes.Compare(item.Key, n.Items[pos].Key); {
		case cmp < 0:
			// The key we are looking for is still smaller than the key of the middle Item that we took from the child,
			// so we can continue following the same direction.
		case cmp > 0:
			// The middle item that we took from the child has a key that is smaller than the one we are looking for,
			// so we need to change our direction.
			pos++
		default:
			// The middle item that we took from the child is the item we are searching for, so just update its value.
			n.Items[pos] = item
			return true
		}

	}
	return n.Children[pos].insert(item)
}

func (n *Node) removeItemAt(pos int) *Item {
	removedItem := n.Items[pos]
	n.Items[pos] = nil

	if lastPos := n.NumItems - 1; pos < lastPos {
		copy(n.Items[pos:lastPos], n.Items[pos+1:lastPos+1])
		n.Items[lastPos] = nil
	}
	n.NumItems--

	return removedItem
}

func (n *Node) removeChildAt(pos int) *Node {
	removedChild := n.Children[pos]
	n.Children[pos] = nil
	if lastPos := n.NumChildren - 1; pos < lastPos {
		copy(n.Children[pos:lastPos], n.Children[pos+1:lastPos+1])
		n.Children[lastPos] = nil
	}
	n.NumChildren--
	return removedChild
}

func (n *Node) fillChildAt(pos int) {
	switch {
	case pos > 0 && n.Children[pos-1].NumItems > minItems:
		left, right := n.Children[pos-1], n.Children[pos]
		copy(right.Items[1:right.NumItems+1], right.Items[:right.NumItems])
		right.Items[0] = n.Items[pos-1]
		right.NumItems++
		if !right.IsLeaf() {
			right.insertChildAt(0, left.removeChildAt(left.NumChildren-1))
		}
		n.Items[pos-1] = left.removeItemAt(left.NumItems - 1)
	case pos < n.NumChildren-1 && n.Children[pos+1].NumItems > minItems:
		left, right := n.Children[pos], n.Children[pos+1]
		left.Items[left.NumItems] = n.Items[pos]
		left.NumItems++
		if !left.IsLeaf() {
			left.insertChildAt(left.NumChildren, right.removeChildAt(0))
		}
		n.Items[pos] = right.removeItemAt(0)
	default:
		if pos >= n.NumItems {
			pos = n.NumItems - 1
		}
		left, right := n.Children[pos], n.Children[pos+1]
		left.Items[left.NumItems] = n.removeItemAt(pos)
		left.NumItems++
		copy(left.Items[left.NumItems:], right.Items[:right.NumItems])
		left.NumItems += right.NumItems
		if !left.IsLeaf() {
			copy(left.Children[left.NumChildren:], right.Children[:right.NumChildren])
			left.NumChildren += right.NumChildren
		}
		n.removeChildAt(pos + 1)
		right = nil
	}
}

func (n *Node) delete(key []byte, isSeekingSuccessor bool) *Item {
	pos, found := n.search(key)
	var next *Node
	if found {
		if n.IsLeaf() {
			return n.removeItemAt(pos)
		}
		next, isSeekingSuccessor = n.Children[pos+1], true
	} else {
		next = n.Children[pos]
	}

	if n.IsLeaf() && isSeekingSuccessor {
		return n.removeItemAt(0)
	}

	if next == nil {
		return nil
	}

	deletedItem := next.delete(key, isSeekingSuccessor)

	if found && isSeekingSuccessor {
		n.Items[pos] = deletedItem
	}

	if next.NumItems < minItems {
		if found && isSeekingSuccessor {
			n.fillChildAt(pos + 1)
		} else {
			n.fillChildAt(pos)
		}
	}
	return deletedItem
}

func (n *Node) IsLeaf() bool {
	return n.NumChildren == 0
}
