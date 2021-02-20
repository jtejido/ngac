package cmd

import (
    "github.com/jtejido/ngac/pip/prohibitions"
)

type AddProhibitionTxCmd struct {
    prohibitions prohibitions.Prohibitions
    prohibition  *prohibitions.Prohibition
}

func NewAddProhibitionTxCmd(p prohibitions.Prohibitions, prohibition *prohibitions.Prohibition) *AddProhibitionTxCmd {
    return &AddProhibitionTxCmd{p, prohibition}
}

func (a *AddProhibitionTxCmd) Commit() error {
    a.prohibitions.Add(a.prohibition)
    return nil
}
