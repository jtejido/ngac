package tx

import (
    "github.com/jtejido/ngac/pip/cmd"
    "github.com/jtejido/ngac/pip/obligations"
    "sync"
)

type TxObligations struct {
    sync.RWMutex
    targetObligations obligations.Obligations
    cmds              []cmd.TxCmd
    txObligations     map[string]obligations.Obligations
}

func NewTxObligations(o obligations.Obligations) *TxObligations {
    return &TxObligations{targetObligations: o, cmds: make([]cmd.TxCmd, 0), txObligations: make(map[string]obligations.Obligations)}
}

func (to *TxObligations) Add(o *obligations.Obligation, enable bool) {
    sync.Lock()

    if to.targetObligations.Get(obligation.Label) != nil {
        panic("obligation already exists with label " + obligation.Label)
    }

    to.cmds = append(to.cmds, cmd.NewAddObligationTxCmd(to.targetObligations, o, enable))
    to.txObligations[obligation.Label] = o
    sync.Unlock()
}

func (to *TxObligations) Get(label string) *obligations.Obligation {
    sync.RLock()
    obligation := to.targetObligations.Get(label)
    if obligation == nil {
        obligation = to.txObligations[label]
    }
    sync.RUnlock()
    return obligation
}

func (to *TxObligations) All() []*obligations.Obligation {
    sync.RLock()
    all := append([]*obligations.Obligation{}, to.targetObligations.All()...)
    for _, v := range to.txObligations {
        all = append(all, v)
    }
    sync.RUnlock()
    return all
}

func (to *TxObligations) Update(label string, o *obligations.Obligation) {
    sync.Lock()
    to.cmds = append(to.cmds, cmd.NewUpdateObligationTxCmd(to.targetObligations, label, o))
    to.txObligations[label] = o
    sync.Unlock()
}

func (to *TxObligations) Remove(label string) {
    sync.Lock()
    to.cmds = append(to.cmds, cmd.NewDeleteObligationTxCmd(to.targetObligations, label))
    delete(to.txObligations, label)
    sync.Unlock()
}

func (to *TxObligations) SetEnable(label string, enabled bool) {
    sync.Lock()
    to.cmds = append(to.cmds, cmd.NewSetEnableTxCmd(to.targetObligations, label, enabled))
    sync.Unlock()
}

func (to *TxObligations) GetEnabled() {
    sync.RLock()
    enabled := to.targetObligations.GetEnabled()
    for _, o := range to.txObligations {
        if o.Enabled {
            enabled = append(enabled, o)
        }
    }
    sync.RUnlock()
    return enabled
}

func (to *TxObligations) Commit() (err error) {
    sync.RLock()
    defer sync.RUnlock()
    for _, txCmd := range to.cmds {
        if err = txCmd.Commit(); err != nil {
            return
        }
    }
}
