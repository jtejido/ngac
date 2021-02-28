package prohibitions

import (
	"strings"
	"sync"
)

var (
	_ Prohibitions = &MemProhibitions{}
)

type MemProhibitions struct {
	prohibitions map[string][]*Prohibition
	sync.RWMutex
}

func NewMemProhibitions() *MemProhibitions {
	return &MemProhibitions{prohibitions: make(map[string][]*Prohibition)}
}

func (mp *MemProhibitions) Add(prohibition *Prohibition) {
	mp.Lock()
	if prohibition == nil {
		panic("a nil prohibition was received when creating a prohibition")
	}

	if len(prohibition.Name) == 0 {
		panic("a nil or empty name was provided when creating a prohibition")
	}

	prohibition = prohibition.Clone()
	subject := prohibition.Subject
	exPros := make([]*Prohibition, 0)
	if v, ok := mp.prohibitions[subject]; ok {
		exPros = v
	}

	exPros = append(exPros, prohibition)

	mp.prohibitions[subject] = exPros
	mp.Unlock()
}

func (mp *MemProhibitions) All() []*Prohibition {
	pros := make([]*Prohibition, 0)
	mp.RLock()
	for _, pt := range mp.prohibitions {
		for _, p := range pt {
			pros = append(pros, p.Clone())
		}
	}
	mp.RUnlock()
	return pros
}

func (mp *MemProhibitions) Get(prohibitionName string) *Prohibition {
	mp.RLock()
	defer mp.RUnlock()
	for _, ps := range mp.prohibitions {
		for _, p := range ps {
			if strings.ToLower(p.Name) == strings.ToLower(prohibitionName) {
				return p.Clone()
			}
		}
	}

	return nil
}

func (mp *MemProhibitions) ProhibitionsFor(subject string) []*Prohibition {
	ret := make([]*Prohibition, 0)
	mp.RLock()
	if pros, ok := mp.prohibitions[subject]; ok {
		for _, p := range pros {
			ret = append(ret, p.Clone())
		}
	}
	mp.RUnlock()
	return ret
}

func (mp *MemProhibitions) Update(prohibitionName string, prohibition *Prohibition) {
	if prohibition == nil {
		panic("a null prohibition was provided when updating a prohibition")
	}
	prohibition.Name = prohibitionName

	mp.Lock()
	mp.Remove(prohibition.Name)
	// add the updated prohibition
	mp.Add(prohibition.Clone())
	mp.Unlock()
}

func (mp *MemProhibitions) Remove(prohibitionName string) {
	mp.Lock()
	for subject, ps := range mp.prohibitions {
		for i := 0; i < len(ps); i++ {
			if ps[i].Name == prohibitionName {
				ps = append(ps[:i], ps[i+1:]...)
				mp.prohibitions[subject] = ps
			}
		}
	}
	mp.Unlock()
}
