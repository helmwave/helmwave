package dependency

const (
	// NodePending is a NodeStatus for pending node.
	NodePending NodeStatus = iota

	// NodeSuccess is a NodeStatus for success node.
	NodeSuccess

	// NodeFailed is a NodeStatus for failed node.
	NodeFailed
)

// NodeStatus is used to code release status - success or failed.
// Please use ReleaseSuccess and ReleaseFailed contants.
type NodeStatus int

// Node is graph node. N stands for data type.
type Node[N any] struct {
	Data         N
	status       NodeStatus
	dependencies []*Node[N]
}

func newNode[N any](data N) *Node[N] {
	return &Node[N]{
		Data:         data,
		status:       NodePending,
		dependencies: make([]*Node[N], 0),
	}
}

func (node Node[N]) SetSucceeded() {
	node.status = NodeSuccess
}

func (node Node[N]) SetFailed() {
	node.status = NodeFailed
}

func (node Node[N]) IsDone() bool {
	return node.status != NodePending
}

func (node Node[N]) IsFailed() bool {
	return node.status == NodeFailed
}

func (node Node[N]) IsReady() bool {
	for _, dependency := range node.dependencies {
		if !dependency.IsDone() {
			return false
		}
	}

	return true
}

func (node Node[N]) addDependency(dependency *Node[N]) {
	node.dependencies = append(node.dependencies, dependency)
}

func (node Node[N]) available() bool {
	return true
}
