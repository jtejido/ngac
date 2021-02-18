package dag

import "fmt"

var (
	dg *DirectedGraph
	_  Graph       = dg
	_  Directed    = dg
	_  NodeAdder   = dg
	_  NodeRemover = dg
	_  EdgeAdder   = dg
	_  EdgeRemover = dg
)

// DirectedGraph implements a generalized directed
type DirectedGraph struct {
	nodes map[int]Node
	from  map[int]map[int]Edge
	to    map[int]map[int]Edge
}

// NewDirectedGraph returns a Directed
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		nodes: make(map[int]Node),
		from:  make(map[int]map[int]Edge),
		to:    make(map[int]map[int]Edge),
	}
}

// AddNode adds n to the  It panics if the added node ID matches an existing node ID.
func (g *DirectedGraph) AddNode(n Node) {
	if _, exists := g.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodes[n.ID()] = n
	g.from[n.ID()] = make(map[int]Edge)
	g.to[n.ID()] = make(map[int]Edge)
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *DirectedGraph) Edge(uid, vid int) Edge {
	edge, ok := g.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

// Edges returns all the edges in the
func (g *DirectedGraph) Edges() EdgeIterator {
	var edges []Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	if len(edges) == 0 {
		return Empty
	}
	return newSafeEdgeIterator(edges)
}

// From returns all nodes in g that can be reached directly from n.
func (g *DirectedGraph) From(id int) NodeIterator {
	if _, ok := g.from[id]; !ok {
		return Empty
	}

	from := make([]Node, len(g.from[id]))
	i := 0
	for vid := range g.from[id] {
		from[i] = g.nodes[vid]
		i++
	}
	if len(from) == 0 {
		return Empty
	}
	return newSafeNodeIterator(from)
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (g *DirectedGraph) HasEdgeBetween(xid, yid int) bool {
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (g *DirectedGraph) HasEdgeFromTo(uid, vid int) bool {
	if _, ok := g.from[uid][vid]; !ok {
		return false
	}
	return true
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *DirectedGraph) Node(id int) Node {
	return g.nodes[id]
}

// Nodes returns all the nodes in the
func (g *DirectedGraph) Nodes() NodeIterator {
	if len(g.nodes) == 0 {
		return Empty
	}
	nodes := make([]Node, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}
	return newSafeNodeIterator(nodes)
}

// RemoveEdge removes the edge with the given end point IDs from the graph, leaving the terminal
// nodes. If the edge does not exist it is a no-op.
func (g *DirectedGraph) RemoveEdge(fid, tid int) {
	if _, ok := g.nodes[fid]; !ok {
		return
	}
	if _, ok := g.nodes[tid]; !ok {
		return
	}

	delete(g.from[fid], tid)
	delete(g.to[tid], fid)
}

// RemoveNode removes the node with the given ID from the graph, as well as any edges attached
// to it. If the node is not in the graph it is a no-op.
func (g *DirectedGraph) RemoveNode(id int) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.from[id] {
		delete(g.to[from], id)
	}
	delete(g.from, id)

	for to := range g.to[id] {
		delete(g.from[to], id)
	}
	delete(g.to, id)
}

// SetEdge adds e, an edge from one node to another. If the nodes do not exist, they are added
// and are set to the nodes of the edge otherwise.
// It will panic if the IDs of the e.From and e.To are equal.
func (g *DirectedGraph) SetEdge(e Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if _, ok := g.nodes[fid]; !ok {
		g.AddNode(from)
	} else {
		g.nodes[fid] = from
	}
	if _, ok := g.nodes[tid]; !ok {
		g.AddNode(to)
	} else {
		g.nodes[tid] = to
	}

	g.from[fid][tid] = e
	g.to[tid][fid] = e
}

// To returns all nodes in g that can reach directly to n.
func (g *DirectedGraph) To(id int) NodeIterator {
	if _, ok := g.from[id]; !ok {
		return Empty
	}

	to := make([]Node, len(g.to[id]))
	i := 0
	for uid := range g.to[id] {
		to[i] = g.nodes[uid]
		i++
	}
	if len(to) == 0 {
		return Empty
	}
	return newSafeNodeIterator(to)
}

func (g *DirectedGraph) IncomingEdgesOf(id int) EdgeIterator {
	if _, ok := g.to[id]; !ok {
		return Empty
	}

	var edges []Edge
	for _, edge := range g.to[id] {
		edges = append(edges, edge)
	}
	if len(edges) == 0 {
		return Empty
	}

	return newSafeEdgeIterator(edges)
}

func (g *DirectedGraph) OutgoingEdgesOf(id int) EdgeIterator {
	if _, ok := g.from[id]; !ok {
		return Empty
	}

	var edges []Edge
	for _, edge := range g.from[id] {
		edges = append(edges, edge)
	}

	if len(edges) == 0 {
		return Empty
	}

	return newSafeEdgeIterator(edges)
}
