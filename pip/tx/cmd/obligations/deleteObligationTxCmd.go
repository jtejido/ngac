package obligations

import (
	"github.com/jtejido/ngac/pip/obligations"
)

type DeleteObligationTxCmd struct {
	obligations obligations.Obligations
	label       string
	obligation  *obligations.Obligation
}

func NewDeleteObligationTxCmd(o obligations.Obligations, label string) *DeleteObligationTxCmd {
	return &DeleteObligationTxCmd{obligations: o, label: label}
}

func (d *DeleteObligationTxCmd) Commit() error {
	d.obligation = d.obligations.Get(label)
	d.obligations.Remove(label)
	return nil
}
