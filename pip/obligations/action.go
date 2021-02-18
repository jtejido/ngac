package obligations

type Action struct {
	Condition        *Condition
	NegatedCondition *NegatedCondition
}

// AssignAction.java
type AssignAction struct {
	Action
	Assignments []*ActionAssignment
}

func NewAssignAction() *AssignAction {
	ans := new(AssignAction)
	ans.Assignments = make([]*ActionAssignment, 0)
	return ans
}

type ActionAssignment struct {
	What, Where *EvrNode
}

// CreateAction.java
type CreateAction struct {
	Action
	CreateNodesList []*ActionCreateNode
	Rules           []*Rule
}

func NewCreateAction() *CreateAction {
	ans := new(CreateAction)
	ans.CreateNodesList = make([]*ActionCreateNode, 0)
	ans.Rules = make([]*Rule, 0)
	return ans
}

type ActionCreateNode struct {
	What, Where *EvrNode
}

type DeleteAction struct {
	Action
	Nodes               []*EvrNode
	Assignments         *AssignAction
	Associations        []*GrantAction
	Prohibitions, Rules []string
}

// DenyAction.java
type DenyAction struct {
	Action
	Label      string
	Subject    *EvrNode
	Operations []string
	Target     *ActionTarget
}

type ActionTarget struct {
	Complement, Intersection bool
	Containers               []*ActionContainer
}

type ActionContainer struct {
	Name, Type string
	Properties map[string]string
	Function   *Function
	Complement bool
}

func NewActionContainerFromFunction(function *Function) *ActionContainer {
	return &ActionContainer{Function: function}
}

func NewActionContainer(name, t string, properties map[string]string) *ActionContainer {
	return &ActionContainer{
		Name:       name,
		Type:       t,
		Properties: properties,
	}

}

type FunctionAction struct {
	Action
	Function *Function
}

func NewFunctionAction(function *Function) *FunctionAction {
	ans := new(FunctionAction)
	ans.Function = function
	return ans
}

type GrantAction struct {
	Action
	Subject    *EvrNode
	Operations []string
	Target     *EvrNode
}

func NewGrantAction() *GrantAction {
	ans := new(GrantAction)
	ans.Operations = make([]string, 0)
	return ans
}
