package epp

import (
	"fmt"
	"github.com/jtejido/ngac/pkg/operations"
	"github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/jtejido/ngac/pkg/pip/obligations"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type EPP interface {
	ProcessEvent(eventCtx EventContext) error
}

func Apply(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, rule *obligations.Rule, obligationLabel string) error {
	// check the response condition
	responsePattern := rule.ResponsePattern
	cc, err := checkCondition(g, p, o, functionEvaluator, responsePattern.Condition, eventCtx)
	if err != nil {
		return err
	}
	nc, err := checkNegatedCondition(g, p, o, functionEvaluator, responsePattern.NegatedCondition, eventCtx)
	if err != nil {
		return err
	}
	if !cc || !nc {
		return nil
	}

	for _, action := range rule.ResponsePattern.Actions {
		cc, err := checkCondition(g, p, o, functionEvaluator, action.Condition(), eventCtx)
		if err != nil {
			return err
		}
		nc, err := checkNegatedCondition(g, p, o, functionEvaluator, action.NegatedCondition(), eventCtx)
		if err != nil {
			return err
		}
		if !cc {
			continue
		} else if !nc {
			continue
		}

		if err := applyAction(g, p, o, functionEvaluator, obligationLabel, eventCtx, action); err != nil {
			return err
		}
	}

	return nil
}

func checkCondition(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, condition *obligations.Condition, eventCtx EventContext) (bool, error) {
	if condition == nil {
		return true, nil
	}

	functions := condition.Condition
	for _, f := range functions {
		result, err := functionEvaluator.Eval(g, p, o, eventCtx, f)
		if err != nil {
			return false, err
		}
		if !result.(bool) {
			return false, nil
		}
	}

	return true, nil
}

/**
 * Return true if the condition is satisfied. A condition is satisfied if all the functions evaluate to false.
 */
func checkNegatedCondition(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, condition *obligations.NegatedCondition, eventCtx EventContext) (bool, error) {
	if condition == nil {
		return true, nil
	}

	functions := condition.Condition
	for _, f := range functions {
		result, err := functionEvaluator.Eval(g, p, o, eventCtx, f)
		if err != nil {
			return false, err
		}
		if result.(bool) {
			return false, err
		}
	}

	return true, nil
}

func applyAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, label string, eventCtx EventContext, action obligations.Action) error {
	if action == nil {
		return nil
	}

	if v, ok := action.(*obligations.AssignAction); ok {
		return applyAssignAction(g, p, o, functionEvaluator, eventCtx, v)
	} else if v, ok := action.(*obligations.CreateAction); ok {
		return applyCreateAction(g, p, o, functionEvaluator, label, eventCtx, v)
	} else if v, ok := action.(*obligations.DeleteAction); ok {
		return applyDeleteAction(g, p, o, functionEvaluator, eventCtx, v)
	} else if v, ok := action.(*obligations.DenyAction); ok {
		return applyDenyAction(g, p, o, functionEvaluator, eventCtx, v)
	} else if v, ok := action.(*obligations.GrantAction); ok {
		return applyGrantAction(g, p, o, functionEvaluator, eventCtx, v)
	} else if v, ok := action.(*obligations.FunctionAction); ok {
		_, err := functionEvaluator.Eval(g, p, o, eventCtx, v.Function)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyGrantAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, action *obligations.GrantAction) error {
	subject := action.Subject
	op := action.Operations
	target := action.Target

	subjectNode, err := toNode(g, p, o, functionEvaluator, eventCtx, subject)
	if err != nil {
		return err
	}
	targetNode, err := toNode(g, p, o, functionEvaluator, eventCtx, target)
	if err != nil {
		return err
	}

	return g.Associate(subjectNode.Name, targetNode.Name, operations.NewOperationSet(op))
}

func applyDenyAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, action *obligations.DenyAction) error {
	subject := action.Subject
	ops := action.Operations
	target := action.Target

	denySubject, err := toDenySubject(g, p, o, functionEvaluator, eventCtx, subject)
	if err != nil {
		return err
	}
	denyNodes, err := toDenyNodes(g, p, o, functionEvaluator, eventCtx, target)
	if err != nil {
		return err
	}

	builder := prohibitions.NewBuilder(action.Label, denySubject, operations.NewOperationSet(ops))
	builder.Intersection = target.Intersection
	for contName, v := range denyNodes {
		builder.AddContainer(contName, v)
	}

	// add the prohibition to the PAP
	p.Add(builder.Build())
	return nil

	// TODO this complement is ignored in the current Prohibition object
	// complement := target.Complement
}

func toDenyNodes(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, target *obligations.ActionTarget) (map[string]bool, error) {
	nodes := make(map[string]bool)
	containers := target.Containers
	for _, container := range containers {
		if container.Function != nil {
			function := container.Function
			result, err := functionEvaluator.Eval(g, p, o, eventCtx, function)
			if err != nil {
				return nil, err
			}

			if cc, ok := result.(*prohibitions.ContainerCondition); ok {
				nodes[cc.Name()] = cc.IsComplement()
			} else {
				return nil, fmt.Errorf("expected function to return a ContainerCondition.")
			}
		} else {
			// get the node
			node, err := g.Node(container.Name)
			if err != nil {
				return nil, err
			}
			nodes[node.Name] = container.Complement
		}
	}

	return nodes, nil
}

