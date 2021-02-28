package pap

import (
	"fmt"
	"github.com/jtejido/ngac/pkg/common"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
)

var (
	_ prohibitions.Prohibitions = &ProhibitionsAdmin{}
)

type ProhibitionsAdmin struct {
	prohibitions prohibitions.Prohibitions
}

func NewProhibitionsAdmin(pip common.FunctionalEntity) *ProhibitionsAdmin {
	return &ProhibitionsAdmin{pip.Prohibitions()}
}

func (pa *ProhibitionsAdmin) Add(prohibition *prohibitions.Prohibition) {
	name := prohibition.Name

	//check that the prohibition name is not null or empty
	if len(name) == 0 {
		panic(fmt.Errorf("a null name was provided when creating a prohibition"))
	}

	//check the prohibitions doesn't already exist
	for _, p := range pa.All() {
		if p.Name == name {
			panic(fmt.Errorf("a prohibition with the name %s already exists", name))
		}
	}

	pa.prohibitions.Add(prohibition)
}

func (pa *ProhibitionsAdmin) All() []*prohibitions.Prohibition {
	return pa.prohibitions.All()
}

func (pa *ProhibitionsAdmin) Get(prohibitionName string) *prohibitions.Prohibition {
	return pa.prohibitions.Get(prohibitionName)
}

func (pa *ProhibitionsAdmin) ProhibitionsFor(subject string) []*prohibitions.Prohibition {
	return pa.prohibitions.ProhibitionsFor(subject)
}

func (pa *ProhibitionsAdmin) Update(prohibitionName string, prohibition *prohibitions.Prohibition) {
	pa.prohibitions.Update(prohibitionName, prohibition)
}

func (pa *ProhibitionsAdmin) Remove(prohibitionName string) {
	pa.prohibitions.Remove(prohibitionName)
}
