package service

import (
	"fmt"
	"github.com/jtejido/ngac/internal/set"
	// "github.com/jtejido/ngac/epp"
	"github.com/jtejido/ngac/operations"
	"github.com/jtejido/ngac/pdp/decider"
	"github.com/jtejido/ngac/pdp/policy"
	"github.com/jtejido/ngac/pip/graph"
)

type pair = graph.PropertyPair

type GraphService struct {
	Service
	// epp         *epp.EPP
	graph       graph.Graph
	superPolicy *policy.SuperPolicy
}

func NewGraphService(g graph.Graph, d decider.Decider) (*GraphService, error) {
	if g == nil {
		return nil, fmt.Errorf("graph cannot be nil")
	}
	gs := new(GraphService)
	gs.graph = g
	if err := gs.configureSuperPolicy(); err != nil {
		return nil, err
	}

	return gs, nil
}

func (gs *GraphService) configureSuperPolicy() error {
	gs.superPolicy = policy.NewSuperPolicy()
	if err := gs.superPolicy.Configure(gs.graph); err != nil {
		return err
	}

	return nil
}

func (gs *GraphService) SuperPolicy() *policy.SuperPolicy {
	return gs.superPolicy
}

func (gs *GraphService) CreatePolicyClass(ctx Context, name string, properties graph.PropertyMap) (*graph.Node, error) {
	// check that the user can create a policy class
	if !gs.hasPermissions(ctx, gs.superPolicy.SuperPolicyClassRep(), operations.CREATE_POLICY_CLASS) {
		return nil, fmt.Errorf("unauthorized permissions to create a policy class")
	}

	if properties == nil {
		properties = graph.NewPropertyMap()
	}

	// create the PC node
	rep := name + "_rep"
	defaultUA := name + "_default_UA"
	defaultOA := name + "_default_OA"

	g := gs.graph
	properties.AddMap(graph.ToProperties(pair{"default_ua", defaultUA}, pair{"default_oa", defaultOA}, pair{graph.REP_PROPERTY, rep}))
	pcNode, err := g.CreatePolicyClass(name, properties)
	if err != nil {
		return nil, err
	}
	// create the PC UA node
	pcUANode, err := g.CreateNode(defaultUA, graph.UA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, name}), pcNode.Name)
	if err != nil {
		return nil, err
	}
	// create the PC OA node
	pcOANode, err := g.CreateNode(defaultOA, graph.OA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, name}), pcNode.Name)
	if err != nil {
		return nil, err
	}

	// assign Super U to PC UA
	// getPAP().getGraphPAP().assign(superPolicy.getSuperU().getID(), pcUANode.getID());
	// assign superUA and superUA2 to PC
	if err := g.Assign(gs.superPolicy.SuperUserAttribute().Name, pcNode.Name); err != nil {
		return nil, err
	}
	if err := g.Assign(gs.superPolicy.SuperUserAttribute2().Name, pcNode.Name); err != nil {
		return nil, err
	}
	// associate Super UA and PC UA
	if err := g.Associate(gs.superPolicy.SuperUserAttribute().Name, pcUANode.Name, operations.NewOperationSet(operations.ALL_OPERATIONS)); err != nil {
		return nil, err
	}
	// associate Super UA and PC OA
	if err := g.Associate(gs.superPolicy.SuperUserAttribute().Name, pcOANode.Name, operations.NewOperationSet(operations.ALL_OPERATIONS)); err != nil {
		return nil, err
	}

	// create an OA that will represent the pc
	if _, err := g.CreateNode(rep, graph.OA, graph.ToProperties(pair{"pc", name}), gs.superPolicy.SuperObjectAttribute().Name); err != nil {
		return nil, err
	}

	return pcNode, nil
}

