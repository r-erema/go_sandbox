package utils

type MAryTree struct {
	value    string
	children []*MAryTree
}

func (node MAryTree) Traverse() []string {
	return node.traverseHelper([]string{})
}

func (node MAryTree) traverseHelper(values []string) []string {
	values = append(values, node.value)

	for _, childNode := range node.children {
		values = childNode.traverseHelper(values)
	}

	return values
}
