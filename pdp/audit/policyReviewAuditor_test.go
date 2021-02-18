package audit

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/operations"
	"github.com/jtejido/ngac/pip/graph"
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
		// graph3(),
		// graph4(),
		// graph5(),
		// graph7(),
		// graph8(),
		// graph9(),
		// graph10(),
		// graph11(),
		// graph12(),
		// graph13(),
		// graph14(),
		// graph15(),
		// graph16(),
		// graph17(),
		// graph18(),
		// graph19(),
		// graph20(),
		// graph21(),
		// graph22(),
		// graph23(),
		// graph24(),
		// graph25(),
	}
}

func graph1(t *testing.T) testCase {
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