func (gs *GraphService) CreateNode(ctx Context, name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (*graph.Node, error) {
	// instantiate the properties map if it's null
	if properties == nil {
		properties = graph.NewPropertyMap()
	}

	parentNode, err := gs.Node(ctx, initialParent)
	if err != nil {
		return nil, err
	}

	// check that the user has the permission to assign to the parent node
	if !gs.hasPermissions(ctx, parentNode, operations.ASSIGN_TO) {
		// if the user cannot assign to the parent node, delete the newly created node
		return nil, fmt.Errorf("unauthorized permission \"%s\" on node %s", operations.ASSIGN_TO, initialParent)
	}

	g := gs.graph

	if parentNode.Type == graph.PC {
		initialParent := gs.PolicyClassDefault(parentNode.Name, t)
		parentNode, err = gs.Node(ctx, initialParent)
		if err != nil {
			return nil, err
		}
	}

	// check any additional parents before assigning
	for i := 0; i < len(additionalParents); i++ {
		additionalParentNode, err := gs.Node(ctx, additionalParents[i])
		if err != nil {
			return nil, err
		}

		if !gs.hasPermissions(ctx, additionalParentNode, operations.ASSIGN_TO) {
			// if the user cannot assign to the parent node, delete the newly created node
			return nil, fmt.Errorf("unauthorized permission \"%s\" on %s", operations.ASSIGN_TO, additionalParentNode.Name)
		}

		if additionalParentNode.Type == graph.PC {
			additionalParents[i] = gs.PolicyClassDefault(additionalParentNode.Name, t)
		}
	}

	//create the node
	node, err := g.CreateNode(name, t, properties, initialParent, additionalParents...)
	if err != nil {
		return nil, err
	}

	// // process event for initial parent
	// gs.epp.ProcessEvent(epp.NewAssignToEvent(parentNode, node), ctx.User(), ctx.Process())
	// // process event for any additional parents
	// for _, parent := range additionalParents {
	// 	parentNode, err = g.Node(parent)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	gs.epp.ProcessEvent(epp.NewAssignToEvent(parentNode, node), ctx.User(), ctx.Process())
	// }

	return node, nil
}

func (gs *GraphService) PolicyClassDefault(pc string, t graph.NodeType) string {
	return pc + "_default_" + t.String()
}

func (gs *GraphService) UpdateNode(ctx Context, name string, properties graph.PropertyMap) error {
	node, err := gs.Node(ctx, name)
	if err != nil {
		return err
	}

	// check that the user can update the node
	if !gs.hasPermissions(ctx, node, operations.UPDATE_NODE) {
		return fmt.Errorf("unauthorized permission %s on node %s", operations.UPDATE_NODE, node.Name)
	}

	return gs.graph.UpdateNode(name, properties)
}

func (gs *GraphService) RemoveNode(ctx Context, name string) error {
	node, err := gs.Node(ctx, name)
	if err != nil {
		return err
	}

	// check the user can deassign the node
	if !gs.hasPermissions(ctx, node, operations.DEASSIGN) {
		return fmt.Errorf("unauthorized permissions on %s: %s", node.Name, operations.DEASSIGN)
	}

	// check that the user can deassign from the node's parents
	parents := gs.graph.Parents(name)
	for parent := range parents.Iter() {
		parentNode, err := gs.Node(ctx, parent.(string))
		if err != nil {
			return err
		}

		if !gs.hasPermissions(ctx, parentNode, operations.DEASSIGN_FROM) {
			return fmt.Errorf("unauthorized permissions on %s: %s", parentNode.Name, operations.DEASSIGN_FROM)
		}

		// gs.epp.ProcessEvent(epp.NewDeassignEvent(node, parentNode), ctx.User(), ctx.Process())
		// gs.epp.ProcessEvent(epp.NewDeassignFromEvent(parentNode, node), ctx.User(), ctx.Process())
	}

	// if it's a PC, delete the rep
	if node.Type == graph.PC {
		if v, found := node.Properties.Get(graph.REP_PROPERTY); found {
			gs.graph.RemoveNode(v.(string))
		}
	}

	gs.graph.RemoveNode(name)

	return nil
}

func (gs *GraphService) Exists(ctx Context, name string) bool {
	if !gs.graph.Exists(name) {
		return false
	}

	node, err := gs.graph.Node(name)
	if err != nil {
		return false
	}

	// if the node is a pc the user must have permission on the rep OA of the PC
	if node.Type == graph.PC {
		pcRep, _ := node.Properties.Get(graph.REP_PROPERTY)
		node, err = gs.graph.Node(pcRep.(string))
		if err != nil {
			return false
		}

		return gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS)
	}

	// node exists, return false if the user does not have access to it.
	return gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS)
}

func (gs *GraphService) Nodes(ctx Context) set.Set {
	nodes := gs.graph.Nodes()
	nodes.Filter(func(v interface{}) bool {
		return !gs.hasPermissions(ctx, v.(*graph.Node), operations.ANY_OPERATIONS)
	})

	return nodes
}

