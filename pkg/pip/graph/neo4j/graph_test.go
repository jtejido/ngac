package neo4j

import (
	"fmt"
	"ngac/pkg/operations"
	g "ngac/pkg/pip/graph"
	"os"
	"testing"
)

var ng *graph = nil

func TestCreateNode(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	pc, err := ng.CreatePolicyClass("pc", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err.Error())
	}

	node, err := ng.CreateNode("oa", g.OA, g.ToProperties(g.PropertyPair{"namespace", "test"}), pc.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}

	node, _ = ng.Node(node.Name)

	if node.Name != "oa" {
		t.Fatalf("failed to lookup node")
	}

	if node.Type != g.OA {
		t.Fatalf("failed to lookup type")
	}

}

func TestUpdateNode(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	node, err := ng.CreatePolicyClass("node", g.ToProperties(g.PropertyPair{"namespace", "test"}))
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}

	if err := ng.UpdateNode("newNodeName", nil); err == nil {
		t.Fatalf("failed to catch an error for non-existing node update")
	}

	if err := ng.UpdateNode("node", g.ToProperties(g.PropertyPair{"newKey", "newValue"})); err != nil {
		t.Fatalf("failed to update node")
	}

	n, _ := ng.Node(node.Name)
	if v, _ := n.Properties["newKey"]; v != "newValue" {
		t.Fatalf("failed to update properties")
	}
}

func TestRemoveNode(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	node, err := ng.CreatePolicyClass("node", g.ToProperties(g.PropertyPair{"namespace", "test"}))
	if err != nil {
		t.Fatalf("failed to create node: %s", err.Error())
	}
	ng.RemoveNode(node.Name)

	if ng.Exists(node.Name) {
		t.Fatalf("node should not exist after deleting")
	}

	if ng.PolicyClasses().Contains(node.Name) {
		t.Fatalf("node should not have policy after deletion of policy node")
	}
}

func TestPolicies(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	ng.CreatePolicyClass("node1", nil)
	ng.CreatePolicyClass("node2", nil)
	ng.CreatePolicyClass("node3", nil)

	if ng.PolicyClasses().Len() != 3 {
		t.Fatalf("node should not have 3 policies")
	}
}

func TestChildren(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parentNode, _ := ng.CreatePolicyClass("parent", nil)

	child1Node, _ := ng.CreateNode("child1", g.OA, nil, "parent")
	child2Node, _ := ng.CreateNode("child2", g.OA, nil, "parent")

	children := ng.Children(parentNode.Name)

	if !children.Contains(child1Node.Name, child2Node.Name) {
		t.Fatalf("failed to lookup child 1 or 2")
	}
}

func TestParents(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := ng.CreatePolicyClass("parent1", nil)
	parent2Node, _ := ng.CreateNode("parent2", g.OA, nil, "parent1")
	child1Node, _ := ng.CreateNode("child1", g.OA, nil, "parent1", "parent2")

	parents := ng.Parents(child1Node.Name)

	if !parents.Contains(parent1Node.Name, parent2Node.Name) {
		t.Fatalf("failed to lookup parent 1 or 2")
	}
}

func TestAssign(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := ng.CreatePolicyClass("parent1", nil)
	child1Node, _ := ng.CreateNode("child1", g.OA, nil, "parent1")
	child2Node, _ := ng.CreateNode("child2", g.OA, nil, "parent1")

	if err := ng.Assign("1241124", "123442141"); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := ng.Assign("1", "12341234"); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := ng.Assign(child1Node.Name, child2Node.Name); err != nil {
		t.Fatalf("failed to assign %s", err.Error())
	}

	if !ng.Children(parent1Node.Name).Contains(child1Node.Name) {
		t.Fatalf("failed to lookup child 1")
	}

	if !ng.Parents(child1Node.Name).Contains(parent1Node.Name) {
		t.Fatalf("failed to lookup parent")
	}
}

func TestDeassign(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	parent1Node, _ := ng.CreatePolicyClass("parent1", nil)
	child1Node, _ := ng.CreateNode("child1", g.OA, nil, "parent1")

	if err := ng.Assign("", ""); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := ng.Assign(child1Node.Name, ""); err == nil {
		t.Fatalf("should not assign non existing node ids")
	}

	if err := ng.Deassign(child1Node.Name, parent1Node.Name); err != nil {
		t.Fatalf("failed to deassign %s", err.Error())
	}

	if ng.Children(parent1Node.Name).Contains(child1Node.Name) {
		t.Fatalf("still able lookup child")
	}

	if ng.Parents(child1Node.Name).Contains(parent1Node.Name) {
		t.Fatalf("still able lookup parent")
	}

}

