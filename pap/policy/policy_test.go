package policy

import (
	"github.com/jtejido/ngac/pip/graph"
	"testing"
)

func TestSuperPolicyWithEmptyGraph(t *testing.T) {
	g := graph.NewMemGraph()
	superPolicy := NewSuperPolicy()
	superPolicy.Configure(g)
	testSuperPolicy(t, g)
}

func testSuperPolicy(t *testing.T, g graph.Graph) {

	if !g.Exists("super_pc") {
		t.Fatalf("super_pc node doesn't exist")
	}
	if !g.Exists("super_pc_rep") {
		t.Fatalf("super_pc_rep node doesn't exist")
	}
	if !g.Exists("super_ua1") {
		t.Fatalf("super_ua1 node doesn't exist")
	}
	if !g.Exists("super_ua2") {
		t.Fatalf("super_ua2 node doesn't exist")
	}
	if !g.Exists("super") {
		t.Fatalf("super node doesn't exist")
	}
	if !g.Exists("super_oa") {
		t.Fatalf("super_oa node doesn't exist")
	}

	if !g.Parents("super").Contains("super_ua1", "super_ua2") {
		t.Fatalf("super node doesn't have super_ua1 and super_ua2 parents")
	}
	if !g.Parents("super_ua1").Contains("super_pc") {
		t.Fatalf("super_ua1 node doesn't have super_pc")
	}
	if !g.Parents("super_ua2").Contains("super_pc") {
		t.Fatalf("super_ua2 node doesn't have super_pc")
	}
	if !g.Parents("super_oa").Contains("super_pc") {
		t.Fatalf("super_oa node doesn't have super_pc")
	}

	if !g.Parents("super_pc_rep").Contains("super_oa") {
		t.Fatalf("super_pc_rep node doesn't have super_oa")
	}

	if v, err := g.SourceAssociations("super_ua1"); err != nil {
		t.Fatalf("error getting super_ua1 source associations")
	} else if _, ok := v["super_oa"]; !ok {
		t.Fatalf("super_oa node not in super_ua1 source associations")
	}

	if v, err := g.SourceAssociations("super_ua2"); err != nil {
		t.Fatalf("error getting super_ua2 source associations")
	} else if _, ok := v["super_ua1"]; !ok {
		t.Fatalf("super_ua1 node not in super_ua2 source associations")
	}
}
