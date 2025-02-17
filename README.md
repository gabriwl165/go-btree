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
## Deleting Node

