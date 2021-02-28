package obligations

import (
	"encoding/json"
	"fmt"
)

type Function struct {
	Name string `json:"name"`
	Args []*Arg `json:"args"`
}

func NewFunction(name string, args []*Arg) *Function {
	return &Function{name, args}
}

func (f *Function) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	json.Unmarshal(b, &raw)
	var name string
	if v, ok := raw["name"]; ok {
		name = v.(string)
	}

	args := make([]*Arg, 0)
	if v, ok := raw["args"]; ok {
		for _, val := range v.([]interface{}) {
			if s, ok := val.(string); ok {
				args = append(args, NewArg(s))
			} else if m, ok := val.(map[string]interface{}); ok {
				if _, ok = m["function"]; ok {
					nf := new(Function)
					b, err := json.Marshal(m)
					if err != nil {
						return err
					}

					err = nf.UnmarshalJSON(b)
					if err != nil {
						return err
					}

					args = append(args, NewArgFromFunction(nf))
				}
			} else {
				return fmt.Errorf("invalid function definition")
			}

		}

	}

	a := NewFunction(name, args)
	*f = *a
	return nil
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
