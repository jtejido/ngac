package pap

type PAP struct {
	graphAdmin        GraphAdmin
	prohibitionsAdmin ProhibitionsAdmin
	obligationsAdmin  ObligationsAdmin
}

func NewPAP(graphAdmin GraphAdmin, prohibitionsAdmin ProhibitionsAdmin, obligationsAdmin ObligationsAdmin) *PAP {
	return &PAP{graphAdmin, prohibitionsAdmin, obligationsAdmin}
}

func (pap *PAP) GraphAdmin() GraphAdmin {
	return pap.graphAdmin
}

func (pap *PAP) ProhibitionsAdmin() ProhibitionsAdmin {
	return pap.prohibitionsAdmin
}

func (pap *PAP) ObligationsAdmin() ObligationsAdmin {
	return pap.obligationsAdmin
}
