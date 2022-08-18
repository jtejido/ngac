package audit

import (
    "github.com/jtejido/ngac/internal/set"
)

type Explain struct {
    Permissions   set.Set
    PolicyClasses map[string]*PolicyClass
}

func NewEmptyExplain() *Explain {
    return NewExplain(set.NewSet(), make(map[string]*PolicyClass))
}

func NewExplain(permissions set.Set, policyClasses map[string]*PolicyClass) *Explain {
    return &Explain{permissions, policyClasses}
}

func (e *Explain) String() string {
    s := "Permissions: "
    var i int
    for p := range e.Permissions.Iter() {
        s += p.(string)
        if i < e.Permissions.Len()-1 {
            s += ", "
        }
        i++
    }

    for pc, policyClass := range e.PolicyClasses {
        s += "\n\t\t"
        s += pc
        s += ": "
        var i int
        for p := range policyClass.Operations.Iter() {
            s += p.(string)
            if i < policyClass.Operations.Len()-1 {
                s += ", "
            }
            i++
        }
    }

    s += "\nPaths:"
    for pc, policyClass := range e.PolicyClasses {
        s += "\n\t\t"
        s += pc
        s += ": "
        var i int
        for p := range policyClass.Operations.Iter() {
            s += p.(string)
            if i < policyClass.Operations.Len()-1 {
                s += ", "
            }
            i++
        }
        paths := policyClass.Paths
        for path := range paths.Iter() {
            s += "\n\t\t\t- "
            s += path.(*Path).String()
        }
    }

    return s
}
