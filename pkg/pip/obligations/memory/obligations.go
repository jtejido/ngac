package memory

import (
    ob "ngac/pkg/pip/obligations"
    "sync"
)

var (
    _ ob.Obligations = &obligations{}
)

type obligations struct {
    obligations map[string]*ob.Obligation
    sync.RWMutex
}

func New() ob.Obligations {
    return &obligations{obligations: make(map[string]*ob.Obligation)}
}

func (o *obligations) Add(obligation *ob.Obligation, enable bool) {
    if obligation == nil {
        panic("a nil obligation was received when creating a obligation")
    }
    o.Lock()
    obligation.Enabled = true
    o.obligations[obligation.Label] = obligation.Clone()
    o.Unlock()
}

func (o *obligations) Get(label string) (ob *ob.Obligation) {
    o.RLock()
    ob, _ = o.obligations[label]
    o.RUnlock()
    return
}

func (o *obligations) All() []*ob.Obligation {
    obs := make([]*ob.Obligation, 0)
    o.Lock()
    for _, o := range o.obligations {
        obs = append(obs, o.Clone())
    }
    o.Unlock()
    return obs
}

func (o *obligations) Update(label string, obligation *ob.Obligation) {
    if obligation == nil {
        panic("a null obligation was provided when updating a obligation")
    }
    updatedLabel := obligation.Label

    o.Lock()
    if updatedLabel != label {
        delete(o.obligations, label)
    }

    o.obligations[label] = obligation.Clone()
    o.Unlock()
}

func (o *obligations) Remove(label string) {
    o.Lock()
    delete(o.obligations, label)
    o.Unlock()
}

func (o *obligations) SetEnable(label string, enabled bool) {
    o.Lock()
    if obligation, ok := o.obligations[label]; ok {
        obligation.Enabled = enabled
        o.Update(label, obligation)
    }
    o.Unlock()
}

func (o *obligations) GetEnabled() []*ob.Obligation {
    obs := make([]*ob.Obligation, 0)
    o.RLock()
    for _, o := range o.obligations {
        if o.Enabled {
            obs = append(obs, o)
        }
    }
    o.RUnlock()
    return obs
}
