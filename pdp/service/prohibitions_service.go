package service

import (
	"fmt"
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/operations"
	"github.com/jtejido/ngac/pdp/decider"
	"github.com/jtejido/ngac/pdp/policy"
	"github.com/jtejido/ngac/pip/prohibitions"
)

type ProhibitionsService struct {
	Service
	prohibitions prohibitions.Prohibitions
}

func NewProhibitionsService(prohibs prohibitions.Prohibitions, policy *policy.SuperPolicy, d decider.Decider) *ProhibitionsService {
	p := new(ProhibitionsService)
	p.prohibitions = prohibs
	p.decider = d
	p.superPolicy = policy

	return p
}

func (ps *ProhibitionsService) AddProhibition(ctx Context, prohibition *prohibitions.Prohibition) error {
	name := prohibition.Name

	//check that the prohibition name is empty
	if len(name) == 0 {
		return fmt.Errorf("a empty name was provided when creating a prohibition")
	}

	//check the prohibitions doesn't already exist
	for p := range ps.All().Iter() {
		if p.(*prohibitions.Prohibition).Name == name {
			return fmt.Errorf("a prohibition with the name %s already exists", name)
		}
	}

	//check the user can create a prohibition on the subject and the nodes
	soaName := ps.superPolicy.SuperObjectAttribute().Name
	if !ps.decider.Check(ctx.User(), ctx.Process(), soaName, operations.CREATE_PROHIBITION) {
		return fmt.Errorf("unauthorized permissions on %s: %s", soaName, operations.CREATE_PROHIBITION)
	}

	//create prohibition in PAP
	return ps.prohibitions.AddProhibition(prohibition)
}

func (ps *ProhibitionsService) All() set.Set {
	return ps.prohibitions.All()
}

func (ps *ProhibitionsService) GetProhibition(prohibitionName string) *prohibitions.Prohibition {
	return ps.prohibitions.GetProhibition(prohibitionName)
}

func (ps *ProhibitionsService) ProhibitionsFor(subject string) set.Set {
	return ps.prohibitions.ProhibitionsFor(subject)
}

func (ps *ProhibitionsService) UpdateProhibition(prohibitionName string, prohibition *prohibitions.Prohibition) error {
	return ps.prohibitions.UpdateProhibition(prohibitionName, prohibition)
}

func (ps *ProhibitionsService) RemoveProhibition(prohibitionName string) {
	ps.prohibitions.RemoveProhibition(prohibitionName)
}

func (ps *ProhibitionsService) Reset(ctx Context) error {
	if !ps.hasPermissions(ctx, ps.superPolicy.SuperPolicyClassRep(), operations.RESET) {
		return fmt.Errorf("unauthorized permissions to reset the graph")
	}

	var names []string
	prohibs := ps.All()
	for prohib := range prohibs.Iter() {
		names = append(names, prohib.(*prohibitions.Prohibition).Name)
	}

	for _, name := range names {
		ps.RemoveProhibition(name)
	}

	return nil
}
