package ngac

import (
    "fmt"
    "github.com/jtejido/ngac/audit"
    "github.com/jtejido/ngac/common"
    "github.com/jtejido/ngac/context"
    "github.com/jtejido/ngac/decider"
    "github.com/jtejido/ngac/epp"
    "github.com/jtejido/ngac/guard"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/operations"
    "github.com/jtejido/ngac/pap/policy"
    "github.com/jtejido/ngac/pip/graph"
)

var _ graph.Graph = &Graph{}

type Graph struct {
    Service
    guard *guard.Graph
}

func NewGraphService(userCtx context.Context, p common.FunctionalEntity, e *EPP, d decider.Decider, a audit.Auditor) *Graph {
    ans := new(Graph)
    ans.userCtx = userCtx
    ans.pap = p
    ans.epp = e
    ans.decider = d
    ans.auditor = a
    ans.guard = guard.NewGraphGuard(p, d)
    return ans
}

func (g *Graph) CreatePolicyClass(name string, properties graph.PropertyMap) (*graph.Node, error) {
    // check user has permission to create a policy class
    g.guard.CheckCreatePolicyClass(g.userCtx)

    // create and return the new policy class
    return g.GraphAdmin().CreatePolicyClass(name, properties)
}

/**
 * Create a node and assign it to the provided parent(s). The name and type must not be null.
 * This method is needed because if a node is created without an initial assignment, it will be impossible
 * to assign the node in the future since no user will have permissions on a node not connected to the graph.
 * In this method we can check the user has the permission to assign to the given parent node and ignore if
 * the user can assign the newly created node.
 *
 * When creating a policy class, a parent node is not required.  The user must have the "create policy class" permission
 * on the super object.  By default the super user will always have this permission. A configuration will be created
 * that grants the user permissions on the policy class' default UA and OA, which will allow the user to delegate admin
 * permissions to other users.
 *
 * @param name the name of the node to create.
 * @param type the type of the node.
 * @param properties properties to add to the node.
 * @param initialParent the name of the node to assign the new node to.
 * @param additionalParents 0 or more node names to assign the new node to.
 * @return the new node.
 */
