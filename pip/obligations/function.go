package obligations

type Function struct {
	Name string
	Args []*Arg
}

func NewFunction(name string, args []*Arg) *Function {
	return &Function{name, args}
}

type Arg struct {
	Value    string
	Function *Function
}

func NewArg(value string) *Arg {
	return &Arg{Value: value}
}

func NewArgFromFunction(function *Function) *Arg {
	return &Arg{Function: function}
}
