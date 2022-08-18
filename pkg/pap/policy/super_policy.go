package policy

import (
    "github.com/jtejido/ngac/pkg/operations"
    "github.com/jtejido/ngac/pkg/pip/graph"
)

const (
    SUPER_USER   = "super"
    SUPER_PC     = "super_pc"
    SUPER_PC_REP = "super_pc_rep"
    SUPER_UA1    = "super_ua1"
    SUPER_UA2    = "super_ua2"
    SUPER_OA     = "super_oa"
)

type pair = graph.PropertyPair

var (
    superUser = &graph.Node{"super", graph.U, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"})}
)

type SuperPolicy struct {
    superUser, superUA1, superUA2, superPolicyClassRep, superOA, superPC *graph.Node
}

func NewSuperPolicy() *SuperPolicy {
    return &SuperPolicy{superUser: superUser}
}

func (sp *SuperPolicy) SuperUser() *graph.Node {
    return sp.superUser
}

func (sp *SuperPolicy) SuperUserAttribute() *graph.Node {
    return sp.superUA1
}

func (sp *SuperPolicy) SuperUserAttribute2() *graph.Node {
    return sp.superUA2
}

func (sp *SuperPolicy) SuperPolicyClassRep() *graph.Node {
    return sp.superPolicyClassRep
}

func (sp *SuperPolicy) SuperObjectAttribute() *graph.Node {
    return sp.superOA
}

func (sp *SuperPolicy) SuperPolicyClass() *graph.Node {
    return sp.superPC
}

func (sp *SuperPolicy) Configure(gr graph.Graph) (err error) {
    superPCRep := SUPER_PC_REP
    if !gr.Exists("super_pc") {
        props := graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}, pair{graph.REP_PROPERTY, superPCRep})
        if sp.superPC, err = gr.CreatePolicyClass("super_pc", props); err != nil {
            return
        }
    } else {
        if sp.superPC, err = gr.Node("super_pc"); err != nil {
            return
        }
        sp.superPC.Properties[graph.REP_PROPERTY] = superPCRep
        if err = gr.UpdateNode(sp.superPC.Name, sp.superPC.Properties); err != nil {
            return
        }
    }

    if !gr.Exists("super_ua1") {
        if sp.superUA1, err = gr.CreateNode("super_ua1", graph.UA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}), sp.superPC.Name); err != nil {
            return
        }
    } else {
        if sp.superUA1, err = gr.Node("super_ua1"); err != nil {
            return
        }

    }

    if !gr.Exists("super_ua2") {
        if sp.superUA2, err = gr.CreateNode("super_ua2", graph.UA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}), sp.superPC.Name); err != nil {
            return
        }
    } else {
        if sp.superUA2, err = gr.Node("super_ua2"); err != nil {
            return
        }
    }

    if !gr.Exists("super") {
        if _, err = gr.CreateNode("super", graph.U, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}), sp.superUA1.Name, sp.superUA2.Name); err != nil {
            return
        }
    }

    if !gr.Exists("super_oa") {
        sp.superOA, err = gr.CreateNode("super_oa", graph.OA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}), sp.superPC.Name)
        if err != nil {
            return
        }
    } else {
        if sp.superOA, err = gr.Node("super_oa"); err != nil {
            return
        }
    }

    if !gr.Exists(superPCRep) {
        if sp.superPolicyClassRep, err = gr.CreateNode(superPCRep, graph.OA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, "super"}, pair{"pc", sp.superPC.Name}), sp.superOA.Name); err != nil {
            return
        }
    } else {
        if sp.superPolicyClassRep, err = gr.Node(superPCRep); err != nil {
            return
        }
    }

    // check super ua1 is assigned to super pc
    children := gr.Children(sp.superPC.Name)
    if !children.Contains(sp.superUA1.Name) {
        if !gr.Children("super_pc_default_UA").Contains(sp.superUA1.Name) {
            if err = gr.Assign(sp.superUA1.Name, sp.superPC.Name); err != nil {
                return
            }
        }
    }

    // check super ua2 is assigned to super pc
    if !children.Contains(sp.superUA2.Name) {
        if !gr.Children("super_pc_default_UA").Contains(sp.superUA2.Name) {
            if err = gr.Assign(sp.superUA2.Name, sp.superPC.Name); err != nil {
                return
            }
        }
    }
    // check super user is assigned to super ua1
    children = gr.Children(sp.superUA1.Name)
    if !children.Contains(sp.superUser.Name) {
        if err = gr.Assign(sp.superUser.Name, sp.superUA1.Name); err != nil {
            return
        }
    }
    // check super user is assigned to super ua2
    children = gr.Children(sp.superUA2.Name)
    if !children.Contains(sp.superUser.Name) {
        if err = gr.Assign(sp.superUser.Name, sp.superUA2.Name); err != nil {
            return
        }
    }
    // check super oa is assigned to super pc
    children = gr.Children(sp.superPC.Name)
    if !children.Contains(sp.superOA.Name) {
        if !gr.Children("super_pc_default_OA").Contains(sp.superOA.Name) {
            if err = gr.Assign(sp.superOA.Name, sp.superPC.Name); err != nil {
                return
            }
        }
    }
    // check super o is assigned to super oa
    children = gr.Children(sp.superOA.Name)
    if !children.Contains(sp.superPolicyClassRep.Name) {
        if err = gr.Assign(sp.superPolicyClassRep.Name, sp.superOA.Name); err != nil {
            return
        }
    }

    // associate super ua to super oa
    if err = gr.Associate(sp.superUA1.Name, sp.superOA.Name, operations.NewOperationSet(operations.ALL_OPS)); err != nil {
        return
    }
    if err = gr.Associate(sp.superUA2.Name, sp.superUA1.Name, operations.NewOperationSet(operations.ALL_OPS)); err != nil {
        return
    }

    if err = gr.Associate(sp.superUA1.Name, sp.superUA2.Name, operations.NewOperationSet(operations.ALL_OPS)); err != nil {
        return
    }

    return sp.configurePolicyClasses(gr)
}