func (g *Graph) CreateNode(name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (node *graph.Node, err error) {
    // check that the user can create the node in each of the parents
    g.guard.CheckCreateNode(g.userCtx, t, initialParent, additionalParents)

    //create the node
    node, err = g.GraphAdmin().CreateNode(name, t, properties, initialParent, additionalParents...)

    // process the event
    g.epp.ProcessEvent(epp.NewCreateNodeEvent(g.userCtx, node, initialParent, additionalParents...))

    return
}

/**
 * Update the node in the database and in the in-memory graph.  If the name is null or empty it is ignored, likewise
 * for properties.
 *
 * @param name the name to give the node.
 * @param properties the properties of the node.
 * @throws PMException if the given node does not exist in the graph.
 * @throws PMAuthorizationException if the user is not authorized to update the node.
 */
func (g *Graph) UpdateNode(name string, properties graph.PropertyMap) error {
    // check that the user can update the node
    g.guard.CheckUpdateNode(g.userCtx, name)

    //update node in the PAP
    return g.GraphAdmin().UpdateNode(name, properties)
}

/**
 * Delete the node with the given name from the PAP.  First check that the current user
 * has the correct permissions to do so. Do this by checking that the user has the permission to deassign from each
 * of the node's parents, and that the user can delete the node.
 * @param name the name of the node to delete.
 * @throws PMException if there is an error accessing the graph through the PAP.
 * @throws PMAuthorizationException if the user is not authorized to delete the node.
 */
func (g *Graph) RemoveNode(name string) {
    node, err := g.GraphAdmin().Node(name)
    if err != nil {
        panic(err)
    }

    // check that the user can delete the node
    g.guard.CheckDeleteNode(g.userCtx, node.Type, name)

    // check that the node does not have any children
    if g.GraphAdmin().Children(name).Len() > 0 {
        panic(fmt.Sprintf("cannot delete %s, nodes are still assigned to it", name))
    }

    // if it's a PC, delete the rep
    if node.Type == graph.PC {
        if v, ok := node.Properties.Get(graph.REP_PROPERTY); ok {
            g.GraphAdmin().RemoveNode(v.(string))
        }
    }

    // get the
    parents := g.GraphAdmin().Parents(name)

    // delete the node
    g.GraphAdmin().RemoveNode(name)

    // process the delete node event
    err = g.epp.ProcessEvent(epp.NewDeleteNodeEvent(g.userCtx, node, parents))
    if err != nil {
        panic(err)
    }
    /*        // process the delete event
              Set<String> parents = graph.getParents(name);
              for(String parent : parents) {
                  Node parentNode = graph.getNode(parent);

                  getEPP().processEvent(new DeassignEvent(userCtx, node, parentNode));
                  getEPP().processEvent(new DeassignFromEvent(userCtx, parentNode, node));
              }

              // delete the node
              graph.deleteNode(name);*/
}

/**
 * Check that a node with the given name exists. This method will return false if the user does not have access to
 * the node.
 * @param name the name of the node to check for.
 * @return true if a node with the given name exists, false otherwise.
 * @throws PMException if there is an error checking if the node exists in the graph through the PAP.
 */
func (g *Graph) Exists(name string) bool {
    exists := g.GraphAdmin().Exists(name)
    if !exists {
        return false
    }
    ok, err := g.guard.CheckExists(g.userCtx, name)
    if err != nil {
        return false
    }
    return ok
}

/**
 * Retrieve the list of all nodes in the graph.  Go to the database to do this, since it is more likely to have
 * all of the node information.
 * @return the set of all nodes in the graph.
 * @throws PMException if there is an error getting the nodes from the PAP.
 */
func (g *Graph) Nodes() set.Set {
    nodes := set.NewSet()
    nodes.AddFrom(g.GraphAdmin().Nodes())
    g.guard.FilterNodes(g.userCtx, nodes)
    res := set.NewSet()
    res.AddFrom(nodes)
    return res
}

/**
 * Get the set of policy classes. This can be performed by the in-memory graph.
 * @return the set of names for the policy classes in the graph.
 * @throws PMException if there is an error getting the policy classes from the PAP.
 */
func (g *Graph) PolicyClasses() set.Set {
    return g.GraphAdmin().PolicyClasses()
}

func (g *Graph) PolicyClassDefault(pc string, t graph.NodeType) string {
    return pc + "_default_" + t.String()
}

/**
 * Get the children of the node from the graph.  Get the children from the database to ensure all node information
 * is present.  Before returning the set of nodes, filter out any nodes that the user has no permissions on.
 *
 * @param name the name of the node to get the children of.
 * @return a set of Node objects, representing the children of the target node.
 * @throws PMException if the target node does not exist.
 * @throws PMException if there is an error getting the children from the PAP.

 */
func (g *Graph) Children(name string) set.Set {
    if !g.Exists(name) {
        panic(fmt.Sprintf("node %s could not be found", name))
    }

    children := g.GraphAdmin().Children(name)
    g.guard.Filter(g.userCtx, children)
    return children
}

/**
 * Get the parents of the node from the graph.  Before returning the set of nodes, filter out any nodes that the user
 * has no permissions on.
 * @param name the name of the node to get the parents of.
 * @return a set of Node objects, representing the parents of the target node.
 * @throws PMException if the target node does not exist.
 * @throws PMException if there is an error getting the parents from the PAP.
 */
func (g *Graph) Parents(name string) set.Set {
    if !g.Exists(name) {
        panic(fmt.Sprintf("node %s could not be found", name))
    }

    parents := g.GraphAdmin().Parents(name)
    g.guard.Filter(g.userCtx, parents)
    return parents
}

/**
 * Create the assignment in both the db and in-memory graphs. First check that the user is allowed to assign the child,
 * and allowed to assign something to the parent.
 * @param child the name of the child node.
 * @param parent the name of the parent node.
 * @throws IllegalArgumentException if the child name is null.
 * @throws IllegalArgumentException if the parent name is null.
 * @throws PMException if the child or parent node does not exist.
 * @throws PMException if the assignment is invalid.
 * @throws PMAuthorizationException if the current user does not have permission to create the assignment.
 */
func (g *Graph) Assign(child, parent string) (err error) {
    // check that the user can make the assignment
    g.guard.CheckAssign(g.userCtx, child, parent)

    // assign in the PAP
    if err = g.GraphAdmin().Assign(child, parent); err != nil {
        return
    }

    // process the assignment as to events - assign and assign to
    var childNode, parentNode *graph.Node
    if childNode, err = g.Node(child); err != nil {
        return
    }
    if parentNode, err = g.Node(parent); err != nil {
        return
    }
    if err = g.epp.ProcessEvent(epp.NewAssignEvent(g.userCtx, childNode, parentNode)); err != nil {
        return
    }

    return g.epp.ProcessEvent(epp.NewAssignToEvent(g.userCtx, parentNode, childNode))
}

/**
 * Create the assignment in both the db and in-memory graphs. First check that the user is allowed to assign the child,
 * and allowed to assign something to the parent.
 * @param child the name of the child of the assignment to delete.
 * @param parent the name of the parent of the assignment to delete.
 * @throws IllegalArgumentException if the child name is null.
 * @throws IllegalArgumentException if the parent name is null.
 * @throws PMException if the child or parent node does not exist.
 * @throws PMAuthorizationException if the current user does not have permission to delete the assignment.
 */
func (g *Graph) Deassign(child, parent string) (err error) {
    // check the user can delete the assignment
    g.guard.CheckDeassign(g.userCtx, child, parent)

    //delete assignment in PAP
    if err = g.GraphAdmin().Deassign(child, parent); err != nil {
        return
    }

    // process the deassign as two events - deassign and deassign from
    var childNode, parentNode *graph.Node
    if childNode, err = g.Node(child); err != nil {
        return
    }
    if parentNode, err = g.Node(parent); err != nil {
        return
    }
    if err = g.epp.ProcessEvent(epp.NewDeassignEvent(g.userCtx, childNode, parentNode)); err != nil {
        return
    }

    return g.epp.ProcessEvent(epp.NewDeassignFromEvent(g.userCtx, parentNode, childNode))
}

func (g *Graph) IsAssigned(child, parent string) bool {
    var childNode, parentNode *graph.Node
    var err error
    if parentNode, err = g.Node(parent); err != nil {
        panic(err)
    }
    if childNode, err = g.Node(child); err != nil {
        panic(err)
    }

    return g.GraphAdmin().IsAssigned(childNode.Name, parentNode.Name)
}

/**
 * Create an association between the user attribute and the target node with the given operations. First, check that
 * the user has the permissions to associate the user attribute and target nodes.  If an association already exists
 * between the two nodes than update the existing association with the provided operations (overwrite).
 *
 * @param ua the name of the user attribute.
 * @param target the name of the target node.
 * @param operations a Set of operations to add to the Association.
 * @throws IllegalArgumentException if the user attribute is null.
 * @throws IllegalArgumentException if the target is null.
 * @throws PMException if the user attribute node does not exist.
 * @throws PMException if the target node does not exist.
 * @throws PMException if the association is invalid.
 * @throws PMAuthorizationException if the current user does not have permission to create the association.
 */
func (g *Graph) Associate(ua, target string, ops operations.OperationSet) (err error) {
    // check that this user can create the association
    g.guard.CheckAssociate(g.userCtx, ua, target)

    //create association in PAP
    if err = g.GraphAdmin().Associate(ua, target, ops); err != nil {
        return
    }

    var n, t *graph.Node
    if n, err = g.GraphAdmin().Node(ua); err != nil {
        return
    }

    if t, err = g.GraphAdmin().Node(target); err != nil {
        return
    }
    return g.epp.ProcessEvent(epp.NewAssociationEvent(g.userCtx, n, t))
}

/**
 * Delete the association between the user attribute and the target node.  First, check that the user has the
 * permission to delete the association.
 *
 * @param ua The name of the user attribute.
 * @param target The name of the target node.
 * @throws IllegalArgumentException If the user attribute is null.
 * @throws IllegalArgumentException If the target node is null.
 * @throws PMException If the user attribute node does not exist.
 * @throws PMException If the target node does not exist.
 * @throws PMAuthorizationException If the current user does not have permission to delete the association.
 */
func (g *Graph) Dissociate(ua, target string) (err error) {
    // check that the user can delete the association
    g.guard.CheckDissociate(g.userCtx, ua, target)

    //create association in PAP
    if err = g.GraphAdmin().Dissociate(ua, target); err != nil {
        return
    }
    var n, t *graph.Node
    if n, err = g.GraphAdmin().Node(ua); err != nil {
        return
    }
    if t, err = g.GraphAdmin().Node(target); err != nil {
        return
    }
    return g.epp.ProcessEvent(epp.NewDeleteAssociationEvent(g.userCtx, n, t))
}

/**
 * Get the associations the given node is the source node of. First, check if the user is allowed to retrieve this
 * information.
 *
 * @param source The name of the source node.
 * @return a map of the target and operations for each association the given node is the source of.
 * @throws PMException If the given node does not exist.
 * @throws PMAuthorizationException If the current user does not have permission to get hte node's associations.
 */
func (g *Graph) SourceAssociations(source string) (sourceAssociations map[string]operations.OperationSet, err error) {
    // check that this user can get the associations of the source node
    g.guard.CheckGetAssociations(g.userCtx, source)

    // get the associations for the source node
    sourceAssociations, err = g.GraphAdmin().SourceAssociations(source)

    // filter out any associations in which the user does not have access to the target attribute
    g.guard.FilterMap(g.userCtx, sourceAssociations)

    return
}

/**
 * Get the associations the given node is the target node of. First, check if the user is allowed to retrieve this
 * information.
 *
 * @param target The name of the source node.
 * @return a map of the source name and operations for each association the given node is the target of.
 * @throws PMException If the given node does not exist.
 * @throws PMAuthorizationException If the current user does not have permission to get hte node's associations.
 */
func (g *Graph) TargetAssociations(target string) (targetAssociations map[string]operations.OperationSet, err error) {
    // check that this user can get the associations of the target node
    g.guard.CheckGetAssociations(g.userCtx, target)

    // get the associations for the target node
    targetAssociations, err = g.GraphAdmin().TargetAssociations(target)

    // filter out any associations in which the user does not have access to the source attribute
    g.guard.FilterMap(g.userCtx, targetAssociations)

    return
}

/**
 * Search the NGAC graph for nodes that match the given parameters. A node must match all non null parameters to be
 * returned in the search.
 *
 * @param type The type of the nodes to search for.
 * @param properties The properties of the nodes to search for.
 * @return a Response with the nodes that match the given search criteria.
 * @throws PMException If the PAP encounters an error with the graph.
 * @throws PMAuthorizationException If the current user does not have permission to get hte node's associations.
 */
func (g *Graph) Search(t graph.NodeType, properties graph.PropertyMap) set.Set {
    search := g.GraphAdmin().Search(t, properties)
    //System.out.println(search);
    g.guard.FilterNodes(g.userCtx, search)
    return search
}

/**
 * Retrieve the node from the graph with the given name.
 *
 * @param name the name of the node to get.
 * @return the Node retrieved from the graph with the given name.
 * @throws PMException If the node does not exist in the graph.
 * @throws PMAuthorizationException if the current user is not authorized to access this node.
 */
func (g *Graph) Node(name string) (node *graph.Node, err error) {
    // get node
    node, err = g.GraphAdmin().Node(name)

    // check user has permissions on the node
    g.guard.CheckExists(g.userCtx, name)

    return
}

func (g *Graph) NodeFromDetails(t graph.NodeType, properties graph.PropertyMap) (node *graph.Node, err error) {
    node, err = g.GraphAdmin().NodeFromDetails(t, properties)

    // check user has permissions on the node
    g.guard.CheckExists(g.userCtx, node.Name)

    return
}

/**
 * Deletes all nodes in the graph
 *
 * @throws PMException if something goes wrong in the deletion process
 */
func (g *Graph) Reset(userCtx context.Context) (err error) {
    g.guard.CheckReset(userCtx)

    nodes := g.GraphAdmin().Nodes()
    names := set.NewSet()
    prohibitions_name := set.NewSet()
    for n := range nodes.Iter() {
        node := n.(*graph.Node)
        names.Add(node.Name)
        if node.Type == graph.UA || node.Type == graph.OA {
            ta, err := g.GraphAdmin().TargetAssociations(node.Name)
            if err != nil {
                return err
            }
            for el, _ := range ta {
                if err = g.GraphAdmin().Dissociate(el, node.Name); err != nil {
                    return err
                }
            }
            if node.Type == graph.UA {
                sa, err := g.GraphAdmin().SourceAssociations(node.Name)
                if err != nil {
                    return err
                }
                for el, _ := range sa {
                    if err = g.GraphAdmin().Dissociate(node.Name, el); err != nil {
                        return err
                    }
                }
            }
        }
        ch := g.GraphAdmin().Children(node.Name)
        for el := range ch.Iter() {
            if g.GraphAdmin().IsAssigned(el.(string), node.Name) {
                if err = g.GraphAdmin().Deassign(el.(string), node.Name); err != nil {
                    return
                }
            }
        }
        parents := g.GraphAdmin().Parents(node.Name)
        for el := range parents.Iter() {
            if g.GraphAdmin().IsAssigned(node.Name, el.(string)) {
                if err = g.GraphAdmin().Deassign(node.Name, el.(string)); err != nil {
                    return
                }
            }
        }

        pf := g.ProhibitionsAdmin().ProhibitionsFor(node.Name)
        for _, el := range pf {
            prohibitions_name.Add(el.Name)
        }
    }
    for prohibition := range prohibitions_name.Iter() {
        g.ProhibitionsAdmin().Remove(prohibition.(string))
    }

    for name := range names.Iter() {
        g.GraphAdmin().RemoveNode(name.(string))
    }
    //setup Super policy in GraphAdmin + copy graph and graph copy
    g.superPolicy = policy.NewSuperPolicy()
    return g.superPolicy.Configure(g.GraphAdmin())
}