func (gs *GraphService) PolicyClasses() set.Set {
	return gs.graph.PolicyClasses()
}

func (gs *GraphService) Children(ctx Context, name string) set.Set {
	if !gs.Exists(ctx, name) {
		return set.NewSet()
	}

	children := gs.graph.Children(name)
	children.Filter(func(v interface{}) bool {
		node, err := gs.graph.Node(v.(string))
		if err != nil {
			return false
		}

		return !gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS)
	})

	return children
}

func (gs *GraphService) Parents(ctx Context, name string) set.Set {
	if !gs.Exists(ctx, name) {
		return set.NewSet()
	}

	parents := gs.graph.Parents(name)
	parents.Filter(func(v interface{}) bool {
		node, err := gs.graph.Node(v.(string))
		if err != nil {
			return false
		}

		return !gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS)
	})

	return parents
}

func (gs *GraphService) Assign(ctx Context, child, parent string) error {
	//check if the assignment is valid
	childNode, err := gs.Node(ctx, child)
	if err != nil {
		return err
	}

	parentNode, err := gs.Node(ctx, parent)
	if err != nil {
		return err
	}

	//check the user can assign the child
	if !gs.hasPermissions(ctx, childNode, operations.ASSIGN) {
		return fmt.Errorf("unauthorized permission %s on node %s", operations.ASSIGN, childNode.Name)
	}

	if err := graph.CheckAssignment(childNode.Type, parentNode.Type); err != nil {
		return err
	}

	// check that the user can assign to the parent node
	if !gs.hasPermissions(ctx, parentNode, operations.ASSIGN_TO) {
		return fmt.Errorf("unauthorized permission %s on node %s", operations.ASSIGN_TO, parentNode.Name)
	}

	if parentNode.Type == graph.PC {
		parent = gs.PolicyClassDefault(parentNode.Name, childNode.Type)
	}

	// assign in the PAP
	if err := gs.graph.Assign(child, parent); err != nil {
		return err
	}

	// gs.epp.ProcessEvent(epp.NewAssignEvent(childNode, parentNode), ctx.User(), ctx.Process())
	// gs.epp.ProcessEvent(epp.NewAssignToEvent(parentNode, childNode), ctx.User(), ctx.Process())
	return nil
}

func (gs *GraphService) Deassign(ctx Context, child, parent string) error {
	//check if the assignment is valid
	childNode, err := gs.Node(ctx, child)
	if err != nil {
		return err
	}

	parentNode, err := gs.Node(ctx, parent)
	if err != nil {
		return err
	}

	//check the user can deassign the child
	if !gs.hasPermissions(ctx, childNode, operations.DEASSIGN) {
		return fmt.Errorf("unauthorized permissions on %s: %s", childNode.Name, operations.DEASSIGN)
	}

	//check that the user can deassign from the parent
	if !gs.hasPermissions(ctx, parentNode, operations.DEASSIGN_FROM) {
		return fmt.Errorf("unauthorized permissions on %s: %s", parentNode.Name, operations.DEASSIGN_FROM)
	}

	//delete assignment in PAP
	if err := gs.graph.Deassign(childNode.Name, parentNode.Name); err != nil {
		return err
	}

	// gs.epp.ProcessEvent(epp.NewDeassignEvent(childNode, parentNode), ctx.User(), ctx.Process())
	// gs.epp.ProcessEvent(epp.NewDeassignFromEvent(parentNode, childNode), ctx.User(), ctx.Process())
	return nil
}

func (gs *GraphService) IsAssigned(ctx Context, child, parent string) bool {
	parentNode, err := gs.Node(ctx, parent)
	if err != nil {
		return false
	}
	childNode, err := gs.Node(ctx, child)
	if err != nil {
		return false
	}

	return gs.graph.IsAssigned(childNode.Name, parentNode.Name)
}

