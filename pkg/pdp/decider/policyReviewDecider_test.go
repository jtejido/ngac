package decider

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pkg/operations"
	"github.com/jtejido/ngac/pkg/pip/graph"
	gm "github.com/jtejido/ngac/pkg/pip/graph/memory"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
	obm "github.com/jtejido/ngac/pkg/pip/prohibitions/memory"
	"testing"
)

var rwe = operations.NewOperationSet("read", "write", "execute")

func TestHasPermission(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o2", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o3", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write", "unknown-op"))
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)

	if !decider.Check(u1.Name, "", o1.Name, "read", "write") {
		t.Fatalf("failed to check permission from source to target node")
	}
}

func TestFilter(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o2, err := g.CreateNode("o2", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o3, err := g.CreateNode("o3", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	nodeIDs := set.NewSet(o1.Name, o2.Name, o3.Name, oa1.Name)

	decider := NewPReviewDecider(g, rwe)

	if !nodeIDs.Contains(decider.Filter(u1.Name, "", nodeIDs, "read").ToSlice()...) {
		t.Fatalf("failed to check filtered node set")
	}
}

func TestChildren(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o2, err := g.CreateNode("o2", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o3, err := g.CreateNode("o3", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)
	children := decider.Children(u1.Name, "", oa1.Name)
	nodeIDs := set.NewSet(o1.Name, o2.Name, o3.Name)

	if !nodeIDs.Contains(children.ToSlice()...) {
		t.Fatalf("failed to get children set")
	}
}

func TestAccessibleNodes(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o2, _ := g.CreateNode("o2", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o3, _ := g.CreateNode("o3", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)
	accessibleNodes := decider.CapabilityList(u1.Name, "")

	if v, found := accessibleNodes[oa1.Name]; !found {
		t.Fatalf("failed to get accessible nodes for %s", oa1.Name)
	} else {
		if set := operations.NewOperationSet("read", "write"); !set.Contains(v.ToSlice()...) {
			t.Fatalf("permissions expected to be read and write for %s", oa1.Name)
		}
	}

	if v, found := accessibleNodes[o1.Name]; !found {
		t.Fatalf("failed to get accessible nodes for %s", o1.Name)
	} else {
		if set := operations.NewOperationSet("read", "write"); !set.Contains(v.ToSlice()...) {
			t.Fatalf("permissions expected to be read and write for %s", o1.Name)
		}
	}

	if v, found := accessibleNodes[o2.Name]; !found {
		t.Fatalf("failed to get accessible nodes for %s", o2.Name)
	} else {
		if set := operations.NewOperationSet("read", "write"); !set.Contains(v.ToSlice()...) {
			t.Fatalf("permissions expected to be read and write for %s", o2.Name)
		}
	}

	if v, found := accessibleNodes[o3.Name]; !found {
		t.Fatalf("failed to get accessible nodes for %s", o3.Name)
	} else {
		if set := operations.NewOperationSet("read", "write"); !set.Contains(v.ToSlice()...) {
			t.Fatalf("permissions expected to be read and write for %s", o3.Name)
		}
	}

}

func TestGraph1(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to be read and write")
	}
}

func TestGraph2(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name, pc2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	ua2, err := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa2, err := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() != 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph3(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to be read and write")
	}
}

func TestGraph4(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	ua2, err := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to be read and write")
	}
}

func TestGraph5(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	ua2, err := g.CreateNode("ua2", graph.UA, nil, pc2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	u1, err := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa2, err := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	o1, err := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read", "write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read") {
		t.Fatalf("permissions expected to be read")
	}
}

func TestGraph6(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc2.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read") {
		t.Fatalf("permissions expected to be read")
	}
}

func TestGraph7(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph8(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))

	decider := NewPReviewDecider(g, rwe)

	l := decider.List(u1.Name, "", o1.Name)
	if !l.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !l.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestGraph9(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))
	g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	l := decider.List(u1.Name, "", o1.Name)
	if !l.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !l.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestGraph10(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc2.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}
}

func TestGraph11(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))

	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph12(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}
}

func TestGraph13(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, ua2.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, oa2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read"))

	decider := NewPReviewDecider(g, rwe)

	l := decider.List(u1.Name, "", o1.Name)
	if !l.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !l.Contains("read") {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestGraph14(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))
	g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("*"))

	decider := NewPReviewDecider(g, rwe)

	l := decider.List(u1.Name, "", o1.Name)
	if !l.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !l.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestGraph15(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, ua2.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, oa2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read"))

	decider := NewPReviewDecider(g, rwe)

	l := decider.List(u1.Name, "", o1.Name)
	if !l.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !l.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestGraph16(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, ua2.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}
}

func TestGraph18(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph19(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))

	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph20(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read") {
		t.Fatalf("permissions expected to contain read")
	}
}

