package tx

import (
    "github.com/jtejido/ngac/pip/cmd"
    "github.com/jtejido/ngac/pip/prohibitions"
    "sync"
)

type Prohibitions interface {
    Add(*Prohibition)
    All() []*Prohibition
    Get(string) *Prohibition
    ProhibitionsFor(string) []*Prohibition
    Update(string, *Prohibition)
    Remove(string)
}

type TxProhibitions struct {
    sync.RWMutex
    targetProhibitions prohibitions.Prohibitions
    prohibitions       []prohibitions.Prohibitions
    cmds               []cmd.TxCmd
}

func NewTxProhibitions(p prohibitions.Prohibitions) *TxProhibitions {
    return &TxObligations{targetProhibitions: p, cmds: make([]cmd.TxCmd, 0), prohibitions: make(map[string]prohibitions.Prohibitions)}
}

func (tp *TxProhibitions) Add(prohibition *prohibitions.Prohibition) {
    sync.Lock()
    tp.cmds = append(tp.cmds, cmd.NewAddProhibitionTxCmd(tp.targetProhibitions, prohibition))
    tp.prohibitions = append(tp.prohibitions, prohibition)
    sync.Unlock()
}

func (tp *TxProhibitions) All() []*prohibitions.Prohibition {
    sync.RLock()
    all := tp.targetProhibitions.All()
    for _, prohibition := range tp.prohibitions {
        all = append(all, prohibition.Clone())
    }
    sync.RUnlock()
    return all
}

func (tp *TxProhibitions) Get(prohibitionName string) *prohibitions.Prohibition {
    sync.RLock()
    prohibition := tp.targetProhibitions.Get(prohibitionName)
    if prohibition == nil {
        for _, p := range tp.prohibitions {
            if p.Name == prohibitionName {
                prohibition = p
            }
        }
    }
    sync.RUnlock()
    return prohibition
}

func (tp *TxProhibitions) ProhibitionsFor(subject string) []*prohibitions.Prohibition {
    sync.RLock()
    ret := tp.targetProhibitions.ProhibitionsFor(subject)
    for _, p := range prohibitions {
        if p.Subject == subject {
            ret = append(ret, p)
        }
    }
    sync.RUnlock()
    return ret
}

func (tp *TxProhibitions) Update(prohibitionName string, prohibition *prohibitions.Prohibition) {
    sync.Lock()
    tp.cmds = append(tp.cmds, cmd.NewUpdateProhibitionTxCmd(tp.targetProhibitions, prohibitionName, prohibition))
    for _, p := range tp.prohibitions {
        if p.Name == prohibitionName {
            p = prohibition
        }
    }
    sync.Unlock()
}

func (tp *TxProhibitions) Remove(prohibitionName string) {
    sync.Lock()
    tp.cmds = append(tp.cmds, cmd.NewDeleteProhibitionTxCmd(tp.targetProhibitions, prohibitionName))
    var idx int
    for i, p := range tp.prohibitions {
        if p.Name == prohibitionName {
            idx = i
        }
    }

    tp.prohibitions = append(tp.prohibitions[:idx], tp.prohibitions[idx+1:]...)
    sync.Unlock()
}

func (tp *TxProhibitions) Commit() error {
    sync.RLock()
    defer sync.RUnlock()
    for _, txCmd := range tp.cmds {
        if err = txCmd.Commit(); err != nil {
            return
        }
    }
}
