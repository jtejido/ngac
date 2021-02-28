package rbac

import (
	"fmt"
	"github.com/jtejido/ngac/pkg/context"
	"github.com/jtejido/ngac/pkg/operations"
	"github.com/jtejido/ngac/pkg/pdp"
	"github.com/jtejido/ngac/pkg/pip/graph"
	"sync"
)

var mu *sync.Mutex

var (
	rbac_pc_name      = "RBAC"
	rbac_users_node   *graph.Node
	rbac_objects_node *graph.Node
	rbac_pc_node      *graph.Node
)

func Name() string {
	return rbac_pc_name
}

func UsersNode() *graph.Node {
	return rbac_users_node
}

func ObjectsNode() *graph.Node {
	return rbac_objects_node
}

func PCNode() *graph.Node {
	return rbac_pc_node
}

/**
 * Utilities for any operation relating to the RBAC NGAC concept
 */
type RBAC struct{}

func New() *RBAC {
	mu = new(sync.Mutex)
	return new(RBAC)
}

/**
 * This sets the RBAC PC for any of the methods in this class.
 * If the given PC already exists it will mark it as the RBAC PC,
 * otherwise it will create and mark it.
 *
 * This will likely be the first call in any method of this class.
 *
 * @param RBACname the name of the RBAC PC, null if you want to use the default name
 * @param pdp PDP of the existing graph
 * @param superUserContext UserContext of the super user
 * @return the RBAC PC
 * @throws PMException
 */
func (policy *RBAC) Configure(RBACname string, p *pdp.PDP, superUserContext context.Context) (rbac_pc_node *graph.Node, err error) {
	g := p.WithUser(superUserContext).Graph()

	// todo: on-boarding methods.

	if RBACname != "" {
		rbac_pc_name = RBACname
	}

	// DAC PC todo: find default pc node's properties
	rbac_pc_node, err = checkAndCreateRBACNode(g, rbac_pc_name, graph.PC)
	if err != nil {
		return
	}
	children := g.Children(rbac_pc_name)
	for child := range children.Iter() {
		childNode, err := g.Node(child.(string))
		if err != nil {
			return nil, err
		}
		if childNode.Type == graph.UA {
			rbac_users_node = childNode
		} else {
			rbac_objects_node = childNode
		}
	}

	return
}

// create user and assign role
func (policy *RBAC) CreateUserAndAssignRole(p *pdp.PDP, superUserContext context.Context, userName, roleName string) (*graph.Node, error) {
	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return nil, fmt.Errorf("Role must exist.")
	}

	// create user
	// choose a role and assign the user to that role
	user, err := g.CreateNode(userName, graph.U, nil, roleName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// assign role
func (policy *RBAC) AssignRole(p *pdp.PDP, superUserContext context.Context, userName, roleName string) error {
	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return fmt.Errorf("Role must exist.")
	}
	if len(userName) == 0 || !g.Exists(userName) {
		return fmt.Errorf("User must exist.")
	}

	// assign user to corresponding role node
	return g.Assign(userName, roleName)
}

// remove role
func (policy *RBAC) DeassignRole(p *pdp.PDP, superUserContext context.Context, userName, roleName string) error {
	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return fmt.Errorf("Role must exist.")
	}
	if len(userName) == 0 || !g.Exists(userName) {
		return fmt.Errorf("User must exist.")
	}

	// de-assign user from corresponding role node
	return g.Deassign(userName, roleName)
}

// create role
func (policy *RBAC) CreateRole(p *pdp.PDP, superUserContext context.Context, roleName string) (*graph.Node, error) {
	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return nil, fmt.Errorf("Role must exist.")
	}

	return g.CreateNode(roleName, graph.UA, nil, rbac_users_node.Name)
}

// create role
func (policy *RBAC) DeleteRole(p *pdp.PDP, superUserContext context.Context, roleName string) error {
	// find role node
	// todo: how to reassign children
	return nil
}

// get user roles
func (policy *RBAC) UserRoles(p *pdp.PDP, superUserContext context.Context, userName string) (s []string, err error) {
	g := p.WithUser(superUserContext).Graph()
	if len(userName) == 0 || !g.Exists(userName) {
		return s, fmt.Errorf("User must exist.")
	}

	// get parents in the RBAC PC
	parents := g.Parents(userName)
	s = make([]string, parents.Len())
	var i int
	for parent := range parents.Iter() {
		s[i] = parent.(string)
		i++
	}
	return
}

// set_role_permissions
func (policy *RBAC) SetRolePermissions(p *pdp.PDP, superUserContext context.Context, roleName string, ops operations.OperationSet, targetName string) error {

	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return fmt.Errorf("Role must exist.")
	}
	if len(targetName) == 0 || !g.Exists(targetName) {
		return fmt.Errorf("Target must exist.")
	}
	if ops == nil || ops.Len() == 0 {
		return fmt.Errorf("Ops must exist.")
	}

	// create association from role to give objects with given permissions
	return g.Associate(roleName, targetName, ops)
}

// get_role_permissions
func (policy *RBAC) RolePermissions(p *pdp.PDP, superUserContext context.Context, roleName string) (map[string]operations.OperationSet, error) {
	g := p.WithUser(superUserContext).Graph()
	if len(roleName) == 0 || !g.Exists(roleName) {
		return nil, fmt.Errorf("Role must exist.")
	}

	// get associations from role
	return g.SourceAssociations(roleName)
}

/********************
 * Helper Functions *
 ********************/

/**
 * Helper Method to check if a RBAC node exists, and other wise create it.
 * It will also set the corresponding property for that DAC node.
 *
 * This methods is specifically for RBAC nodes, and not meant to be used elsewhere
 */
type pair = graph.PropertyPair

func checkAndCreateRBACNode(g graph.Graph, name string, t graph.NodeType) (RBAC *graph.Node, err error) {

	if !g.Exists(name) {
		if t == graph.PC {
			return g.CreatePolicyClass(name, graph.ToProperties(pair{"ngac_type", "RBAC"}))
		} else {
			return g.CreateNode(name, t, graph.ToProperties(pair{"ngac_type", "RBAC"}), rbac_pc_name)
		}
	} else {
		RBAC, err = g.Node(name)
		if err != nil {
			return
		}

		// add ngac_type=RBAC to properties
		mu.Lock()
		defer mu.Unlock()
		typeValue, ok := RBAC.Properties["ngac_type"]
		if !ok {
			RBAC.Properties["ngac_type"] = "RBAC"
		} else if typeValue != "RBAC" {
			err = fmt.Errorf("Node cannot have property key of ngac_type")
			return
		}
	}
	return
}
