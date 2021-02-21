package tx

import (
    "github.com/jtejido/ngac/pip/prohibitions"
    "sync"
)

type TxProhibitions struct {
    sync.RWMutex
    targetProhibitions prohibitions.Prohibitions
    prohibitions       []*prohibitions.Prohibition
    cmds               []Committer
}

func NewTxProhibitions(p prohibitions.Prohibitions) *TxProhibitions {
    return &TxProhibitions{targetProhibitions: p, cmds: make([]Committer, 0), prohibitions: make([]*prohibitions.Prohibition, 0)}
}

func (tp *TxProhibitions) Add(prohibition *prohibitions.Prohibition) {
    tp.Lock()
    tp.cmds = append(tp.cmds, func() error {
        tp.targetProhibitions.Add(prohibition)
        return nil
    })
    tp.prohibitions = append(tp.prohibitions, prohibition)
    tp.Unlock()
}

func (tp *TxProhibitions) All() []*prohibitions.Prohibition {
    tp.RLock()
    all := tp.targetProhibitions.All()
    for _, prohibition := range tp.prohibitions {
        all = append(all, prohibition.Clone())
    }
    tp.RUnlock()
    return all
}

func (tp *TxProhibitions) Get(prohibitionName string) *prohibitions.Prohibition {
    tp.RLock()
    prohibition := tp.targetProhibitions.Get(prohibitionName)
    if prohibition == nil {
        for _, p := range tp.prohibitions {
            if p.Name == prohibitionName {
                prohibition = p
            }
        }
    }
    tp.RUnlock()
    return prohibition
}

func (tp *TxProhibitions) ProhibitionsFor(subject string) []*prohibitions.Prohibition {
    tp.RLock()
    ret := tp.targetProhibitions.ProhibitionsFor(subject)
    for _, p := range tp.prohibitions {
        if p.Subject == subject {
            ret = append(ret, p)
        }
    }
    tp.RUnlock()
    return ret
}

func (tp *TxProhibitions) Update(prohibitionName string, prohibition *prohibitions.Prohibition) {
    tp.Lock()
    tp.cmds = append(tp.cmds, func() error {
        tp.targetProhibitions.Update(prohibitionName, prohibition)
        return nil
    })
    for i := 0; i < len(tp.prohibitions); i++ {
        if tp.prohibitions[i].Name == prohibitionName {
            tp.prohibitions[i] = prohibition
        }
    }
    tp.Unlock()
}

func (tp *TxProhibitions) Remove(prohibitionName string) {
    tp.Lock()
    tp.cmds = append(tp.cmds, func() error {
        tp.targetProhibitions.Remove(prohibitionName)
        return nil
    })
    var idx int
    for i, p := range tp.prohibitions {
        if p.Name == prohibitionName {
            idx = i
        }
    }

    tp.prohibitions = append(tp.prohibitions[:idx], tp.prohibitions[idx+1:]...)
    tp.Unlock()
}

func (tp *TxProhibitions) Commit() (err error) {
    tp.RLock()
    defer tp.RUnlock()
    for _, txCmd := range tp.cmds {
        if err = txCmd(); err != nil {
            return
        }
    }
    return nil
}
