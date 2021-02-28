package obligations

import (
	"encoding/json"
	"fmt"
	"github.com/jtejido/ngac/pip/graph"
)

type Obligation struct {
	User    string
	Enabled bool
	Label   string  `json:"label" yaml:"label"`
	Rules   []*Rule `json:"rules, omitempty" yaml:"rules, omitempty"`
	Source  string
}

func NewObligation(user string) *Obligation {
	return &Obligation{User: user, Rules: make([]*Rule, 0)}
}

func (ob *Obligation) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	if v, ok := raw["label"]; ok {
		if v.(string) == "" {
			return fmt.Errorf("no label specified for obligation")
		}

		ob.Label = v.(string)
	} else {
		return fmt.Errorf("no label specified for obligation")
	}

	if v, ok := raw["rules"]; ok {
		ob.Rules = make([]*Rule, len(v.([]interface{})))
		for i, _ := range v.([]interface{}) {
			b, err := json.Marshal(v.([]interface{})[i])
			if err != nil {
				return err
			}
			ob.Rules[i] = new(Rule)
			err = ob.Rules[i].UnmarshalJSON(b)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ob *Obligation) Clone() *Obligation {
	return &Obligation{ob.User, ob.Enabled, ob.Label, append([]*Rule{}, ob.Rules...), ob.Source}
}

type Rule struct {
	Label           string           `json:"label" yaml:"label"`
	EventPattern    *EventPattern    `json:"event" yaml:"event"`
	ResponsePattern *ResponsePattern `json:"response" yaml:"response"`
}

func (r *Rule) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["label"]; ok {
		if v.(string) == "" {
			return fmt.Errorf("no label provided for rule")
		}

		r.Label = v.(string)
	} else {
		return fmt.Errorf("no label provided for rule")
	}

	if v, ok := raw["event"]; ok {
		r.EventPattern = new(EventPattern)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = r.EventPattern.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no event provided for rule")
	}

	if v, ok := raw["response"]; ok {
		r.ResponsePattern = new(ResponsePattern)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = r.ResponsePattern.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no response provided for rule")
	}

	return nil
}

func NewRule() *Rule {
	return &Rule{"", new(EventPattern), NewResponsePattern()}
}

type EventPattern struct {
	Subject     *Subject     `json:"subject, omitempty" yaml:"subject, omitempty"`
	PolicyClass *PolicyClass `json:"policyClass, omitempty" yaml:"policyClass, omitempty"`
	Operations  []string     `json:"operations, omitempty" yaml:"operations, omitempty"`
	Target      *Target      `json:"target, omitempty" yaml:"target, omitempty"`
}

