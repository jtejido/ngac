package memory

import (
	"fmt"
	"ngac/internal/set"
	"ngac/pkg/operations"
	g "ngac/pkg/pip/graph"
)

var _ g.Graph = &graph{}

const (
	node_not_found_msg = "node %s does not exist in the graph"
)

// This is an in-memory dag implementation.
type graph struct {
	nodes map[string]*g.Node // contains all nodes
	from  map[string]map[string]g.Edge
	to    map[string]map[string]g.Edge
	pcs   set.Set // contains all policies
}

func New() *graph {
	return &graph{
		nodes: make(map[string]*g.Node),
		from:  make(map[string]map[string]g.Edge),
		to:    make(map[string]map[string]g.Edge),
		pcs:   set.NewSet(),
	}
}

func (mg *graph) addNode(n *g.Node) {
	if _, exists := mg.nodes[n.Name]; exists {
		panic(fmt.Sprintf("simple: node collision: %s", n.Name))
	}

	mg.nodes[n.Name] = n
	mg.from[n.Name] = make(map[string]g.Edge)
	mg.to[n.Name] = make(map[string]g.Edge)

}

func (mg *graph) node(name string) (n *g.Node) {
	n, _ = mg.nodes[name]
	return
}

func (mg *graph) setEdge(e g.Edge) error {
	var (
		sid = e.From()
		tid = e.To()
	)

	if sid == tid {
		return fmt.Errorf("adding self edge")
	}
	var found1, found2 bool
	_, found1 = mg.nodes[sid]
	_, found2 = mg.nodes[tid]
	if !found1 {
		return fmt.Errorf("source vertex not in the g.")
	}

	if !found2 {
		return fmt.Errorf("target vertex not in the g.")
	}

	mg.from[sid][tid] = e
	mg.to[tid][sid] = e

	return nil
}

func (mg *graph) removeNode(name string) {
	var found bool
	_, found = mg.nodes[name]
	if !found {
		return
	}

	delete(mg.nodes, name)

	for from := range mg.from[name] {
		delete(mg.to[from], name)
	}
	delete(mg.from, name)

	for to := range mg.to[name] {
		delete(mg.from[to], name)
	}
	delete(mg.to, name)

}

func (mg *graph) incomingEdgesOf(name string) []g.Edge {
	var edges []g.Edge
	if _, ok := mg.to[name]; !ok {
		return []g.Edge{}
	}

	for _, edge := range mg.to[name] {
		edges = append(edges, edge)
	}
	if len(edges) == 0 {
		return []g.Edge{}
	}

	return edges
}

func (mg *graph) outgoingEdgesOf(name string) []g.Edge {
	var edges []g.Edge
	if _, ok := mg.from[name]; !ok {
		return []g.Edge{}
	}

	for _, edge := range mg.from[name] {
		edges = append(edges, edge)
	}

	if len(edges) == 0 {
		return []g.Edge{}
	}

	return edges
}

func (mg *graph) hasEdgeFromTo(u, v string) bool {
	var found bool
	_, found = mg.from[u][v]

	return found
}

func (mg *graph) removeEdge(fid, tid string) error {
	if _, ok := mg.nodes[fid]; !ok {
		return fmt.Errorf("source vertex not in the g.")
	}
	if _, ok := mg.nodes[tid]; !ok {
		return fmt.Errorf("target vertex not in the g.")
	}

	delete(mg.from[fid], tid)
	delete(mg.to[tid], fid)

	return nil
}

func (mg *graph) CreatePolicyClass(name string, properties g.PropertyMap) (*g.Node, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if mg.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	// add the pc's name to the pc set and to the graph
	mg.pcs.Add(name)

	// create the node
	if properties == nil {
		properties = g.NewPropertyMap()
	}

	node := &g.Node{name, g.PC, properties}
	mg.addNode(node)

	return node, nil
}

func (mg *graph) CreateNode(name string, t g.NodeType, properties g.PropertyMap, initialParent string, additionalParents ...string) (*g.Node, error) {
	//check for null values

	if t == g.PC {
		return nil, fmt.Errorf("use CreatePolicyClass to create a policy class node")
	} else if len(name) == 0 {
		return nil, fmt.Errorf("no name was provided when creating a node in the in-memory graph")
	} else if mg.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	//store the node in the map
	if properties == nil {
		properties = g.NewPropertyMap()
	}

	node := &g.Node{name, t, properties}

	mg.addNode(node)

	// assign the new node the to given parent nodes
	if err := mg.Assign(name, initialParent); err != nil {
		return nil, err
	}

	for _, parent := range additionalParents {
		if err := mg.Assign(name, parent); err != nil {
			return nil, err
		}
	}
	//return the Node
	return node, nil
}

func (mg *graph) UpdateNode(name string, properties g.PropertyMap) error {
	n, exists := mg.nodes[name]
	if !exists {
		return fmt.Errorf("node with the name %s could not be found to update", name)
	}

	// update the properties
	if properties != nil {
		n.Properties = properties
	}

	mg.nodes[name] = n // don't change the stored edges
	return nil
}

