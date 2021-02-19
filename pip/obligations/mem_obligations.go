package obligations

import (
    "sync"
)

var (
    _ Obligations = &MemObligations{}
)

type MemObligations struct {
    obligations map[string]*Obligation
    sync.RWMutex
}

func NewMemObligations() *MemObligations {
    return &MemObligations{obligations: make(map[string]*Obligation)}
}

func (mo *MemObligations) Add(obligation *Obligation, enable bool) {
    if obligation == nil {
        panic("a nil obligation was received when creating a obligation")
    }
    mo.Lock()
    obligation.Enabled = true
    mo.obligations[obligation.Label] = obligation.Clone()
    mo.Unlock()
}

func (mo *MemObligations) Get(label string) (ob *Obligation) {
    mo.RLock()
    ob, _ = mo.obligations[label]
    mo.RUnlock()
    return
}

func (mo *MemObligations) All() []*Obligation {
    obs := make([]*Obligation, 0)
    mo.Lock()
    for _, o := range mo.obligations {
        obs = append(obs, o.Clone())
    }
    mo.Unlock()
    return obs
}

func (mo *MemObligations) Update(label string, obligation *Obligation) {
    if obligation == nil {
        panic("a null obligation was provided when updating a obligation")
    }
    updatedLabel := obligation.Label

    mo.Lock()
    if updatedLabel != label {
        delete(mo.obligations, label)
    }

    mo.obligations[label] = obligation.Clone()
    mo.Unlock()
}

func (mo *MemObligations) Remove(label string) {
    mo.Lock()
    delete(mo.obligations, label)
    mo.Unlock()
}

func (mo *MemObligations) SetEnable(label string, enabled bool) {
    mo.Lock()
    if obligation, ok := mo.obligations[label]; ok {
        obligation.Enabled = enabled
        mo.Update(label, obligation)
    }
    mo.Unlock()
}

func (mo *MemObligations) GetEnabled() []*Obligation {
    obs := make([]*Obligation, 0)
    mo.RLock()
    for _, o := range mo.obligations {
        if o.Enabled {
            obs = append(obs, o)
        }
    }
    mo.RUnlock()
    return obs
}
