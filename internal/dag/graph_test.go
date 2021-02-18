package dag_test

import (
	"github.com/jtejido/ngac/common/dag"
	// "math/rand"
	"testing"
)

func TestDAG(t *testing.T) {
	d := dag.NewDirectedGraph()

	if d.Order() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %d", d.Order())
	}
}

func TestDAG_AddNode(t *testing.T) {
	dag1 := dag.NewDirectedGraph()
	vertex1 := 1
	dag1.AddNode(vertex1)

	if dag1.Order() != 1 {
		t.Fatalf("DAG number of vertices expected to be 1 but got %d", dag1.Order())
	}
}

func TestDAG_DeleteNode(t *testing.T) {
	dag1 := dag.NewDirectedGraph()
	vertex1 := 1

	dag1.AddNode(vertex1)

	if dag1.Order() != 1 {
		t.Fatalf("DAG number of vertices expected to be 1 but got %d", dag1.Order())
	}

	dag1.RemoveNode(vertex1)

	if dag1.Order() != 0 {
		t.Fatalf("DAG number of vertices expected to be 0 but got %d", dag1.Order())
	}
}

func TestDAG_AddEdge(t *testing.T) {
	dag1 := dag.NewDirectedGraph()

	vertex1 := 1
	vertex2 := 2

	dag1.AddNode(1)
	dag1.AddNode(2)

	dag1.SetEdge(&edge{vertex1, vertex2})
	if !dag1.HasEdgeFromTo(vertex1, vertex2) {
		t.Fatalf("Not valid edge")
	}

}

// // EdgeAdder is a graph.EdgeAdder graph.
// type EdgeAdder interface {
// 	dag.Directed
// 	dag.EdgeAdder
// }

// func addEdges(t *testing.T, n int, g EdgeAdder, newNode func(id int) dag.Node) {
// 	defer func() {
// 		r := recover()
// 		if r != nil {
// 			t.Errorf("unexpected panic: %v", r)
// 		}
// 	}()

// 	type altNode struct {
// 		dag.Node
// 	}

// 	rnd := rand.New(rand.NewSource(1))
// 	for i := 0; i < n; i++ {
// 		u := newNode(rnd.Intn(int(n)))
// 		var v dag.Node
// 		for {
// 			v = newNode(rnd.Intn(int(n)))
// 			if u.ID() != v.ID() {
// 				break
// 			}
// 		}

// 		g.AddEdge(u, v)
// 		if !g.HasEdgeFromTo(u, v) {
// 			t.Fatalf("SetEdge failed to find edge. from: %#v, to: %#v", u.ID(), v.ID())
// 		}
// 		if g.Node(u.ID()) == nil {
// 			t.Fatalf("SetEdge failed to add from node: %#v", u)
// 		}
// 		if g.Node(v.ID()) == nil {
// 			t.Fatalf("SetEdge failed to add to node: %#v", v)
// 		}

// 		g.AddEdge(altNode{u}, altNode{v})
// 		if nu := g.Node(u.ID()); nu == u {
// 			t.Fatalf("SetEdge failed to update from node: u=%#v nu=%#v", u, nu)
// 		}
// 		if nv := g.Node(v.ID()); nv == v {
// 			t.Fatalf("SetEdge failed to update to node: v=%#v nv=%#v", v, nv)
// 		}
// 	}
// }

func TestDAG_RemoveEdge(t *testing.T) {
	dag1 := dag.NewDirectedGraph()

	vertex1 := 1
	vertex2 := 2

	dag1.AddNode(vertex1)
	dag1.AddNode(vertex2)

	dag1.SetEdge(&edge{vertex1, vertex2})
	if !dag1.HasEdgeFromTo(vertex1, vertex2) {
		t.Fatalf("Can't add edge to DAG")
	}

}

// func TestDAG_Node(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)

// 	v1 := dag1.Node(1)
// 	v2 := dag1.Node(2)

// 	if v1.ID() != 1 {
// 		t.Fatalf("Expected value1 to be %d but got %v.", 1, v1.ID())
// 	}
// 	if v2.ID() != 2 {
// 		t.Fatalf("Expected value2 to be %d but got %v.", 2, v2.ID())
// 	}
// }

// func TestDAG_Order(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	expected_order := 0
// 	order := dag1.Order()
// 	if order != expected_order {
// 		t.Fatalf("Expected order to be %d but got %d", expected_order, order)
// 	}

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)
// 	vertex3 := Node(3)

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)
// 	dag1.AddNode(vertex3)

// 	expected_order = 3
// 	order = dag1.Order()
// 	if order != expected_order {
// 		t.Fatalf("Expected order to be %d but got %d", expected_order, order)
// 	}
// }

// func TestDAG_Size(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	expected_size := 0
// 	size := dag1.Size()
// 	if size != expected_size {
// 		t.Fatalf("Expected size to be %d but got %d", expected_size, size)
// 	}

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)
// 	vertex3 := Node(3)
// 	vertex4 := Node(4)

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)
// 	dag1.AddNode(vertex3)
// 	dag1.AddNode(vertex4)

// 	expected_size = 0
// 	size = dag1.Size()
// 	if size != expected_size {
// 		t.Fatalf("Expected size to be %d but got %d", expected_size, size)
// 	}

