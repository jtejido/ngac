package prohibitions

import (
    "github.com/jtejido/ngac/pip/prohibitions"
)

type DeleteProhibitionTxCmd struct {
    prohibitions    prohibitions.Prohibitions
    prohibitionName string
}

func NewDeleteProhibitionTxCmd(p prohibitions.Prohibitions, prohibitionName string) *DeleteProhibitionTxCmd {
    return &DeleteProhibitionTxCmd{p, prohibitionName}
}

func (d *DeleteProhibitionTxCmd) Commit() error {
    d.prohibitions.Remove(d.prohibitionName)
    return nil
}