func TestGraph21(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	pc2, err := g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name, ua2.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	g.Associate(ua2.Name, oa2.Name, operations.NewOperationSet("write"))

	decider := NewPReviewDecider(g, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph22(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	g.CreatePolicyClass("pc2", nil)
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}
}

func TestGraph23WithProhibition(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa3, _ := g.CreateNode("oa3", graph.OA, nil, pc1.Name)
	oa4, _ := g.CreateNode("oa4", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, oa3.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, oa4.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa3.Name, operations.NewOperationSet("read", "write", "execute"))

	prohibs := obm.New()
	prohibition := prohibitions.NewBuilder("deny", ua1.Name, operations.NewOperationSet("read"))
	prohibition.AddContainer(oa1.Name, false)
	prohibition.AddContainer(oa2.Name, false)
	prohibition.Intersection = true

	prohibs.Add(prohibition.Build())

	prohibition = prohibitions.NewBuilder("deny2", u1.Name, operations.NewOperationSet("write"))
	prohibition.Intersection = true
	prohibition.AddContainer(oa3.Name, false)

	prohibs.Add(prohibition.Build())

	decider := NewPReviewDeciderWithProhibitions(g, prohibs, rwe)
	list := decider.List(u1.Name, "", o1.Name)

	if list.Len() != 1 {
		t.Fatalf("incorrect list size")
	}

	if !list.Contains("execute") {
		t.Fatalf("permissions expected to contain execute")
	}
}

func TestGraph24WithProhibition(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)
	o2, _ := g.CreateNode("o2", graph.O, nil, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))

	prohibs := obm.New()
	prohibition := prohibitions.NewBuilder("deny", ua1.Name, operations.NewOperationSet("read"))
	prohibition.AddContainer(oa1.Name, false)
	prohibition.AddContainer(oa2.Name, true)
	prohibition.Intersection = true
	prohibs.Add(prohibition.Build())

	decider := NewPReviewDeciderWithProhibitions(g, prohibs, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read") {
		t.Fatalf("permissions expected to contain read")
	}

	if decider.List(u1.Name, "", o2.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}

	g.Associate(ua1.Name, oa2.Name, operations.NewOperationSet("read"))

	prohibition = prohibitions.NewBuilder("deny-process", "1234", operations.NewOperationSet("read"))
	prohibition.AddContainer(oa1.Name, false)
	prohibs.Add(prohibition.Build())

	if decider.List(u1.Name, "1234", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph25WithProhibition(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, oa1.Name)
	oa3, _ := g.CreateNode("oa3", graph.OA, nil, oa1.Name)
	oa4, _ := g.CreateNode("oa4", graph.OA, nil, oa3.Name)
	oa5, _ := g.CreateNode("oa5", graph.OA, nil, oa2.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa4.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	prohibs := obm.New()
	prohibition := prohibitions.NewBuilder("deny", ua1.Name, operations.NewOperationSet("read", "write"))
	prohibition.AddContainer(oa4.Name, true)
	prohibition.AddContainer(oa1.Name, false)
	prohibition.Intersection = true
	prohibs.Add(prohibition.Build())

	decider := NewPReviewDeciderWithProhibitions(g, prohibs, rwe)

	if !decider.List(u1.Name, "", o1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}

	if decider.List(u1.Name, "", oa5.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestGraph25WithProhibition2(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	u1, _ := g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	o1, _ := g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read", "write"))

	prohibs := obm.New()
	prohibition := prohibitions.NewBuilder("deny", ua1.Name, operations.NewOperationSet("read", "write"))
	prohibition.AddContainer(oa1.Name, false)
	prohibition.AddContainer(oa2.Name, false)
	prohibition.Intersection = true
	prohibs.Add(prohibition.Build())

	decider := NewPReviewDeciderWithProhibitions(g, prohibs, rwe)

	if decider.List(u1.Name, "", o1.Name).Len() > 0 {
		t.Fatalf("permissions expected to be empty")
	}
}

func TestDeciderWithUA(t *testing.T) {
	g := gm.New()

	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua2, _ := g.CreateNode("ua2", graph.UA, nil, pc1.Name)
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, ua2.Name)
	g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	oa2, _ := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	g.CreateNode("o1", graph.O, nil, oa1.Name, oa2.Name)
	g.CreateNode("o2", graph.O, nil, oa2.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("read"))
	g.Associate(ua2.Name, oa1.Name, operations.NewOperationSet("write"))

	decider := NewPReviewDecider(g, rwe)

	if !decider.List(ua1.Name, "", oa1.Name).Contains("read", "write") {
		t.Fatalf("permissions expected to contain read and write")
	}
}

func TestProhibitionsAllCombinations(t *testing.T) {
	g := gm.New()
	g.CreatePolicyClass("pc1", nil)
	g.CreateNode("oa1", graph.OA, nil, "pc1")
	g.CreateNode("oa2", graph.OA, nil, "pc1")
	g.CreateNode("oa3", graph.OA, nil, "pc1")
	g.CreateNode("oa4", graph.OA, nil, "pc1")
	g.CreateNode("o1", graph.O, nil, "oa1", "oa2", "oa3")
	g.CreateNode("o2", graph.O, nil, "oa1", "oa4")
	g.CreateNode("ua1", graph.UA, nil, "pc1")
	g.CreateNode("u1", graph.U, nil, "ua1")
	g.CreateNode("u2", graph.U, nil, "ua1")
	g.CreateNode("u3", graph.U, nil, "ua1")
	g.CreateNode("u4", graph.U, nil, "ua1")

	g.Associate("ua1", "oa1", operations.NewOperationSet(operations.WRITE, operations.READ))

	prohibs := obm.New()
	prohibition := prohibitions.NewBuilder("p1", "u1", operations.NewOperationSet(operations.WRITE))
	prohibition.AddContainer("oa1", false)
	prohibition.AddContainer("oa2", false)
	prohibition.AddContainer("oa3", false)
	prohibition.Intersection = true
	prohibs.Add(prohibition.Build())

	prohibition = prohibitions.NewBuilder("p1", "u2", operations.NewOperationSet(operations.WRITE))
	prohibition.AddContainer("oa1", false)
	prohibition.AddContainer("oa2", false)
	prohibition.AddContainer("oa3", false)
	prohibs.Add(prohibition.Build())

	prohibition = prohibitions.NewBuilder("p1", "u3", operations.NewOperationSet(operations.WRITE))
	prohibition.AddContainer("oa1", false)
	prohibition.AddContainer("oa2", true)
	prohibition.Intersection = true
	prohibs.Add(prohibition.Build())

	prohibition = prohibitions.NewBuilder("p1", "u4", operations.NewOperationSet(operations.WRITE))
	prohibition.AddContainer("oa1", false)
	prohibition.AddContainer("oa2", true)
	prohibs.Add(prohibition.Build())

	prohibition = prohibitions.NewBuilder("p1", "u4", operations.NewOperationSet(operations.WRITE))
	prohibition.AddContainer("oa2", true)
	prohibs.Add(prohibition.Build())

	decider := NewPReviewDeciderWithProhibitions(g, prohibs, rwe)

	list := decider.List("u1", "", "o1")
	if !list.Contains("read") && list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}

	list = decider.List("u1", "", "o2")
	if !list.Contains("read") && !list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}

	list = decider.List("u2", "", "o2")
	if !list.Contains("read") && list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}

	list = decider.List("u3", "", "o2")
	if !list.Contains("read") && list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}

	list = decider.List("u4", "", "o1")
	if !list.Contains("read") && list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}

	list = decider.List("u4", "", "o2")
	if !list.Contains("read") && list.Contains("write") {
		t.Fatalf("permissions expected to contain read and NOT write")
	}
}

