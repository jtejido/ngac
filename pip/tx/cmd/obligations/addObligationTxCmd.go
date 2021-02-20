package obligations

import (
    "github.com/jtejido/ngac/pip/obligations"
)

type AddObligationTxCmd struct {
    obligations obligations.Obligations
    obligation  *obligations.Obligation
    enabled     bool
}

func NewAddObligationTxCmd(o obligations.Obligations, obligation *obligations.Obligation, enabled bool) *AddObligationTxCmd {
    return &AddObligationTxCmd{o, obligation, enabled}
}

func (a *AddObligationTxCmd) Commit() error {
    a.obligations.Add(a.obligation, a.enabled)
}
