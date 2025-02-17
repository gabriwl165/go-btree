package entities_test

import (
	"testing"

	"github.com/gabriwl165/go-btree/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNodeIsLeaf(t *testing.T) {
	newNode := &entities.Node{
		NumChildren: 0,
	}
	assert.Equal(t, true, newNode.IsLeaf(), "they should be a leaf")
}

func TestNodeInsertItemAt(t *testing.T) {
	newNode := &entities.Node{
		NumChildren: 0,
		NumItems:    0,
	}
	newItem := &entities.Item{
		Key: []byte("key"),
		Val: []byte("val"),
	}
	newNode.InsertItemAt(0, newItem)
	assert.NotEqual(t, nil, newNode.Items[0])
}
