package decider

import "ngac/internal/set"

type Decider interface {
	/**
	 * Check if the subject has the permissions on the target node. Use '*' as the permission to check
	 * if the subject has any permissions on the node.
	 */
	Check(subject, process, target string, perms ...interface{}) bool

	/**
	 * List the permissions that the subject has on the target node.
	 */
	List(subject, process, target string) set.Set

	/**
	 * Given a list of nodes filter out any nodes that the given subject does not have the given permissions on. To filter
	 * based on any permissions use Operations.ANY as the permission to check for.
	 */
	Filter(subject, process string, nodes set.Set, perms ...interface{}) set.Set

	/**
	 * Get the children of the target node that the subject has the given permissions on.
	 */
	Children(subject, process, target string, perms ...interface{}) set.Set

	/**
	 * Given a subject ID, return every node the subject has access to and the permissions they have on each.
	 */
	CapabilityList(subject, process string) map[string]set.Set

	/**
	 * Given an Object Attribute ID, returns the id of every user (long), and what permissions(Set<String>) it has on it
	 */
	GenerateACL(target, process string) map[string]set.Set
}
