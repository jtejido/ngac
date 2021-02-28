package epp

type EPPOptions struct {
	executors []FunctionExecutor
}

func NewEPPOptions(executors ...FunctionExecutor) *EPPOptions {
	eo := new(EPPOptions)
	eo.executors = executors
	return eo
}

func (eo *EPPOptions) Executors() []FunctionExecutor {
	return eo.executors
}