func (sp *SuperPolicy) configurePolicyClasses(gr graph.Graph) error {
    policyClasses := gr.PolicyClasses()
    for p := range policyClasses.Iter() {
        // configure default nodes
        pc := p.(string)
        rep := pc + "_rep"
        defaultUA := pc + "_default_UA"
        defaultOA := pc + "_default_OA"

        if !gr.Exists(defaultOA) {
            if _, err := gr.CreateNode(defaultOA, graph.OA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, pc}), pc); err != nil {
                return err
            }
        }
        if !gr.Exists(defaultUA) {
            if _, err := gr.CreateNode(defaultUA, graph.UA, graph.ToProperties(pair{graph.NAMESPACE_PROPERTY, pc}), pc); err != nil {
                return err
            }
        }

        // update pc node if necessary
        node, err := gr.Node(pc)
        if err != nil {
            return err
        }
        props := node.Properties
        props["default_ua"] = defaultUA
        props["default_oa"] = defaultOA
        props[graph.REP_PROPERTY] = rep
        if err := gr.UpdateNode(pc, props); err != nil {
            return err
        }
        //remove potential parents of super uas
        if gr.Parents(sp.superUA1.Name).Contains("super_pc_default_UA") {
            if err := gr.Deassign(sp.superUA1.Name, "super_pc_default_UA"); err != nil {
                return err
            }
        }
        if gr.Parents(sp.superUA2.Name).Contains("super_pc_default_UA") {
            if err := gr.Deassign(sp.superUA2.Name, "super_pc_default_UA"); err != nil {
                return err
            }
        }
        // assign both super uas if not already
        if !gr.IsAssigned(sp.superUA1.Name, pc) {
            if err := gr.Assign(sp.superUA1.Name, pc); err != nil {
                return err
            }
        }
        if !gr.IsAssigned(sp.superUA2.Name, pc) {
            if err := gr.Assign(sp.superUA2.Name, pc); err != nil {
                return err
            }
        }

        // associate super ua 1 with pc default node
        if err := gr.Associate(sp.superUA1.Name, defaultUA, operations.NewOperationSet(operations.ALL_OPS)); err != nil {
            return err
        }
        if err := gr.Associate(sp.superUA1.Name, defaultOA, operations.NewOperationSet(operations.ALL_OPS)); err != nil {
            return err
        }

        // create the rep
        if !gr.Exists(rep) {
            if _, err := gr.CreateNode(rep, graph.OA, graph.ToProperties(pair{"pc", pc}), sp.superOA.Name); err != nil {
                return err
            }
        } else {
            // check that the rep is assigned to the super OA
            if !gr.IsAssigned(rep, sp.superOA.Name) {
                if err := gr.Assign(rep, sp.superOA.Name); err != nil {
                    return err
                }
            }
        }

    }

    return nil

}
