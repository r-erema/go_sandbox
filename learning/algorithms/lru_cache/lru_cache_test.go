package lrucache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	t.Parallel()

	cache := Constructor(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	assert.Equal(t, 1, cache.Get(1))
	cache.Put(3, 3)
	assert.Equal(t, -1, cache.Get(2))
	cache.Put(4, 4)
	assert.Equal(t, -1, cache.Get(1))
	assert.Equal(t, 3, cache.Get(3))
	assert.Equal(t, 4, cache.Get(4))
}

type node struct {
	key, value int
	next, prev *node
}
type LRUCache struct {
	capacity            int
	nodesMap            map[int]*node
	lastNode, firstNode *node
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		nodesMap: make(map[int]*node, capacity),
	}
}

// Time O(1), since get the node using the map and renew needed links with constant time
// Space O(n), since we use a map.
func (c *LRUCache) Get(key int) int {
	curr, ok := c.nodesMap[key]
	if !ok {
		return -1
	}

	c.moveNodeToLast(curr)

	return curr.value
}

func (c *LRUCache) moveNodeToLast(curr *node) {
	if c.lastNode == curr {
		return
	}

	var head, tail *node
	if curr.next != nil {
		head = curr.next
		head.prev = nil
		curr.next = nil

		if c.firstNode == curr {
			c.firstNode = head
		}
	}

	if curr.prev != nil {
		tail = curr.prev
		tail.next = head

		if head != nil {
			head.prev = tail
		}
	}

	curr.prev = c.lastNode
	c.lastNode.next = curr
	c.lastNode = curr
}

// Time O(1), since put the node to the end of the list and renew needed links with constant time
// Space O(n), since we put into a map.
func (c *LRUCache) Put(key, value int) {
	if curr, ok := c.nodesMap[key]; ok {
		curr.value = value
		c.moveNodeToLast(curr)

		return
	}

	tail := c.lastNode
	c.lastNode = &node{key: key, value: value}

	if tail != nil {
		c.lastNode.prev = tail
		tail.next = c.lastNode
	}

	if len(c.nodesMap) == c.capacity {
		delete(c.nodesMap, c.firstNode.key)
		head := c.firstNode.next
		head.prev = nil
		c.firstNode = head
	}

	c.nodesMap[key] = c.lastNode
	if len(c.nodesMap) == 1 {
		c.firstNode = c.lastNode
	}
}
