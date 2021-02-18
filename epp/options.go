package epp

import (
	"github.com/jtejido/ngac/epp/functions"
)

type EPPOptions struct {
	executors []functions.FunctionExecutor
}

func NewEPPOptions(executors ...functions.FunctionExecutor) *EPPOptions {
	eo := new(EPPOptions)
	eo.executors = executors
}

func (eo *EPPOptions) Executors() []functions.FunctionExecutor {
	return eo.executors
}