func TestAssociate(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	ng.CreatePolicyClass("pc", nil)
	uaNode, _ := ng.CreateNode("subject", g.UA, nil, "pc")
	targetNode, _ := ng.CreateNode("target", g.OA, nil, "pc")

	err := ng.Associate(uaNode.Name, targetNode.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("error thrown at associate %s", err.Error())
	}

	associations, err := ng.SourceAssociations(uaNode.Name)
	if err != nil {
		t.Fatalf("error thrown at getting source associations")
	}

	if _, ok := associations[targetNode.Name]; !ok {
		t.Fatalf("failed to get association for id: %s", targetNode.Name)
	}

	if !associations[targetNode.Name].Contains("read", "write") {
		t.Fatalf("failed to get right associations for source:  read/write")
	}

	associations, err = ng.TargetAssociations(targetNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting target associations")
	}

	if _, ok := associations[uaNode.Name]; !ok {
		t.Fatalf("failed to get association for id: %s", uaNode.Name)
	}

	if !associations[uaNode.Name].Contains("read", "write") {
		t.Fatalf("failed to get right associations for target:  read/write")
	}

	_, err = ng.CreateNode("test", g.UA, nil, "subject")
	if err != nil {
		t.Fatalf("error thrown at create node %s", err.Error())
	}
	err = ng.Associate("test", "subject", operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("error thrown at associate %s", err.Error())
	}
	associations, err = ng.SourceAssociations("test")
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
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	ng.CreatePolicyClass("pc", nil)
	uaNode, _ := ng.CreateNode("subject", g.UA, nil, "pc")
	targetNode, _ := ng.CreateNode("target", g.OA, nil, "pc")

	ng.Associate(uaNode.Name, targetNode.Name, operations.NewOperationSet("read", "write"))
	ng.Dissociate(uaNode.Name, targetNode.Name)

	associations, err := ng.SourceAssociations(uaNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting source associations")
	}

	if _, ok := associations[targetNode.Name]; ok {
		t.Fatalf("able to get association for target id: %s", targetNode.Name)
	}

	associations, err = ng.TargetAssociations(targetNode.Name)

	if err != nil {
		t.Fatalf("error thrown at getting target associations")
	}

	if _, ok := associations[uaNode.Name]; ok {
		t.Fatalf("able to get association for source id: %s", uaNode.Name)
	}
}

func TestSearch(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	ng.CreatePolicyClass("pc", nil)
	ng.CreateNode("oa1", g.OA, g.ToProperties(g.PropertyPair{"namespace", "test"}), "pc")
	ng.CreateNode("oa2", g.OA, g.ToProperties(g.PropertyPair{"key1", "value1"}), "pc")
	ng.CreateNode("oa3", g.OA, g.ToProperties(g.PropertyPair{"key1", "value1"}, g.PropertyPair{"key2", "value2"}), "pc")

	// name and type no properties
	nodes := ng.Search(g.OA, nil)
	if nodes.Len() != 3 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// one property
	nodes = ng.Search(-1, g.ToProperties(g.PropertyPair{"key1", "value1"}))
	if nodes.Len() != 2 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// just namespace
	nodes = ng.Search(-1, g.ToProperties(g.PropertyPair{"namespace", "test"}))

	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	// name, type, namespace
	nodes = ng.Search(g.OA, g.ToProperties(g.PropertyPair{"namespace", "test"}))

	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	nodes = ng.Search(g.OA, g.ToProperties(g.PropertyPair{"namespace", "test"}))
	if nodes.Len() != 1 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}

	nodes = ng.Search(g.OA, nil)
	if nodes.Len() != 3 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
	nodes = ng.Search(g.OA, g.ToProperties(g.PropertyPair{"key1", "value1"}))
	if nodes.Len() != 2 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
	nodes = ng.Search(-1, nil)
	if nodes.Len() != 4 {
		t.Fatalf("incorrect length after search: %d", nodes.Len())
	}
}

func TestNodes(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}

	ng.CreatePolicyClass("pc", nil)
	ng.CreateNode("node1", g.OA, nil, "pc")
	ng.CreateNode("node2", g.OA, nil, "pc")
	ng.CreateNode("node3", g.OA, nil, "pc")
	// name and type no properties

	if ng.Nodes().Len() != 4 {
		t.Fatalf("incorrect length : %d", ng.Nodes().Len())
	}
}

func TestNode(t *testing.T) {
	if err := ng.reset(); err != nil {
		t.Fatalf("failed to reset graph: %s", err.Error())
	}
	_, err := ng.Node("123")

	if err == nil {
		t.Fatalf("no node expected")
	}

	node, _ := ng.CreatePolicyClass("pc", nil)

	// name and type no properties
	n, _ := ng.Node(node.Name)
	if n.Name != "pc" {
		t.Fatalf("incorrect node name")
	}

	if n.Type != g.PC {
		t.Fatalf("incorrect node type")
	}
}

func TestMain(m *testing.M) {
	gg, err := New(`test_config.yaml`)
	ng = gg.(*graph)
	if err != nil {
		panic(fmt.Sprintf("failed to create graph: %s", err.Error()))
	}

	if err = ng.Start(); err != nil {
		panic(fmt.Sprintf("failed to start graph: %s", err.Error()))
	}

	defer ng.Close()

	os.Exit(m.Run())
}
