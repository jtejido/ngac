package dag

type NodeIterator interface {
	HasNext() bool
	Node() interface{}
	Reset()
}

type safeNodeIterator struct {
	nodes []interface{}
	index int
}

func newSafeNodeIterator(nodes []interface{}) *safeNodeIterator {
	return &safeNodeIterator{nodes, -1}
}

func (i *safeNodeIterator) Node() interface{} {
	i.index++
	return i.nodes[i.index]
}

func (i *safeNodeIterator) HasNext() bool {
	return i.index < (len(i.nodes) - 1)
}

func (i *safeNodeIterator) Reset() {
	i.index = -1
}

type EdgeIterator interface {
	HasNext() bool
	Edge() Edge
	Reset()
}

type safeEdgeIterator struct {
	edges []Edge
	index int
}

func newSafeEdgeIterator(edges []Edge) *safeEdgeIterator {
	return &safeEdgeIterator{edges, -1}
}

func (i *safeEdgeIterator) Edge() Edge {
	i.index++
	return i.edges[i.index]
}

func (i *safeEdgeIterator) HasNext() bool {
	return i.index < (len(i.edges) - 1)
}

func (i *safeEdgeIterator) Reset() {
	i.index = -1
}
