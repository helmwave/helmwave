package dependency

import (
	"fmt"
)

// Graph is dependencies graph. K stands for map keys type (e.g. string names), N for data type.
type Graph[K comparable, N any] struct {
	Nodes        map[K]*Node[N]
	dependencies [][2]K
}

// NewGraph returns empty graph.
func NewGraph[K comparable, N any]() *Graph[K, N] {
	return &Graph[K, N]{
		Nodes:        make(map[K]*Node[N]),
		dependencies: make([][2]K, 0),
	}
}

func (graph *Graph[K, N]) NewNode(key K, data N) error {
	if _, ok := graph.Nodes[key]; ok {
		return fmt.Errorf("key %v already exists", key)
	}

	node := newNode(data)
	graph.Nodes[key] = node

	return nil
}

// AddDependency adds lazy dependency. It will be evaluated only in `Build` method.
func (graph *Graph[K, N]) AddDependency(dependant, dependency K) {
	graph.dependencies = append(graph.dependencies, [2]K{dependant, dependency})
}

func (graph *Graph[K, N]) Reverse() (*Graph[K, N], error) {
	newDependenciesGraph := NewGraph[K, N]()

	for key, node := range graph.Nodes {
		err := newDependenciesGraph.NewNode(key, node.Data)
		if err != nil {
			return nil, err
		}
	}

	for _, dep := range graph.dependencies {
		dependantKey := dep[1]
		dependencyKey := dep[0]
		newDependenciesGraph.AddDependency(dependantKey, dependencyKey)
	}

	err := newDependenciesGraph.Build()
	if err != nil {
		return nil, err
	}

	return newDependenciesGraph, nil
}

func (graph *Graph[K, N]) Build() error {
	for _, dep := range graph.dependencies {
		dependantKey := dep[0]
		dependencyKey := dep[1]

		dependant, ok := graph.Nodes[dependantKey]
		if !ok {
			return fmt.Errorf("dependant key %v does not exist", dependantKey)
		}
		dependency, ok := graph.Nodes[dependencyKey]
		if !ok {
			return fmt.Errorf("dependency key %v does not exist", dependencyKey)
		}

		dependant.addDependency(dependency)
	}

	if err := graph.detectCycle(); err != nil {
		return err
	}

	return nil
}

func (graph *Graph[K, N]) detectCycle() error {
	visited := make(map[*Node[N]]int)

	for _, node := range graph.Nodes {
		err := graph.dfs(node, visited)
		if err != nil {
			return err
		}
	}

	return nil
}

// dfs is Depth First Search.
func (graph *Graph[K, N]) dfs(node *Node[N], visited map[*Node[N]]int) error {
	// This means that during recursion we hit node that is already being dfs'd
	if visited[node] == -1 {
		return fmt.Errorf("graph loop detected (starts with %v)", node)
	}

	if visited[node] == 1 {
		return nil
	}

	visited[node] = -1

	for _, dep := range node.dependencies {
		err := graph.dfs(dep, visited)
		if err != nil {
			return err
		}
	}

	visited[node] = 1

	return nil
}

func (graph *Graph[K, N]) runChan(ch chan<- *Node[N]) {
	nodes := make(map[K]*Node[N])
	for key, node := range graph.Nodes {
		nodes[key] = node
	}

	for len(nodes) > 0 {
		for key, node := range nodes {
			// In case some release failed because it's dependency failed
			if node.IsDone() {
				delete(nodes, key)

				continue
			}

			if node.IsReady() {
				ch <- node
				delete(nodes, key)

				continue
			}
		}
	}

	close(ch)
}

// Run returns channel for data and runs goroutine that handles dependency graph
// and populates channel with ready to install releases.
func (graph *Graph[K, N]) Run() <-chan *Node[N] {
	ch := make(chan *Node[N], len(graph.Nodes))
	go graph.runChan(ch)

	return ch
}
