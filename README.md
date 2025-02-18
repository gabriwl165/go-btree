# BTree Implementation In Go

A B-tree is a self-balancing tree data structure that maintains sorted data and allows searches, sequential access, insertions, and deletions in logarithmic time.

![Logarithmic Time Example](assets/logarithmic_time_complexity.jpg)


### Introduction

B-Tree has some initials definiitions, such as:
- The minimun amout of items that a node can have is half of the maximum, except for the root node

Let's assume the maximum items in node is 4 (four)

![Creating the root node](assets/examples_insert/image3.png)

So, what we need to do now, that there's 5 items within this node? We should split them up! We basically get the mid item, and make him the current root node, all from the left is going to point to him, and the same from the right.

![Splitting the root node](assets/examples_insert/image4.png)

So, with this introduction, we should be able to start our code!

```go
const (
	degree      = 5
	maxChildren = 2 * degree
	maxItems    = maxChildren - 1
	minItems    = degree - 1
)

type Item struct {
	Key []byte
	Val []byte
}

type Node struct {
	Items       [maxItems]*Item
	Children    [maxChildren]*Node
	NumItems    int
	NumChildren int
}

type BTree struct {
	root *Node
}
```

- Item is our smallest unit of data, it will store whatever you want
- Node is responsible to store our items
- BTree is de Data Structure that will arrange all nodes

![BTree Definition](assets/examples_insert/btree_definition.png)

## Searching Node

Before we dive into how to insert and delete node, there's one step earlier, that is implement how we're going to iterate over the btree.

```go
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
```

This method performs a binary search on a sorted list of items inside a Node. It finds the position of a given key or determines where it should be inserted.

## Inserting Node

Let's start by create the implementation for our btree, that is kind simple, for now.

```go
func (b *BTree) Insert(key []byte, val []byte) bool {
	item := &Item{
		Key: key,
		Val: val,
	}

    // If there's no node, so it's the first insertion
	if b.root == nil {
		b.root = &Node{}
	}

	hasSuccess := b.root.insert(item)
	return hasSuccess
}

func (n *Node) insert(item *Item) bool {
	// Search into the files to see if found the item with the key provided
	pos, found := n.search(item.Key)

	// Found an item with the same key, so we must update this position with the new item
	if found {
		n.Items[pos] = item
		return false
	}
	return true
}
```

For the `Insert` BTree method is just instantiating the `Item` Object and Node `insert` method is appending into the root node. But this is just the beggning, what about the current Node is already fullfill? we need to split them up!

Let's begin with the current node has reached his limit:
```go
func (n *Node) insert(item *Item) bool {
	// Search into the files to see if found the item with the key provided
	pos, found := n.search(item.Key)

	// Found an item with the same key, so we must update this position with the new item
	if found {
		n.Items[pos] = item
		return false
	}

	// The node is already full
	if n.Children[pos].NumItems >= maxItems {
		// ....
	}

	return true
}
```

Right here we must implement our logic that will separate our node!
```go
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
```

- `mid` and `midItem` is getting the mid item within the current node.
- after `&Node{}` we `copy` half right of the item to inside the new node.
- If the node is not a leaf it copies the right half of the child pointers to newNode.Children

![BTree Definition](assets/btree/image2.png)

## Deleting Node

