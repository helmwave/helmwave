package dependency

func NewNode[N any](data N) *Node[N] {
	return newNode(data)
}

func (node *Node[N]) AddDependency(dependency *Node[N]) {
	node.addDependency(dependency)
}
