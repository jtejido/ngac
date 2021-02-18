package service

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pap/policy"
	"github.com/jtejido/ngac/pdp/audit"
	"github.com/jtejido/ngac/pdp/decider"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
)

type Service struct {
	pap         common.FunctionalEntity
	superPolicy *policy.SuperPolicy
	UserCtx     Context
	decider     decider.Decider
	auditor     audit.Auditor
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
