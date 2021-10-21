package audit

import (
	"ngac/internal/set"
)

type PolicyClass struct {
	Operations, Paths set.Set
}

func NewEmptyPolicyClass() *PolicyClass {
	return NewPolicyClass(set.NewSet(), set.NewSet())
}

func NewPolicyClass(operations, paths set.Set) *PolicyClass {
	return &PolicyClass{operations, paths}
}