func (mg *graph) RemoveNode(name string) {
	_, exists := mg.nodes[name]
	if !exists {
		return
	}

	//remove the node from the graph
	mg.removeNode(name)

	//remove the node from the policies if it is a policy class
	mg.pcs.Remove(name)
}

func (mg *graph) Exists(name string) bool {
	_, exists := mg.nodes[name]
	return exists
}

func (mg *graph) PolicyClasses() set.Set {
	return mg.pcs
}

func (mg *graph) Nodes() set.Set {
	s := set.NewSet()
	for _, v := range mg.nodes {
		s.Add(v)
	}

	return s
}

func (mg *graph) Node(name string) (*g.Node, error) {
	node := mg.node(name)
	if node == nil {
		return nil, fmt.Errorf("a node with the name %s does not exist", name)
	}

	return node, nil
}

func (mg *graph) NodeFromDetails(t g.NodeType, properties g.PropertyMap) (*g.Node, error) {
	search := mg.Search(t, properties).Iterator()
	if !search.HasNext() {
		return nil, fmt.Errorf("a node matching the criteria (%s, %v) does not exist", t.String(), properties)
	}

	return search.Next().(*g.Node), nil
}

func (mg *graph) Search(t g.NodeType, properties g.PropertyMap) set.Set {
	if properties == nil {
		properties = g.NewPropertyMap()
	}

	results := set.NewSet()
	// iterate over the nodes to find ones that match the search parameters
	for _, node := range mg.nodes {
		if node.Type != t && t != g.NOOP {
			continue
		}

		match := true
		for k, v := range properties {
			if node.Properties[k] != v {
				match = false
			}
		}

		if match {
			results.Add(node)
		}
	}

	return results
}

func (mg *graph) Children(name string) set.Set {
	if !mg.Exists(name) {
		panic(fmt.Errorf(node_not_found_msg, name))
	}

	children := set.NewSet()
	for _, rel := range mg.incomingEdgesOf(name) {
		if _, ok := rel.(*g.Association); ok {
			continue
		}

		children.Add(rel.From())
	}

	return children
}

func (mg *graph) Parents(name string) set.Set {
	if !mg.Exists(name) {
		panic(fmt.Errorf(node_not_found_msg, name))
	}

	parents := set.NewSet()
	for _, rel := range mg.outgoingEdgesOf(name) {
		if _, ok := rel.(*g.Association); ok {
			continue
		}

		parents.Add(rel.To())
	}

	return parents
}

func (mg *graph) Assign(child, parent string) error {
	if !mg.Exists(child) {
		return fmt.Errorf(node_not_found_msg, child)
	} else if !mg.Exists(parent) {
		return fmt.Errorf(node_not_found_msg, parent)
	}

	if mg.hasEdgeFromTo(child, parent) {
		return fmt.Errorf("%s is already assigned to %s", parent, child)
	}

	c := mg.node(child)
	p := mg.node(parent)

	if err := g.CheckAssignment(c.Type, p.Type); err != nil {
		return err
	}

	a := new(g.Assignment)
	a.Source = child
	a.Target = parent

	return mg.setEdge(a)
}

func (mg *graph) Deassign(child, parent string) error {
	return mg.removeEdge(child, parent)
}

func (mg *graph) IsAssigned(child, parent string) bool {
	return mg.hasEdgeFromTo(child, parent)
}

func (mg *graph) Associate(ua, target string, ops operations.OperationSet) error {
	if !mg.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !mg.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	uaNode := mg.node(ua)
	targetNode := mg.node(target)

	// check that the association is valid
	if err := g.CheckAssociation(uaNode.Type, targetNode.Type); err != nil {
		return err
	}

	// if no edge exists create an association
	// if an assignment exists create a new edge for the association
	// if an association exists update it
	edge, found := mg.from[ua][target]
	var isAssign bool
	if found {
		_, isAssign = edge.(*g.Assignment)
	}

	if !found || isAssign {
		e := new(g.Association)
		e.Source = ua
		e.Target = target
		e.Operations = ops
		if err := mg.setEdge(e); err != nil {
			return err
		}
	} else if assoc, ok := edge.(*g.Association); ok {
		assoc.Operations = ops
	}

	return nil
}

func (mg *graph) Dissociate(ua, target string) error {
	return mg.removeEdge(ua, target)
}

func (mg *graph) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
	if !mg.Exists(source) {
		return nil, fmt.Errorf(node_not_found_msg, source)
	}

	assocs := make(map[string]operations.OperationSet)
	for _, rel := range mg.outgoingEdgesOf(source) {
		if assoc, ok := rel.(*g.Association); ok {
			assocs[assoc.Target] = operations.NewOperationSetFromSet(assoc.Operations)
		}
	}
	return assocs, nil
}

func (mg *graph) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
	if !mg.Exists(target) {
		return nil, fmt.Errorf(node_not_found_msg, target)
	}

	assocs := make(map[string]operations.OperationSet)
	for _, rel := range mg.incomingEdgesOf(target) {
		if assoc, ok := rel.(*g.Association); ok {
			assocs[assoc.Source] = operations.NewOperationSetFromSet(assoc.Operations)
		}
	}

	return assocs, nil
}
