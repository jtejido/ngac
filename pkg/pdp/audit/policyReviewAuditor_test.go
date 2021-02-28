package audit

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pkg/operations"
	"github.com/jtejido/ngac/pkg/pip/graph"
	"sort"
	"strings"
	"testing"
)

var (
	rw    = operations.NewOperationSet("read", "write")
	r     = operations.NewOperationSet("read")
	w     = operations.NewOperationSet("write")
	noops = operations.NewOperationSet()
)

type testCase struct {
	name          string
	graph         graph.Graph
	expectedPaths map[string][]string
	expectedOps   set.Set
}

func getTests(t *testing.T) []testCase {
	return []testCase{
		graph1(t),
		graph2(t),
		graph3(t),
		graph4(t),
		graph5(t),
		graph7(t),
		graph8(t),
		graph9(t),
		graph10(t),
		graph11(t),
		graph12(t),
		graph13(t),
		graph14(t),
		graph15(t),
		graph16(t),
		graph17(t),
		graph18(t),
		graph19(t),
		graph20(t),
		graph21(t),
		graph22(t),
		graph23(t),
		graph24(t),
		// graph25(t), // you'll find yourself with a resulting string with different ops order and will fail string match
	}
}

func graph1(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph1", g, expectedPaths, rw}
}

func graph2(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa1", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua2(UA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph2", g, expectedPaths, rw}
}

func graph3(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa1", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua1(UA)-ua2(UA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph3", g, expectedPaths, rw}
}

func graph4(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{}
	return testCase{"graph4", g, expectedPaths, noops}
}

func graph5(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{}
	return testCase{"graph5", g, expectedPaths, noops}
}

// removed graph 6 because of change to Graph interface -- requiring parent nodes on creation prevents floating nodes

func graph7(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua2", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	expectedPaths["pc2"] = []string{"u1(U)-ua2(UA)-oa2(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph7", g, expectedPaths, r}
}

func graph8(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	expectedPaths["pc2"] = []string{"u1(U)-ua2(UA)-oa2(OA)-o1(O) ops=[write]"}
	return testCase{"graph8", g, expectedPaths, noops}
}

func graph9(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph9", g, expectedPaths, rw}
}

func graph10(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	expectedPaths["pc2"] = []string{}
	return testCase{"graph10", g, expectedPaths, noops}
}

func graph11(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph11", g, expectedPaths, rw}
}

func graph12(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa1", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua2(UA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph12", g, expectedPaths, rw}
}

func graph13(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	expectedPaths["pc2"] = []string{"u1(U)-ua2(UA)-oa2(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph13", g, expectedPaths, r}
}

func graph14(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	expectedPaths["pc2"] = []string{"u1(U)-ua2(UA)-oa2(OA)-o1(O) ops=[read]"}
	return testCase{"graph14", g, expectedPaths, r}
}

func graph15(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	expectedPaths["pc2"] = []string{}
	return testCase{"graph15", g, expectedPaths, noops}
}

func graph16(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	return testCase{"graph16", g, expectedPaths, r}
}

func graph17(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	return testCase{"graph17", g, expectedPaths, r}
}

func graph18(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreatePolicyClass("pc2", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa2", graph.OA, nil, "pc2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc2"] = []string{"u1(U)-ua1(UA)-ua2(UA)-oa2(OA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]"}
	return testCase{"graph18", g, expectedPaths, noops}
}

func graph19(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa2", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", rw)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa2", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write]", "u1(U)-ua1(UA)-ua2(UA)-oa2(OA)-oa1(OA)-o1(O) ops=[read]"}
	return testCase{"graph19", g, expectedPaths, rw}
}

func graph20(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua2", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua2(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph20", g, expectedPaths, rw}
}

func graph21(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}

	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua2", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1", "ua2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", w)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	err = g.Associate("ua2", "oa1", r)
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua2(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph21", g, expectedPaths, rw}
}

func graph22(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1", "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	return testCase{"graph22", g, expectedPaths, r}
}

func graph23(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua1", "oa2", operations.NewOperationSet("write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]"}
	return testCase{"graph23", g, expectedPaths, r}
}

