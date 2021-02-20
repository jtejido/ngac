package graph

import (
	"fmt"
	"github.com/jtejido/ngac/internal/omap"
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/operations"
)

const (
	node_not_found_msg = "node %s does not exist in the graph"
)

// This is an in-memory dag implementation.
type MemGraph struct {
	nodes map[string]*Node // contains all nodes
	from  map[string]map[string]Edge
	to    map[string]map[string]Edge
	pcs   set.Set // contains all policies
}

func NewMemGraph() *MemGraph {
	return &MemGraph{
		nodes: make(map[string]*Node),
		from:  make(map[string]map[string]Edge),
		to:    make(map[string]map[string]Edge),
		pcs:   set.NewSet(),
	}
}

func (mg *MemGraph) addNode(n *Node) {
	if _, exists := mg.nodes[n.Name]; exists {
		panic(fmt.Sprintf("simple: node collision: %s", n.Name))
	}

	mg.nodes[n.Name] = n
	mg.from[n.Name] = make(map[string]Edge)
	mg.to[n.Name] = make(map[string]Edge)

}

func (mg *MemGraph) node(name string) (n *Node) {
	n, _ = mg.nodes[name]
	return
}

func (mg *MemGraph) setEdge(e Edge) error {
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
		return fmt.Errorf("source vertex not in the graph.")
	}

	if !found2 {
		return fmt.Errorf("target vertex not in the graph.")
	}

	mg.from[sid][tid] = e
	mg.to[tid][sid] = e

	return nil
}

