package prohibitions

type Prohibitions interface {
	Add(*Prohibition)
	All() []*Prohibition
	Get(string) *Prohibition
	ProhibitionsFor(string) []*Prohibition
	Update(string, *Prohibition)
	Remove(string)
}
