package obligations

import (
    "github.com/jtejido/ngac/internal/omap"
)

var (
    _ Obligations = &MemObligations{}
)

type MemObligations struct {
    omap.OrderedMap
}

func NewMemObligations() *MemObligations {
    return &MemObligations{omap.NewOrderedMap()}
}

func (mo *MemObligations) AddObligation(obligation *Obligation, enable bool) error {
    obligation.Enabled = true
    mo.Add(obligation.Label, obligation.Clone())
    return nil
}

func (mo *MemObligations) GetObligation(label string) *Obligation {
    if v, ok := mo.Get(label); ok {
        return v.(*Obligation).Clone()
    }

    return nil
}

func (mo *MemObligations) All() []*Obligation {
    obs := make([]*Obligation, 0)
    for _, o := range mo.Values() {
        obs = append(obs, o.(*Obligation).Clone())
    }
    return obs
}

func (mo *MemObligations) UpdateObligation(label string, obligation *Obligation) {
    updatedLabel := obligation.Label
    if updatedLabel != label {
        mo.Remove(label)
    }

    mo.Add(label, obligation.Clone())
}

func (mo *MemObligations) RemoveObligation(label string) {
    mo.Remove(label)
}

func (mo *MemObligations) SetEnable(label string, enabled bool) {
    if obligation, ok := mo.Get(label); ok {
        obligation.(*Obligation).Enabled = enabled
        mo.UpdateObligation(label, obligation.(*Obligation))
    }
}

func (mo *MemObligations) GetEnabled() []*Obligation {
    obs := make([]*Obligation, 0)
    for _, o := range mo.Values() {
        if o.(*Obligation).Enabled {
            obs = append(obs, o.(*Obligation))
        }
    }

    return obs
}
