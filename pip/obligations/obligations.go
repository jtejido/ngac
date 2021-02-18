package obligations

type Obligations interface {
	AddObligation(*Obligation, bool) error
	GetObligation(string) *Obligation
	All() []*Obligation
	UpdateObligation(string, *Obligation)
	RemoveObligation(string)
	SetEnable(string, bool)
	GetEnabled() []*Obligation
}
