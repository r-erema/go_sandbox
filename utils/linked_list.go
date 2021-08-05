package utils

type LinkedListNode struct {
	value      float32
	prev, next *LinkedListNode
}

func Build(nodes []*LinkedListNode) *LinkedListNode {
	var startNode *LinkedListNode

	iterationThreshold := len(nodes) - 1

	for i := range nodes {
		var currentNode, prevNode, nextNode *LinkedListNode
		currentNode = nodes[i]

		if i != 0 {
			prevNode = nodes[i-1]
			prevNode.next = currentNode
		} else {
			startNode = nodes[i]
		}

		if i < iterationThreshold {
			nextNode = nodes[i+1]
			nextNode.prev = currentNode
		}

		currentNode.next = nextNode
		currentNode.prev = prevNode
	}

	return startNode
}

func (node LinkedListNode) IsHead() bool {
	return node.prev == nil
}

func (node *LinkedListNode) IsTail() bool {
	return node.next == nil
}

func (node *LinkedListNode) Head() *LinkedListNode {
	currentNode := node

	for {
		if currentNode.prev == nil {
			return currentNode
		}

		currentNode = currentNode.prev
	}
}

func (node *LinkedListNode) Tail() *LinkedListNode {
	currentNode := node

	for {
		if currentNode.next == nil {
			return currentNode
		}

		currentNode = currentNode.next
	}
}

func (node *LinkedListNode) Append(nodeToAppend *LinkedListNode) {
	tail := node.Tail()
	tail.next = nodeToAppend
	nodeToAppend.prev = tail
}

func (node *LinkedListNode) Prepend(nodeToPrepend *LinkedListNode) {
	head := node.Head()
	head.prev = nodeToPrepend
	nodeToPrepend.next = head
}

func (node *LinkedListNode) Search(needle float32) *LinkedListNode {
	currentNode := node

	for {
		if currentNode.value == needle {
			return currentNode
		}

		if currentNode.next == nil {
			break
		}

		currentNode = currentNode.next
	}

	if currentNode = node.prev; currentNode != nil {
		for {
			if currentNode.value == needle {
				return currentNode
			}

			if currentNode.prev == nil {
				break
			}

			currentNode = currentNode.prev
		}
	}

	return nil
}

func (node *LinkedListNode) Remove(needle float32) {
	nodeToRemove := node.Search(needle)
	if nodeToRemove == nil {
		return
	}

	prev, next := nodeToRemove.prev, nodeToRemove.next

	if prev != nil {
		prev.next = next
	}

	if next != nil {
		next.prev = prev
	}
}

func (node *LinkedListNode) InsertAfter(nodeToInsert *LinkedListNode, mountPoint float32) {
	startNode := node.Search(mountPoint)
	if startNode == nil {
		return
	}

	next := startNode.next

	startNode.next = nodeToInsert
	nodeToInsert.prev = startNode
	nodeToInsert.next = next

	if next != nil {
		next.prev = nodeToInsert
	}
}