// 	dag1.AddEdge(vertex1, vertex2)
// 	dag1.AddEdge(vertex2, vertex3)
// 	dag1.AddEdge(vertex2, vertex4)

// 	expected_size = 3
// 	size = dag1.Size()
// 	if size != expected_size {
// 		t.Fatalf("Expected size to be %d but got %d", expected_size, size)
// 	}
// }

// func TestDAG_Descendants(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)
// 	dag1.AddEdge(vertex1, vertex2)

// 	successors := dag1.Descendants(vertex1)
// 	got := successors.Next().ID()
// 	if got != 2 {
// 		t.Fatalf("Successor vertex expected to be 2 but got %q", got)
// 	}
// }

// func TestDAG_Ancestors(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)

// 	dag1.AddEdge(vertex1, vertex2)

// 	predecessors := dag1.Ancestors(vertex2)
// 	got := predecessors.Next().ID()
// 	if got != 1 {
// 		t.Fatalf("Predecessor vertex expected to be 1 but got %q", got)
// 	}
// }

// // https://www.cs.hmc.edu/~keller/courses/cs60/s98/examples/acyclic/
// func TestDAG_Acyclic(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)
// 	vertex3 := Node(3)
// 	vertex4 := Node(4)
// 	vertex5 := Node(5)
// 	vertex6 := Node(6)

// 	dag1.AddEdge(vertex1, vertex2)
// 	dag1.AddEdge(vertex2, vertex3)
// 	dag1.AddEdge(vertex2, vertex4)
// 	dag1.AddEdge(vertex4, vertex5)
// 	dag1.AddEdge(vertex4, vertex6)
// 	dag1.AddEdge(vertex5, vertex6)
// 	dag1.AddEdge(vertex6, vertex3)

// 	if dag1.IsCyclic() {
// 		t.Fatalf("Expected dag to be acyclic but got %v.", dag1.IsCyclic())
// 	}
// }

// // https://www.cs.hmc.edu/~keller/courses/cs60/s98/examples/acyclic/
// func TestDAG_Cyclic(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := Node(1)
// 	vertex2 := Node(2)
// 	vertex3 := Node(3)
// 	vertex4 := Node(4)
// 	vertex5 := Node(5)
// 	vertex6 := Node(6)

// 	dag1.AddEdge(vertex1, vertex2)
// 	dag1.AddEdge(vertex2, vertex3)
// 	dag1.AddEdge(vertex2, vertex4)
// 	dag1.AddEdge(vertex4, vertex5)
// 	dag1.AddEdge(vertex6, vertex3)
// 	dag1.AddEdge(vertex5, vertex6)
// 	dag1.AddEdge(vertex6, vertex4)

// 	if !dag1.IsCyclic() {
// 		t.Fatalf("Expected dag to be cyclic but got %v.", dag1.IsCyclic())
// 	}
// }

// func TestDAG_TopologicalSort(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex0 := Node(0)
// 	vertex1 := Node(1)
// 	vertex2 := Node(2)
// 	vertex3 := Node(3)
// 	vertex4 := Node(4)
// 	vertex5 := Node(5)

// 	dag1.AddEdge(vertex5, vertex2)
// 	dag1.AddEdge(vertex5, vertex0)
// 	dag1.AddEdge(vertex4, vertex0)
// 	dag1.AddEdge(vertex4, vertex1)
// 	dag1.AddEdge(vertex2, vertex3)
// 	dag1.AddEdge(vertex3, vertex1)

// 	res := dag1.TopologicalSort()

// 	exp := []int{5, 4, 2, 3, 1, 0}

// 	var i int
// 	for res.HasNext() {
// 		id := res.Next().ID()
// 		if id != exp[i] {
// 			t.Fatalf("Expected %d but got %v.", exp[i], id)
// 		}

// 		i++
// 	}
// }

type edge struct {
	f, t interface{}
}

func (e *edge) From() dag.Node {
	return e.f
}

func (e *edge) To() dag.Node {
	return e.t
}

// https://www.hackerearth.com/practice/algorithms/graphs/topological-sort/tutorial/
// func TestDAG_TopologicalSort2(t *testing.T) {
// 	dag1 := dag.NewDirectedGraph()

// 	vertex1 := 1
// 	vertex2 := 2
// 	vertex3 := 3
// 	vertex4 := 4
// 	vertex5 := 5

// 	dag1.AddNode(vertex1)
// 	dag1.AddNode(vertex2)
// 	dag1.AddNode(vertex3)
// 	dag1.AddNode(vertex4)
// 	dag1.AddNode(vertex5)

// 	dag1.SetEdge(&edge{vertex1, vertex2})
// 	dag1.SetEdge(&edge{vertex2, vertex3})
// 	dag1.SetEdge(&edge{vertex2, vertex4})
// 	dag1.SetEdge(&edge{vertex3, vertex4})
// 	dag1.SetEdge(&edge{vertex3, vertex5})
// 	dag1.SetEdge(&edge{vertex1, vertex3})

// 	res := dag1.TopologicalSort()

// 	exp := []int{1, 2, 3, 4, 5}

// 	var i int
// 	for res.HasNext() {
// 		id := res.Next()
// 		if id != exp[i] {
// 			t.Fatalf("Expected %d but got %v.", exp[i], id)
// 		}

// 		i++
// 	}
// }
