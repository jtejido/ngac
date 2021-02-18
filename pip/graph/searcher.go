package graph

import (
	"container/list"
	"github.com/jtejido/ngac/internal/set"
)

type Direction int

const (
	CHILDREN Direction = iota
	PARENTS
)

type Searcher interface {
	Traverse(start *Node, propagator Propagator, visitor Visitor)
}

type Propagator func(parent, child *Node) error

type Visitor func(*Node) error

type BFS struct {
	graph Graph
	Seen  set.Set
}

func NewBFS(gr Graph) *BFS {
	return &BFS{gr, nil}
}

func (bfs *BFS) Traverse(start *Node, propagator Propagator, visitor Visitor, direction Direction) error {
	// set up a queue to ensure FIFO
	queue := list.New()
	// set up a set to ensure nodes are only visited once
	bfs.Seen = set.NewSet()
	queue.PushBack(start)
	bfs.Seen.Add(start.Name)

	for queue.Len() > 0 {
		qval := queue.Front()
		node := qval.Value.(*Node)
		if err := visitor(node); err != nil {
			return err
		}

		nextLevel := bfs.nextLevel(node.Name, direction)
		for s := range nextLevel.Iter() {
			n, err := bfs.graph.Node(s.(string))
			if err != nil {
				return err
			}

			// if this node has already been seen, we don't need to se it again
			if bfs.Seen.Contains(n.Name) {
				continue
			}

			// add the node to the queue and the seen set
			queue.PushBack(n)
			bfs.Seen.Add(n.Name)

			// propagate from the nextLevel to the current node
			if err := propagator(node, n); err != nil {
				return err
			}
		}

		queue.Remove(qval)
	}

	return nil
}

func (bfs *BFS) nextLevel(node string, direction Direction) set.Set {
	if direction == PARENTS {
		return bfs.graph.Parents(node)
	}

	return bfs.graph.Children(node)

}

type DFS struct {
	graph   Graph
	visited set.Set
}

func NewDFS(gr Graph) *DFS {
	return &DFS{gr, set.NewSet()}
}

func (dfs *DFS) Traverse(start *Node, propagator Propagator, visitor Visitor, direction Direction) error {
	if dfs.visited.Contains(start.Name) {
		return nil
	}

	// mark the node as visited
	dfs.visited.Add(start.Name)

	var nodes set.Set
	if direction == PARENTS {
		nodes = dfs.graph.Parents(start.Name)
	} else {
		nodes = dfs.graph.Children(start.Name)
	}

	it := nodes.Iterator()
	for it.HasNext() {
		n := it.Next().(string)
		node, err := dfs.graph.Node(n)
		if err != nil {
			return err
		}

		// traverse from the node
		dfs.Traverse(node, propagator, visitor, direction)

		// propagate from the node to the start node
		if err := propagator(node, start); err != nil {
			return err
		}
	}

	// after processing the parents, visit the start node
	return visitor(start)
}