func graph24(t *testing.T) testCase {
	g := graph.NewMemGraph()
	var err error
	_, err = g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	_, err = g.CreateNode("oa2", graph.OA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("oa1", graph.OA, nil, "oa2")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("ua1", graph.UA, nil, "pc1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, "oa1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, "ua1")
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate("ua1", "oa1", operations.NewOperationSet("read"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}
	err = g.Associate("ua1", "oa2", operations.NewOperationSet("write"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read]", "u1(U)-ua1(UA)-oa2(OA)-oa1(OA)-o1(O) ops=[write]"}
	return testCase{"graph24", g, expectedPaths, rw}
}

func graph25(t *testing.T) testCase {
	g := graph.NewMemGraph()
	pc1, err := g.CreatePolicyClass("pc1", nil)
	if err != nil {
		t.Fatalf("failed to create policy class: %s", err)
	}
	ua1, err := g.CreateNode("ua1", graph.UA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("u1", graph.U, nil, ua1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa2, err := g.CreateNode("oa2", graph.OA, nil, pc1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	oa1, err := g.CreateNode("oa1", graph.OA, nil, oa2.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}
	_, err = g.CreateNode("o1", graph.O, nil, oa1.Name)
	if err != nil {
		t.Fatalf("failed to create node: %s", err)
	}

	err = g.Associate(ua1.Name, oa1.Name, operations.NewOperationSet("*r"))
	if err != nil {
		t.Fatalf("failed to associate: %s", err)
	}

	expectedPaths := make(map[string][]string)
	expectedPaths["pc1"] = []string{"u1(U)-ua1(UA)-oa1(OA)-o1(O) ops=[read, write, execute]"}
	return testCase{"graph25", g, expectedPaths, rw}
}

func TestExplain(t *testing.T) {
	for _, tc := range getTests(t) {
		auditor := NewPReviewAuditor(tc.graph, operations.NewOperationSet("read", "write", "execute"))
		explain, err := auditor.Explain("u1", "o1")
		if err != nil {
			t.Fatalf("failed to explain: %s", err)
		}

		if !explain.Permissions.Contains(tc.expectedOps.ToSlice()...) {
			t.Errorf("%s expected ops %q but got %q", tc.name, tc.expectedOps.ToSlice(), explain.Permissions.ToSlice())
		}

		for pcName, expectedPaths := range tc.expectedPaths {
			if expectedPaths == nil {
				t.Errorf("%s should not be nil", tc.name)
			}

			pc := explain.PolicyClasses[pcName]
			if len(expectedPaths) != pc.Paths.Len() {
				t.Errorf("should be equal. want: %d, got: %d", len(expectedPaths), pc.Paths.Len())
			}

			for _, exPathStr := range expectedPaths {
				var match bool
				for resPath := range pc.Paths.Iter() {
					if pathsMatch(exPathStr, resPath.(*Path).String()) {
						match = true
						break
					}
				}
				if !match {
					t.Errorf("%s expected path \"%s\" but it was not in the results \"%q\"", tc.name, exPathStr, pc.Paths.ToSlice())
				}
			}
		}
	}
}

func pathsMatch(expectedStr, actualStr string) bool {
	expectedArr := strings.Split(expectedStr, "-")
	actualArr := strings.Split(actualStr, "-")

	if len(expectedArr) != len(actualArr) {
		return false
	}

	for i := 0; i < len(expectedArr); i++ {
		ex := expectedArr[i]
		res := actualArr[i]
		// if the element has brackets, it's a list of permissions
		if strings.HasPrefix(ex, "[") && strings.HasPrefix(res, "[") {
			// trim the brackets from the strings
			ex = ex[1:]
			res = res[1:]

			// split both into an array of strings
			exOps := strings.Split(ex, ",")
			resOps := strings.Split(res, ",")

			sort.Strings(exOps)
			sort.Strings(resOps)

			if len(exOps) != len(resOps) {
				return false
			}
			for j := 0; j < len(exOps); j++ {
				if exOps[j] != resOps[j] {
					return false
				}
			}
		} else if ex != actualArr[i] {
			return false
		}
	}

	return true
}