func (gs *GraphService) Associate(ctx Context, ua, target string, ops operations.OperationSet) error {
	sourceNode, err := gs.Node(ctx, ua)
	if err != nil {
		return err
	}

	targetNode, err := gs.Node(ctx, target)
	if err != nil {
		return err
	}

	err = graph.CheckAssociation(sourceNode.Type, targetNode.Type)
	if err != nil {
		return err
	}

	//check the user can associate the source and target nodes
	if !gs.hasPermissions(ctx, sourceNode, operations.ASSOCIATE) {
		return fmt.Errorf("unauthorized permissions on %s: %s", sourceNode.Name, operations.ASSOCIATE)
	}
	if !gs.hasPermissions(ctx, targetNode, operations.ASSOCIATE) {
		return fmt.Errorf("unauthorized permissions on %s: %s", targetNode.Name, operations.ASSOCIATE)
	}

	//create association in PAP
	return gs.graph.Associate(sourceNode.Name, targetNode.Name, ops)
}

func (gs *GraphService) Dissociate(ctx Context, ua, target string) error {
	sourceNode, err := gs.Node(ctx, ua)
	if err != nil {
		return err
	}

	targetNode, err := gs.Node(ctx, target)
	if err != nil {
		return err
	}

	//check the user can associate the source and target nodes
	if !gs.hasPermissions(ctx, sourceNode, operations.DISASSOCIATE) {
		return fmt.Errorf("unauthorized permissions on %s: %s", sourceNode.Name, operations.DISASSOCIATE)
	}
	if !gs.hasPermissions(ctx, targetNode, operations.DISASSOCIATE) {
		return fmt.Errorf("unauthorized permissions on %s: %s", targetNode.Name, operations.DISASSOCIATE)
	}

	//create association in PAP
	return gs.graph.Dissociate(sourceNode.Name, targetNode.Name)
}

func (gs *GraphService) SourceAssociations(ctx Context, source string) (map[string]operations.OperationSet, error) {
	sourceNode, err := gs.Node(ctx, source)
	if err != nil {
		return nil, err
	}

	//check the user can get the associations of the source node
	if !gs.hasPermissions(ctx, sourceNode, operations.GET_ASSOCIATIONS) {
		return nil, fmt.Errorf("unauthorized permissions on %s: %s", sourceNode.Name, operations.GET_ASSOCIATIONS)
	}

	return gs.graph.SourceAssociations(sourceNode.Name)
}

func (gs *GraphService) TargetAssociations(ctx Context, target string) (map[string]operations.OperationSet, error) {
	targetNode, err := gs.Node(ctx, target)
	if err != nil {
		return nil, err
	}

	//check the user can get the associations of the source node
	if !gs.hasPermissions(ctx, targetNode, operations.GET_ASSOCIATIONS) {
		return nil, fmt.Errorf("unauthorized permissions on %s: %s", targetNode.Name, operations.GET_ASSOCIATIONS)
	}

	return gs.graph.TargetAssociations(targetNode.Name)
}

func (gs *GraphService) Search(ctx Context, t graph.NodeType, properties graph.PropertyMap) set.Set {
	search := gs.graph.Search(t, properties)
	search.Filter(func(v interface{}) bool {
		return !gs.hasPermissions(ctx, v.(*graph.Node), operations.ANY_OPERATIONS)
	})

	return search
}

func (gs *GraphService) Node(ctx Context, name string) (*graph.Node, error) {
	if !gs.Exists(ctx, name) {
		return nil, fmt.Errorf("node %s could not be found", name)
	}

	node, _ := gs.graph.Node(name)

	if !gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS) {
		return nil, fmt.Errorf("unauthorized permissions on %s: %s", node.Name, operations.ANY_OPERATIONS)
	}

	return node, nil
}

func (gs *GraphService) NodeFromDetails(ctx Context, t graph.NodeType, properties graph.PropertyMap) (*graph.Node, error) {
	node, err := gs.graph.NodeFromDetails(t, properties)
	if err != nil {
		return nil, err
	}

	if !gs.hasPermissions(ctx, node, operations.ANY_OPERATIONS) {
		return nil, fmt.Errorf("node (%s, %v) could not be found", t.String(), properties)
	}

	return node, nil
}

func (gs *GraphService) Reset(ctx Context) error {
	// check that the user can reset the graph
	if !gs.hasPermissions(ctx, gs.superPolicy.SuperPolicyClassRep(), operations.RESET) {
		return fmt.Errorf("unauthorized permissions to reset the graph")
	}

	nodes := gs.graph.Nodes()
	names := []string{}
	for node := range nodes.Iter() {
		names = append(names, node.(*graph.Node).Name)
	}

	for _, name := range names {
		gs.graph.RemoveNode(name)
	}

	return nil
}
