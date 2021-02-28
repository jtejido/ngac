package ngac

import (
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/pkg/context"
    "github.com/jtejido/ngac/pkg/epp"
    "github.com/jtejido/ngac/pkg/operations"
    "github.com/jtejido/ngac/pkg/pap"
    . "github.com/jtejido/ngac/pkg/pdp"
    "github.com/jtejido/ngac/pkg/pdp/audit"
    "github.com/jtejido/ngac/pkg/pdp/decider"
    "github.com/jtejido/ngac/pkg/pip"
    "github.com/jtejido/ngac/pkg/pip/graph"
    "github.com/jtejido/ngac/pkg/pip/obligations"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
    "testing"
)

type testContext struct {
    pdp                   *PDP
    u1, ua1, o1, oa1, pc1 *graph.Node
}

func testCtx(t *testing.T) testContext {
    ops := operations.NewOperationSet("read", "write", "execute")
    functionalEntity := pip.NewPIP(graph.NewMemGraph(), prohibitions.NewMemProhibitions(), obligations.NewMemObligations())
    p, err := pap.NewPAP(functionalEntity)
    if err != nil {
        t.Fatalf("%s", err)
    }
    pdp := NewPDP(
        p,
        epp.NewEPPOptions(),
        decider.NewPReviewDeciderWithProhibitions(functionalEntity.Graph(), functionalEntity.Prohibitions(), ops),
        audit.NewPReviewAuditor(functionalEntity.Graph(), ops),
    )
    ctx, _ := context.NewUserContext("super")
    g := pdp.WithUser(ctx).Graph()
    pc1, err := g.CreatePolicyClass("pc1", nil)
    if err != nil {
        t.Fatalf("%s", err)
    }
    oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
    if err != nil {
        t.Fatalf("%s", err)
    }
    o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
    if err != nil {
        t.Fatalf("%s", err)
    }
    ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
    if err != nil {
        t.Fatalf("%s", err)
    }
    u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
    if err != nil {
        t.Fatalf("%s", err)
    }
    err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
    if err != nil {
        t.Fatalf("%s", err)
    }

    return testContext{pdp, u1, ua1, o1, oa1, pc1}
}

func TestChildOfAssignExecutor(t *testing.T) {
    tctx := testCtx(t)
    ctx, _ := context.NewUserContext(tctx.u1.Name)
    executor := new(epp.ChildOfAssignExecutor)
    eventContext := epp.NewAssignEvent(ctx, tctx.o1, tctx.oa1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), nil)

    superUser, _ := context.NewUserContext("super")
    node, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(), pdp.WithUser(superUser).Obligations(), eventContext, function, epp.NewFunctionEvaluator())

    if err != nil {
        t.Fatalf("%s", err)
    }

    if node == nil {
        t.Errorf("node should not be nil")
    }

    if tctx.o1 != node.(*graph.Node) {
        t.Errorf("o1 should not be same as node")
    }

}

func TestCreateNodeExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.CreateNodeExecutor)
    var eventContext epp.EventContext
    pdp := tctx.pdp
    function := obligations.NewFunction(
        executor.Name(),
        []*obligations.Arg{
            obligations.NewArg("oa1"),
            obligations.NewArg("OA"),
            obligations.NewArg("testNode"),
            obligations.NewArg("OA"),
            obligations.NewArgFromFunction(obligations.NewFunction("to_props", []*obligations.Arg{obligations.NewArg("k=v")})),
        },
    )

    superUser, _ := context.NewUserContext("super")
    n, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }

    if n == nil {
        t.Errorf("node should not be nil")
    }

    if "testNode" != n.(*graph.Node).Name {
        t.Errorf("node name should be testNode")
    }

    if graph.OA != n.(*graph.Node).Type {
        t.Errorf("node type should be OA")
    }
    if n.(*graph.Node).Properties == nil {
        t.Errorf("properties should not be nil")
    }
    v, ok := n.(*graph.Node).Properties["k"]
    if !ok || v != "v" {
        t.Errorf("v should be present on properties")
    }
}

func TestCurrentProcessExecutorTest(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.CurrentProcessExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignEvent(ctx, tctx.o1, tctx.oa1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), nil)

    superUser, _ := context.NewUserContext("super")
    result, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }

    if result == nil {
        t.Errorf("result should not be nil")
    }
    if result.(string) != "1234" {
        t.Errorf("result should 1234")
    }
}

func TestCurrentTargetExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.CurrentTargetExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignEvent(ctx, tctx.o1, tctx.oa1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), nil)

    superUser, _ := context.NewUserContext("super")
    target, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if target == nil {
        t.Errorf("target should not be nil")
    }
    if tctx.o1 != target.(*graph.Node) {
        t.Errorf("target should be same as o1")
    }
}

func TestCurrentUserExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.CurrentUserExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignEvent(ctx, tctx.o1, tctx.oa1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), nil)

    superUser, _ := context.NewUserContext("super")
    userNode, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if userNode == nil {
        t.Errorf("userNode should not be nil")
    }
    if tctx.u1.Name != userNode.(*graph.Node).Name {
        t.Errorf("userNode name should be same as u1")
    }
}

func TestGetChildrenExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.GetChildrenExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignToEvent(ctx, tctx.oa1, tctx.o1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), []*obligations.Arg{obligations.NewArg("oa1"), obligations.NewArg("OA")})

    superUser, _ := context.NewUserContext("super")
    children, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if children == nil {
        t.Errorf("children should not be nil")
    }
    if !children.(set.Set).Contains(tctx.o1.Name) {
        t.Errorf("children should contain o1")
    }
}

func TestGetNodeExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.GetNodeExecutor)

    var eventContext epp.EventContext
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), []*obligations.Arg{obligations.NewArg("oa1"), obligations.NewArg("OA")})

    superUser, _ := context.NewUserContext("super")
    node, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if node == nil {
        t.Errorf("node should not be nil")
    }
    if !tctx.oa1.Equals(node.(*graph.Node)) {
        t.Errorf("oa1 should be same as node")
    }
}

func TestGetNodeNameExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.GetNodeNameExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignToEvent(ctx, tctx.oa1, tctx.o1)
    pdp := tctx.pdp

    function := obligations.NewFunction(executor.Name(), []*obligations.Arg{obligations.NewArgFromFunction(obligations.NewFunction("get_node", []*obligations.Arg{obligations.NewArg("oa1"), obligations.NewArg("OA")}))})

    superUser, _ := context.NewUserContext("super")
    name, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if name.(string) == "" {
        t.Errorf("name should not be empty")
    }
    if name.(string) != "oa1" {
        t.Errorf("name should be oa1")
    }
}

func TestIsNodeContainedInExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.IsNodeContainedInExecutor)
    ctx, _ := context.NewUserContextWithProcess(tctx.u1.Name, "1234")
    eventContext := epp.NewAssignToEvent(ctx, tctx.oa1, tctx.o1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(),
        []*obligations.Arg{
            obligations.NewArgFromFunction(
                obligations.NewFunction("get_node", []*obligations.Arg{obligations.NewArg("o1"), obligations.NewArg("O")}),
            ),
            obligations.NewArgFromFunction(
                obligations.NewFunction("get_node", []*obligations.Arg{obligations.NewArg("oa1"), obligations.NewArg("OA")}),
            ),
        })
    superUser, _ := context.NewUserContext("super")
    isContained, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if !isContained.(bool) {
        t.Errorf("isContained should be true")
    }

    function = obligations.NewFunction(executor.Name(),
        []*obligations.Arg{
            obligations.NewArgFromFunction(
                obligations.NewFunction("get_node", []*obligations.Arg{obligations.NewArg("u1"), obligations.NewArg("U")}),
            ),
            obligations.NewArgFromFunction(
                obligations.NewFunction("get_node", []*obligations.Arg{obligations.NewArg("oa1"), obligations.NewArg("OA")}),
            ),
        })
    isContained, err = executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if isContained.(bool) {
        t.Errorf("isContained should be false")
    }
}

func TestParentOfAssignExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.ParentOfAssignExecutor)
    ctx, _ := context.NewUserContext(tctx.u1.Name)
    eventContext := epp.NewAssignEvent(ctx, tctx.o1, tctx.oa1)
    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), nil)

    superUser, _ := context.NewUserContext("super")
    node, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if node == nil {
        t.Errorf("node should not be nil")
    }

    if node.(*graph.Node) != tctx.oa1 {
        t.Errorf("node should be equal to oa1")
    }
}

func TestToPropertiesExecutor(t *testing.T) {
    tctx := testCtx(t)
    executor := new(epp.ToPropertiesExecutor)
    var eventContext epp.EventContext

    pdp := tctx.pdp
    function := obligations.NewFunction(executor.Name(), []*obligations.Arg{obligations.NewArg("k=v"), obligations.NewArg("k1=v1"), obligations.NewArg("k2=v2")})

    superUser, _ := context.NewUserContext("super")
    props, err := executor.Exec(pdp.WithUser(superUser).Graph(), pdp.WithUser(superUser).Prohibitions(),
        pdp.WithUser(superUser).Obligations(),
        eventContext, function, epp.NewFunctionEvaluator())
    if err != nil {
        t.Fatalf("%s", err)
    }
    if props == nil {
        t.Errorf("props should not be nil")
    }
    if len(props.(graph.PropertyMap)) != 3 {
        t.Errorf("props size should  be 3")
    }
    v, ok := props.(graph.PropertyMap)["k"]
    if !ok || v != "v" {
        t.Errorf("k should be present and should equals v")
    }
    v, ok = props.(graph.PropertyMap)["k1"]
    if !ok || v != "v1" {
        t.Errorf("k1 should be present and should equals v1")
    }
    v, ok = props.(graph.PropertyMap)["k2"]
    if !ok || v != "v2" {
        t.Errorf("k2 should be present and should equals v2")
    }
}
