package tx

import (
    "ngac/pkg/pip/obligations"
    "sync"
)

type TxObligations struct {
    sync.RWMutex
    targetObligations obligations.Obligations
    cmds              []Committer
    txObligations     map[string]*obligations.Obligation
}

func NewTxObligations(o obligations.Obligations) *TxObligations {
    return &TxObligations{targetObligations: o, cmds: make([]Committer, 0), txObligations: make(map[string]*obligations.Obligation)}
}

func (to *TxObligations) Add(o *obligations.Obligation, enable bool) {
    to.Lock()

    if to.targetObligations.Get(o.Label) != nil {
        panic("obligation already exists with label " + o.Label)
    }

    to.cmds = append(to.cmds, func() error {
        to.targetObligations.Add(o, enable)
        return nil
    })
    to.txObligations[o.Label] = o
    to.Unlock()
}

func (to *TxObligations) Get(label string) *obligations.Obligation {
    to.RLock()
    obligation := to.targetObligations.Get(label)
    if obligation == nil {
        obligation = to.txObligations[label]
    }
    to.RUnlock()
    return obligation
}

func (to *TxObligations) All() []*obligations.Obligation {
    to.RLock()
    all := append([]*obligations.Obligation{}, to.targetObligations.All()...)
    for _, v := range to.txObligations {
        all = append(all, v)
    }
    to.RUnlock()
    return all
}

func (to *TxObligations) Update(label string, o *obligations.Obligation) {
    to.Lock()
    to.cmds = append(to.cmds, func() error {
        to.targetObligations.Update(label, o)
        return nil
    })
    to.txObligations[label] = o
    to.Unlock()
}

func (to *TxObligations) Remove(label string) {
    to.Lock()
    to.cmds = append(to.cmds, func() error {
        // obligation := to.targetObligations.Get(label)
        to.targetObligations.Remove(label)
        return nil
    })
    delete(to.txObligations, label)
    to.Unlock()
}

func (to *TxObligations) SetEnable(label string, enabled bool) {
    to.Lock()
    to.cmds = append(to.cmds, func() error {
        to.targetObligations.SetEnable(label, enabled)
        return nil
    })
    to.Unlock()
}

func (to *TxObligations) GetEnabled() []*obligations.Obligation {
    to.RLock()
    enabled := to.targetObligations.GetEnabled()
    for _, o := range to.txObligations {
        if o.Enabled {
            enabled = append(enabled, o)
        }
    }
    to.RUnlock()
    return enabled
}

func (to *TxObligations) Commit() (err error) {
    to.RLock()
    defer to.RUnlock()
    for _, txCmd := range to.cmds {
        if err = txCmd(); err != nil {
            return
        }
    }
    return nil
}
