package obligations

import (
	"github.com/jtejido/ngac/pip/graph"
)

type Obligation struct {
	User    string
	Enabled bool
	Label   string
	Rules   []*Rule
	Source  string
}

func NewObligation(user string) *Obligation {
	return &Obligation{User: user, Rules: make([]*Rule, 0)}
}

func (ob *Obligation) Clone() *Obligation {
	return &Obligation{ob.User, ob.Enabled, ob.Label, append([]*Rule{}, ob.Rules...), ob.Source}
}

type Rule struct {
	Label           string
	EventPattern    *EventPattern
	ResponsePattern *ResponsePattern
}

func NewRule() *Rule {
	return &Rule{"", new(EventPattern), NewResponsePattern()}
}

type EventPattern struct {
	Subject     *Subject
	PolicyClass *PolicyClass
	Operations  []string
	Target      *Target
}

type ResponsePattern struct {
	Condition        *Condition
	NegatedCondition *NegatedCondition
	Actions          []Action
}

func NewResponsePattern() *ResponsePattern {
	return &ResponsePattern{Actions: make([]Action, 0)}
}

type Condition struct {
	Condition []*Function
}

type NegatedCondition struct {
	Condition []*Function
}

type Subject struct {
	User    string
	AnyUser []string
	Process *EvrProcess
}

func NewSubject(user string) *Subject {
	return &Subject{User: user}
}

func NewSubjectFromUsers(users []string) *Subject {
	return &Subject{AnyUser: users}
}

func NewSubjectFromProcess(process *EvrProcess) *Subject {
	return &Subject{Process: process}
}

type EvrProcess struct {
	Value    string
	Function *Function
}

func NewEvrProcess(process string) *EvrProcess {
	return &EvrProcess{Value: process}
}

func (evr *EvrProcess) Equals(i interface{}) bool {
	if v, ok := i.(*EvrProcess); ok {
		if len(v.Value) == 0 {
			return false
		}

		return v.Value == evr.Value
	}

	return false
}

type PolicyClass struct {
	AnyOf  []string
	EachOf []string
}

func NewPolicyClass() *PolicyClass {
	return &PolicyClass{make([]string, 0), make([]string, 0)}
}

type Target struct {
	PolicyElements []*EvrNode
	Containers     []*EvrNode
}

type EvrNode struct {
	Name       string
	Type       string
	Properties graph.PropertyMap
	Function   *Function
	Process    *EvrProcess
}

func NewEvrNodeFromFunction(function *Function) *EvrNode {
	return &EvrNode{
		Function: function,
	}
}

func NewEvrNode(name, t string, properties graph.PropertyMap) *EvrNode {
	return &EvrNode{
		Name:       name,
		Type:       t,
		Properties: properties,
	}
}

func NewEvrNodeFromProcess(process *EvrProcess) *EvrNode {
	return &EvrNode{
		Process: process,
	}
}

func (evr *EvrNode) Equals(i interface{}) bool {
	if v, ok := i.(*EvrNode); ok {
		return evr.Name == v.Name
	}

	return false
}

type Containers struct {
	AnyOf  []*EvrNode
	EachOf []*EvrNode
}
