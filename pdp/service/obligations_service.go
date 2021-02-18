package service

import (
	"fmt"
	"github.com/jtejido/ngac/operations"
	"github.com/jtejido/ngac/pdp/decider"
	"github.com/jtejido/ngac/pdp/policy"
	"github.com/jtejido/ngac/pip/obligations"
)

type ObligationsService struct {
	Service
	obligations obligations.Obligations
}

func NewObligationsService(obligs obligations.Obligations, policy *policy.SuperPolicy, d decider.Decider) *ObligationsService {
	p := new(ObligationsService)
	p.obligations = obligs
	p.decider = d
	p.superPolicy = policy

	return p
}

func (os *ObligationsService) AddObligation(obligation *obligations.Obligation, enable bool) error {
	return os.obligations.AddObligation(obligation, enable)
}

func (os *ObligationsService) GetObligation(label string) *obligations.Obligation {
	return os.obligations.GetObligation(label)
}

func (os *ObligationsService) All() []*obligations.Obligation {
	return os.obligations.All()
}

func (os *ObligationsService) UpdateObligation(label string, obligation *obligations.Obligation) {
	os.obligations.UpdateObligation(label, obligation)
}

func (os *ObligationsService) RemoveObligation(label string) {
	os.obligations.RemoveObligation(label)
}

func (os *ObligationsService) SetEnable(label string, enabled bool) {
	os.obligations.SetEnable(label, enabled)
}

func (os *ObligationsService) GetEnabled() []*obligations.Obligation {
	return os.obligations.GetEnabled()
}

func (os *ObligationsService) Reset(ctx Context) error {
	// check that the user can reset the graph
	if !os.hasPermissions(ctx, os.superPolicy.SuperPolicyClassRep(), operations.RESET) {
		return fmt.Errorf("unauthorized permissions to reset obligations")
	}

	obligs := os.All()
	var labels []string
	for _, oblig := range obligs {
		labels = append(labels, oblig.Label)
	}

	for _, label := range labels {
		os.RemoveObligation(label)
	}

	return nil
}
