package dependency

import "sync"

const (
	// NodePending is a NodeStatus for pending node.
	NodePending NodeStatus = iota

	// NodeSuccess is a NodeStatus for success node.
	NodeSuccess

	// NodeFailed is a NodeStatus for failed node.
	NodeFailed
)

// NodeStatus is used to code release status - success or failed.
// Please use ReleaseSuccess and ReleaseFailed constants.
type NodeStatus int

// Node is graph node. N stands for data type.
type Node[N any] struct {
	Data         N
	dependencies []*Node[N]
	lock         sync.RWMutex
	status       NodeStatus
}

func newNode[N any](data N) *Node[N] {
	return &Node[N]{
		Data:         data,
		status:       NodePending,
		dependencies: make([]*Node[N], 0),
	}
}

func (node *Node[N]) SetSucceeded() {
	node.lock.Lock()
	defer node.lock.Unlock()

	node.status = NodeSuccess
}

func (node *Node[N]) SetFailed() {
	node.lock.Lock()
	defer node.lock.Unlock()

	node.status = NodeFailed
}

func (node *Node[N]) IsDone() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()

	return node.status != NodePending
}

func (node *Node[N]) IsFailed() bool {
	node.lock.RLock()
	defer node.lock.RUnlock()

	return node.status == NodeFailed
}

func (node *Node[N]) IsReady() bool {
	node.lock.RLock()
	deps := node.dependencies
	node.lock.RUnlock()

	for _, dependency := range deps {
		if !dependency.IsDone() {
			return false
		}

		if dependency.IsFailed() {
			node.SetFailed()

			return false
		}
	}

	return true
}

func (node *Node[N]) addDependency(dependency *Node[N]) {
	node.lock.Lock()
	defer node.lock.Unlock()

	node.dependencies = append(node.dependencies, dependency)
}
