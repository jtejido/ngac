package pap

import (
	"fmt"
	"ngac/internal/set"
	"ngac/pkg/common"
	"ngac/pkg/operations"
	"ngac/pkg/pap/policy"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

var (
	_ graph.Graph = &GraphAdmin{}
)

type GraphAdmin struct {
	pip         common.PolicyStore
	graph       graph.Graph
	superPolicy *policy.SuperPolicy
}

type pair = graph.PropertyPair

func NewGraphAdmin(pip common.PolicyStore) (*GraphAdmin, error) {
	ans := &GraphAdmin{pip, pip.Graph(), policy.NewSuperPolicy()}
	err := ans.superPolicy.Configure(ans.graph)
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func (ga *GraphAdmin) PolicyClassDefault(pcName string, t graph.NodeType) string {
	return pcName + "_default_" + t.String()
}

func (ga *GraphAdmin) CreatePolicyClass(name string, properties graph.PropertyMap) (*graph.Node, error) {
	nodeProps := graph.NewPropertyMap()
	if properties != nil {
		nodeProps = properties
	}

	rep := name + "_rep"
	defaultUA := ga.PolicyClassDefault(name, graph.UA)
	defaultOA := ga.PolicyClassDefault(name, graph.OA)

	nodeProps["default_ua"] = defaultUA
	nodeProps["default_oa"] = defaultOA
	nodeProps[graph.REP_PROPERTY] = rep

	pcNode := graph.NewNode()

	if err := ga.pip.RunTx(func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error {
		// create the pc node
		node, err := g.CreatePolicyClass(name, nodeProps)
		if err != nil {
			return err
		}
		pcNode.Name = node.Name
		pcNode.Type = node.Type
		pcNode.Properties = node.Properties

		// create the PC UA node
		pcUANode, err := g.CreateNode(defaultUA, graph.UA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, name}), pcNode.Name)
		if err != nil {
			return err
		}
		// create the PC OA node
		pcOANode, err := g.CreateNode(defaultOA, graph.OA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, name}), pcNode.Name)
		if err != nil {
			return err
		}

		// assign Super U to PC UA
		// getPAP().getGraphPAP().assign(superPolicy.getSuperU().getID(), pcUANode.getID());
		// assign superUA and superUA2 to PC
		err = g.Assign(ga.superPolicy.SuperUserAttribute().Name, pcNode.Name)
		if err != nil {
			return err
		}
		err = g.Assign(ga.superPolicy.SuperUserAttribute2().Name, pcNode.Name)
		if err != nil {
			return err
		}
		// associate Super UA and PC UA
		err = g.Associate(ga.superPolicy.SuperUserAttribute().Name, pcUANode.Name, operations.NewOperationSet(operations.ALL_OPS))
		if err != nil {
			return err
		}
		// associate Super UA and PC OA
		err = g.Associate(ga.superPolicy.SuperUserAttribute().Name, pcOANode.Name, operations.NewOperationSet(operations.ALL_OPS))
		if err != nil {
			return err
		}
		// create an OA that will represent the pc
		_, err = g.CreateNode(rep, graph.OA, graph.ToProperties(pair{"pc", name}), ga.superPolicy.SuperObjectAttribute().Name)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return pcNode, nil
}

