package prohibitions

import (
    "github.com/jtejido/ngac/operations"
)

type Prohibition struct {
    Name         string
    Subject      string
    containers   map[string]bool
    Operations   operations.OperationSet
    Intersection bool
}

func NewProhibition(name, subject string, containers map[string]bool, ops operations.OperationSet, intersection bool) *Prohibition {
    p := new(Prohibition)
    p.Name = name
    p.Subject = subject

    if containers == nil {
        containers = make(map[string]bool)
    }

    p.containers = containers

    if ops == nil {
        ops = operations.NewOperationSet()
    }

    p.Operations = ops
    p.Intersection = intersection

    return p
}

func (p *Prohibition) Clone() *Prohibition {
    containers := make(map[string]bool)
    for key, value := range p.containers {
        containers[key] = value
    }
    return &Prohibition{
        Name:         p.Name,
        Subject:      p.Subject,
        containers:   containers,
        Operations:   p.Operations.Clone(),
        Intersection: p.Intersection,
    }
}

func (p *Prohibition) Containers() map[string]bool {
    return p.containers
}

func (p *Prohibition) AddContainer(name string, complement bool) {
    p.containers[name] = complement
}

func (p *Prohibition) RemoveContainerCondition(name string) {
    delete(p.containers, name)
}

func (p *Prohibition) Equals(i interface{}) bool {
    if v, ok := i.(*Prohibition); ok {
        return p.Name == v.Name
    }

    return false
}

type Builder struct {
    name, subject string
    containers    map[string]bool
    operations    operations.OperationSet
    Intersection  bool
}

func NewBuilder(name, subject string, operations operations.OperationSet) *Builder {
    b := new(Builder)
    b.name = name
    b.subject = subject
    b.containers = make(map[string]bool)
    b.operations = operations
    b.Intersection = false

    return b
}

func (b *Builder) AddContainer(container string, complement bool) {
    b.containers[container] = complement
}

func (b *Builder) Build() *Prohibition {
    return NewProhibition(b.name, b.subject, b.containers, b.operations, b.Intersection)
}
