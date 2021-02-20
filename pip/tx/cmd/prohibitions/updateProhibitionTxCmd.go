package prohibitions

import (
    "github.com/jtejido/ngac/pip/prohibitions"
)

type UpdateProhibitionTxCmd struct {
    prohibitions    prohibitions.Prohibitions
    prohibitionName string
    prohibition     *prohibitions.Prohibition
}

func NewUpdateProhibitionTxCmd(p Prohibitions, prohibitionName string, prohibition *prohibitions.Prohibition) *UpdateProhibitionTxCmd {
    return &UpdateProhibitionTxCmd{p, prohibitionName, prohibition}
}

func (u *UpdateProhibitionTxCmd) Commit() error {
    u.prohibitions.Update(u.prohibitionName, u.prohibition)
    return nil
}
