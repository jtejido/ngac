package obligations

type Obligations interface {
	Add(*Obligation, bool)
	Get(string) *Obligation
	All() []*Obligation
	Update(string, *Obligation)
	Remove(string)
	SetEnable(string, bool)
	GetEnabled() []*Obligation
}