func toDenySubject(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, subject *obligations.EvrNode) (string, error) {
	var denySubject string

	if subject.Function != nil {
		function := subject.Function
		t, err := functionEvaluator.Eval(g, p, o, eventCtx, function)
		if err != nil {
			return "", err
		}
		denySubject = t.(string)
	} else if subject.Process != nil {
		denySubject = subject.Process.Value
	} else {
		// if (subject.getName() != null) {
		if len(subject.Name) > 0 {
			t, err := g.Node(subject.Name)
			if err != nil {
				return "", err
			}
			denySubject = t.Name
		} else {
			t, err := g.NodeFromDetails(graph.ToNodeType(subject.Type), subject.Properties)
			if err != nil {
				return "", err
			}
			denySubject = t.Name
		}
	}

	return denySubject, nil
}

func applyDeleteAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, action *obligations.DeleteAction) error {
	nodes := action.Nodes
	if nodes != nil {
		for _, evrNode := range nodes {
			node, err := toNode(g, p, o, functionEvaluator, eventCtx, evrNode)
			if err != nil {
				return err
			}
			g.RemoveNode(node.Name)
		}
	}

	assignAction := action.Assignments
	if assignAction != nil {
		for _, assignment := range assignAction.Assignments {
			what, err := toNode(g, p, o, functionEvaluator, eventCtx, assignment.What)
			if err != nil {
				return err
			}
			where, err := toNode(g, p, o, functionEvaluator, eventCtx, assignment.Where)
			if err != nil {
				return err
			}
			if err := g.Deassign(what.Name, where.Name); err != nil {
				return err
			}
		}
	}

	associations := action.Associations
	if associations != nil {
		for _, grantAction := range associations {
			subject, err := toNode(g, p, o, functionEvaluator, eventCtx, grantAction.Subject)
			if err != nil {
				return err
			}
			target, err := toNode(g, p, o, functionEvaluator, eventCtx, grantAction.Target)
			if err != nil {
				return err
			}
			if err := g.Dissociate(subject.Name, target.Name); err != nil {
				return err
			}
		}
	}

	actionProhibitions := action.Prohibitions
	if p != nil {
		for _, label := range actionProhibitions {
			p.Remove(label)
		}
	}

	rules := action.Rules
	if rules != nil {
		for _, label := range rules {
			allObligations := o.All()
			for _, obligation := range allObligations {
				oblRules := obligation.Rules
				for i := 0; i < len(oblRules); i++ {
					if oblRules[i].Label == label {
						oblRules = append(oblRules[:i], oblRules[i+1:]...)
					}
				}
			}
		}
	}

	return nil
}

func toNode(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, evrNode *obligations.EvrNode) (node *graph.Node, err error) {
	if evrNode.Function != nil {
		n, err := functionEvaluator.Eval(g, p, o, eventCtx, evrNode.Function)
		if err != nil {
			return nil, err
		}
		node = n.(*graph.Node)
	} else {
		if len(evrNode.Name) > 0 {
			node, err = g.Node(evrNode.Name)
			if err != nil {
				return
			}
		} else {
			node, err = g.NodeFromDetails(graph.ToNodeType(evrNode.Type), evrNode.Properties)
			if err != nil {
				return
			}
		}
	}
	return node, nil
}

func applyCreateAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, label string, eventCtx EventContext, action *obligations.CreateAction) (err error) {
	rules := action.Rules
	if rules != nil {
		for _, rule := range rules {
			createRule(g, p, o, label, eventCtx, rule)
		}
	}

	createNodesList := action.CreateNodesList
	if createNodesList != nil {
		for _, createNode := range createNodesList {
			what := createNode.What
			where := createNode.Where
			whereNode, err := toNode(g, p, o, functionEvaluator, eventCtx, where)
			if err != nil {
				return err
			}
			_, err = g.CreateNode(what.Name, graph.ToNodeType(what.Type), what.Properties, whereNode.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createRule(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, obligationLabel string, eventCtx EventContext, rule *obligations.Rule) {
	// add the rule to the obligation
	obligation := o.Get(obligationLabel)
	rules := obligation.Rules
	rules = append(rules, rule)
	obligation.Rules = rules
	o.Update(obligationLabel, obligation)
}

func applyAssignAction(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations, functionEvaluator *FunctionEvaluator, eventCtx EventContext, action *obligations.AssignAction) error {
	assignments := action.Assignments
	if assignments != nil {
		for _, assignment := range assignments {
			what := assignment.What
			where := assignment.Where

			whatNode, err := toNode(g, p, o, functionEvaluator, eventCtx, what)
			if err != nil {
				return err
			}
			whereNode, err := toNode(g, p, o, functionEvaluator, eventCtx, where)
			if err != nil {
				return err
			}
			if err := g.Assign(whatNode.Name, whereNode.Name); err != nil {
				return err
			}
		}
	}

	return nil
}
