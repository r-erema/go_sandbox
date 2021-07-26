package utils

type BST struct {
	value       float32
	left, right *BST
}

func (node *BST) Left() *BST {
	return node.left
}

func (node *BST) Right() *BST {
	return node.right
}

func (node *BST) Value() float32 {
	return node.value
}

func NewBST(value float32) *BST {
	return &BST{value: value, left: nil, right: nil}
}

func (node *BST) InsertRecursively(nodeToAdd *BST) *BST {
	if nodeToAdd.value < node.value {
		node.leftBehaviourRecursively(nodeToAdd)
	} else {
		node.rightBehaviourRecursively(nodeToAdd)
	}

	return node
}

func (node *BST) leftBehaviourRecursively(nodeToAdd *BST) {
	if node.left == nil {
		node.left = nodeToAdd
	} else {
		node.left.InsertRecursively(nodeToAdd)
	}
}

func (node *BST) rightBehaviourRecursively(nodeToAdd *BST) {
	if node.right == nil {
		node.right = nodeToAdd
	} else {
		node.right.InsertRecursively(nodeToAdd)
	}
}

func (node *BST) InsertIteratively(nodeToAdd *BST) *BST {
	currentNode := node

	for {
		if nodeToAdd.value < currentNode.value {
			if currentNode.Left() == nil {
				currentNode.left = nodeToAdd

				break
			}

			currentNode = currentNode.Left()
		} else {
			if currentNode.Right() == nil {
				currentNode.right = nodeToAdd

				break
			}

			currentNode = currentNode.Right()
		}
	}

	return node
}

func (node *BST) TraversePreorder() chan *BST {
	ch := make(chan *BST)

	go func() {
		ch <- node

		if node.left != nil {
			for v := range node.left.TraversePreorder() {
				ch <- v
			}
		}

		if node.right != nil {
			for v := range node.right.TraversePreorder() {
				ch <- v
			}
		}

		close(ch)
	}()

	return ch
}

func (node *BST) TraverseInorder() chan *BST {
	ch := make(chan *BST)

	go func() {
		if node.left != nil {
			for v := range node.left.TraverseInorder() {
				ch <- v
			}
		}

		ch <- node

		if node.right != nil {
			for v := range node.right.TraverseInorder() {
				ch <- v
			}
		}

		close(ch)
	}()

	return ch
}

func (node *BST) find(value float32) *BST {
	if node == nil || node.value == value {
		return node
	}

	if node.value < value {
		return node.right.find(value)
	}

	return node.left.find(value)
}

func (node *BST) inorderSuccessor(valueOfNode float32) (targetNode, inorderSuccessor *BST) {
	nextIsSuccessor := false
	for n := range node.TraverseInorder() {
		if nextIsSuccessor {
			inorderSuccessor = n

			return
		}

		if valueOfNode == n.value {
			nextIsSuccessor = true
			targetNode = n
		}
	}

	return
}

func (node *BST) Delete(value float32) (newBST *BST) {
	tmpBSTClone := *node
	tmpBST := &tmpBSTClone
	onlySingleChild := func(node *BST) (child *BST) {
		if node.left != nil && node.right == nil {
			child = node.left
		}

		if node.right != nil && node.left == nil {
			child = node.right
		}

		return
	}
	targetNode, inorderSuccessor := tmpBST.inorderSuccessor(value)

	if targetNode.IsLeaf() {
		newBST = leafModeDeletion(tmpBST, targetNode)
	} else if child := onlySingleChild(targetNode); child != nil {
		newBST = onlySingleChildModeDeletion(tmpBST, targetNode, child)
	} else {
		newBST = multiChildModeDeletion(tmpBST, targetNode, inorderSuccessor)
	}

	*node = *newBST

	return
}

func leafModeDeletion(tmpBST, targetNode *BST) (newBST *BST) {
	newBST = NewBST(tmpBST.value)

	for n := range tmpBST.TraversePreorder() {
		if n != targetNode && n != tmpBST {
			newBST.InsertIteratively(NewBST(n.value))
		}
	}

	return
}

func onlySingleChildModeDeletion(tmpBST, targetNode, child *BST) (newBST *BST) {
	newBST = NewBST(tmpBST.value)
	targetNode.value = child.value

	for n := range tmpBST.TraversePreorder() {
		if n != child && n != tmpBST {
			newBST.InsertIteratively(NewBST(n.value))
		}
	}

	return
}

func multiChildModeDeletion(tmpBST, targetNode, inorderSuccessor *BST) (newBST *BST) {
	targetNode.value = inorderSuccessor.value
	newBST = NewBST(tmpBST.value)

	for n := range tmpBST.TraversePreorder() {
		if n != inorderSuccessor && n != tmpBST {
			newBST.InsertIteratively(NewBST(n.value))
		}
	}

	return
}

func (node *BST) IsLeaf() bool {
	return node != nil && node.left == nil && node.right == nil
}

func (node *BST) BranchVectors() [][]float32 {
	return branchVectorsHelper(node, []float32{}, [][]float32{})
}

func branchVectorsHelper(node *BST, branch []float32, branches [][]float32) [][]float32 {
	branch = append(branch, node.value)

	if node.IsLeaf() {
		return append(branches, branch)
	}

	branchTempCopy := make([]float32, len(branch))
	copy(branchTempCopy, branch)

	if node.left != nil {
		branches = branchVectorsHelper(node.left, branchTempCopy, branches)
	}

	if node.right != nil {
		branches = branchVectorsHelper(node.right, branchTempCopy, branches)
	}

	return branches
}
