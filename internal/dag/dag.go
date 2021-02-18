package dag

// Edge is a graph edge. In directed graphs, the direction of the
// edge is given from -> to, otherwise the edge is semantically
// unordered.
type Edge interface {
	// From returns the from node of the edge.
	From() interface{}

	// To returns the to node of the edge.
	To() interface{}
}

// Graph is a generalized graph.
type Graph interface {
	// Node returns the node with the given ID if it exists
	// in the graph, and nil otherwise.
	Node(id interface{}) interface{}

	// Nodes returns all the nodes in the graph.
	//
	// Nodes must not return nil.
	Nodes() NodeIterator

	// From returns all nodes that can be reached directly
	// from the node with the given ID.
	//
	// From must not return nil.
	From(id interface{}) NodeIterator

	// HasEdgeBetween returns whether an edge exists between
	// nodes with IDs xid and yid without considering direction.
	HasEdgeBetween(xid, yid interface{}) bool

	// Edge returns the edge from u to v, with IDs uid and vid,
	// if such an edge exists and nil otherwise. The node v
	// must be directly reachable from u as defined by the
	// From method.
	Edge(uid, vid interface{}) Edge
}

// Directed is a directed graph.
type Directed interface {
	Graph

	// HasEdgeFromTo returns whether an edge exists
	// in the graph from u to v with IDs uid and vid.
	HasEdgeFromTo(uid, vid interface{}) bool

	// To returns all nodes that can reach directly
	// to the node with the given ID.
	//
	// To must not return nil.
	To(id interface{}) NodeIterator
}

// NodeAdder is an interface for adding arbitrary nodes to a graph.
type NodeAdder interface {

	// AddNode adds a node to the graph. AddNode panics if
	// the added node ID matches an existing node ID.
	AddNode(interface{})
}

// NodeRemover is an interface for removing nodes from a graph.
type NodeRemover interface {
	// RemoveNode removes the node with the given ID
	// from the graph, as well as any edges attached
	// to it. If the node is not in the graph it is
	// a no-op.
	RemoveNode(id interface{})
}

// EdgeAdder is an interface for adding edges to a graph.
type EdgeAdder interface {
	// SetEdge adds an edge from one node to another.
	// If the graph supports node addition the nodes
	// will be added if they do not exist, otherwise
	// SetEdge will panic.
	// The behavior of an EdgeAdder when the IDs
	// returned by e.From() and e.To() are equal is
	// implementation-dependent.
	// Whether e, e.From() and e.To() are stored
	// within the graph is implementation dependent.
	SetEdge(e Edge)
}

// EdgeRemover is an interface for removing nodes from a graph.
type EdgeRemover interface {
	// RemoveEdge removes the edge with the given end
	// IDs, leaving the terminal nodes. If the edge
	// does not exist it is a no-op.
	RemoveEdge(fid, tid interface{})
}

// Builder is a graph that can have nodes and edges added.
type Builder interface {
	NodeAdder
	EdgeAdder
}

// DirectedBuilder is a directed graph builder.
type DirectedBuilder interface {
	Directed
	Builder
}

var (
	_ NodeIterator = Empty
	_ EdgeIterator = Empty
)

const Empty = nothing

const nothing = empty(true)

type empty bool

func (empty) HasNext() bool     { return false }
func (empty) Reset()            {}
func (empty) Node() interface{} { return nil }
func (empty) Edge() Edge        { return nil }