func (mg *MemGraph) removeNode(name string) {
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

func (mg *MemGraph) incomingEdgesOf(name string) []Edge {
	var edges []Edge
	if _, ok := mg.to[name]; !ok {
		return []Edge{}
	}

	for _, edge := range mg.to[name] {
		edges = append(edges, edge)
	}
	if len(edges) == 0 {
		return []Edge{}
	}

	return edges
}

func (mg *MemGraph) outgoingEdgesOf(name string) []Edge {
	var edges []Edge
	if _, ok := mg.from[name]; !ok {
		return []Edge{}
	}

	for _, edge := range mg.from[name] {
		edges = append(edges, edge)
	}

	if len(edges) == 0 {
		return []Edge{}
	}

	return edges
}

func (mg *MemGraph) hasEdgeFromTo(u, v string) bool {
	var found bool
	_, found = mg.from[u][v]

	return found
}

func (mg *MemGraph) removeEdge(fid, tid string) error {
	if _, ok := mg.nodes[fid]; !ok {
		return fmt.Errorf("source vertex not in the graph.")
	}
	if _, ok := mg.nodes[tid]; !ok {
		return fmt.Errorf("target vertex not in the graph.")
	}

	delete(mg.from[fid], tid)
	delete(mg.to[tid], fid)

	return nil
}

func (mg *MemGraph) CreatePolicyClass(name string, properties PropertyMap) (*Node, error) {
	if mg.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	// add the pc's name to the pc set and to the graph
	mg.pcs.Add(name)

	// create the node
	if properties == nil {
		properties = NewPropertyMap()
	}

	node := &Node{name, PC, properties}
	mg.addNode(node)

	return node, nil
}

func (mg *MemGraph) CreateNode(name string, t NodeType, properties PropertyMap, initialParent string, additionalParents ...string) (*Node, error) {
	//check for null values
	if t == PC {
		return nil, fmt.Errorf("use CreatePolicyClass to create a policy class node")
	} else if mg.Exists(name) {
		return nil, fmt.Errorf("the name %s already exists in the graph", name)
	}

	//store the node in the map
	if properties == nil {
		properties = NewPropertyMap()
	}

	node := &Node{name, t, properties}

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

func (mg *MemGraph) UpdateNode(name string, properties PropertyMap) error {
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

func (mg *MemGraph) RemoveNode(name string) {
	_, exists := mg.nodes[name]
	if !exists {
		return
	}

	//remove the node from the graph
	mg.removeNode(name)

	//remove the node from the policies if it is a policy class
	mg.pcs.Remove(name)
}

func (mg *MemGraph) Exists(name string) bool {
	_, exists := mg.nodes[name]
	return exists
}

func (mg *MemGraph) PolicyClasses() set.Set {
	return mg.pcs
}

func (mg *MemGraph) Nodes() set.Set {
	s := set.NewSet()
	for _, v := range mg.nodes {
		s.Add(v)
	}

	return s
}

func (mg *MemGraph) Node(name string) (*Node, error) {
	node := mg.node(name)
	if node == nil {
		return nil, fmt.Errorf("a node with the name %s does not exist", name)
	}

	return node, nil
}

func (mg *MemGraph) NodeFromDetails(t NodeType, properties PropertyMap) (*Node, error) {
	search := mg.Search(t, properties).Iterator()
	if !search.HasNext() {
		return nil, fmt.Errorf("a node matching the criteria (%s, %v) does not exist", t.String(), properties)
	}

	return search.Next().(*Node), nil
}

func (mg *MemGraph) Search(t NodeType, properties PropertyMap) set.Set {
	if properties == nil {
		properties = omap.NewOrderedMap().(PropertyMap)
	}

	results := set.NewSet()
	// iterate over the nodes to find ones that match the search parameters
	for _, node := range mg.nodes {
		if node.Type != t && t != ALL {
			continue
		}

		match := true
		for _, key := range properties.Keys() {
			checkValue, _ := properties.Get(key)
			foundValue, _ := node.Properties.Get(key)
			if checkValue != foundValue {
				match = false
			}
		}

		if match {
			results.Add(node)
		}
	}

	return results
}

func (mg *MemGraph) Children(name string) set.Set {
	if !mg.Exists(name) {
		return set.NewSet()
	}

	children := set.NewSet()
	for _, rel := range mg.incomingEdgesOf(name) {
		if _, ok := rel.(*Association); ok {
			continue
		}

		children.Add(rel.From())
	}

	return children
}

func (mg *MemGraph) Parents(name string) set.Set {
	if !mg.Exists(name) {
		return set.NewSet()
	}

	parents := set.NewSet()
	for _, rel := range mg.outgoingEdgesOf(name) {
		if _, ok := rel.(*Association); ok {
			continue
		}

		parents.Add(rel.To())
	}

	return parents
}

func (mg *MemGraph) Assign(child, parent string) error {
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

	if err := CheckAssignment(c.Type, p.Type); err != nil {
		return err
	}

	a := new(Assignment)
	a.Source = child
	a.Target = parent

	return mg.setEdge(a)
}

func (mg *MemGraph) Deassign(child, parent string) error {
	return mg.removeEdge(child, parent)
}

func (mg *MemGraph) IsAssigned(child, parent string) bool {
	return mg.hasEdgeFromTo(child, parent)
}

func (mg *MemGraph) Associate(ua, target string, ops operations.OperationSet) error {
	if !mg.Exists(ua) {
		return fmt.Errorf(node_not_found_msg, ua)
	} else if !mg.Exists(target) {
		return fmt.Errorf(node_not_found_msg, target)
	}

	uaNode := mg.node(ua)
	targetNode := mg.node(target)

	// check that the association is valid
	if err := CheckAssociation(uaNode.Type, targetNode.Type); err != nil {
		return err
	}

	// if no edge exists create an association
	// if an assignment exists create a new edge for the association
	// if an association exists update it
	edge, found := mg.from[ua][target]
	var isAssign bool
	if found {
		_, isAssign = edge.(*Assignment)
	}

	if !found || isAssign {
		e := new(Association)
		e.Source = ua
		e.Target = target
		e.Operations = ops
		if err := mg.setEdge(e); err != nil {
			return err
		}
	} else if assoc, ok := edge.(*Association); ok {
		assoc.Operations = ops
	}

	return nil
}

func (mg *MemGraph) Dissociate(ua, target string) error {
	return mg.removeEdge(ua, target)
}

func (mg *MemGraph) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
	if !mg.Exists(source) {
		return nil, fmt.Errorf(node_not_found_msg, source)
	}

	assocs := make(map[string]operations.OperationSet)
	for _, rel := range mg.outgoingEdgesOf(source) {
		if assoc, ok := rel.(*Association); ok {
			assocs[assoc.Target] = operations.NewOperationSetFromSet(assoc.Operations)
		}
	}
	return assocs, nil
}

func (mg *MemGraph) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
	if !mg.Exists(target) {
		return nil, fmt.Errorf(node_not_found_msg, target)
	}

	assocs := make(map[string]operations.OperationSet)
	for _, rel := range mg.incomingEdgesOf(target) {
		if assoc, ok := rel.(*Association); ok {
			assocs[assoc.Source] = operations.NewOperationSetFromSet(assoc.Operations)
		}
	}

	return assocs, nil
}
