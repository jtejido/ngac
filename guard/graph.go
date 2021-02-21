package guard

import (
    "fmt"
    "github.com/jtejido/ngac/common"
    "github.com/jtejido/ngac/context"
    "github.com/jtejido/ngac/decider"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/operations"
    "github.com/jtejido/ngac/pap/policy"
    "github.com/jtejido/ngac/pip/graph"
)

type Graph struct {
    Guard
}

func NewGraphGuard(p common.FunctionalEntity, d decider.Decider) *Graph {
    ans := new(Graph)
    ans.pap = p
    ans.decider = d
    return ans
}

func (g *Graph) CheckCreatePolicyClass(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := g.hasPermissions(userCtx, policy.SUPER_PC_REP, operations.CREATE_POLICY_CLASS)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to create a policy class")
    }

    return nil
}

func (g *Graph) CheckCreateNode(userCtx context.Context, nodeType graph.NodeType, initialParent string, additionalParents []string) error {
    var op string
    switch nodeType {
    case graph.OA:
        op = operations.CREATE_OBJECT_ATTRIBUTE
        break
    case graph.UA:
        op = operations.CREATE_USER_ATTRIBUTE
        break
    case graph.O:
        op = operations.CREATE_OBJECT
        break
    case graph.U:
        op = operations.CREATE_USER
        break
    default:
        op = operations.CREATE_POLICY_CLASS
    }
    ok, err := g.hasPermissions(userCtx, initialParent, op)
    if err != nil {
        return err
    }
    // check that the user has the permission to assign to the parent node
    if !ok {
        // if the user cannot assign to the parent node, delete the newly created node
        return fmt.Errorf("unauthorized permission \"%s\" on node %s", op, initialParent)
    }

    // check any additional parents
    for _, parent := range additionalParents {
        ok, err := g.hasPermissions(userCtx, parent, op)
        if err != nil {
            return err
        }
        if !ok {
            // if the user cannot assign to the parent node, delete the newly created node
            return fmt.Errorf("unauthorized permission \"%s\" on %s", op, parent)
        }
    }

    return nil
}

func (g *Graph) CheckUpdateNode(userCtx context.Context, name string) error {
    // check that the user can update the node
    ok, err := g.hasPermissions(userCtx, name, operations.UPDATE_NODE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permission %s on node %s", operations.UPDATE_NODE, name)
    }

    return nil
}

func (g *Graph) CheckDeleteNode(userCtx context.Context, nodeType graph.NodeType, node string) error {
    // check that the user can delete a policy class if that is the type
    if nodeType == graph.PC {
        ok, err := g.hasPermissions(userCtx, policy.SUPER_PC_REP, operations.DELETE_POLICY_CLASS)
        if err != nil {
            return err
        }
        if !ok {
            return fmt.Errorf("unauthorized permissions to delete a policy class")
        } else {
            return nil
        }
    }

    var op string
    switch nodeType {
    case graph.OA:
        op = operations.DELETE_OBJECT_ATTRIBUTE
        break
    case graph.UA:
        op = operations.DELETE_USER_ATTRIBUTE
        break
    case graph.O:
        op = operations.DELETE_OBJECT
        break
    case graph.U:
        op = operations.DELETE_USER
        break
    default:
        op = operations.DELETE_POLICY_CLASS
    }

    // check the user can delete the node
    ok, err := g.hasPermissions(userCtx, node, operations.DELETE_NODE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", node, operations.DELETE_NODE)
    }

    // check that the user can delete the node from the node's parents
    parents := g.pap.Graph().Parents(node)
    for parent := range parents.Iter() {
        ok, err := g.hasPermissions(userCtx, parent.(string), op)
        if err != nil {
            return err
        }
        if !ok {
            return fmt.Errorf("unauthorized permissions on %s: %s", parent.(string), op)
        }
    }
    return nil
}

func (g *Graph) CheckExists(userCtx context.Context, name string) (bool, error) {
    // a user only needs one permission to know a node exists
    // however receiving an unauthorized exception would let the user know it exists
    // therefore, false is returned if they don't have permissions on the node
    return g.hasPermissions(userCtx, name)
}

func (g *Graph) Filter(userCtx context.Context, nodes set.Set) {
    nodes.Filter(func(node interface{}) bool {
        ok, err := g.hasPermissions(userCtx, node.(string))
        if err != nil {
            return true
        }

        return !ok
    })
}

func (g *Graph) FilterNodes(userCtx context.Context, nodes set.Set) {
    nodes.Filter(func(node interface{}) bool {
        ok, err := g.hasPermissions(userCtx, node.(*graph.Node).Name)
        if err != nil {
            return true
        }

        return !ok
    })
}

func (g *Graph) FilterMap(userCtx context.Context, m map[string]operations.OperationSet) {
    for key := range m {
        ok, err := g.hasPermissions(userCtx, key)
        if err != nil {
            delete(m, key)
        }
        if !ok {
            delete(m, key)
        }
    }
}

func (g *Graph) CheckAssign(userCtx context.Context, child, parent string) error {
    //check the user can assign the child
    ok, err := g.hasPermissions(userCtx, child, operations.ASSIGN)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permission %s on node %s", operations.ASSIGN, child)
    }
    ok, err = g.hasPermissions(userCtx, parent, operations.ASSIGN_TO)
    if err != nil {
        return err
    }
    // check that the user can assign to the parent node
    if !ok {
        return fmt.Errorf("unauthorized permission %s on node %s", operations.ASSIGN_TO, parent)
    }

    return nil
}

func (g *Graph) CheckDeassign(userCtx context.Context, child, parent string) error {
    //check the user can deassign the child
    ok, err := g.hasPermissions(userCtx, child, operations.DEASSIGN)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", child, operations.DEASSIGN)
    }
    ok, err = g.hasPermissions(userCtx, parent, operations.DEASSIGN_FROM)
    if err != nil {
        return err
    }
    //check that the user can deassign from the parent
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", parent, operations.DEASSIGN_FROM)
    }

    return nil
}

func (g *Graph) CheckAssociate(userCtx context.Context, ua, target string) error {
    //check the user can associate the source and target nodes
    ok, err := g.hasPermissions(userCtx, ua, operations.ASSOCIATE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", ua, operations.ASSOCIATE)
    }
    ok, err = g.hasPermissions(userCtx, target, operations.ASSOCIATE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", target, operations.ASSOCIATE)
    }

    return nil
}

func (g *Graph) CheckDissociate(userCtx context.Context, ua, target string) error {
    //check the user can associate the source and target nodes
    ok, err := g.hasPermissions(userCtx, ua, operations.DISASSOCIATE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", ua, operations.DISASSOCIATE)
    }
    ok, err = g.hasPermissions(userCtx, target, operations.DISASSOCIATE)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", target, operations.DISASSOCIATE)
    }

    return nil
}

func (g *Graph) CheckGetAssociations(userCtx context.Context, node string) error {
    //check the user can get the associations of the source node
    ok, err := g.hasPermissions(userCtx, node, operations.GET_ASSOCIATIONS)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions on %s: %s", node, operations.GET_ASSOCIATIONS)
    }

    return nil
}
