package obligations

import (
	"encoding/json"
	"fmt"
	"github.com/jtejido/ngac/pip/graph"
)

type Action interface {
	Condition() *Condition
	SetCondition(condition *Condition)
	NegatedCondition() *NegatedCondition
	SetNegatedCondition(*NegatedCondition)
	UnmarshalJSON(b []byte) error
}

type action struct {
	condition        *Condition
	negatedCondition *NegatedCondition
}

func NewAction() *action {
	return new(action)
}

func (a *action) Condition() *Condition {
	return a.condition
}

func (a *action) SetCondition(c *Condition) {
	a.condition = c
}

func (a *action) NegatedCondition() *NegatedCondition {
	return a.negatedCondition
}

func (a *action) SetNegatedCondition(n *NegatedCondition) {
	a.negatedCondition = n
}

// AssignAction.java
type AssignAction struct {
	action
	Assignments []*ActionAssignment
}

func NewAssignAction() *AssignAction {
	ans := new(AssignAction)
	ans.Assignments = make([]*ActionAssignment, 0)
	return ans
}

func (a *AssignAction) UnmarshalJSON(b []byte) error {

	return nil
}

type ActionAssignment struct {
	What, Where *EvrNode
}

// CreateAction.java
type CreateAction struct {
	action
	CreateNodesList []*ActionCreateNode
	Rules           []*Rule
}

func NewCreateAction() *CreateAction {
	ans := new(CreateAction)
	ans.CreateNodesList = make([]*ActionCreateNode, 0)
	ans.Rules = make([]*Rule, 0)
	return ans
}

func (a *CreateAction) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	if v, ok := raw["condition"]; ok {
		a.condition = new(Condition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = a.condition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["not_condition"]; ok {
		a.negatedCondition = new(NegatedCondition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = a.negatedCondition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	// if v, ok := raw["create"]; ok {
	// 	// TO-DO
	// }

	return fmt.Errorf("invalid action received")
}

type ActionCreateNode struct {
	What, Where *EvrNode
}

type DeleteAction struct {
	action
	Nodes               []*EvrNode
	Assignments         *AssignAction
	Associations        []*GrantAction
	Prohibitions, Rules []string
}

func (a *DeleteAction) UnmarshalJSON(b []byte) error {

	return nil
}

// DenyAction.java
type DenyAction struct {
	action
	Label      string
	Subject    *EvrNode
	Operations []string
	Target     *ActionTarget
}

func (a *DenyAction) UnmarshalJSON(b []byte) error {

	return nil
}

type ActionTarget struct {
	Complement, Intersection bool
	Containers               []*ActionContainer
}

type ActionContainer struct {
	Name, Type string
	Properties graph.PropertyMap
	Function   *Function
	Complement bool
}

func NewActionContainerFromFunction(function *Function) *ActionContainer {
	return &ActionContainer{Function: function}
}

func NewActionContainer(name, t string, properties graph.PropertyMap) *ActionContainer {
	return &ActionContainer{
		Name:       name,
		Type:       t,
		Properties: properties,
	}

}

type FunctionAction struct {
	action
	Function *Function
}

func NewFunctionAction(function *Function) *FunctionAction {
	ans := new(FunctionAction)
	ans.Function = function
	return ans
}

func (a *FunctionAction) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	if v, ok := raw["condition"]; ok {
		a.condition = new(Condition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = a.condition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["not_condition"]; ok {
		a.negatedCondition = new(NegatedCondition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = a.negatedCondition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["function"]; ok {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		a.Function = new(Function)
		err = a.Function.UnmarshalJSON(b)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("invalid action received")
}

type GrantAction struct {
	action
	Subject    *EvrNode
	Operations []string
	Target     *EvrNode
}

func NewGrantAction() *GrantAction {
	ans := new(GrantAction)
	ans.Operations = make([]string, 0)
	return ans
}

func (a *GrantAction) UnmarshalJSON(b []byte) error {

	return nil
}
