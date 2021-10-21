package ngac

import (
    "ngac/pkg/context"
    "ngac/pkg/operations"
    "ngac/pkg/pap"
    . "ngac/pkg/pdp"
    "ngac/pkg/pdp/audit"
    "ngac/pkg/pdp/decider"
    "ngac/pkg/pip"
    "ngac/pkg/pip/graph"
    gm "ngac/pkg/pip/graph/memory"
    obm "ngac/pkg/pip/obligations/memory"
    pm "ngac/pkg/pip/prohibitions/memory"
    "testing"
)

func TestPolicyClassReps(t *testing.T) {
    var g graph.Graph
    g = gm.New()
    mp := pm.New()
    ops := operations.NewOperationSet("read", "write", "execute")
    functionalEntity := pip.NewPIP(g, mp, obm.New())
    p, err := pap.NewPAP(functionalEntity)
    if err != nil {
        t.Fatalf("%s", err)
    }
    pdp := NewPDP(
        p,
        nil,
        decider.NewPReviewDeciderWithProhibitions(g, mp, ops),
        audit.NewPReviewAuditor(g, ops))
    ctx, _ := context.NewUserContext("super")
    g = pdp.WithUser(ctx).Graph()

    test, err := g.CreatePolicyClass("test", nil)
    if err != nil {
        t.Fatalf("%s", err)
    }
    defUA, ok := test.Properties["default_ua"]
    if !ok {
        t.Fatalf("default_ua should be present")
    }
    defOA, ok := test.Properties["default_oa"]
    if !ok {
        t.Fatalf("default_oa should be present")
    }
    repProp, ok := test.Properties[graph.REP_PROPERTY]
    if !ok {
        t.Fatalf("%s should be present", graph.REP_PROPERTY)
    }

    if !g.Exists(defUA) {
        t.Errorf("default_ua should exist")
    }
    if !g.Exists(defOA) {
        t.Errorf("default_oa should exist")
    }
    if !g.Exists(repProp) {
        t.Errorf("%s should exist", graph.REP_PROPERTY)
    }
}
