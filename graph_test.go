package ngac

import (
    "github.com/jtejido/ngac/audit"
    "github.com/jtejido/ngac/context"
    "github.com/jtejido/ngac/decider"
    "github.com/jtejido/ngac/operations"
    "github.com/jtejido/ngac/pap"
    "github.com/jtejido/ngac/pip"
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
    "testing"
)

func TestPolicyClassReps(t *testing.T) {
    var g graph.Graph
    g = graph.NewMemGraph()
    mp := prohibitions.NewMemProhibitions()
    ops := operations.NewOperationSet("read", "write", "execute")
    functionalEntity := pip.NewPIP(g, mp, obligations.NewMemObligations())
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
    defUA, ok := test.Properties.Get("default_ua")
    if !ok {
        t.Fatalf("default_ua should be present")
    }
    defOA, ok := test.Properties.Get("default_oa")
    if !ok {
        t.Fatalf("default_oa should be present")
    }
    repProp, ok := test.Properties.Get(graph.REP_PROPERTY)
    if !ok {
        t.Fatalf("%s should be present", graph.REP_PROPERTY)
    }

    if !g.Exists(defUA.(string)) {
        t.Errorf("default_ua should exist")
    }
    if !g.Exists(defOA.(string)) {
        t.Errorf("default_oa should exist")
    }
    if !g.Exists(repProp.(string)) {
        t.Errorf("%s should exist", graph.REP_PROPERTY)
    }
}
