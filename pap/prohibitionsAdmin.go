package pap

import (
	"fmt"
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pip/prohibitions"
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

func (pa *ProhibitionsAdmin) AddProhibition(prohibition *prohibitions.Prohibition) error {
	name := prohibition.Name

	//check that the prohibition name is not null or empty
	if len(name) == 0 {
		return fmt.Errorf("a null name was provided when creating a prohibition")
	}

	//check the prohibitions doesn't already exist
	for p := range pa.All().Iter() {
		if p.(*prohibitions.Prohibition).Name == name {
			return fmt.Errorf("a prohibition with the name %s already exists", name)
		}
	}

	return pa.prohibitions.AddProhibition(prohibition)
}

func (pa *ProhibitionsAdmin) All() set.Set {
	return pa.prohibitions.All()
}

func (pa *ProhibitionsAdmin) GetProhibition(prohibitionName string) (*prohibitions.Prohibition, error) {
	return pa.prohibitions.GetProhibition(prohibitionName)
}

func (pa *ProhibitionsAdmin) ProhibitionsFor(subject string) set.Set {
	return pa.prohibitions.ProhibitionsFor(subject)
}

func (pa *ProhibitionsAdmin) UpdateProhibition(prohibitionName string, prohibition *prohibitions.Prohibition) error {
	return pa.prohibitions.UpdateProhibition(prohibitionName, prohibition)
}

func (pa *ProhibitionsAdmin) RemoveProhibition(prohibitionName string) {
	pa.prohibitions.RemoveProhibition(prohibitionName)
}
