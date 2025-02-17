package examples

import (
	"strconv"
	"testing"

	"github.com/gabriwl165/go-btree/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestInsertOneNode(t *testing.T) {
	btree := entities.BTree{}
	hasSuccess := btree.Insert(
		[]byte("key"),
		[]byte("val"),
	)
	assert.Equal(t, true, hasSuccess)
}

func TestInsertFourNode(t *testing.T) {
	btree := entities.BTree{}
	var hadSuccess [4]bool
	for i := 0; i < 4; i++ {
		key := "key" + strconv.Itoa(i)
		hasSuccess := btree.Insert(
			[]byte(key),
			[]byte("val"),
		)
		hadSuccess[i] = hasSuccess
	}
	for _, success := range hadSuccess {
		assert.Equal(t, true, success)
	}
}
