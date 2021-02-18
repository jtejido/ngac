package graph

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/operations"
)

// Interface for maintaining an NGAC graph.
// Every Dao for a graph implementation should follow this
type Graph interface {

	/**
	 * Create a policy class in the graph.
	 */
	CreatePolicyClass(name string, properties PropertyMap) (*Node, error)

	/**
	 * Create a new node with the given name, type and properties and add it to the graph. Node names must be unique.
	 *
	 */
	CreateNode(name string, t NodeType, properties PropertyMap, initialParent string, additionalParents ...string) (*Node, error)

	/**
	 * Update the properties of the node with the given name. The given properties overwrite any existing properties.
	 */
	UpdateNode(name string, properties PropertyMap) error

	/**
	 * Delete the node with the given name from the graph.
	 */
	RemoveNode(name string)

	/**
	 * Check that a node with the given name exists in the graph.
	 */
	Exists(name string) bool

	/**
	 * Get the set of policy classes.  This operation is run every time a decision is made, so a separate
	 * method is needed to improve efficiency. The returned set is just the names of each policy class.
	 */
	PolicyClasses() set.Set

	/**
	 * Retrieve the set of all nodes in the graph.
	 */
	Nodes() set.Set

	/**
	 * Retrieve the node with the given name.
	 */
	Node(name string) (*Node, error)

	/**
	 * Search the graph for a node that matches the given parameters. A node must
	 * contain all properties provided to be returned.
	 * To get a node that has a specific property key with any value use "*" as the value in the parameter.
	 * (i.e. {key=*})
	 * If more than one node matches the criteria, only one will be returned.
	 */
	NodeFromDetails(t NodeType, properties PropertyMap) (*Node, error)

	/**
	 * Search the graph for nodes matching the given parameters. A node must
	 * contain all properties provided to be returned.
	 * To get all the nodes that have a specific property key with any value use "*" as the value in the parameter.
	 * (i.e. {key=*})
	 */
	Search(t NodeType, properties PropertyMap) set.Set

	/**
	 * Get the set of nodes that are assigned to the node with the given name.
	 */
	Children(name string) set.Set

	/**
	 * Get the set of nodes that the node with the given name is assigned to.
	 */
	Parents(name string) set.Set

	/**
	 * Assign the child node to the parent node. The child and parent nodes must both already exist in the graph,
	 * and the types must make a valid assignment. An example of a valid assignment is assigning o1, an object, to oa1,
	 * an object attribute.  o1 is the child (objects can never be the parent in an assignment), and oa1 is the parent.
	 */
	Assign(child, parent string) error

	/**
	 * Remove the Assignment between the child and parent nodes.
	 */
	Deassign(child, parent string) error

	/**
	 * Returns true if the child is assigned to the parent.
	 */
	IsAssigned(child, parent string) bool

	/**
	 * Create an Association between the user attribute and the Target node with the provided operations. If an association
	 * already exists between these two nodes, overwrite the existing operations with the ones provided.  Associations
	 * can only begin at a user attribute but can point to either an Object or user attribute
	 */
	Associate(ua, target string, operations operations.OperationSet) error

	/**
	 * Delete the Association between the user attribute and Target node.
	 */
	Dissociate(ua, target string) error

	/**
	 * Retrieve the associations the given node is the source of.  The source node of an association is always a
	 * user attribute and this method will throw an exception if an invalid node is provided.  The returned Map will
	 * contain the target and operations of each association.
	 */
	SourceAssociations(source string) (map[string]operations.OperationSet, error)

	/**
	 * Retrieve the associations the given node is the target of.  The target node can be an Object Attribute or a User
	 * Attribute. This method will throw an exception if a node of any other type is provided.  The returned Map will
	 * contain the source node names and the operations of each association.
	 */
	TargetAssociations(target string) (map[string]operations.OperationSet, error)
}
