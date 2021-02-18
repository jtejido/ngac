package prohibitions

import (
	"fmt"
	"github.com/jtejido/ngac/internal/omap"
	"github.com/jtejido/ngac/internal/set"
	"strings"
)

var (
	_ Prohibitions = &MemProhibitions{}
)

// embeds a regular map interface
type MemProhibitions struct {
	omap.OrderedMap
}

func NewMemProhibitions() *MemProhibitions {
	return &MemProhibitions{omap.NewOrderedMap()}
}

func (mp *MemProhibitions) AddProhibition(prohibition *Prohibition) error {
	if prohibition == nil {
		return fmt.Errorf("a nil prohibition was received when creating a prohibition")
	}

	if len(prohibition.Name) == 0 {
		return fmt.Errorf("a nil or empty name was provided when creating a prohibition")
	}

	prohibition = prohibition.Clone()
	subject := prohibition.Subject
	var exPros set.Set
	if v, ok := mp.Get(subject); ok {
		exPros = v.(set.Set)
	} else {
		exPros = set.NewSet()
	}

	exPros.Add(prohibition)

	mp.Add(subject, exPros)

	return nil
}

func (mp *MemProhibitions) All() set.Set {
	pros := set.NewSet()
	for _, pList := range mp.Values() {
		pt := pList.(set.Set)
		for p := range pt.Iter() {
			pros.Add(p.(*Prohibition).Clone())
		}
	}

	return pros
}

func (mp *MemProhibitions) GetProhibition(prohibitionName string) (*Prohibition, error) {
	for _, ps := range mp.Values() {
		for _, p := range ps.([]*Prohibition) {
			if strings.ToLower(p.Name) == strings.ToLower(prohibitionName) {
				return p.Clone(), nil
			}
		}
	}

	return nil, fmt.Errorf("a prohibition does not exist with the name %s", prohibitionName)
}

func (mp *MemProhibitions) ProhibitionsFor(subject string) set.Set {
	ret := set.NewSet()
	if v, ok := mp.Get(subject); ok {
		pros := v.(set.Set)
		for p := range pros.Iter() {
			ret.Add(p.(*Prohibition).Clone())
		}
	}

	return ret
}

func (mp *MemProhibitions) UpdateProhibition(prohibitionName string, prohibition *Prohibition) error {
	if prohibition == nil {
		return fmt.Errorf("a null prohibition was provided when updating a prohibition")
	}

	prohibition.Name = prohibitionName
	mp.RemoveProhibition(prohibition.Name)
	// add the updated prohibition
	return mp.AddProhibition(prohibition.Clone())
}

func (mp *MemProhibitions) RemoveProhibition(prohibitionName string) {
	for subject := range mp.Keys() {
		if v, ok := mp.Get(subject); ok {
			ps := v.(set.Set).ToSlice()
			for i, p := range ps {
				if p.(*Prohibition).Name == prohibitionName {
					ps = append(ps[:i], ps[i+1:]...)
					mp.Add(subject, ps)
				}
			}
		}
	}
}
