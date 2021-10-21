package neo4j

import (
	"fmt"
	"ngac/pkg/operations"
	"ngac/pkg/pip/graph"
	"os"
	"testing"
)

var g *Graph = nil

func TestCreateNode(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	pc, err := g.CreatePolicyClass("pc", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err.Error())
	}

	node, err := g.CreateNode("oa", graph.OA, graph.ToProperties(graph.PropertyPair{"namespace", "test"}), pc.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}

	node, _ = g.Node(node.Name)

	if node.Name != "oa" {
		t.Fatalf("failed to lookup node")
	}

	if node.Type != graph.OA {
		t.Fatalf("failed to lookup type")
	}

}

func TestUpdateNode(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	node, err := g.CreatePolicyClass("node", graph.ToProperties(graph.PropertyPair{"namespace", "test"}))
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}

	if err := g.UpdateNode("newNodeName", nil); err == nil {
		t.Fatalf("failed to catch an error for non-existing node update")
	}

	if err := g.UpdateNode("node", graph.ToProperties(graph.PropertyPair{"newKey", "newValue"})); err != nil {
		t.Fatalf("failed to update node")
	}

	n, _ := g.Node(node.Name)
	if v, _ := n.Properties["newKey"]; v != "newValue" {
		t.Fatalf("failed to update properties")
	}
}

func TestRemoveNode(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	node, err := g.CreatePolicyClass("node", graph.ToProperties(graph.PropertyPair{"namespace", "test"}))
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}
	g.RemoveNode(node.Name)

	if g.Exists(node.Name) {
		t.Fatalf("node should not exist after deleting")
	}

	if g.PolicyClasses().Contains(node.Name) {
		t.Fatalf("node should not have policy after deletion of policy node")
	}
}

func TestPolicies(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	g.CreatePolicyClass("node1", nil)
	g.CreatePolicyClass("node2", nil)
	g.CreatePolicyClass("node3", nil)

	if g.PolicyClasses().Len() != 3 {
		t.Fatalf("node should not have 3 policies")
	}
}

func TestChildren(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parentNode, _ := g.CreatePolicyClass("parent", nil)

	child1Node, _ := g.CreateNode("child1", graph.OA, nil, "parent")
	child2Node, _ := g.CreateNode("child2", graph.OA, nil, "parent")

	children := g.Children(parentNode.Name)

	if !children.Contains(child1Node.Name, child2Node.Name) {
		t.Fatalf("failed to lookup child 1 or 2")
	}
}

func TestParents(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := g.CreatePolicyClass("parent1", nil)
	parent2Node, _ := g.CreateNode("parent2", graph.OA, nil, "parent1")
	child1Node, _ := g.CreateNode("child1", graph.OA, nil, "parent1", "parent2")

	parents := g.Parents(child1Node.Name)

	if !parents.Contains(parent1Node.Name, parent2Node.Name) {
		t.Fatalf("failed to lookup parent 1 or 2")
	}
}

func TestAssign(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := g.CreatePolicyClass("parent1", nil)
	child1Node, _ := g.CreateNode("child1", graph.OA, nil, "parent1")
	child2Node, _ := g.CreateNode("child2", graph.OA, nil, "parent1")

	if err := g.Assign("1241124", "123442141"); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := g.Assign("1", "12341234"); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := g.Assign(child1Node.Name, child2Node.Name); err != nil {
		t.Fatalf("failed to assign %s", err.Error())
	}

	if !g.Children(parent1Node.Name).Contains(child1Node.Name) {
		t.Fatalf("failed to lookup child 1")
	}

	if !g.Parents(child1Node.Name).Contains(parent1Node.Name) {
		t.Fatalf("failed to lookup parent")
	}
}

func TestDeassign(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := g.CreatePolicyClass("parent1", nil)
	child1Node, _ := g.CreateNode("child1", graph.OA, nil, "parent1")

	if err := g.Assign("", ""); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := g.Assign(child1Node.Name, ""); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := g.Deassign(child1Node.Name, parent1Node.Name); err != nil {
		t.Fatalf("failed to deassign %s", err.Error())
	}

	if g.Children(parent1Node.Name).Contains(child1Node.Name) {
		t.Fatalf("still able lookup child")
	}

	if g.Parents(child1Node.Name).Contains(parent1Node.Name) {
		t.Fatalf("still able lookup parent")
	}

}

