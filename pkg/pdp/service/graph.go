package service

import (
    "fmt"
    "ngac/internal/set"
    "ngac/pkg/common"
    "ngac/pkg/context"
    "ngac/pkg/epp"
    "ngac/pkg/operations"
    "ngac/pkg/pap/policy"
    "ngac/pkg/pdp/audit"
    "ngac/pkg/pdp/decider"
    "ngac/pkg/pdp/service/guard"
    "ngac/pkg/pip/graph"
)

var _ graph.Graph = &Graph{}

type Graph struct {
    Service
    guard *guard.Graph
}

func NewGraphService(userCtx context.Context, p common.PolicyStore, e epp.EPP, d decider.Decider, a audit.Auditor) *Graph {
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
    if err := g.guard.CheckCreatePolicyClass(g.userCtx); err != nil {
        return nil, err
    }

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
 */
func (g *Graph) CreateNode(name string, t graph.NodeType, properties graph.PropertyMap, initialParent string, additionalParents ...string) (node *graph.Node, err error) {
    // check that the user can create the node in each of the parents
    if err := g.guard.CheckCreateNode(g.userCtx, t, initialParent, additionalParents); err != nil {
        return nil, err
    }

    //create the node
    if node, err = g.GraphAdmin().CreateNode(name, t, properties, initialParent, additionalParents...); err != nil {
        return
    }

    // process the event
    err = g.epp.ProcessEvent(epp.NewCreateNodeEvent(g.userCtx, node, initialParent, additionalParents...))
    return
}

/**
 * Update the node in the database and in the in-memory graph.  If the name is null or empty it is ignored, likewise
 * for properties.
 */
func (g *Graph) UpdateNode(name string, properties graph.PropertyMap) error {
    // check that the user can update the node
    if err := g.guard.CheckUpdateNode(g.userCtx, name); err != nil {
        return err
    }

    //update node in the PAP
    return g.GraphAdmin().UpdateNode(name, properties)
}

/**
 * Delete the node with the given name from the PAP.  First check that the current user
 * has the correct permissions to do so. Do this by checking that the user has the permission to deassign from each
 * of the node's parents, and that the user can delete the node.
 */
func (g *Graph) RemoveNode(name string) {
    node, err := g.GraphAdmin().Node(name)
    if err != nil {
        panic(err)
    }

    // check that the user can delete the node
    if err := g.guard.CheckDeleteNode(g.userCtx, node.Type, name); err != nil {
        panic(err)
    }

    // check that the node does not have any children
    if g.GraphAdmin().Children(name).Len() > 0 {
        panic(fmt.Sprintf("cannot delete %s, nodes are still assigned to it", name))
    }

    // if it's a PC, delete the rep
    if node.Type == graph.PC {
        if v, ok := node.Properties[graph.REP_PROPERTY]; ok {
            g.GraphAdmin().RemoveNode(v)
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

}

/**
 * Check that a node with the given name exists. This method will return false if the user does not have access to
 * the node.
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
 */
func (g *Graph) Assign(child, parent string) (err error) {
    // check that the user can make the assignment
    if err = g.guard.CheckAssign(g.userCtx, child, parent); err != nil {
        return
    }

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
 */
func (g *Graph) Deassign(child, parent string) (err error) {
    // check the user can delete the assignment
    if err = g.guard.CheckDeassign(g.userCtx, child, parent); err != nil {
        return
    }

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
 */
func (g *Graph) Associate(ua, target string, ops operations.OperationSet) (err error) {
    // check that this user can create the association
    if err = g.guard.CheckAssociate(g.userCtx, ua, target); err != nil {
        return
    }

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
 */
func (g *Graph) Dissociate(ua, target string) (err error) {
    // check that the user can delete the association
    if err = g.guard.CheckDissociate(g.userCtx, ua, target); err != nil {
        return
    }

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
 */
func (g *Graph) SourceAssociations(source string) (sourceAssociations map[string]operations.OperationSet, err error) {
    // check that this user can get the associations of the source node
    if err = g.guard.CheckGetAssociations(g.userCtx, source); err != nil {
        return
    }

    // get the associations for the source node
    sourceAssociations, err = g.GraphAdmin().SourceAssociations(source)

    // filter out any associations in which the user does not have access to the target attribute
    g.guard.FilterMap(g.userCtx, sourceAssociations)

    return
}

/**
 * Get the associations the given node is the target node of. First, check if the user is allowed to retrieve this
 * information.
 */
func (g *Graph) TargetAssociations(target string) (targetAssociations map[string]operations.OperationSet, err error) {
    // check that this user can get the associations of the target node
    if err = g.guard.CheckGetAssociations(g.userCtx, target); err != nil {
        return
    }

    // get the associations for the target node
    targetAssociations, err = g.GraphAdmin().TargetAssociations(target)

    // filter out any associations in which the user does not have access to the source attribute
    g.guard.FilterMap(g.userCtx, targetAssociations)

    return
}

/**
 * Search the NGAC graph for nodes that match the given parameters. A node must match all non null parameters to be
 * returned in the search.
 */
func (g *Graph) Search(t graph.NodeType, properties graph.PropertyMap) set.Set {
    search := g.GraphAdmin().Search(t, properties)
    //System.out.println(search);
    g.guard.FilterNodes(g.userCtx, search)
    return search
}

/**
 * Retrieve the node from the graph with the given name.
 */
func (g *Graph) Node(name string) (node *graph.Node, err error) {
    // get node
    if node, err = g.GraphAdmin().Node(name); err != nil {
        return
    }

    // check user has permissions on the node
    if _, err = g.guard.CheckExists(g.userCtx, name); err != nil {
        return
    }

    return
}

func (g *Graph) NodeFromDetails(t graph.NodeType, properties graph.PropertyMap) (node *graph.Node, err error) {
    node, err = g.GraphAdmin().NodeFromDetails(t, properties)

    // check user has permissions on the node
    if _, err = g.guard.CheckExists(g.userCtx, node.Name); err != nil {
        return
    }

    return
}

/**
 * Deletes all nodes in the graph
 */
func (g *Graph) Reset(userCtx context.Context) (err error) {
    if err = g.guard.CheckReset(userCtx); err != nil {
        return
    }

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
