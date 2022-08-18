package memory

import (
	p "github.com/jtejido/ngac/pkg/pip/prohibitions"
	"strings"
	"sync"
)

var (
	_ p.Prohibitions = &prohibitions{}
)

type prohibitions struct {
	prohibitions map[string][]*p.Prohibition
	sync.RWMutex
}

func New() p.Prohibitions {
	return &prohibitions{prohibitions: make(map[string][]*p.Prohibition)}
}

func (mp *prohibitions) Add(prohibition *p.Prohibition) {
	mp.Lock()
	if prohibition == nil {
		panic("a nil prohibition was received when creating a prohibition")
	}

	if len(prohibition.Name) == 0 {
		panic("a nil or empty name was provided when creating a prohibition")
	}

	prohibition = prohibition.Clone()
	subject := prohibition.Subject
	exPros := make([]*p.Prohibition, 0)
	if v, ok := mp.prohibitions[subject]; ok {
		exPros = v
	}

	exPros = append(exPros, prohibition)

	mp.prohibitions[subject] = exPros
	mp.Unlock()
}

func (mp *prohibitions) All() []*p.Prohibition {
	pros := make([]*p.Prohibition, 0)
	mp.RLock()
	for _, pt := range mp.prohibitions {
		for _, p := range pt {
			pros = append(pros, p.Clone())
		}
	}
	mp.RUnlock()
	return pros
}

func (mp *prohibitions) Get(prohibitionName string) *p.Prohibition {
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

func (mp *prohibitions) ProhibitionsFor(subject string) []*p.Prohibition {
	ret := make([]*p.Prohibition, 0)
	mp.RLock()
	if pros, ok := mp.prohibitions[subject]; ok {
		for _, p := range pros {
			ret = append(ret, p.Clone())
		}
	}
	mp.RUnlock()
	return ret
}

func (mp *prohibitions) Update(prohibitionName string, prohibition *p.Prohibition) {
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

func (mp *prohibitions) Remove(prohibitionName string) {
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