func (ga *GraphAdmin) CreateNode(name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (*graph.Node, error) {
	if t == graph.PC {
		return nil, fmt.Errorf("use CreatePolicyClass to create a policy class node")
	}

	defaultType := graph.UA
	if t == graph.OA || t == graph.O {
		defaultType = graph.OA
	}

	node := graph.NewNode()
	if err := ga.pip.RunTx(func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error {
		// if the parent is a PC get the PC default
		parentNode, err := g.Node(initialParent)
		if err != nil {
			return err
		}
		var pcInitParent string
		noname := true
		if parentNode.Type == graph.PC {
			pcInitParent = ga.PolicyClassDefault(parentNode.Name, defaultType)
			noname = false
		}

		for i := 0; i < len(additionalParents); i++ {
			parent := additionalParents[i]

			// if the parent is a PC get the PC default attribute
			additionalParentNode, err := g.Node(parent)
			if err != nil {
				return err
			}
			if additionalParentNode.Type == graph.PC {
				additionalParents[i] = ga.PolicyClassDefault(additionalParentNode.Name, defaultType)
			}
		}
		var nn string
		if noname {
			nn = initialParent
		} else {
			nn = pcInitParent
		}

		n, err := g.CreateNode(name, t, properties, nn, additionalParents...)
		if err != nil {
			return err
		}
		node.Name = n.Name
		node.Type = n.Type
		node.Properties = n.Properties
		return nil
	}); err != nil {
		return nil, err
	}

	return node, nil
}

func (ga *GraphAdmin) UpdateNode(name string, properties graph.PropertyMap) error {
	return ga.graph.UpdateNode(name, properties)
}
func (ga *GraphAdmin) RemoveNode(name string) {
	if ga.graph.Children(name).Len() != 0 {
		panic(fmt.Sprintf("cannot delete %s, nodes are still assigned to it", name))
	}
	ga.graph.RemoveNode(name)
}
func (ga *GraphAdmin) Exists(name string) bool {
	return ga.graph.Exists(name)
}

func (ga *GraphAdmin) PolicyClasses() set.Set {
	return ga.graph.PolicyClasses()
}
func (ga *GraphAdmin) Nodes() set.Set {
	return ga.graph.Nodes()
}

func (ga *GraphAdmin) Node(name string) (*graph.Node, error) {
	if !ga.Exists(name) {
		return nil, fmt.Errorf("node %s could not be found", name)
	}

	return ga.graph.Node(name)
}
func (ga *GraphAdmin) NodeFromDetails(t graph.NodeType, properties graph.PropertyMap) (*graph.Node, error) {
	return ga.graph.NodeFromDetails(t, properties)
}
func (ga *GraphAdmin) Search(t graph.NodeType, properties graph.PropertyMap) set.Set {
	return ga.graph.Search(t, properties)
}

func (ga *GraphAdmin) Children(name string) set.Set {
	return ga.graph.Children(name)
}

func (ga *GraphAdmin) Parents(name string) set.Set {
	return ga.graph.Parents(name)
}

func (ga *GraphAdmin) Assign(child, parent string) (err error) {
	if !ga.Exists(child) {
		return fmt.Errorf("child node %s does not exist", child)
	} else if !ga.Exists(parent) {
		return fmt.Errorf("parent node %s does not exist", parent)
	}

	//check if the assignment is valid
	var childNode, parentNode *graph.Node
	if childNode, err = ga.Node(child); err != nil {
		return
	}
	if parentNode, err = ga.Node(parent); err != nil {
		return
	}
	if err = graph.CheckAssignment(childNode.Type, parentNode.Type); err != nil {
		return
	}

	if parentNode.Type == graph.PC {
		parent = ga.PolicyClassDefault(parentNode.Name, childNode.Type)
	}

	return ga.graph.Assign(child, parent)
}

func (ga *GraphAdmin) Deassign(child, parent string) error {
	if !ga.Exists(child) {
		return fmt.Errorf("child node %s could not be found when deassigning", child)
	} else if !ga.Exists(parent) {
		return fmt.Errorf("parent node %s could not be found when deassigning", parent)
	}

	return ga.graph.Deassign(child, parent)
}

func (ga *GraphAdmin) IsAssigned(child, parent string) bool {
	return ga.graph.IsAssigned(child, parent)
}

func (ga *GraphAdmin) Associate(ua, target string, operations operations.OperationSet) error {
	if !ga.Exists(ua) {
		return fmt.Errorf("node %s could not be found when creating an association", ua)
	} else if !ga.Exists(target) {
		return fmt.Errorf("node %s could not be found when creating an association", target)
	}

	return ga.graph.Associate(ua, target, operations)
}

func (ga *GraphAdmin) Dissociate(ua, target string) error {
	if !ga.Exists(ua) {
		return fmt.Errorf("node %s could not be found when deleting an association", ua)
	} else if !ga.Exists(target) {
		return fmt.Errorf("node %s could not be found when deleting an association", target)
	}

	return ga.graph.Dissociate(ua, target)
}

func (ga *GraphAdmin) SourceAssociations(source string) (map[string]operations.OperationSet, error) {
	if !ga.Exists(source) {
		return nil, fmt.Errorf("node %s could not be found", source)
	}

	return ga.graph.SourceAssociations(source)
}

func (ga *GraphAdmin) TargetAssociations(target string) (map[string]operations.OperationSet, error) {
	if !ga.Exists(target) {
		return nil, fmt.Errorf("node %s could not be found", target)
	}

	return ga.graph.TargetAssociations(target)
}
