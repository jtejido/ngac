package service

import (
	"github.com/jtejido/ngac/pkg/common"
	"github.com/jtejido/ngac/pkg/context"
	"github.com/jtejido/ngac/pkg/epp"
	"github.com/jtejido/ngac/pkg/pap/policy"
	"github.com/jtejido/ngac/pkg/pdp/audit"
	"github.com/jtejido/ngac/pkg/pdp/decider"
	"github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/jtejido/ngac/pkg/pip/obligations"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type Service struct {
	userCtx     context.Context
	pap         common.FunctionalEntity
	epp         epp.EPP
	decider     decider.Decider
	auditor     audit.Auditor
	superPolicy *policy.SuperPolicy
}

func (s *Service) GraphAdmin() graph.Graph {
	return s.pap.Graph()
}

func (s *Service) ProhibitionsAdmin() prohibitions.Prohibitions {
	return s.pap.Prohibitions()
}

func (s *Service) ObligationsAdmin() obligations.Obligations {
	return s.pap.Obligations()
}

func (s *Service) Decider() decider.Decider {
	return s.decider
}

func (s *Service) Auditor() audit.Auditor {
	return s.auditor
}
