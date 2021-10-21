package graph

import (
	"container/list"
	"ngac/internal/set"
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
	seen  map[string]struct{}
}

func NewBFS(gr Graph) *BFS {
	return &BFS{gr, nil}
}

func (bfs *BFS) Traverse(start *Node, propagator Propagator, visitor Visitor, direction Direction) error {
	// set up a queue to ensure FIFO
	queue := list.New()
	// set up a set to ensure nodes are only visited once
	bfs.seen = make(map[string]struct{})
	queue.PushBack(start)
	bfs.seen[start.Name] = struct{}{}

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
			if _, ok := bfs.seen[n.Name]; ok {
				continue
			}

			// add the node to the queue and the seen set
			queue.PushBack(n)

			bfs.seen[n.Name] = struct{}{}
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
	visited map[string]struct{}
}

func NewDFS(gr Graph) *DFS {
	return &DFS{gr, make(map[string]struct{})}
}

func (dfs *DFS) Traverse(start *Node, propagator Propagator, visitor Visitor, direction Direction) error {
	if _, ok := dfs.visited[start.Name]; ok {
		return nil
	}

	// mark the node as visited
	dfs.visited[start.Name] = struct{}{}

	nextLevel := dfs.nextLevel(start.Name, direction)
	for s := range nextLevel.Iter() {
		node, err := dfs.graph.Node(s.(string))
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

func (dfs *DFS) nextLevel(node string, direction Direction) set.Set {
	if direction == PARENTS {
		return dfs.graph.Parents(node)
	}

	return dfs.graph.Children(node)
}

type IDS struct {
	graph   Graph
	limit   int
	visited map[string]struct{}
}

func NewIDS(gr Graph, limit int) *IDS {
	return &IDS{gr, limit, make(map[string]struct{})}
}

func (ids *IDS) Traverse(start *Node, propagator Propagator, visitor Visitor, direction Direction) error {
	var depth int
	node := start
	var err error
	for depth < ids.limit && node != nil {
		node, err = ids.dls(node, propagator, visitor, direction, depth)
		if err != nil {
			return err
		}

		depth++
	}

	// after processing the parents, visit the start node
	return visitor(start)
}

func (ids *IDS) dls(n *Node, propagator Propagator, visitor Visitor, direction Direction, depth int) (*Node, error) {
	if _, ok := ids.visited[n.Name]; ok {
		return nil, nil
	}
	// mark the node as visited
	ids.visited[n.Name] = struct{}{}

	nextLevel := ids.nextLevel(n.Name, direction)

	for s := range nextLevel.Iter() {
		node, err := ids.graph.Node(s.(string))
		if err != nil {
			return nil, err
		}

		// traverse from the node
		ids.dls(node, propagator, visitor, direction, depth-1)

		if err := propagator(node, n); err != nil {
			return nil, err
		}
	}

	return nil, visitor(n)
}

func (ids *IDS) nextLevel(node string, direction Direction) set.Set {
	if direction == PARENTS {
		return ids.graph.Parents(node)
	}

	return ids.graph.Children(node)
}
