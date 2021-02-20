package audit

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pip/graph"
)

type Path struct {
	Operations set.Set
	Nodes      []*graph.Node
}

func NewEmptyPath() *Path {
	return NewPath(set.NewSet(), make([]*graph.Node, 0))
}

func NewPath(operations set.Set, nodes []*graph.Node) *Path {
	return &Path{operations, nodes}
}

func (p *Path) Equals(o interface{}) bool {
	if v, ok := o.(*Path); ok {
		for i, n := range v.Nodes {
			if n.Equals(p.Nodes[i]) {
				return p.Operations.Equal(v.Operations)
			}
		}
	}
	return false
}

func (p *Path) String() string {
	if len(p.Nodes) == 0 {
		return ""
	}

	var s string
	var i int
	for _, node := range p.Nodes {
		s += node.Name
		s += "("
		s += node.Type.String()
		s += ")"
		if i < len(p.Nodes)-1 {
			s += "-"
		}
		i++
	}
	s += " ops=["
	i = 0
	for op := range p.Operations.Iter() {
		s += op.(string)
		if i < p.Operations.Len()-1 {
			s += ", "
		}
		i++
	}
	s += "]"
	return s
}
