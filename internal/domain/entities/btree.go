package entities

import "errors"

type BTree struct {
	root *Node
}

func (t *BTree) Find(key []byte) ([]byte, error) {
	for next := t.root; next != nil; {
		pos, found := next.search(key)

		if found {
			return next.Items[pos].Val, nil
		}

		next = next.Children[pos]
	}
	return nil, errors.New("key not found")
}

func (t *BTree) splitRoot() {
	newRoot := &Node{}
	midItem, newNode := t.root.split()
	newRoot.InsertItemAt(0, midItem)
	newRoot.insertChildAt(0, t.root)
	newRoot.insertChildAt(1, newNode)
	t.root = newRoot
}

func (t *BTree) Insert(Key, Val []byte) bool {
	i := &Item{Key, Val}

	if t.root == nil {
		t.root = &Node{}
	}

	if t.root.NumItems >= maxItems {
		t.splitRoot()
	}

	hasSuccess := t.root.insert(i)
	return hasSuccess
}

func (t *BTree) Delete(key []byte) bool {
	if t.root == nil {
		return false
	}
	deletedItem := t.root.delete(key, false)

	if t.root.NumItems == 0 {
		if t.root.IsLeaf() {
			t.root = nil
		} else {
			t.root = t.root.Children[0]
		}
	}

	if deletedItem != nil {
		return true
	}
	return false
}