func TestAssociate(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	g.CreatePolicyClass("pc", nil)
	uaNode, _ := g.CreateNode("subject", graph.UA, nil, "pc")
	targetNode, _ := g.CreateNode("target", graph.OA, nil, "pc")

	err := g.Associate(uaNode.Name, targetNode.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("error thrown at associate %s", err.Error())
	}

	associations, err := g.SourceAssociations(uaNode.Name)
	if err != nil {
		t.Fatalf("error thrown at getting source associations")
	}

	if _, ok := associations[targetNode.Name]; !ok {
		t.Fatalf("failed to get association for id: %s", targetNode.Name)
	}

	if !associations[targetNode.Name].Contains("read", "write") {
		t.Fatalf("failed to get right associations for source:  read/write")
	}

	associations, err = g.TargetAssociations(targetNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting target associations")
	}

	if _, ok := associations[uaNode.Name]; !ok {
		t.Fatalf("failed to get association for id: %s", uaNode.Name)
	}

	if !associations[uaNode.Name].Contains("read", "write") {
		t.Fatalf("failed to get right associations for target:  read/write")
	}

	_, err = g.CreateNode("test", graph.UA, nil, "subject")
	if err != nil {
		t.Fatalf("error thrown at create node %s", err.Error())
	}
	err = g.Associate("test", "subject", operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("error thrown at associate %s", err.Error())
	}
	associations, err = g.SourceAssociations("test")
	if err != nil {
		t.Fatalf("error thrown at getting source associations")
	}

	if _, ok := associations["subject"]; !ok {
		t.Fatalf("failed to get association for id: subject")
	}

	if !associations["subject"].Contains("read") {
		t.Fatalf("failed to get right associations for source:  read")
	}

}

func TestDissociate(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	g.CreatePolicyClass("pc", nil)
	uaNode, _ := g.CreateNode("subject", graph.UA, nil, "pc")
	targetNode, _ := g.CreateNode("target", graph.OA, nil, "pc")

	g.Associate(uaNode.Name, targetNode.Name, operations.NewOperationSet("read", "write"))
	g.Dissociate(uaNode.Name, targetNode.Name)

	associations, err := g.SourceAssociations(uaNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting source associations")
	}

	if _, ok := associations[targetNode.Name]; ok {
		t.Fatalf("able to get association for target id: %s", targetNode.Name)
	}

	associations, err = g.TargetAssociations(targetNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting target associations")
	}

	if _, ok := associations[uaNode.Name]; ok {
		t.Fatalf("able to get association for source id: %s", uaNode.Name)
	}
}

func TestSearch(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	g.CreatePolicyClass("pc", nil)
	g.CreateNode("oa1", graph.OA, graph.ToProperties(graph.PropertyPair{"namespace", "test"}), "pc")
	g.CreateNode("oa2", graph.OA, graph.ToProperties(graph.PropertyPair{"key1", "value1"}), "pc")
	g.CreateNode("oa3", graph.OA, graph.ToProperties(graph.PropertyPair{"key1", "value1"}, graph.PropertyPair{"key2", "value2"}), "pc")

	// name and type no properties
	nodes := g.Search(graph.OA, nil)
	if nodes.Len() != 3 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// one property
	nodes = g.Search(-1, graph.ToProperties(graph.PropertyPair{"key1", "value1"}))
	if nodes.Len() != 2 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// just namespace
	nodes = g.Search(-1, graph.ToProperties(graph.PropertyPair{"namespace", "test"}))

	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// name, type, namespace
	nodes = g.Search(graph.OA, graph.ToProperties(graph.PropertyPair{"namespace", "test"}))

	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	nodes = g.Search(graph.OA, graph.ToProperties(graph.PropertyPair{"namespace", "test"}))
	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	nodes = g.Search(graph.OA, nil)
	if nodes.Len() != 3 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
	nodes = g.Search(graph.OA, graph.ToProperties(graph.PropertyPair{"key1", "value1"}))
	if nodes.Len() != 2 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
	nodes = g.Search(-1, nil)
	if nodes.Len() != 4 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
}

func TestNodes(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	g.CreatePolicyClass("pc", nil)
	g.CreateNode("node1", graph.OA, nil, "pc")
	g.CreateNode("node2", graph.OA, nil, "pc")
	g.CreateNode("node3", graph.OA, nil, "pc")
	// name and type no properties

	if g.Nodes().Len() != 4 {
		t.Fatalf("incorrect length : %d", g.Nodes().Len())
	}
}

func TestNode(t *testing.T) {
	if err := g.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	_, err := g.Node("123")

	if err == nil {
		t.Fatalf("no node expected")
	}

	node, _ := g.CreatePolicyClass("pc", nil)

	// name and type no properties
	n, _ := g.Node(node.Name)
	if n.Name != "pc" {
		t.Fatalf("incorrect node name")
	}

	if n.Type != graph.PC {
		t.Fatalf("incorrect node type")
	}
}

func TestMain(m *testing.M) {
	var err error
	g, err = New(`test_config.yaml`)
	if err != nil {
		panic(fmt.Sprintf("failed to create graph: %s", err.Error()))
	}

	if err = g.Start(); err != nil {
		panic(fmt.Sprintf("failed to start graph: %s", err.Error()))
	}

	defer g.Close()

	os.Exit(m.Run())
}
