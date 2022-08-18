package guard

import (
	"fmt"
	"ngac/pkg/common"
	"ngac/pkg/context"
	"ngac/pkg/operations"
	"ngac/pkg/pap/policy"
	"ngac/pkg/pdp/decider"
	"ngac/pkg/pip/graph"
)

type Guard struct {
	pap         common.PolicyStore
	decider     decider.Decider
	resourceOps operations.OperationSet
}

func assertUserCtx(userCtx context.Context) {
	if userCtx == nil {
		panic("no user context provided to the PDP")
	}
}

func (g *Guard) hasPermissions(userCtx context.Context, target string, permissions ...string) (bool, error) {
	// assert that the user context is not null
	assertUserCtx(userCtx)

	// if checking the permissions on a PC, check the permissions on the rep node for the PC
	node, err := g.pap.Graph().Node(target)
	if err != nil {
		return false, err
	}
	if node.Type == graph.PC {
		t, found := node.Properties[graph.REP_PROPERTY]
		if !found {
			return false, fmt.Errorf("unable to check permissions for policy class %s, rep property not set", node.Name)
		}

		target = t
	}

	// check for permissions
	allowed := g.decider.List(userCtx.User(), userCtx.Process(), target)
	if len(permissions) == 0 {
		return allowed.Len() > 0, nil
	}

	for _, perm := range permissions {
		if !allowed.Contains(perm) {
			return false, nil
		}
	}

	return true, nil
}

func (g *Guard) CheckReset(userCtx context.Context) error {
	// check that the user can reset the graph
	ok, err := g.hasPermissions(userCtx, policy.SUPER_PC_REP, operations.RESET)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unauthorized permissions to reset the prohibitions")
	}

	return nil
}

func (g *Guard) ResourceOps() operations.OperationSet {
	return g.resourceOps
}

func (g *Guard) SetResourceOps(resourceOps operations.OperationSet) {
	g.resourceOps = resourceOps
	if d, ok := g.decider.(*decider.PReviewDecider); ok {
		d.ResourceOps = resourceOps
	}
}
