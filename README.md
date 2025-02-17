# BTree Implementation In Go

A B-tree is a self-balancing tree data structure that maintains sorted data and allows searches, sequential access, insertions, and deletions in logarithmic time.

![Logarithmic Time Example](assets/logarithmic_time_complexity.jpg)

## Insert Node

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
	key []byte
	val []byte
}

type Node struct {
	items       [maxItems]*Item
	children    [maxChildren]*Node
	numItems    int
	numChildren int
}

type BTree struct {
	root *Node
}
```

- Item is our smallest unit of data, it will store whatever you want
- Node is responsible to store our items
- BTree is de Data Structure that will arrange all nodes

![BTree Definition](assets/examples_insert/btree_definition.png)

## Deleting Node