func (e *EventPattern) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["subject"]; ok {
		e.Subject = new(Subject)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = e.Subject.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["policyClass"]; ok {
		e.PolicyClass = new(PolicyClass)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = e.PolicyClass.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["operations"]; ok {
		e.Operations = make([]string, len(v.([]interface{})))
		for i, ops := range v.([]interface{}) {
			e.Operations[i] = ops.(string)
		}
	}

	if v, ok := raw["target"]; ok {
		e.Target = new(Target)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = e.Target.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	return nil
}

type ResponsePattern struct {
	Condition        *Condition        `json:"condition, omitempty" yaml:"condition, omitempty"`
	NegatedCondition *NegatedCondition `json:"not_condition, omitempty" yaml:"condition!, omitempty"`
	Actions          []Action          `json:"actions, omitempty" yaml:"actions, omitempty"`
}

func NewResponsePattern() *ResponsePattern {
	return &ResponsePattern{Actions: make([]Action, 0)}
}

func (r *ResponsePattern) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["condition"]; ok {
		r.Condition = new(Condition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = r.Condition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["not_condition"]; ok {
		r.NegatedCondition = new(NegatedCondition)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = r.NegatedCondition.UnmarshalJSON(b)
		if err != nil {
			return err
		}
	}

	if v, ok := raw["actions"]; ok {
		r.Actions = make([]Action, len(v.([]interface{})))
		for i, act := range v.([]interface{}) {
			if _, ok2 := act.(map[string]interface{})["function"]; ok2 {
				r.Actions[i] = new(FunctionAction)
			} else if _, ok2 = act.(map[string]interface{})["create"]; ok2 {
				r.Actions[i] = new(CreateAction)
			} else if _, ok2 = act.(map[string]interface{})["assign"]; ok2 {
				r.Actions[i] = new(AssignAction)
			} else if _, ok2 = act.(map[string]interface{})["deny"]; ok2 {
				r.Actions[i] = new(DenyAction)
			} else if _, ok2 = act.(map[string]interface{})["grant"]; ok2 {
				r.Actions[i] = new(GrantAction)
			} else if _, ok2 = act.(map[string]interface{})["delete"]; ok2 {
				r.Actions[i] = new(DeleteAction)
			}

			b, err := json.Marshal(act)
			if err != nil {
				return err
			}

			err = r.Actions[i].UnmarshalJSON(b)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Condition struct {
	Condition []*Function `json:"function"`
}

func (c *Condition) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	c.Condition = make([]*Function, 0)
	if v, ok := raw["function"]; ok {
		for i, _ := range v.([]interface{}) {
			b, err := json.Marshal(v.([]interface{})[i])
			if err != nil {
				return err
			}
			nf := new(Function)
			err = nf.UnmarshalJSON(b)
			if err != nil {
				return err
			}
			c.Condition = append(c.Condition, nf)
		}
	}

	return nil
}

type NegatedCondition struct {
	Condition []*Function `json:"function"`
}

func (c *NegatedCondition) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	c.Condition = make([]*Function, 0)
	if v, ok := raw["function"]; ok {
		for i, _ := range v.([]interface{}) {
			b, err := json.Marshal(v.([]interface{})[i])
			if err != nil {
				return err
			}
			nf := new(Function)
			err = nf.UnmarshalJSON(b)
			if err != nil {
				return err
			}
			c.Condition = append(c.Condition, nf)
		}
	}

	return nil
}

type Subject struct {
	User    string      `json:"user, omitempty" yaml:"user, omitempty"`
	AnyUser []string    `json:"anyUser, omitempty" yaml:"anyUser, omitempty"`
	Process *EvrProcess `json:"process, omitempty" yaml:"process, omitempty"`
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

func (s *Subject) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["user"]; ok {
		s.User = v.(string)
		return nil
	} else if v, ok = raw["anyUser"]; ok {
		s.AnyUser = make([]string, len(v.([]interface{})))
		for i, u := range v.([]interface{}) {
			s.AnyUser[i] = u.(string)
		}
		return nil
	} else if v, ok = raw["user"]; ok {
		s.Process = NewEvrProcess(v.(string))
		return nil
	} else {
		return fmt.Errorf("invalid subject specification")
	}

	return nil
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
	AnyOf  []string `json:"anyOf" yaml:"anyOf"`
	EachOf []string `json:"eachOf" yaml:"eachOf"`
}

func NewPolicyClass() *PolicyClass {
	return &PolicyClass{make([]string, 0), make([]string, 0)}
}

func (pc *PolicyClass) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["anyOf"]; ok {
		pc.AnyOf = make([]string, len(v.([]interface{})))
		for i, u := range v.([]interface{}) {
			pc.AnyOf[i] = u.(string)
		}

		return nil
	}

	if v, ok := raw["eachOf"]; ok {
		pc.EachOf = make([]string, len(v.([]interface{})))
		for i, u := range v.([]interface{}) {
			pc.EachOf[i] = u.(string)
		}

		return nil
	}

	return fmt.Errorf("expected one of (anyOf, eachOf)")
}

type Target struct {
	PolicyElements []*EvrNode `json:"policyElements" yaml:"policyElements"`
	Containers     []*EvrNode `json:"containers" yaml:"containers"`
}

func (t *Target) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["policyElements"]; ok {
		t.PolicyElements = make([]*EvrNode, len(v.([]interface{})))
		for i, _ := range v.([]interface{}) {
			b, err := json.Marshal(v.([]interface{})[i])
			if err != nil {
				return err
			}
			t.PolicyElements[i] = new(EvrNode)
			err = t.PolicyElements[i].UnmarshalJSON(b)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if v, ok := raw["containers"]; ok {
		t.Containers = make([]*EvrNode, len(v.([]interface{})))
		for i, _ := range v.([]interface{}) {
			b, err := json.Marshal(v.([]interface{})[i])
			if err != nil {
				return err
			}
			t.Containers[i] = new(EvrNode)
			err = t.Containers[i].UnmarshalJSON(b)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// should there be an error when no target is specified?
	return nil
}

type EvrNode struct {
	Name       string `json:"name" yaml:"name"`
	Type       string `json:"type" yaml:"type"`
	Properties graph.PropertyMap
	Function   *Function `json:"function" yaml:"function"`
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

func (evr *EvrNode) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)

	if v, ok := raw["function"]; ok {
		evr.Function = new(Function)
		b, err := json.Marshal(v.(interface{}))
		if err != nil {
			return err
		}
		err = evr.Function.UnmarshalJSON(b)
		if err != nil {
			return err
		}

		return nil
	}

	if v, ok := raw["name"]; ok {
		if len(v.(string)) == 0 {
			return fmt.Errorf("name cannot be empty")
		}

		evr.Name = v.(string)
		if v, ok = raw["type"]; ok {
			if len(v.(string)) == 0 {
				return fmt.Errorf("type cannot be empty")
			}

			evr.Type = v.(string)
		} else {
			return fmt.Errorf("type cannot be empty")
		}

		return nil
	}

	return fmt.Errorf("invalid EVR node")
}

type Containers struct {
	AnyOf  []*EvrNode `json:"anyOf" yaml:"anyOf"`
	EachOf []*EvrNode `json:"eachOf" yaml:"eachOf"`
}
