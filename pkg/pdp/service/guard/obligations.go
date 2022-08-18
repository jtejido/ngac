package guard

import (
    "fmt"
    "ngac/pkg/common"
    "ngac/pkg/context"
    "ngac/pkg/operations"
    "ngac/pkg/pap/policy"
    "ngac/pkg/pdp/decider"
)

type Obligations struct {
    Guard
}

func NewObligationsGuard(p common.PolicyStore, d decider.Decider) *Obligations {
    ans := new(Obligations)
    ans.pap = p
    ans.decider = d
    return ans
}

func (o *Obligations) CheckAdd(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := o.hasPermissions(userCtx, policy.SUPER_OA, operations.ADD_OBLIGATION)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to add an obligation")
    }

    return nil
}

func (o *Obligations) CheckGet(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := o.hasPermissions(userCtx, policy.SUPER_OA, operations.GET_OBLIGATION)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to get obligation")
    }
    return nil
}

func (o *Obligations) CheckUpdate(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := o.hasPermissions(userCtx, policy.SUPER_OA, operations.UPDATE_OBLIGATION)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to update obligation")
    }
    return nil
}

func (o *Obligations) CheckDelete(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := o.hasPermissions(userCtx, policy.SUPER_OA, operations.DELETE_OBLIGATION)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to delete obligation")
    }
    return nil
}

func (o *Obligations) CheckEnable(userCtx context.Context) error {
    // check that the user can create a policy class
    ok, err := o.hasPermissions(userCtx, policy.SUPER_OA, operations.ENABLE_OBLIGATION)
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unauthorized permissions to enable obligation")
    }
    return nil
}
