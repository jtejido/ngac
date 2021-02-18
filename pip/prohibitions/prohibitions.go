package prohibitions

import "github.com/jtejido/ngac/internal/set"

type Prohibitions interface {
	AddProhibition(*Prohibition) error
	All() set.Set
	GetProhibition(string) (*Prohibition, error)
	ProhibitionsFor(string) set.Set
	UpdateProhibition(string, *Prohibition) error
	RemoveProhibition(string)
}
