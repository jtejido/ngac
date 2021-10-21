package guard

import (
    "fmt"
    "ngac/pkg/common"
    "ngac/pkg/context"
    "ngac/pkg/operations"
    "ngac/pkg/pdp/decider"
    "ngac/pkg/pip/prohibitions"
)

type Prohibitions struct {
    Guard
}

func NewProhibitionsGuard(p common.FunctionalEntity, d decider.Decider) *Prohibitions {
    ans := new(Prohibitions)
    ans.pap = p
    ans.decider = d
    return ans
}

func (p *Prohibitions) check(userCtx context.Context, prohibition *prohibitions.Prohibition, permission string) error {
    subject := prohibition.Subject
    containers := prohibition.Containers()

    // check prohibition subject
    if p.pap.Graph().Exists(subject) {
        ok, err := p.hasPermissions(userCtx, subject, permission)
        if err != nil {
            return err
        }
        if !ok {
            return fmt.Errorf("unauthorized permission %s on %s", permission, subject)
        }
    }

    // check each container in prohibition
    for container := range containers {
        ok, err := p.hasPermissions(userCtx, container, permission)
        if err != nil {
            return err
        }
        if !ok {
            return fmt.Errorf("unauthorized permission %s on %s", permission, container)
        }
    }

    return nil
}

func (p *Prohibitions) CheckAdd(userCtx context.Context, prohibition *prohibitions.Prohibition) error {
    return p.check(userCtx, prohibition, operations.CREATE_PROHIBITION)
}

func (p *Prohibitions) CheckGet(userCtx context.Context, prohibition *prohibitions.Prohibition) error {
    return p.check(userCtx, prohibition, operations.VIEW_PROHIBITION)
}

func (p *Prohibitions) CheckUpdate(userCtx context.Context, prohibition *prohibitions.Prohibition) error {
    return p.check(userCtx, prohibition, operations.UPDATE_PROHIBITION)
}

func (p *Prohibitions) CheckDelete(userCtx context.Context, prohibition *prohibitions.Prohibition) error {
    return p.check(userCtx, prohibition, operations.DELETE_PROHIBITION)
}

func (p *Prohibitions) Filter(userCtx context.Context, prohibs []*prohibitions.Prohibition) {
    for i := 0; i < len(prohibs); i++ {
        if err := p.CheckGet(userCtx, prohibs[i]); err != nil {
            prohibs = append(prohibs[:i], prohibs[i+1:]...)
        }
    }
}