func TestPermissions(t *testing.T) {
	g := gm.New()
	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, _ := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	g.CreateNode("u1", graph.U, nil, ua1.Name)
	oa1, _ := g.CreateNode("oa1", graph.OA, nil, pc1.Name)
	g.CreateNode("o1", graph.O, nil, oa1.Name)

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet(operations.ALL_OPS))
	decider := NewPReviewDecider(g, rwe)
	list := decider.List("u1", "", "o1")
	if !list.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if !list.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet(operations.ALL_ADMIN_OPS))
	list = decider.List("u1", "", "o1")
	if !list.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to contain admin_ops")
	}

	if list.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to not contain rwe")
	}

	g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet(operations.ALL_RESOURCE_OPS))
	list = decider.List("u1", "", "o1")
	if list.Contains(operations.AdminOps().ToSlice()...) {
		t.Fatalf("permissions expected to not contain admin_ops")
	}

	if !list.Contains(rwe.ToSlice()...) {
		t.Fatalf("permissions expected to contain rwe")
	}
}

func TestPermissionInOnlyOnePC(t *testing.T) {
	g := gm.New()
	g.CreatePolicyClass("pc1", nil)
	g.CreatePolicyClass("pc2", nil)
	g.CreateNode("ua3", graph.UA, nil, "pc1")
	g.CreateNode("ua2", graph.UA, nil, "ua3")
	g.CreateNode("u1", graph.UA, nil, "ua2")

	g.CreateNode("oa1", graph.UA, nil, "pc1")
	g.CreateNode("oa3", graph.UA, nil, "pc2")
	g.Assign("oa3", "oa1")
	g.CreateNode("o1", graph.UA, nil, "oa3")

	g.Associate("ua3", "oa1", operations.NewOperationSet("read"))

	decider := NewPReviewDecider(g, operations.NewOperationSet("read"))

	if decider.List("u1", "", "o1").Len() > 0 {
		t.Fatalf("permissions should be empty")
	}

}
