package service

import (
	"ngac/pkg/common"
	"ngac/pkg/context"
	"ngac/pkg/epp"
	"ngac/pkg/pap/policy"
	"ngac/pkg/pdp/audit"
	"ngac/pkg/pdp/decider"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
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
